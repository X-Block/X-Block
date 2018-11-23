package message

import (
	"XBlock/common"
	"XBlock/common/log"
	"XBlock/common/serialization"
	"XBlock/core/ledger"
	. "XBlock/net/protocol"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
)

type headersReq struct {
	hdr msgHdr
	p   struct {
		len       uint8
		hashStart [HASHLEN]byte
		hashEnd   [HASHLEN]byte
	}
}

type blkHeader struct {
	hdr    msgHdr
	cnt    uint32
	blkHdr []ledger.Header
}

func NewHeadersReq() ([]byte, error) {
	var h headersReq

	h.p.len = 1
	buf := ledger.DefaultLedger.Store.GetCurrentHeaderHash()
	copy(h.p.hashEnd[:], reverse(buf[:]))

	p := new(bytes.Buffer)
	err := binary.Write(p, binary.LittleEndian, &(h.p))
	if err != nil {
		log.Error("Binary Write failed at new headersReq")
		return nil, err
	}

	s := checkSum(p.Bytes())
	h.hdr.init("getheaders", s, uint32(len(p.Bytes())))

	m, err := h.Serialization()
	return m, err
}

func (msg headersReq) Verify(buf []byte) error {
	err := msg.hdr.Verify(buf)
	return err
}

func (msg blkHeader) Verify(buf []byte) error {
	err := msg.hdr.Verify(buf)
	return err
}

func (msg headersReq) Serialization() ([]byte, error) {
	var buf bytes.Buffer

	err := binary.Write(&buf, binary.LittleEndian, msg)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), err
}

func (msg *headersReq) Deserialization(p []byte) error {
	buf := bytes.NewBuffer(p)
	err := binary.Read(buf, binary.LittleEndian, &msg)
	return err
}

func (msg blkHeader) Serialization() ([]byte, error) {
	hdrBuf, err := msg.hdr.Serialization()
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(hdrBuf)
	binary.Write(buf, binary.LittleEndian, msg.cnt)

	for _, header := range msg.blkHdr {
		header.Serialize(buf)
	}
	return buf.Bytes(), err
}

func (msg *blkHeader) Deserialization(p []byte) error {
	buf := bytes.NewBuffer(p)
	err := binary.Read(buf, binary.LittleEndian, &(msg.hdr))
	if err != nil {
		return err
	}

	err = binary.Read(buf, binary.LittleEndian, &(msg.cnt))
	if err != nil {
		return err
	}

	for i := 0; i < int(msg.cnt); i++ {
		var headers ledger.Header
		err := (&headers).Deserialize(buf)
		msg.blkHdr = append(msg.blkHdr, headers)
		if err != nil {
			log.Debug("blkHeader Deserialization failed")
			goto blkHdrErr
		}
	}

blkHdrErr:
	return err
}

func (msg headersReq) Handle(node Noder) error {
	log.Debug()
	var startHash [HASHLEN]byte
	var stopHash [HASHLEN]byte
	startHash = msg.p.hashStart
	stopHash = msg.p.hashEnd
	headers, cnt, err := GetHeadersFromHash(startHash, stopHash)
	if err != nil {
		return err
	}
	buf, err := NewHeaders(headers, cnt)
	if err != nil {
		return err
	}
	go node.Tx(buf)
	return nil
}

