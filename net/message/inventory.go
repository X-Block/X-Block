package message

import (
	. "XBlock/common"
	"XBlock/common/log"
	"XBlock/common/serialization"
	"XBlock/core/ledger"
	. "XBlock/net/protocol"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
)

type blocksReq struct {
	msgHdr
	p struct {
		HeaderHashCount uint8
		hashStart       [HASHLEN]byte
		hashStop        [HASHLEN]byte
	}
}

type InvPayload struct {
	InvType InventoryType
	Cnt     uint32
	Blk     []byte
}

type Inv struct {
	Hdr msgHdr
	P   InvPayload
}

func NewBlocksReq(n Noder) ([]byte, error) {
	var h blocksReq
	log.Debug("request block hash")
	h.p.HeaderHashCount = 1
	buf := ledger.DefaultLedger.Blockchain.CurrentBlockHash()

	copy(h.p.hashStart[:], reverse(buf[:]))

	p := new(bytes.Buffer)
	err := binary.Write(p, binary.LittleEndian, &(h.p))
	if err != nil {
		log.Error("Binary Write failed at new blocksReq")
		return nil, err
	}

	s := checkSum(p.Bytes())
	h.msgHdr.init("getblocks", s, uint32(len(p.Bytes())))

	m, err := h.Serialization()

	return m, err
}

func (msg blocksReq) Verify(buf []byte) error {


	err := msg.msgHdr.Verify(buf)
	return err
}

func (msg blocksReq) Handle(node Noder) error {
	log.Debug()
	log.Debug("handle blocks request")
	var starthash Uint256
	var stophash Uint256
	starthash = msg.p.hashStart
	stophash = msg.p.hashStop
	inv, err := GetInvFromBlockHash(starthash, stophash)
	if err != nil {
		return err
	}
	buf, err := NewInv(inv)
	if err != nil {
		return err
	}
	go node.Tx(buf)
	return nil
}

func (msg blocksReq) Serialization() ([]byte, error) {
	var buf bytes.Buffer

	err := binary.Write(&buf, binary.LittleEndian, msg)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), err
}

func (msg *blocksReq) Deserialization(p []byte) error {
	buf := bytes.NewBuffer(p)
	err := binary.Read(buf, binary.LittleEndian, &msg)
	return err
}

func (msg Inv) Verify(buf []byte) error {

	err := msg.Hdr.Verify(buf)
	return err
}

func (msg Inv) Handle(node Noder) error {
	log.Debug()
	var id Uint256
	str := hex.EncodeToString(msg.P.Blk)
	log.Debug(fmt.Sprintf("The inv type: 0x%x block len: %d, %s\n",
		msg.P.InvType, len(msg.P.Blk), str))

	invType := InventoryType(msg.P.InvType)
	switch invType {
	case TRANSACTION:
		log.Debug("RX TRX message")
		id.Deserialize(bytes.NewReader(msg.P.Blk[:32]))
		if !node.ExistedID(id) {
			reqTxnData(node, id)
		}
	case BLOCK:
		log.Debug("RX block message")
		var i uint32
		count := msg.P.Cnt
		log.Debug("RX inv-block message, hash is ", msg.P.Blk)
		for i = 0; i < count; i++ {
			id.Deserialize(bytes.NewReader(msg.P.Blk[HASHLEN*i:]))
			if !ledger.DefaultLedger.Store.BlockInCache(id) &&
				!ledger.DefaultLedger.BlockInLedger(id) {
				log.Info("inv request block hash: ", id)
				ReqBlkData(node, id)
			}

		}
	case CONSENSUS:
		log.Debug("RX consensus message")
		id.Deserialize(bytes.NewReader(msg.P.Blk[:32]))
		reqConsensusData(node, id)
	default:
		log.Warn("RX unknown inventory message")
	}
	return nil
}

func (msg Inv) Serialization() ([]byte, error) {
	hdrBuf, err := msg.Hdr.Serialization()
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(hdrBuf)
	msg.P.Serialization(buf)

	return buf.Bytes(), err
}

func (msg *Inv) Deserialization(p []byte) error {
	err := msg.Hdr.Deserialization(p)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(p[MSGHDRLEN:])
	invType, err := serialization.ReadUint8(buf)
	if err != nil {
		return err
	}
	msg.P.InvType = InventoryType(invType)
	msg.P.Cnt, err = serialization.ReadUint32(buf)
	if err != nil {
		return err
	}

	msg.P.Blk = make([]byte, msg.P.Cnt*HASHLEN)
	err = binary.Read(buf, binary.LittleEndian, &(msg.P.Blk))

	return err
}

func (msg Inv) invType() InventoryType {
	return msg.P.InvType
}

func GetInvFromBlockHash(starthash Uint256, stophash Uint256) (*InvPayload, error) {
	var count uint32 = 0
	var i uint32
	var empty Uint256
	var startheight uint32
	var stopheight uint32
	curHeight := ledger.DefaultLedger.GetLocalBlockChainHeight()
	if starthash == empty {
		if stophash == empty {
			if curHeight > MAXBLKHDRCNT {
				count = MAXBLKHDRCNT
			} else {
				count = curHeight
			}
		} else {
			bkstop, err := ledger.DefaultLedger.GetBlockWithHash(stophash)
			if err != nil {
				return nil, err
			}
			stopheight = bkstop.Blockdata.Height
			count = curHeight - stopheight
			if curHeight > MAXINVHDRCNT {
				count = MAXINVHDRCNT
			}
		}
	} else {
		bkstart, err := ledger.DefaultLedger.GetBlockWithHash(starthash)
		if err != nil {
			return nil, err
		}
		startheight = bkstart.Blockdata.Height
		if stophash != empty {
			bkstop, err := ledger.DefaultLedger.GetBlockWithHash(stophash)
			if err != nil {
				return nil, err
			}
			stopheight = bkstop.Blockdata.Height
			count = startheight - stopheight
			if count >= MAXINVHDRCNT {
				count = MAXINVHDRCNT
				stopheight = startheight + MAXINVHDRCNT
			}
		} else {

			if startheight > MAXINVHDRCNT {
				count = MAXINVHDRCNT
			} else {
				count = startheight
			}
		}
	}
	tmpBuffer := bytes.NewBuffer([]byte{})
	for i = 1; i <= count; i++ {
		hash, _ := ledger.DefaultLedger.Store.GetBlockHash(stopheight + i)
		log.Debug("GetInvFromBlockHash i is ", i, " , hash is ", hash)
		hash.Serialize(tmpBuffer)
	}
	log.Debug("GetInvFromBlockHash hash is ", tmpBuffer.Bytes())
	return NewInvPayload(BLOCK, count, tmpBuffer.Bytes()), nil
}

func NewInvPayload(invType InventoryType, count uint32, msg []byte) *InvPayload {
	return &InvPayload{
		InvType: invType,
		Cnt:     count,
		Blk:     msg,
	}
}

