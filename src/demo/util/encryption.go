package util

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"unsafe"
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

// 封装API隐藏调用和打乱汇编的方法
func CallAPI(apiName string, args []string) error {
	// 将API名称和参数转换为随机字符串
	randStrs := make([]string, len(args)+1)
	randStrs[0] = generateRandomString(len(apiName))
	for i := 1; i < len(randStrs); i++ {
		randStrs[i] = generateRandomString(len(args[i-1]))
	}

	// 内联汇编调用API
	asm := fmt.Sprintf("call %s", randStrs[0])
	for i := 1; i < len(randStrs); i++ {
		asm += fmt.Sprintf(", %s", randStrs[i])
	}
	asmCode := []byte(asm)
	asmFunc := func() int { return 0 }
	asmFuncPtr := uintptr(unsafe.Pointer(&asmFunc))
	asmPtr := unsafe.Pointer(asmFuncPtr)
	_ = *(*func())(asmPtr)

	// 输出汇编代码和API名称和参数
	fmt.Println("Assembly code:", asmCode)
	fmt.Println("API name:", apiName)
	fmt.Println("Arguments:", args)

	// TODO: 调用实际的API，并返回结果或错误
	return nil
}

// 生成指定长度的随机字符串
func generateRandomString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	randStr := make([]byte, length)
	for i := range randStr {
		randStr[i] = letters[rand.Intn(len(letters))]
	}

	return string(randStr)
}
