package network

import (
	"encoding/binary"
	"fmt"
	"github.com/erDong01/micro-kit/tools"
	"log"
)

const (
	PACKET_LEN_BYTE  = 1
	PACKET_LEN_WORD  = 2
	PACKET_LEN_DWORD = 4
)

type (
	PacketParser struct {
		PacketLen       int
		MaxPacketLen    int
		LittleEndian    bool
		MaxPacketBuffer []byte //max receive buff
		PacketFunc      HandlePacket
	}

	PacketConfig struct {
		MaxPacketLen int
		Func         HandlePacket
	}
)

func NewPacketParser(conf PacketConfig) PacketParser {
	p := PacketParser{}
	p.PacketLen = PACKET_LEN_DWORD
	p.MaxPacketLen = tools.MAX_PACKET
	p.LittleEndian = true
	if conf.Func != nil {
		p.PacketFunc = conf.Func
	} else {
		p.PacketFunc = func(buff []byte) {
		}
	}
	return p
}
func (this *PacketParser) readLen(buff []byte) (bool, int) {
	nLen := len(buff)
	if nLen < this.PacketLen {
		return false, 0
	}
	bufMsgLen := buff[:this.PacketLen]
	var msgLen int

	switch this.PacketLen {
	case PACKET_LEN_BYTE:
		msgLen = int(bufMsgLen[0])
	case PACKET_LEN_WORD:
		if this.LittleEndian {
			msgLen = int(binary.LittleEndian.Uint16(bufMsgLen))
		} else {
			msgLen = int(binary.BigEndian.Uint16(bufMsgLen))
		}
	case PACKET_LEN_DWORD:
		if this.LittleEndian {
			msgLen = int(binary.LittleEndian.Uint32(bufMsgLen))
		} else {
			msgLen = int(binary.BigEndian.Uint32(bufMsgLen))
		}

	}
	if msgLen+this.PacketLen <= nLen {
		return true, msgLen + this.PacketLen
	}
	return false, 0
}

func (this *PacketParser) Read(dat []byte) bool {
	buff := append(this.MaxPacketBuffer, dat...)
	this.MaxPacketBuffer = []byte{}
	nCurSize := 0
ParsePacket:
	nPacketSize := 0
	nBufferSize := len(buff[nCurSize:])
	bFindFlag := false
	bFindFlag, nPacketSize = this.readLen(buff[nCurSize:])

	if bFindFlag {
		if nBufferSize == nPacketSize { //完整包
			this.PacketFunc(buff[nCurSize+this.PacketLen : nCurSize+nPacketSize])
			nCurSize += nPacketSize
		} else if nBufferSize > nPacketSize {
			this.PacketFunc(buff[nCurSize+this.PacketLen : nCurSize+nPacketSize])
			nCurSize += nPacketSize
			goto ParsePacket
		}
	} else if nBufferSize < this.MaxPacketLen {
		this.MaxPacketBuffer = buff[nCurSize:]
	} else {
		log.Println("超出最大包限制，丢弃该包")
		return false
	}

	return true
}

func (this *PacketParser) Write(dat []byte) []byte {
	msgLen := len(dat)

	if msgLen+this.PacketLen > tools.MAX_PACKET {
		fmt.Println("write over base.MAX_PACKET")
	}
	msg := make([]byte, this.PacketLen+msgLen)

	switch this.PacketLen {
	case PACKET_LEN_BYTE:
		msg[0] = byte(msgLen)
	case PACKET_LEN_WORD:
		if this.LittleEndian {
			binary.LittleEndian.PutUint16(msg, uint16(msgLen))
		} else {
			binary.BigEndian.PutUint16(msg, uint16(msgLen))
		}
	case PACKET_LEN_DWORD:
		if this.LittleEndian {
			binary.LittleEndian.PutUint32(msg, uint32(msgLen))
		} else {
			binary.BigEndian.PutUint32(msg, uint32(msgLen))
		}
	}
	copy(msg[this.PacketLen:], dat)
	return msg
}
