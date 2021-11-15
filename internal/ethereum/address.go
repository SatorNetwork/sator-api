package ethereum

import "encoding/hex"

const EthereumAddressLen = 42

func IsEthereumAddress(addr string) bool {
	if len(addr) != EthereumAddressLen {
		return false
	}

	if _, err := hex.DecodeString(addr); err != nil {
		return false
	}

	return true
}
