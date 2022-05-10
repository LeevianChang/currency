package tool

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

//
//const (
//	aesKey = "HQECux7Tt6UrGOUl"
//	gsmIV  = "000000010000010000000010"
//)
//
//func GcmEncrypt(key, plaintext string) (string, error) {
//	keyByte := []byte(key)
//	plainByte := []byte(plaintext)
//	block, err := aes.NewCipher(keyByte)
//	if err != nil {
//		return "", err
//	}
//	aesGcm, err := cipher.NewGCM(block)
//	if err != nil {
//		return "", err
//	}
//	nonce := make([]byte, 12)
//	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
//		return "", err
//	}
//	seal := aesGcm.Seal(nonce, nonce, plainByte, nil)
//	return base64.URLEncoding.EncodeToString(seal), nil
//}
//
//func GcmDecrypt(key, cipherText string) (string, error) {
//	cipherByte, err := base64.URLEncoding.DecodeString(cipherText)
//	if err != nil {
//		return "", err
//	}
//	nonce, cipherByte := cipherByte[:12], cipherByte[12:]
//	keyByte := []byte(key)
//	block, err := aes.NewCipher(keyByte)
//	if err != nil {
//		return "", err
//	}
//	aesGcm, err := cipher.NewGCM(block)
//	if err != nil {
//		return "", err
//	}
//	plainByte, err := aesGcm.Open(nil, nonce, cipherByte, nil)
//	if err != nil {
//		return "", err
//	}
//	return string(plainByte), nil
//}

//cbc建议使用PKCS#7或PKCS#5，这里使用PKCS#7
func pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func pkcs7UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	if length == 0 {
		return nil, errors.New("wrong encryption parameters")
	} else {
		unPadding := int(origData[length-1])
		return origData[:(length - unPadding)], nil
	}
}

func aes128EncryptPKCS7UnPadding(origData []byte, key []byte, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	origData = pkcs7Padding(origData, blockSize)

	//使用cbc
	aesGcm, _ := cipher.NewGCM(block)
	//encrypted := make([]byte, len(origData))
	//nonce := make([]byte, aesGcm.NonceSize())
	//if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
	//}
	cipherText := aesGcm.Seal(nil, nonce, origData, nil)
	return cipherText, nil
}

func aes128DecryptPKCS7UnPadding(cypted []byte, key []byte, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesGcm, _ := cipher.NewGCM(block)
	//cipherText, _ := hex.DecodeString(string(cypted))

	//nonce := make([]byte, aesGcm.NonceSize())
	//if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
	//}

	//nonce, cypted := cypted[:12], cypted[12:]
	plaintext, err := aesGcm.Open(nil, nonce, cypted, nil)
	if err != nil {
		return nil, err
	}
	origData, err := pkcs7UnPadding(plaintext)
	if err != nil {
		return nil, err
	}
	return origData, err
}

//cbc+PKCS7+16位key+16位偏移量（直接用key）
func GcmEncrypt(input, key string, iv []byte) (string, error) {
	result, err := aes128EncryptPKCS7UnPadding([]byte(input), []byte(key), iv)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(result), err
}

func GcmDecrypt(input, key string, iv []byte) (string, error) {
	pwdByte, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", err
	}
	res, err := aes128DecryptPKCS7UnPadding(pwdByte, []byte(key), iv)
	return string(res), err
}
