package reurl

import (
    "gorm.io/gorm"
	"fmt"
	"errors"
)

const base62Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
var base62Rev = [128]int{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, -1, -1, -1, -1, -1, -1, -1, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, -1, -1, -1, -1, -1, -1, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, -1, -1, -1, -1, -1}
const maxGeneratedKeyLength = 3
const invalidGeneratedKeyVal = 62*62*62

func base62CharToValue(c byte) uint {
	var ret = base62Rev[c]
	if ret == -1 {
		fmt.Printf("invalid base62 character: %c\n", c)
		return invalidGeneratedKeyVal
	}
	return uint(ret)
}

func base62ValueToChar(v uint) byte {
	if v >= 62 {
		fmt.Printf("invalid base62 value: %d\n", v)
		return base62Alphabet[0]
	}
	return base62Alphabet[v]
}

// base62Encode encodes an unsigned integer into base62 string (no padding).
func base62Encode(n uint) string {
    if n == 0 {
        return string(base62ValueToChar(0))
    }
    s := ""
    for n > 0 {
        rem := n % 62
        s = string(base62ValueToChar(rem)) + s
        n = n / 62
    }
    return s
}

func base62Decode(s string) uint {
	var n uint = 0
	for i := 0; i < len(s); i++ {
		if base62CharToValue(s[i]) == invalidGeneratedKeyVal {
			fmt.Printf("invalid base62 string: %s\n", s)
			return invalidGeneratedKeyVal
		}
		n = n*62 + base62CharToValue(s[i])
	}
	return n
}

// 經過思考後決定使用暴力搜尋法，因為判斷對於流量較低的網站而言，這麼做是比較快的
func generateKey(db *gorm.DB) (string, error) {
	var existing []string
	// use LENGTH() which works for common DBs (MySQL, SQLite, Postgres)
	if err := db.Raw("SELECT `key` FROM reurls WHERE LENGTH(`key`) <= ?", maxGeneratedKeyLength).Scan(&existing).Error; err != nil {
		return "", err
	}

	exists := make(map[uint]struct{}, len(existing))
	for _, k := range existing {
		exists[base62Decode(k)] = struct{}{}
	}

	var keyVal uint = invalidGeneratedKeyVal
	for i := 0; i < invalidGeneratedKeyVal; i++ {
		if _, found := exists[uint(i)]; !found {
			keyVal = uint(i)
			break
		}
	}

	if keyVal == invalidGeneratedKeyVal {
		return "", errors.New("no available keys")
	}

	return base62Encode(keyVal), nil
}
