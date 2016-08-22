// encrypt
package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

//---------------------crypto-----------------------
func PKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

// AES
func AesEncrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	paddingData := PKCS5Padding(data, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	cryptData := make([]byte, len(paddingData))
	blockMode.CryptBlocks(cryptData, paddingData)
	return cryptData, nil
}

func AesDecrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	paddingData := make([]byte, len(data))
	blockMode.CryptBlocks(paddingData, data)
	originData := PKCS5UnPadding(paddingData)
	return originData, nil
}

// base64.
func UrlBase64Encode(value []byte) string {
	encoding := base64.URLEncoding
	encoded := make([]byte, encoding.EncodedLen(len(value)))
	encoding.Encode(encoded, value)
	return string(encoded)
}

func UrlBase64Decode(value []byte) ([]byte, error) {
	pad := len(value) % 4
	if pad != 0 {
		for i := 0; i < 4-pad; i++ {
			value = append(value, '=')
		}
	}

	encoding := base64.URLEncoding
	decoded := make([]byte, encoding.DecodedLen(len(value)))
	b, err := encoding.Decode(decoded, value)
	if err != nil {
		return nil, err
	}
	return decoded[:b], nil
}

// md5
func MD5(data []byte) []byte {
	m := md5.New()
	m.Write(data)
	return m.Sum(nil)
}

func MD5Hex(data []byte) string {
	out := MD5(data)
	return hex.EncodeToString(out)
}
