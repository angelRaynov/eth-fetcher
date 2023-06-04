package helper

import (
	"fmt"
	"math/big"
	"strings"
)

func DecodeHexBigInt(valueHex string) (*big.Int, error) {
	valueHex = strings.TrimPrefix(valueHex, "0x")
	value := new(big.Int)
	value, success := value.SetString(valueHex, 16)
	if !success {
		return nil, fmt.Errorf("failed to decode value")
	}

	return value, nil
}
