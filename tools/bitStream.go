package tools

const (
	Bit8              = 8
	Bit16             = 16
	Bit32             = 32
	Bit64             = 64
	Bit128            = 128
	MAX_PACKET        = 1 * 1024 * 1024 // 1MB
	MAX_CLIENT_PACKET = 10 * 1024       // 10KB
)

type (
	BitStream struct {
		dataPtr        []byte
		bitNum         int
		flagNum        int
		tailFlag       bool
		bufSize        int
		bitsLimite     int
		error          bool
		maxReadBitNum  int
		maxWriteBitNum int
	}

	IBitStream interface {
		BuildPacketStream([]byte, int) bool
		setBuffer([]byte, int, int)
		GetBuffer() []byte
		GetBytePtr() []byte
		GetReadByteSize() int
		GetCurPos() int
		GetPosition() int
		GetStreamSize() int
		SetPosition(int) bool
		clear()
		resize() bool

		WriteBits([]byte, int)
		ReadBits(int) []byte
		WriteInt(int, int)
		ReadInt(int) int
		ReadFlag() bool
		WriteFlag(bool) bool
		WriteString(string)
		ReadString() string

		WriteInt64(int64, int)
		ReadInt64(int) int64
		WriteFloat(float32)
		ReadFloat() float32
		WriteFloat64(float64)
		ReadFloat64() float64
	}
)

func (this BitStream) BuildPacketStream(buffer []byte, writeSize int) bool {
	if writeSize <= 0 {
		return false
	}
	this.setBuffer(buffer, writeSize, -1)
	this.SetPosition(0)
	return true
}
func (this *BitStream) SetPosition(pos int) bool {
	Assert(pos == 0 || this.flagNum == 0, "不正确的setPosition调用")
	if pos != 0 && this.flagNum != 0 {
		return false
	}

	this.bitNum = pos << 3
	this.flagNum = 0
	return true
}
func (this BitStream) setBuffer(bufPtr []byte, size int, maxSize int) {
	this.dataPtr = bufPtr
	this.bitNum = 0
	this.flagNum = 0
	this.tailFlag = false
	this.bufSize = size
	this.maxReadBitNum = size << 3
	if maxSize < 0 {
		maxSize = size
	}
	this.maxWriteBitNum = maxSize << 3
	this.bitsLimite = size
	this.error = false
}
func (this *BitStream) WriteInt(value int, bitCount int) {
	this.WriteBits(IntToBytes(value), bitCount)
}
func (this *BitStream) WriteBits(bitPtr []byte, bitCount int) {
	if bitCount == 0 {
		return
	}

	if this.tailFlag {
		this.error = true
		Assert(false, "Out of range write")
		return
	}

	if (bitCount & 0x7) != 0 {
		bitCount = (bitCount & ^0x7) + 8
	}

	for bitCount+this.bitNum > this.maxWriteBitNum {
		if !this.resize() {
			this.error = true
			Assert(false, "Out of range write")
			return
		}
	}

	bitNum := this.bitNum >> 3
	byteCount := (bitCount + 7) >> 3
	copy(this.dataPtr[bitNum:], bitPtr[:byteCount])
	/*for i, v := range bitPtr[:byteCount] {
		this.dataPtr[bitNum+i] = v
	}*/
	this.bitNum += bitCount
}
func (this *BitStream) resize() bool {
	//fmt.Println("BitStream Resize")
	this.dataPtr = append(this.dataPtr, make([]byte, this.bitsLimite)...)
	size := this.bitsLimite * 2
	if size <= 0 || size >= MAX_PACKET*2 {
		return false
	}
	this.bufSize = size
	this.maxReadBitNum = size << 3
	this.maxWriteBitNum = size << 3
	this.bitsLimite = size
	return true
}

func (this *BitStream) GetBuffer() []byte {
	return this.dataPtr[0:this.GetPosition()]
}

func (this *BitStream) GetPosition() int {
	return (this.bitNum + 7) >> 3
}

func NewBitStream(buf []byte, nLen int) *BitStream {
	var bitstream BitStream
	bitstream.BuildPacketStream(buf, nLen)
	return &bitstream
}
