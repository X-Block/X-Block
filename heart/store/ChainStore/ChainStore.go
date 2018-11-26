package ChainStore

import (
	. "XBlock/common"
	"XBlock/common/log"
	"XBlock/common/serialization"
	. "XBlock/core/asset"
	"XBlock/core/contract/program"
	. "XBlock/core/ledger"
	. "XBlock/core/store"
	. "XBlock/core/store/LevelDBStore"
	tx "XBlock/core/transaction"
	"XBlock/core/transaction/payload"
	"XBlock/core/validation"
	"XBlock/crypto"
	. "XBlock/errors"
	"XBlock/events"
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"XBlock/core/account"
)


var (
	ErrDBNotFound = errors.New("leveldb: not found")
)

type ChainStore struct {
	st IStore

	headerIndex map[uint32]Uint256
	blockCache  map[Uint256]*Block
	headerCache map[Uint256]*Header

	currentBlockHeight uint32
	storedHeaderCount  uint32

	mu sync.RWMutex

	disposed bool
}

func init() {
}

func NewStore() IStore {
	ldbs, _ := NewLevelDBStore("Chain")

	return ldbs
}

func NewLedgerStore() ILedgerStore {
	cs, _ := NewChainStore("Chain")

	return cs
}

func NewChainStore(file string) (*ChainStore, error) {

	return &ChainStore{
		st:                 NewStore(),
		headerIndex:        map[uint32]Uint256{},
		blockCache:         map[Uint256]*Block{},
		headerCache:        map[Uint256]*Header{},
		currentBlockHeight: 0,
		storedHeaderCount:  0,
		disposed:           false,
	}, nil
}

func (bd *ChainStore) InitLedgerStoreWithGenesisBlock(genesisBlock *Block, defaultBookKeeper []*crypto.PubKey) (uint32, error) {

	hash := genesisBlock.Hash()
	bd.headerIndex[0] = hash
	log.Debug(fmt.Sprintf("listhash genesis: %x\n", hash))

	prefix := []byte{byte(IX_Version)}
	version, err := bd.st.Get(prefix)
	if err != nil {
		version = []byte{0x00}
	}

	if version[0] == 0x01 {
		currentBlockPrefix := []byte{byte(SYS_CurrentBlock)}
		data, err := bd.st.Get(currentBlockPrefix)
		if err != nil {
			return 0, err
		}

		r := bytes.NewReader(data)
		var blockHash Uint256
		blockHash.Deserialize(r)
		bd.currentBlockHeight, err = serialization.ReadUint32(r)
		current_Header_Height := bd.currentBlockHeight


		var headerHash Uint256
		currentHeaderPrefix := []byte{byte(SYS_CurrentHeader)}
		data, err = bd.st.Get(currentHeaderPrefix)
		if err == nil {
			r = bytes.NewReader(data)
			headerHash.Deserialize(r)

			headerHeight, err_get := serialization.ReadUint32(r)
			if err_get != nil {
				return 0, err_get
			}

			current_Header_Height = headerHeight
		}

		log.Debug(fmt.Sprintf("blockHash: %x\n", blockHash.ToArray()))
		log.Debug(fmt.Sprintf("blockheight: %d\n", current_Header_Height))

		var listHash Uint256
		iter := bd.st.NewIterator([]byte{byte(IX_HeaderHashList)})
		for iter.Next() {
			rk := bytes.NewReader(iter.Key())
			_, _ = serialization.ReadBytes(rk, 1)
			startNum, err := serialization.ReadUint32(rk)
			if err != nil {
				return 0, err
			}
			log.Debug(fmt.Sprintf("start index: %d\n", startNum))

			r = bytes.NewReader(iter.Value())
			listNum, err := serialization.ReadVarUint(r, 0)
			if err != nil {
				return 0, err
			}

			for i := 0; i < int(listNum); i++ {
				listHash.Deserialize(r)
				bd.headerIndex[startNum+uint32(i)] = listHash
				bd.storedHeaderCount++
			}
		}

		if bd.storedHeaderCount == 0 {
			iter = bd.st.NewIterator([]byte{byte(DATA_BlockHash)})
			for iter.Next() {
				rk := bytes.NewReader(iter.Key())				_, _ = serialization.ReadBytes(rk, 1)
				listheight, err := serialization.ReadUint32(rk)
				if err != nil {
					return 0, err
				}

				r := bytes.NewReader(iter.Value())
				listHash.Deserialize(r)


				bd.headerIndex[listheight] = listHash
			}
		} else if current_Header_Height >= bd.storedHeaderCount {
			hash = headerHash
			for {
				if hash == bd.headerIndex[bd.storedHeaderCount-1] {
					break
				}

				header, err := bd.GetHeader(hash)
				if err != nil {
					return 0, err
				}


				bd.headerIndex[header.Blockdata.Height] = hash
				hash = header.Blockdata.PrevBlockHash
			}
		}

		return current_Header_Height, nil

	} else {

		bd.st.NewBatch()
		iter := bd.st.NewIterator(nil)
		for iter.Next() {
			bd.st.BatchDelete(iter.Key())
		}
		iter.Release()

		err := bd.st.BatchCommit()
		if err != nil {
			return 0, err
		}

		sort.Sort(crypto.PubKeySlice(defaultBookKeeper))


		bkListKey := bytes.NewBuffer(nil)
		bkListKey.WriteByte(byte(SYS_CurrentBookKeeper))


		bkListValue := bytes.NewBuffer(nil)
		serialization.WriteUint8(bkListValue, uint8(len(defaultBookKeeper)))
		for k := 0; k < len(defaultBookKeeper); k++ {
			defaultBookKeeper[k].Serialize(bkListValue)
		}


		serialization.WriteUint8(bkListValue, uint8(len(defaultBookKeeper)))
		for k := 0; k < len(defaultBookKeeper); k++ {
			defaultBookKeeper[k].Serialize(bkListValue)
		}

		bd.persist(genesisBlock)
		err = bd.st.Put(prefix, []byte{0x01})
		if err != nil {
			return 0, err
		}

		return 0, nil
	}
}

