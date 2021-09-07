package utils

import (
	"encoding/hex"
	"errors"
	"strings"
)

func CleanUpHex(s string) string {
	s = strings.Replace(strings.TrimPrefix(s, "0x"), " ", "", -1)

	return strings.ToLower(s)
}

func ValidateAccount(accountAddress string) (string, error) {
	accountAddress = CleanUpHex(accountAddress)
	// check account length
	if len(accountAddress) != 40 {
		return "", errors.New("invalid account address")
	}

	_, err := hex.DecodeString(accountAddress)
	if err != nil {
		return "", errors.New("invalid account address")
	}

	return NormalizeAddress(accountAddress), nil
}

func NormalizeAddress(addr string) string {
	return "0x" + Trim0x(strings.ToLower(addr))
}

func NormalizeAddresses(addrs []string) []string {
	for k, v := range addrs {
		addrs[k] = NormalizeAddress(v)
	}

	return addrs
}

// Trim0x removes the "0x" prefix of hexes if it exists
func Trim0x(str string) string {
	return strings.TrimPrefix(str, "0x")
}
