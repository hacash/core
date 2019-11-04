package account

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/big"
	"reflect"
	"fmt"
)

const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// EncodeHexString encodes the given version and data to a base58check encoded string
func EncodeHexString(version, data string) (string, error) {
	verbytes, e1 := hex.DecodeString(version)
	if e1 != nil {
		return "", e1
	}
	databytes, e2 := hex.DecodeString(data)
	if e2 != nil {
		return "", e2
	}
	return EncodeWithVersion(verbytes, databytes), nil
}

func EncodeWithVersion(version []byte, data []byte) (string) {
	prefix := make([]byte, len(version))
	copy(prefix, version)
	dataBytes := make([]byte, 0)
	dataBytes = append(prefix, data...)
	return Encode(dataBytes)
}

func Encode(dataBytes []byte) (string) {

	// Performing SHA256 twice
	sha256hash := sha256.New()
	sha256hash.Write(dataBytes)
	middleHash := sha256hash.Sum(nil)
	sha256hash = sha256.New()
	sha256hash.Write(middleHash)
	hash := sha256hash.Sum(nil)

	checksum := hash[:4]
	//fmt.Println("checksum", checksum)
	stuff := append([]byte{}, dataBytes...)
	stuff = append(stuff, checksum...)


	// For all the "00" versions or any prepended zeros as base58 removes them
	zeroCount := 0
	for i := 0; i < len(stuff); i++ {
		if stuff[i] == 0 {
			zeroCount++
		} else {
			break
		}
	}

	// Performing base58 encoding
	encoded := b58encode(stuff)

	for i := 0; i < zeroCount; i++ {
		encoded = "1" + encoded
	}

	return encoded
}



// DecodeHexString decodes the given base58check encoded string and returns the version prepended decoded string
func Decode(encoded string) ([]byte, error) {
	zeroCount := 0
	for i := 0; i < len(encoded); i++ {
		if encoded[i] == 49 {
			zeroCount++
		} else {
			break
		}
	}

	dataBytes, err := b58decode(encoded)
	if err != nil {
		return nil, err
	}

	dataBytesLen := len(dataBytes)
	if dataBytesLen <= 4 {
		return nil, errors.New("base58check data cannot be less than 4 bytes")
	}

	data, checksum := dataBytes[:dataBytesLen-4], dataBytes[dataBytesLen-4:]

	for i := 0; i < zeroCount; i++ {
		data = append([]byte{0}, data...)
	}

	// Performing SHA256 twice to validate checksum
	sha256hash := sha256.New()
	sha256hash.Write(data)
	middleHash := sha256hash.Sum(nil)
	sha256hash = sha256.New()
	sha256hash.Write(middleHash)
	hash := sha256hash.Sum(nil)

	if !reflect.DeepEqual(checksum, hash[:4]) {
		return nil, errors.New("Data and checksum don't match")
	}

	return data, nil
}


// DecodeHexString decodes the given base58check encoded string and returns the version prepended decoded string
func DecodeTestPrint(encoded string) ([]byte, error) {
	zeroCount := 0
	for i := 0; i < len(encoded); i++ {
		if encoded[i] == 49 {
			zeroCount++
		} else {
			break
		}
	}

	dataBytes, err := b58decode(encoded)
	if err != nil {
		return nil, err
	}

	dataBytesLen := len(dataBytes)
	if dataBytesLen <= 4 {
		return nil, errors.New("base58check data cannot be less than 4 bytes")
	}

	data, checksum := dataBytes[:dataBytesLen-4], dataBytes[dataBytesLen-4:]

	for i := 0; i < zeroCount; i++ {
		data = append([]byte{0}, data...)
	}

	// Performing SHA256 twice to validate checksum
	sha256hash := sha256.New()
	sha256hash.Write(data)
	middleHash := sha256hash.Sum(nil)
	sha256hash = sha256.New()
	sha256hash.Write(middleHash)
	hash := sha256hash.Sum(nil)

	if !reflect.DeepEqual(checksum, hash[:4]) {

		fmt.Println("----data----", hex.EncodeToString(data), data)
		fmt.Println("right addr =", Encode(data))

		return nil, errors.New("Data and checksum don't match")
	}

	fmt.Println(encoded, " address ok")

	return data, nil
}




func b58encode(data []byte) string {
	var encoded string
	decimalData := new(big.Int)
	decimalData.SetBytes(data)
	divisor, zero := big.NewInt(58), big.NewInt(0)

	for decimalData.Cmp(zero) > 0 {
		mod := new(big.Int)
		decimalData.DivMod(decimalData, divisor, mod)
		encoded = string(alphabet[mod.Int64()]) + encoded
		//fmt.Println(encoded)
	}

	return encoded
}

func b58decode(data string) ([]byte, error) {
	decimalData := new(big.Int)
	alphabetBytes := []byte(alphabet)
	multiplier := big.NewInt(58)

	for _, value := range data {
		pos := bytes.IndexByte(alphabetBytes, byte(value))
		if pos == -1 {
			return nil, errors.New("Character not found in alphabet")
		}
		decimalData.Mul(decimalData, multiplier)
		decimalData.Add(decimalData, big.NewInt(int64(pos)))
	}

	return decimalData.Bytes(), nil
}