func (bd *ChainStore) InitLedgerStore(l *Ledger) error {
	return nil
}

func (bd *ChainStore) IsDoubleSpend(tx *tx.Transaction) bool {
	if len(tx.UTXOInputs) == 0 {
		return false
	}

	unspentPrefix := []byte{byte(IX_Unspent)}
	for i := 0; i < len(tx.UTXOInputs); i++ {
		txhash := tx.UTXOInputs[i].ReferTxID
		unspentValue, err_get := bd.st.Get(append(unspentPrefix, txhash.ToArray()...))
		if err_get != nil {
			return true
		}

		unspents, _ := GetUint16Array(unspentValue)
		findFlag := false
		for k := 0; k < len(unspents); k++ {
			if unspents[k] == tx.UTXOInputs[i].ReferTxOutputIndex {
				findFlag = true
				break
			}
		}

		if !findFlag {
			return true
		}
	}

	return false
}

func (bd *ChainStore) GetBlockHash(height uint32) (Uint256, error) {

	if height >= 0 {
		queryKey := bytes.NewBuffer(nil)
		queryKey.WriteByte(byte(DATA_BlockHash))
		err := serialization.WriteUint32(queryKey, height)

		if err == nil {
			blockHash, err_get := bd.st.Get(queryKey.Bytes())
			if err_get != nil {
				return Uint256{}, err_get
			} else {
				blockHash256, err_parse := Uint256ParseFromBytes(blockHash)
				if err_parse == nil {
					return blockHash256, nil
				} else {
					return Uint256{}, err_parse
				}

			}
		} else {
			return Uint256{}, err
		}
	} else {
		return Uint256{}, NewDetailErr(errors.New("[LevelDB]: GetBlockHash error param height < 0."), ErrNoCode, "")
	}
}

func (bd *ChainStore) GetCurrentBlockHash() Uint256 {
	bd.mu.RLock()
	defer bd.mu.RUnlock()

	return bd.headerIndex[bd.currentBlockHeight]
}

func (bd *ChainStore) GetContract(hash []byte) ([]byte, error) {
	prefix := []byte{byte(DATA_Contract)}
	bData, err_get := bd.st.Get(append(prefix, hash...))
	if err_get != nil {
		return nil, err_get
	}

	log.Debug("GetContract Data: ", bData)

	return bData, nil
}

func (bd *ChainStore) GetHeaderWithCache(hash Uint256) *Header {
	if _, ok := bd.headerCache[hash]; ok {
		return bd.headerCache[hash]
	}

	header, _ := bd.GetHeader(hash)

	return header
}

func (bd *ChainStore) containsBlock(hash Uint256) bool {
	header := bd.GetHeaderWithCache(hash)
	if header != nil {
		return header.Blockdata.Height <= bd.currentBlockHeight
	} else {
		return false
	}
}

func (bd *ChainStore) VerifyHeader(header *Header) bool {
	if bd.containsBlock(header.Blockdata.Hash()) {
		return true
	}

	prevHeader := bd.GetHeaderWithCache(header.Blockdata.PrevBlockHash)

	if prevHeader == nil {
		log.Error("[VerifyHeader] failed, not found prevHeader.")
		return false
	}

	if prevHeader.Blockdata.Height+1 != header.Blockdata.Height {
		log.Error("[VerifyHeader] failed, prevHeader.Height + 1 != header.Height")
		return false
	}

	if prevHeader.Blockdata.Timestamp >= header.Blockdata.Timestamp {
		log.Error("[VerifyHeader] failed, prevHeader.Timestamp >= header.Timestamp")
		return false
	}

	flag, err := validation.VerifySignableData(header.Blockdata)
	if flag == false || err != nil {
		log.Error("[VerifyHeader] failed, VerifySignableData failed.")
		log.Error(err)
		return false
	}

	return true
}

