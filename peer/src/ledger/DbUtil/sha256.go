package DbUtil

import (
	"crypto/sha256"
	"encoding/hex"
)

func Str_Sha256(workInfo []string) []byte {
	str_to_hash := ""
	for _, value := range workInfo {
		str_to_hash = str_to_hash + value

	}
	hash := sha256.New()
	hash.Write([]byte(str_to_hash))
	md := hash.Sum(nil)
	//	mdStr := hex.EncodeToString(md)
	return md
}
func Str_Sha256_String(workInfo []string) string {
	str_to_hash := ""
	for _, value := range workInfo {
		str_to_hash = str_to_hash + value

	}
	hash := sha256.New()
	hash.Write([]byte(str_to_hash))
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)
	return mdStr
}
