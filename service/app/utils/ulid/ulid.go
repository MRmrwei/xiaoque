// Package ulid 生成 26 位 ULID，用于页面翻译组 ID 等业务标识。
package ulid

import (
	"crypto/rand"
	"fmt"
	"time"
)

// encoding 是 Crockford Base32 字母表，与标准 ULID 一致。
const encoding = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

// New 生成一个 26 位 ULID：前 10 位是毫秒时间戳，后 16 位是 80 位随机数。
// 同一毫秒内随机数不保证严格单调，但足够避免翻译组 ID 冲突。
func New() (string, error) {
	var entropy [10]byte
	if _, err := rand.Read(entropy[:]); err != nil {
		return "", fmt.Errorf("生成 ULID 失败：%w", err)
	}

	var dst [26]byte

	// 48 位毫秒时间戳编码为 10 个 5 位字符，最高位恒为 0，因此首字符在 0-7 之间。
	timestamp := uint64(time.Now().UnixMilli())
	for i := 9; i >= 0; i-- {
		dst[i] = encoding[timestamp&0x1f]
		timestamp >>= 5
	}

	// 80 位随机数编码为 16 个 5 位字符。
	var buffer uint64
	var bits uint
	index := 10
	for _, value := range entropy {
		buffer = (buffer << 8) | uint64(value)
		bits += 8
		for bits >= 5 {
			bits -= 5
			dst[index] = encoding[(buffer>>bits)&0x1f]
			index++
		}
	}

	return string(dst[:]), nil
}
