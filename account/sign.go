package account

import (
	"fmt"
	"github.com/hacash/core/crypto/btcec"
)

func CheckSignByHash32(hash32 []byte, publicKeyBytes33 []byte, signatureBytes64 []byte) (bool, error) {
	if len(hash32) != 32 {
		return false, fmt.Errorf("Hash length is not 32")
	}
	sigobj, e3 := btcec.ParseSignatureByte64(signatureBytes64)
	if e3 != nil {
		return false, e3
	}
	pubKey, e4 := btcec.ParsePubKey(publicKeyBytes33, btcec.S256())
	if e4 != nil {
		return false, e4
	}
	verok := sigobj.Verify(hash32, pubKey)
	if !verok {
		address := NewAddressFromPublicKeyV0(publicKeyBytes33)
		return false, fmt.Errorf("Address %s verify signature fail.", Base58CheckEncode(address))
	}
	// ok

	return true, nil
}
