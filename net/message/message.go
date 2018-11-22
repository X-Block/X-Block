package message

import (
	"XBlock/common/log"
	. "XBlock/net/protocol"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
)

type Messager interface {
	Verify([]byte) error
	Serialization() ([]byte, error)
	Deserialization([]byte) error
	Handle(Noder) error
}

type msgHdr struct {
	Magic uint32
	
	CMD      [MSGCMDLEN]byte 
	Length   uint32
	Checksum [CHECKSUMLEN]byte
}


type msgCont struct {
	hdr msgHdr
	p   interface{}
}

type varStr struct {
	len uint
	buf []byte
}

type filteradd struct {
	msgHdr
	
}

type filterclear struct {
	msgHdr
	
}

type filterload struct {
	msgHdr
	
}

func AllocMsg(t string, length int) Messager {
	switch t {
	case "msgheader":
		var msg msgHdr
		return &msg
	case "version":
		var msg version
		copy(msg.Hdr.CMD[0:len(t)], t)
		return &msg
	case "verack":
		var msg verACK
		copy(msg.msgHdr.CMD[0:len(t)], t)
		return &msg
	case "getheaders":
		var msg headersReq
		copy(msg.hdr.CMD[0:len(t)], t)
		return &msg
	case "headers":
		var msg blkHeader
		copy(msg.hdr.CMD[0:len(t)], t)
		return &msg
	case "getaddr":
		var msg addrReq
		copy(msg.Hdr.CMD[0:len(t)], t)
		return &msg
	case "addr":
		var msg addr
		copy(msg.hdr.CMD[0:len(t)], t)
		return &msg
	case "inv":
		var msg Inv
		copy(msg.Hdr.CMD[0:len(t)], t)
		msg.P.Blk = make([]byte, length-MSGHDRLEN-1)
		return &msg
	case "getdata":
		var msg dataReq
		copy(msg.msgHdr.CMD[0:len(t)], t)
		return &msg
	case "block":
		var msg block
		copy(msg.msgHdr.CMD[0:len(t)], t)
		return &msg
	case "tx":
		var msg trn
		copy(msg.msgHdr.CMD[0:len(t)], t)
	
		return &msg
	case "consensus":
		var msg consensus
		copy(msg.msgHdr.CMD[0:len(t)], t)
		return &msg
	case "filteradd":
		var msg filteradd
		copy(msg.msgHdr.CMD[0:len(t)], t)
		return &msg
	case "filterclear":
		var msg filterclear
		copy(msg.msgHdr.CMD[0:len(t)], t)
		return &msg
	case "filterload":
		var msg filterload
		copy(msg.msgHdr.CMD[0:len(t)], t)
		return &msg
	case "getblocks":
		var msg blocksReq
		copy(msg.msgHdr.CMD[0:len(t)], t)
		return &msg
	case "txnpool":
		var msg txnPool
		copy(msg.msgHdr.CMD[0:len(t)], t)
		return &msg
	case "alert":
		log.Warn("Not supported message type - alert")
		return nil
	case "merkleblock":
		log.Warn("Not supported message type - merkleblock")
		return nil
	case "notfound":
		log.Warn("Not supported message type - notfound")
		return nil
	case "ping":
		var msg ping
		copy(msg.msgHdr.CMD[0:len(t)], t)
		return &msg
	case "pong":
		var msg pong
		copy(msg.msgHdr.CMD[0:len(t)], t)
		return &msg
	case "reject":
		log.Warn("Not supported message type - reject")
		return nil
	default:
		log.Warn("Unknown message type")
		return nil
	}
}

func MsgType(buf []byte) (string, error) {
	cmd := buf[CMDOFFSET : CMDOFFSET+MSGCMDLEN]
	n := bytes.IndexByte(cmd, 0)
	if n < 0 || n >= MSGCMDLEN {
		return "", errors.New("Unexpected length of CMD command")
	}
	s := string(cmd[:n])
	return s, nil
}


func NewMsg(t string, n Noder) ([]byte, error) {
	switch t {
	case "version":
		return NewVersion(n)
	case "verack":
		return NewVerack()
	case "getheaders":
		return NewHeadersReq()
	case "getaddr":
		return newGetAddr()

	default:
		return nil, errors.New("Unknown message type")
	}
}


func HandleNodeMsg(node Noder, buf []byte, len int) error {
	if len < MSGHDRLEN {
		log.Warn("Unexpected size of received message")
		return errors.New("Unexpected size of received message")
	}

	str := hex.EncodeToString(buf[:len])
	log.Debug("Received data len: ", len, "\n", str)

	s, err := MsgType(buf)
	if err != nil {
		log.Error("Message type parsing error")
		return err
	}

	msg := AllocMsg(s, len)
	if msg == nil {
		log.Error(fmt.Sprintf("Allocation message %s failed", s))
		return errors.New("Allocation message failed")
	}
	msg.Deserialization(buf[:len])
	msg.Verify(buf[MSGHDRLEN:len])

	return msg.Handle(node)
}

