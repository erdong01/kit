package json

import (
	"encoding/binary"
	"fmt"

	"github.com/erdong01/kit/base"
	"github.com/erdong01/kit/network"
)

const (
	PACKET_LEN_BYTE  = 1
	PACKET_LEN_WORD  = 2
	PACKET_LEN_DWORD = 4
)

// --------------
// | len | data |
// --------------
type (
	PacketParser struct {
		packetLen       int
		maxPacketLen    int
		littleEndian    bool
		maxPacketBuffer []byte //max receive buff
		packetFunc      network.HandlePacket
	}

	PacketConfig struct {
		MaxPacketLen *int
		Func         network.HandlePacket
	}
)

func NewPacketParser(conf PacketConfig) PacketParser {
	p := PacketParser{}
	p.packetLen = PACKET_LEN_DWORD
	p.maxPacketLen = base.MAX_PACKET
	p.littleEndian = true
	if conf.Func != nil {
		p.packetFunc = conf.Func
	} else {
		p.packetFunc = func(buff []byte) {
		}
	}
	return p
}

func (p *PacketParser) readLen(buff []byte) (bool, int) {
	nLen := len(buff)
	if nLen < p.packetLen {
		return false, 0
	}

	bufMsgLen := buff[:p.packetLen]
	// parse len
	var msgLen int
	switch p.packetLen {
	case PACKET_LEN_BYTE:
		msgLen = int(bufMsgLen[0])
	case PACKET_LEN_WORD:
		if p.littleEndian {
			msgLen = int(binary.LittleEndian.Uint16(bufMsgLen))
		} else {
			msgLen = int(binary.BigEndian.Uint16(bufMsgLen))
		}
	case PACKET_LEN_DWORD:
		if p.littleEndian {
			msgLen = int(binary.LittleEndian.Uint32(bufMsgLen))
		} else {
			msgLen = int(binary.BigEndian.Uint32(bufMsgLen))
		}
	}

	if msgLen+p.packetLen <= nLen {
		return true, msgLen + p.packetLen
	}

	return false, 0
}

func (p *PacketParser) Read(buff []byte) bool {
	//fmt.Println(p.maxPacketBuffer)
	//fmt.Println(bFindFlag, nPacketSize, nBufferSize)
			p.packetFunc(buff)
	return true
}

func (p *PacketParser) Write(dat []byte) []byte {
	// get len
	msgLen := len(dat)
	// check len
	if msgLen+p.packetLen > base.MAX_PACKET {
		fmt.Println("write over base.MAX_PACKET")
	}

	msg := make([]byte, p.packetLen+msgLen)
	// write len
	switch p.packetLen {
	case PACKET_LEN_BYTE:
		msg[0] = byte(msgLen)
	case PACKET_LEN_WORD:
		if p.littleEndian {
			binary.LittleEndian.PutUint16(msg, uint16(msgLen))
		} else {
			binary.BigEndian.PutUint16(msg, uint16(msgLen))
		}
	case PACKET_LEN_DWORD:
		if p.littleEndian {
			binary.LittleEndian.PutUint32(msg, uint32(msgLen))
		} else {
			binary.BigEndian.PutUint32(msg, uint32(msgLen))
		}
	}

	copy(msg[p.packetLen:], dat)
	return msg
}
