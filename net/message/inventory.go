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

