package util

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

func AesDecrypt(decodeStr string, key string,iv string) ([]byte, error) {
	//先解密base64
	decodeBytes, err := base64.StdEncoding.DecodeString(decodeStr)
	if err != nil {
		return nil, err
	}

	baseKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(baseKey)
	if err != nil {
		return nil, err
	}

	ivBase4, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, ivBase4)
	origData := make([]byte, len(decodeBytes))

	blockMode.CryptBlocks(origData, decodeBytes)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
