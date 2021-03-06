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

func (bd *ChainStore) AddHeaders(headers []Header, ledger *Ledger) error {
	bd.mu.Lock()
	defer bd.mu.Unlock()

	for i := 0; i < len(headers); i++ {
		if headers[i].Blockdata.Height >= (uint32(len(bd.headerIndex)) + 1) {
			break
		}

		if headers[i].Blockdata.Height < uint32(len(bd.headerIndex)) {
			continue
		}


		if !bd.VerifyHeader(&headers[i]) {
			log.Error("Verify header failed")
			break
		}

		bd.st.NewBatch()

		bd.addHeader(&headers[i])

		err := bd.st.BatchCommit()
		if err != nil {
			return err
		}


		bd.headerCache[headers[i].Blockdata.Hash()] = &headers[i]
	}

	for k, _ := range bd.headerCache {
		delete(bd.headerCache, k)
	}

	return nil
}

func (bd *ChainStore) GetHeader(hash Uint256) (*Header, error) {
	var h *Header = new(Header)

	h.Blockdata = new(Blockdata)
	h.Blockdata.Program = new(program.Program)

	prefix := []byte{byte(DATA_Header)}
	log.Debug("GetHeader Data:", hash.ToArray())
	data, err_get := bd.st.Get(append(prefix, hash.ToArray()...))
	if err_get != nil {

		return nil, err_get
	}

	r := bytes.NewReader(data)
	systemfee, err := serialization.ReadUint64(r)
	if err != nil {
		return nil, err
	}
	log.Debug(fmt.Sprintf("systemfee: %d\n", systemfee))


	err = h.Deserialize(r)
	if err != nil {
		return nil, err
	}

	return h, err
}

