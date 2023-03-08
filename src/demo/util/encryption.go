package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
)

// EncryptString 使用 AES 加密算法对字符串进行加密
func EncryptString(key, plaintext string) (string, error) {
	// 将密钥转换为字节数组
	keyBytes := []byte(key)

	// 创建一个 AES 加密算法实例
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// 对字符串进行加密
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))
	encrypted := base64.URLEncoding.EncodeToString(ciphertext)

	// 返回加密后的结果
	return encrypted, nil
}

// DecryptString 使用 AES 解密算法对字符串进行解密
func DecryptString(key, encrypted string) (string, error) {
	// 将密钥转换为字节数组
	keyBytes := []byte(key)

	// 解密字符串
	decoded, err := base64.URLEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}
	if len(decoded) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := decoded[:aes.BlockSize]
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(decoded[aes.BlockSize:], decoded[aes.BlockSize:])
	decrypted := string(decoded[aes.BlockSize:])

	// 返回解密后的结果
	return decrypted, nil
}
