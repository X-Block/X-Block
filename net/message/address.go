package message

import (
	"XBlock/common/log"
	. "XBlock/net/protocol"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
)

type addrReq struct {
	Hdr msgHdr
	
}

type addr struct {
	hdr       msgHdr
	nodeCnt   uint64
	nodeAddrs []NodeAddr
}

const (
	NODEADDRSIZE = 30
)

func newGetAddr() ([]byte, error) {
	var msg addrReq
	var sum []byte
	sum = []byte{0x23, 0x72, 0x35, 0x24}
	msg.Hdr.init("getaddr", sum, 0)

	buf, err := msg.Serialization()
	if err != nil {
		return nil, err
	}

	str := hex.EncodeToString(buf)
	log.Debug("The message get addr length is: ", len(buf), " ", str)

	return buf, err
}

func NewAddrs(nodeaddrs []NodeAddr, count uint64) ([]byte, error) {
	var msg addr
	msg.nodeAddrs = nodeaddrs
	msg.nodeCnt = count
	msg.hdr.Magic = NETMAGIC
	cmd := "addr"
	copy(msg.hdr.CMD[0:7], cmd)
	p := new(bytes.Buffer)
	err := binary.Write(p, binary.LittleEndian, msg.nodeCnt)
	if err != nil {
		log.Error("Binary Write failed at new Msg: ", err.Error())
		return nil, err
	}

	err = binary.Write(p, binary.LittleEndian, msg.nodeAddrs)
	if err != nil {
		log.Error("Binary Write failed at new Msg: ", err.Error())
		return nil, err
	}
	s := sha256.Sum256(p.Bytes())
	s2 := s[:]
	s = sha256.Sum256(s2)
	buf := bytes.NewBuffer(s[:4])
	binary.Read(buf, binary.LittleEndian, &(msg.hdr.Checksum))
	msg.hdr.Length = uint32(len(p.Bytes()))
	log.Debug("The message payload length is ", msg.hdr.Length)

	m, err := msg.Serialization()
	if err != nil {
		log.Error("Error Convert net message ", err.Error())
		return nil, err
	}

	return m, nil
}

func (msg addrReq) Verify(buf []byte) error {
	err := msg.Hdr.Verify(buf)
	return err
}

func (msg addrReq) Handle(node Noder) error {
	log.Debug()
	var addrstr []NodeAddr
	var count uint64
	addrstr, count = node.LocalNode().GetNeighborAddrs()
	buf, err := NewAddrs(addrstr, count)
	if err != nil {
		return err
	}
	go node.Tx(buf)
	return nil
}

func (msg addrReq) Serialization() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, msg)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), err
}

func (msg *addrReq) Deserialization(p []byte) error {
	buf := bytes.NewBuffer(p)
	err := binary.Read(buf, binary.LittleEndian, msg)
	return err
}