func (bd *ChainStore) SaveAsset(assetId Uint256, asset *Asset) error {
	w := bytes.NewBuffer(nil)

	asset.Serialize(w)


	assetKey := bytes.NewBuffer(nil)

	assetKey.WriteByte(byte(ST_Info))

	assetId.Serialize(assetKey)

	log.Debug(fmt.Sprintf("asset key: %x\n", assetKey))
	err := bd.st.Put(assetKey.Bytes(), w.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (bd *ChainStore) GetAsset(hash Uint256) (*Asset, error) {
	log.Debug(fmt.Sprintf("GetAsset Hash: %x\n", hash))

	asset := new(Asset)

	prefix := []byte{byte(ST_Info)}
	data, err_get := bd.st.Get(append(prefix, hash.ToArray()...))

	log.Debug(fmt.Sprintf("GetAsset Data: %x\n", data))
	if err_get != nil {

		return nil, err_get
	}

	r := bytes.NewReader(data)
	asset.Deserialize(r)

	return asset, nil
}

func (bd *ChainStore) GetTransaction(hash Uint256) (*tx.Transaction, error) {
	log.Debug()
	log.Debug(fmt.Sprintf("GetTransaction Hash: %x\n", hash))

	t := new(tx.Transaction)
	err := bd.getTx(t, hash)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (bd *ChainStore) getTx(tx *tx.Transaction, hash Uint256) error {
	prefix := []byte{byte(DATA_Transaction)}
	tHash, err_get := bd.st.Get(append(prefix, hash.ToArray()...))
	if err_get != nil {

		return err_get
	}

	r := bytes.NewReader(tHash)


	_, err := serialization.ReadUint32(r)
	if err != nil {
		return err
	}


	err = tx.Deserialize(r)

	return err
}

func (bd *ChainStore) SaveTransaction(tx *tx.Transaction, height uint32) error {

	txhash := bytes.NewBuffer(nil)

	txhash.WriteByte(byte(DATA_Transaction))

	txHashValue := tx.Hash()
	txHashValue.Serialize(txhash)
	log.Debug(fmt.Sprintf("transaction header + hash: %x\n", txhash))


	w := bytes.NewBuffer(nil)
	serialization.WriteUint32(w, height)
	tx.Serialize(w)
	log.Debug(fmt.Sprintf("transaction tx data: %x\n", w))
	err := bd.st.Put(txhash.Bytes(), w.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (bd *ChainStore) GetBlock(hash Uint256) (*Block, error) {
	var b *Block = new(Block)

	b.Blockdata = new(Blockdata)
	b.Blockdata.Program = new(program.Program)

	prefix := []byte{byte(DATA_Header)}
	bHash, err_get := bd.st.Get(append(prefix, hash.ToArray()...))
	if err_get != nil {

		return nil, err_get
	}

	r := bytes.NewReader(bHash)


	_, err := serialization.ReadUint64(r)
	if err != nil {
		return nil, err
	}


	err = b.FromTrimmedData(r)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(b.Transactions); i++ {
		err = bd.getTx(b.Transactions[i], b.Transactions[i].Hash())
		if err != nil {
			return nil, err
		}
	}

	return b, nil
}

func (self *ChainStore) GetBookKeeperList() ([]*crypto.PubKey, []*crypto.PubKey, error) {
	prefix := []byte{byte(SYS_CurrentBookKeeper)}
	bkListValue, err_get := self.st.Get(prefix)
	if err_get != nil {
		return nil, nil, err_get
	}

	r := bytes.NewReader(bkListValue)

	currCount, err := serialization.ReadUint8(r)
	if err != nil {
		return nil, nil, err
	}

	var currBookKeeper = make([]*crypto.PubKey, currCount)
	for i := uint8(0); i < currCount; i++ {
		bk := new(crypto.PubKey)
		err := bk.DeSerialize(r)
		if err != nil {
			return nil, nil, err
		}

		currBookKeeper[i] = bk
	}

	nextCount, err := serialization.ReadUint8(r)
	if err != nil {
		return nil, nil, err
	}

	var nextBookKeeper = make([]*crypto.PubKey, nextCount)
	for i := uint8(0); i < nextCount; i++ {
		bk := new(crypto.PubKey)
		err := bk.DeSerialize(r)
		if err != nil {
			return nil, nil, err
		}

		nextBookKeeper[i] = bk
	}

	return currBookKeeper, nextBookKeeper, nil
}

func (bd *ChainStore) persist(b *Block) error {
	unspents := make(map[Uint256][]uint16)
	quantities := make(map[Uint256]Fixed64)

	unspentPrefix := []byte{byte(IX_Unspent)}
	accounts := make(map[Uint160]*account.AccountState, 0)

	bd.st.NewBatch()

	bhhash := bytes.NewBuffer(nil)

	bhhash.WriteByte(byte(DATA_Header))

	blockHash := b.Hash()
	blockHash.Serialize(bhhash)
	log.Debug(fmt.Sprintf("block header + hash: %x\n", bhhash))


	bd.st.BatchPut(bhhash.Bytes(), w.Bytes())

	bhash := bytes.NewBuffer(nil)
	bhash.WriteByte(byte(DATA_BlockHash))
	err := serialization.WriteUint32(bhash, b.Blockdata.Height)
	if err != nil {
		return err
	}
	log.Debug(fmt.Sprintf("DATA_BlockHash table key: %x\n", bhash))

	hashWriter := bytes.NewBuffer(nil)
	hashValue := b.Blockdata.Hash()
	hashValue.Serialize(hashWriter)
	log.Debug(fmt.Sprintf("DATA_BlockHash table value: %x\n", hashValue))

	needUpdateBookKeeper := false
	currBookKeeper, nextBookKeeper, err := bd.GetBookKeeperList()

	if len(currBookKeeper) != len(nextBookKeeper) {
		needUpdateBookKeeper = true
	} else {
		for i, _ := range currBookKeeper {
			if currBookKeeper[i].X.Cmp(nextBookKeeper[i].X) != 0 ||
				currBookKeeper[i].Y.Cmp(nextBookKeeper[i].Y) != 0 {
				needUpdateBookKeeper = true
				break
			}
		}
	}
	if needUpdateBookKeeper {
		currBookKeeper = make([]*crypto.PubKey, len(nextBookKeeper))
		for i := 0; i < len(nextBookKeeper); i++ {
			currBookKeeper[i] = new(crypto.PubKey)
			currBookKeeper[i].X = new(big.Int).Set(nextBookKeeper[i].X)
			currBookKeeper[i].Y = new(big.Int).Set(nextBookKeeper[i].Y)
		}
	}

	bd.st.BatchPut(bhash.Bytes(), hashWriter.Bytes())

	nLen := len(b.Transactions)

	for i := 0; i < nLen; i++ {

		if b.Transactions[i].TxType == tx.RegisterAsset ||
			b.Transactions[i].TxType == tx.IssueAsset ||
			b.Transactions[i].TxType == tx.TransferAsset ||
			b.Transactions[i].TxType == tx.Record ||
			b.Transactions[i].TxType == tx.BookKeeper ||
			b.Transactions[i].TxType == tx.PrivacyPayload ||
			b.Transactions[i].TxType == tx.BookKeeping {
			err = bd.SaveTransaction(b.Transactions[i], b.Blockdata.Height)
			if err != nil {
				return err
			}
		}
		if b.Transactions[i].TxType == tx.RegisterAsset {
			ar := b.Transactions[i].Payload.(*payload.RegisterAsset)
			err = bd.SaveAsset(b.Transactions[i].Hash(), ar.Asset)
			if err != nil {
				return err
			}
		}

		if b.Transactions[i].TxType == tx.IssueAsset {
			results := b.Transactions[i].GetMergedAssetIDValueFromOutputs()
			for assetId, value := range results {
				if _, ok := quantities[assetId]; !ok {
					quantities[assetId] += value
				} else {
					quantities[assetId] = value
				}
			}
		}
		
		for index := 0; index < len(b.Transactions[i].Outputs); index++ {
			output := b.Transactions[i].Outputs[index]
			programHash := output.ProgramHash
			assetId := output.AssetID
			if value, ok := accounts[programHash]; ok {
				value.Balances[assetId] += output.Value
			} else {
				accountState, err := bd.GetAccount(programHash)
				if err != nil && err.Error() != ErrDBNotFound.Error() { return err }
				if accountState != nil {
					accountState.Balances[assetId] += output.Value
				} else {
					balances := make(map[Uint256]Fixed64, 0)
					balances[assetId] = output.Value
					accountState = account.NewAccountState(programHash, balances)
				}
				accounts[programHash] = accountState
			}
		}

		for index := 0; index < len(b.Transactions[i].UTXOInputs); index++ {
			input := b.Transactions[i].UTXOInputs[index]
			transaction, err := bd.GetTransaction(input.ReferTxID)
			if err != nil { return err }
			index := input.ReferTxOutputIndex
			output := transaction.Outputs[index]
			programHash := output.ProgramHash
			assetId := output.AssetID
			if value, ok := accounts[programHash]; ok {
				value.Balances[assetId] -= output.Value
			}else {
				accountState, err := bd.GetAccount(programHash)
				if err != nil { return err }
				accountState.Balances[assetId] -= output.Value
				accounts[programHash] = accountState
			}
			if accounts[programHash].Balances[assetId] < 0 {
				return errors.New(fmt.Sprintf("account programHash:%v, assetId:%v insufficient of balance", programHash, assetId))
			}
		}

		txhash := b.Transactions[i].Hash()
		for index := 0; index < len(b.Transactions[i].Outputs); index++ {
			unspents[txhash] = append(unspents[txhash], uint16(index))
		}

		for index := 0; index < len(b.Transactions[i].UTXOInputs); index++ {
			txhash := b.Transactions[i].UTXOInputs[index].ReferTxID


			if _, ok := unspents[txhash]; !ok {
				unspentValue, err_get := bd.st.Get(append(unspentPrefix, txhash.ToArray()...))

				if err_get != nil {
					return err_get
				}

				unspents[txhash], err_get = GetUint16Array(unspentValue)
				if err_get != nil {
					return err_get
				}
			}

Transactions[i].UTXOInputs[index].ReferTxOutputIndex and delete it
			for k := 0; k < len(unspents[txhash]); k++ {
				if unspents[txhash][k] == uint16(b.Transactions[i].UTXOInputs[index].ReferTxOutputIndex) {
					unspents[txhash] = append(unspents[txhash], unspents[txhash][:k]...)
					unspents[txhash] = append(unspents[txhash], unspents[txhash][k+1:]...)
					break
				}
			}
		}


		if b.Transactions[i].TxType == tx.BookKeeper {
			bk := b.Transactions[i].Payload.(*payload.BookKeeper)

			switch bk.Action {
			case payload.BookKeeperAction_ADD:
				findflag := false
				for k := 0; k < len(nextBookKeeper); k++ {
					if bk.PubKey.X.Cmp(nextBookKeeper[k].X) == 0 && bk.PubKey.Y.Cmp(nextBookKeeper[k].Y) == 0 {
						findflag = true
						break
					}
				}

				if !findflag {
					needUpdateBookKeeper = true
					nextBookKeeper = append(nextBookKeeper, bk.PubKey)
					sort.Sort(crypto.PubKeySlice(nextBookKeeper))
				}
			case payload.BookKeeperAction_SUB:
				ind := -1
				for k := 0; k < len(nextBookKeeper); k++ {
					if bk.PubKey.X.Cmp(nextBookKeeper[k].X) == 0 && bk.PubKey.Y.Cmp(nextBookKeeper[k].Y) == 0 {
						ind = k
						break
					}
				}

				if ind != -1 {
					needUpdateBookKeeper = true
					nextBookKeeper = append(nextBookKeeper[:ind], nextBookKeeper[ind+1:]...)
				}
			}

		}

	}

	if needUpdateBookKeeper {
		bkListKey := bytes.NewBuffer(nil)
		bkListKey.WriteByte(byte(SYS_CurrentBookKeeper))

		bkListValue := bytes.NewBuffer(nil)

		serialization.WriteUint8(bkListValue, uint8(len(currBookKeeper)))
		for k := 0; k < len(currBookKeeper); k++ {
			currBookKeeper[k].Serialize(bkListValue)
		}

		serialization.WriteUint8(bkListValue, uint8(len(nextBookKeeper)))
		for k := 0; k < len(nextBookKeeper); k++ {
			nextBookKeeper[k].Serialize(bkListValue)
		}

		bd.st.BatchPut(bkListKey.Bytes(), bkListValue.Bytes())

	}
	for txhash, value := range unspents {
		unspentKey := bytes.NewBuffer(nil)
		unspentKey.WriteByte(byte(IX_Unspent))
		txhash.Serialize(unspentKey)

		if len(value) == 0 {
			bd.st.BatchDelete(unspentKey.Bytes())
		} else {
			unspentArray := ToByteArray(value)
			bd.st.BatchPut(unspentKey.Bytes(), unspentArray)
		}
	}

	for assetId, value := range quantities {
		quantityKey := bytes.NewBuffer(nil)
		quantityKey.WriteByte(byte(ST_QuantityIssued))
		assetId.Serialize(quantityKey)

		qt, err := bd.GetQuantityIssued(assetId)
		if err != nil {
			return err
		}

		qt = qt + value

		quantityArray := bytes.NewBuffer(nil)
		qt.Serialize(quantityArray)

		bd.st.BatchPut(quantityKey.Bytes(), quantityArray.Bytes())
		log.Debug(fmt.Sprintf("quantityKey: %x\n", quantityKey.Bytes()))
		log.Debug(fmt.Sprintf("quantityArray: %x\n", quantityArray.Bytes()))
	}

	for programHash, value := range accounts {
		accountKey := new(bytes.Buffer)
		accountKey.WriteByte(byte(ST_ACCOUNT))
		programHash.Serialize(accountKey)

		accountValue := new(bytes.Buffer)
		value.Serialize(accountValue)

		bd.st.BatchPut(accountKey.Bytes(), accountValue.Bytes())
	}


	bd.currentBlockHeight = b.Blockdata.Height

	currentBlockKey := bytes.NewBuffer(nil)
	currentBlockKey.WriteByte(byte(SYS_CurrentBlock))

	currentBlock := bytes.NewBuffer(nil)
	blockHash.Serialize(currentBlock)
	serialization.WriteUint32(currentBlock, b.Blockdata.Height)

	bd.st.BatchPut(currentBlockKey.Bytes(), currentBlock.Bytes())

	err = bd.st.BatchCommit()

	if err != nil {
		return err
	}

	return nil
}

func (bd *ChainStore) addHeader(header *Header) {
	log.Debug(fmt.Sprintf("addHeader(), Height=%d\n", header.Blockdata.Height))

	hash := header.Blockdata.Hash()
	bd.headerIndex[header.Blockdata.Height] = hash
	for header.Blockdata.Height-bd.storedHeaderCount >= HeaderHashListCount {
		hashBuffer := new(bytes.Buffer)
		serialization.WriteVarUint(hashBuffer, uint64(HeaderHashListCount))
		var hashArray []byte
		for i := 0; i < HeaderHashListCount; i++ {
			index := bd.storedHeaderCount + uint32(i)
			thash := bd.headerIndex[index]
			thehash := thash.ToArray()
			hashArray = append(hashArray, thehash...)
		}
		hashBuffer.Write(hashArray)

		hhlPrefix := bytes.NewBuffer(nil)

		hhlPrefix.WriteByte(byte(IX_HeaderHashList))
		serialization.WriteUint32(hhlPrefix, bd.storedHeaderCount)

		bd.st.BatchPut(hhlPrefix.Bytes(), hashBuffer.Bytes())
		bd.storedHeaderCount += HeaderHashListCount
	}

	headerKey := bytes.NewBuffer(nil)
	headerKey.WriteByte(byte(DATA_Header))
	blockHash := header.Blockdata.Hash()
	blockHash.Serialize(headerKey)
	log.Debug(fmt.Sprintf("header key: %x\n", headerKey))

	w := bytes.NewBuffer(nil)
	var systemfee uint64 = 0xFFFFFFFFFFFFFFFF
	serialization.WriteUint64(w, systemfee)
	header.Serialize(w)
	log.Debug(fmt.Sprintf("header data: %x\n", w))

	bd.st.BatchPut(headerKey.Bytes(), w.Bytes())

	currentHeaderKey := bytes.NewBuffer(nil)
	currentHeaderKey.WriteByte(byte(SYS_CurrentHeader))

	currentHeader := bytes.NewBuffer(nil)
	blockHash.Serialize(currentHeader)
	serialization.WriteUint32(currentHeader, header.Blockdata.Height)

	bd.st.BatchPut(currentHeaderKey.Bytes(), currentHeader.Bytes())

	log.Debug("[addHeader]: finish, header height:", header.Blockdata.Height)
}

func (bd *ChainStore) persistBlocks(ledger *Ledger) {
	bd.mu.Lock()
	defer bd.mu.Unlock()
	for !bd.disposed {
		if uint32(len(bd.headerIndex)) < bd.currentBlockHeight+1 {
			log.Warn("[persistBlocks]: warn, headerIndex.count < currentBlockHeight + 1")
			break
		}

		hash := bd.headerIndex[bd.currentBlockHeight+1]

		block, ok := bd.blockCache[hash]
		if !ok {
			log.Warn("[persistBlocks]: warn, blockCache not contain key hash.")
			break
		}
		bd.persist(block)

		ledger.Blockchain.BCEvents.Notify(events.EventBlockPersistCompleted, block)
		ledger.Blockchain.BlockHeight = block.Blockdata.Height
		log.Trace("The latest block height:", block.Blockdata.Height)

		delete(bd.blockCache, hash)
	}

}

func (bd *ChainStore) SaveBlock(b *Block, ledger *Ledger) error {
	log.Debug("SaveBlock()")

	bd.mu.Lock()
	defer bd.mu.Unlock()

	if bd.blockCache[b.Hash()] == nil {
		bd.blockCache[b.Hash()] = b
	}

	if b.Blockdata.Height >= (uint32(len(bd.headerIndex)) + 1) {
headerIndex.count:%d",b.Blockdata.Height, uint32(len(bd.headerIndex)) )),ErrDuplicatedBlock,"")
		return errors.New(fmt.Sprintf("WARNING: [SaveBlock] block height - headerIndex.count >= 1, block height:%d, headerIndex.count:%d", b.Blockdata.Height, uint32(len(bd.headerIndex))))
	}

	if b.Blockdata.Height == uint32(len(bd.headerIndex)) {

		err := validation.VerifyBlock(b, ledger, false)
		if err != nil {
			log.Debug("VerifyBlock() error!")
			return err
		}

		bd.st.NewBatch()
		h := new(Header)
		h.Blockdata = b.Blockdata
		bd.addHeader(h)
		err = bd.st.BatchCommit()
		if err != nil {
			return err
		}
	}

	if b.Blockdata.Height < uint32(len(bd.headerIndex)) {
		go bd.persistBlocks(ledger)
	}

	return nil
}

func (bd *ChainStore) BlockInCache(hash Uint256) bool {
	bd.mu.RLock()
	defer bd.mu.RUnlock()
	if _, ok := bd.blockCache[hash]; ok {
		return true
	}
	return false
}

func (bd *ChainStore) GetQuantityIssued(assetId Uint256) (Fixed64, error) {
	log.Debug(fmt.Sprintf("GetQuantityIssued Hash: %x\n", assetId))

	prefix := []byte{byte(ST_QuantityIssued)}
	data, err_get := bd.st.Get(append(prefix, assetId.ToArray()...))
	log.Debug(fmt.Sprintf("GetQuantityIssued Data: %x\n", data))

	var quantity Fixed64
	if err_get != nil {
		quantity = Fixed64(0)
	} else {
		r := bytes.NewReader(data)
		quantity.Deserialize(r)
	}

	return quantity, nil
}

func (bd *ChainStore) GetUnspent(txid Uint256, index uint16) (*tx.TxOutput, error) {
	if ok, _ := bd.ContainsUnspent(txid, index); ok {
		Tx, err := bd.GetTransaction(txid)
		if err != nil {
			return nil, err
		}

		return Tx.Outputs[index], nil
	}

	return nil, errors.New("[GetUnspent] NOT ContainsUnspent.")
}

func (bd *ChainStore) ContainsUnspent(txid Uint256, index uint16) (bool, error) {
	unspentPrefix := []byte{byte(IX_Unspent)}
	unspentValue, err_get := bd.st.Get(append(unspentPrefix, txid.ToArray()...))

	if err_get != nil {
		return false, err_get
	}

	unspentArray, err_get := GetUint16Array(unspentValue)
	if err_get != nil {
		return false, err_get
	}

	for i := 0; i < len(unspentArray); i++ {
		if unspentArray[i] == index {
			return true, nil
		}
	}

	return false, nil
}

func (bd *ChainStore) GetCurrentHeaderHash() Uint256 {
	bd.mu.RLock()
	defer bd.mu.RUnlock()
	return bd.headerIndex[uint32(len(bd.headerIndex)-1)]
}

func (bd *ChainStore) GetHeaderHashByHeight(height uint32) Uint256 {
	bd.mu.RLock()
	defer bd.mu.RUnlock()

	return bd.headerIndex[height]
}

func (bd *ChainStore) GetHeaderHeight() uint32 {
	bd.mu.RLock()
	defer bd.mu.RUnlock()

	return uint32(len(bd.headerIndex) - 1)
}

func (bd *ChainStore) GetHeight() uint32 {
	bd.mu.RLock()
	defer bd.mu.RUnlock()

	return bd.currentBlockHeight
}

func (bd *ChainStore) GetAccount(programHash Uint160) (*account.AccountState, error) {
	accountPrefix := []byte{byte(ST_ACCOUNT)}

	state, err := bd.st.Get(append(accountPrefix, programHash.ToArray()...))

	if err != nil { return nil, err }

	accountState := new(account.AccountState)
	accountState.Deserialize(bytes.NewBuffer(state))

	return accountState, nil
}
