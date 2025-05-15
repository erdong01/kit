package uuid

import (
	"math/rand/v2"
	"strings"

	"github.com/google/uuid"
)

func New() (uid string) {
	id, err := uuid.NewV7()
	if err != nil {
		return ""
	}
	return strings.ReplaceAll(id.String(), "-", "")
}

func GenerateNumber(l int) string {
	return GenerateNo(l)
}

// 生成编号
func GenerateNo(l int, head ...string) string {
	strDigits := "0123456789" // 用于从计算值映射到数字字符
	sb := strings.Builder{}
	sb.Grow(l) // 预分配容量
	if len(head) > 0 {
		sb.WriteString(head[0])
	}
	// 如果head本身就达到或超过了所需长度l
	if sb.Len() >= l {
		return sb.String()[:l]
	}

	// 获取基于时间的UUID v7字符串
	uuidStr := New()
	if uuidStr == "" {
		// UUID生成失败的后备方案：用随机数填充剩余部分
		// 注意：这部分的ID将不是B+树友好的
		remainingLen := l - sb.Len()
		for i := 0; i < remainingLen; i++ {
			sb.WriteByte(strDigits[rand.IntN(len(strDigits))])
		}
		return sb.String()
	}

	uuidChars := []byte(uuidStr) // UUID字符串的字节表示
	uuidLen := len(uuidChars)
	uuidCharIndex := uuidLen // 当前从uuidChars中取字符的索引

	// 填充ID直到达到长度l
	for sb.Len() < l {
		uuidCharIndex--
		// 从UUID字符（十六进制）派生一个0-9的数字。
		// 直接取字符的字节值对10取模。
		// 例如：'0' (ASCII 48) % 10 = 8
		//       'f' (ASCII 102) % 10 = 2
		// 这种方法虽然不保证数字的均匀分布，但重要的是它保留了源UUID字符的顺序性。
		// （即，如果uuid_char1 < uuid_char2，派生的数字通常也会保持类似的顺序关系，尤其是在高位）
		charByte := uuidChars[uuidCharIndex]
		digit := strDigits[charByte%10] // charByte % 10 的结果范围是0-9

		sb.WriteByte(digit)

		if uuidCharIndex <= 0 {
			// 如果UUID的字符已经用完，但仍需更多数字，则从头开始循环使用UUID字符。
			// 这比纯随机数更能保持一定的顺序性。
			uuidStr = New()
			uuidChars = []byte(uuidStr)
			uuidLen = len(uuidChars)
			uuidCharIndex = uuidLen
		}
	}

	return sb.String()
}
