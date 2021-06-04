package tools

import (
	"encoding/binary"
	"hash/crc32"
	"math"
	"os"
	"reflect"
	"strings"
	"time"
)

const (
	INT_MAX       = int(2147483647)
	TCP_HEAD_SIZE = 4     //解决tpc粘包半包,包头固定长度
	TCP_END       = "💞♡" //解决tpc粘包半包,特殊结束标志,pb采用Varint编码高位有特殊含义
)

var (
	SEVERNAME      string
	TCP_END_LENGTH = len([]byte(TCP_END)) //tcp结束标志长度
)

// IntToBytes
// @Description: 整形转换成字节
// @param val int
// @return []byte
func IntToBytes(val int) []byte {
	tmp := uint32(val)
	buff := make([]byte, 4)
	binary.LittleEndian.PutUint32(buff, tmp)
	return buff
}

// BytesToInt
// @Description: 字节转换成整形
// @param data []byte
// @return int
func BytesToInt(data []byte) int {
	buff := make([]byte, 4)
	copy(buff, data)
	tmp := int32(binary.LittleEndian.Uint32(buff))
	return int(tmp)
}

// Float64ToByte
// @Description: 转化float64
// @param val float64
// @return []byte
func Float64ToByte(val float64) []byte {
	tmp := math.Float64bits(val)
	buff := make([]byte, 8)
	binary.LittleEndian.PutUint64(buff, tmp)
	return buff
}

func BytesToFloat64(data []byte) float64 {
	buff := make([]byte, 8)
	copy(buff, data)
	tmp := binary.LittleEndian.Uint64(buff)
	return math.Float64frombits(tmp)
}

// Int16ToBytes
// @Description: 整形16转换成字节
// @param val int16
// @return []byte
func Int16ToBytes(val int16) []byte {
	tmp := uint16(val)
	buff := make([]byte, 2)
	binary.LittleEndian.PutUint16(buff, tmp)
	return buff
}

// BytesToInt16
// @author chenqiaojie
// @Description: 字节转换成为int16
// @param data []byte
// @return int16
func BytesToInt16(data []byte) int16 {
	buff := make([]byte, 2)
	copy(buff, data)
	tmp := binary.LittleEndian.Uint16(buff)
	return int16(tmp)
}

// Int64ToBytes
// @Description: 转化64位
// @param val int64
// @return []byte
func Int64ToBytes(val int64) []byte {
	tmp := uint64(val)
	buff := make([]byte, 8)
	binary.LittleEndian.PutUint64(buff, tmp)
	return buff
}

func BytesToInt64(data []byte) int64 {
	buff := make([]byte, 8)
	copy(buff, data)
	tmp := binary.LittleEndian.Uint64(buff)
	return int64(tmp)
}

// Float32ToByte
// @Description: 转化float
// @param val float32
// @return []byte
func Float32ToByte(val float32) []byte {
	tmp := math.Float32bits(val)
	buff := make([]byte, 4)
	binary.LittleEndian.PutUint32(buff, tmp)
	return buff
}

func BytesToFloat32(data []byte) float32 {
	buff := make([]byte, 4)
	copy(buff, data)
	tmp := binary.LittleEndian.Uint32(buff)
	return math.Float32frombits(tmp)
}

func GetDBTime(strTime string) *time.Time {
	DefaultTimeLoc := time.Local
	loginTime, _ := time.ParseInLocation("2006-01-02 15:04:05", strTime, DefaultTimeLoc)
	return &loginTime
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func Clamp(val, low, high int) int {
	return int(math.Max(math.Min(float64(val), float64(high)), float64(low)))
}

func GetClassName(param interface{}) string {
	sType := strings.ToLower(reflect.ValueOf(param).Type().String())
	index := strings.Index(sType, ".")
	if index != -1 {
		sType = sType[index+1:]
	}
	return sType
}

func ToHash(str string) uint32 {
	return crc32.ChecksumIEEE([]byte(str))
}
