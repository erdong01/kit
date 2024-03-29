package network

import (
	"encoding/binary"
	"log"

	"github.com/erdong01/kit/base"
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
		packetFunc      HandlePacket
	}

	PacketConfig struct {
		MaxPacketLen *int
		Func         HandlePacket
	}

	IPacketParser interface {
		readLen(buff []byte) (bool, int)
		Read(dat []byte) bool
		Write(dat []byte) []byte
		GetMaxPacketLen() int
		SetMaxPacketLen(val int)
	}
)

func NewPacketParser(conf PacketConfig) *PacketParser {
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
	return &p
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

func (p *PacketParser) Read(dat []byte) bool {
	buff := append(p.maxPacketBuffer, dat...)
	p.maxPacketBuffer = []byte{}
	nCurSize := 0
	//fmt.Println(p.maxPacketBuffer)
ParsePacekt:
	nPacketSize := 0
	nBufferSize := len(buff[nCurSize:])
	bFindFlag := false
	bFindFlag, nPacketSize = p.readLen(buff[nCurSize:])
	//fmt.Println(bFindFlag, nPacketSize, nBufferSize)
	if bFindFlag {
		if nBufferSize == nPacketSize { //完整包
			p.packetFunc(buff[nCurSize+p.packetLen : nCurSize+nPacketSize])
			nCurSize += nPacketSize
		} else if nBufferSize > nPacketSize {
			p.packetFunc(buff[nCurSize+p.packetLen : nCurSize+nPacketSize])
			nCurSize += nPacketSize
			goto ParsePacekt
		}
	} else if nBufferSize < p.maxPacketLen {
		p.maxPacketBuffer = buff[nCurSize:]
	} else {
		log.Println("超出最大包限制，丢弃该包")
		return false
	}
	return true
}

func (p *PacketParser) Write(dat []byte) []byte {
	// get len
	msgLen := len(dat)
	// check len
	if msgLen+p.packetLen > base.MAX_PACKET {
		log.Println("write over base.MAX_PACKET")
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

func (p *PacketParser) GetMaxPacketLen() int {
	return p.maxPacketLen
}
func (p *PacketParser) SetMaxPacketLen(val int) {
	p.maxPacketLen = val
}
