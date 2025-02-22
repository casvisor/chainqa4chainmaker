package service

import (
	"errors"
	"fmt"

	// AES
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

// AES加解密
/**
 * @Description: AES解密
 * @author: jjq
 * @param ciphertext 密文
 * @param key 密钥
 * @return string 明文
 * @return error 错误信息
 */
func AesDecrypt(ciphertext, key string) (string, error) {
	byteKey := []byte(key)
	decodedCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(byteKey)
	if err != nil {
		return "", err
	}

	if len(decodedCiphertext) < aes.BlockSize {
		return "", errors.New("AES解密失败，密文太短。")
	}

	iv := decodedCiphertext[:aes.BlockSize]
	decodedCiphertext = decodedCiphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(decodedCiphertext, decodedCiphertext)

	// 移除PKCS#7填充
	paddingLen := int(decodedCiphertext[len(decodedCiphertext)-1])
	if paddingLen > len(decodedCiphertext) || paddingLen > aes.BlockSize {
		return "", errors.New("AES解密失败，无效的填充长度。")
	}
	for i := 0; i < paddingLen; i++ {
		if decodedCiphertext[len(decodedCiphertext)-1-i] != byte(paddingLen) {
			return "", errors.New("AES解密失败，无效的填充。")
		}
	}
	decodedCiphertext = decodedCiphertext[:len(decodedCiphertext)-paddingLen]

	return string(decodedCiphertext), nil
}

/**
 * @Description: AES加密
 * @author: jjq
 * @param plaintext 明文
 * @param key 密钥
 * @return string 密文
 * @return error 错误信息
 */
// AesEncrypt 使用AES算法和CFB模式加密明文字符串。
func AesEncrypt(plaintext, key string) (string, error) {
	byteKey := []byte(key)

	// 检查密钥长度是否合适。
	if len(byteKey) != 16 && len(byteKey) != 24 && len(byteKey) != 32 {
		return "", errors.New("密钥长度必须是16、24或32字节")
	}

	// 创建一个新的AES加密器。
	block, err := aes.NewCipher(byteKey)
	if err != nil {
		return "", err
	}

	// 为IV生成随机值。
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// 对明文进行填充，以确保其长度是AES块大小的倍数。
	padding := aes.BlockSize - len(plaintext)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	plaintextPadded := plaintext + string(padtext)

	// 将明文转换为字节切片。
	plaintextBytes := []byte(plaintextPadded)

	// 创建一个CFB加密器并加密明文。
	ciphertext := make([]byte, len(plaintextBytes))
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext, plaintextBytes)

	// 将IV附加到密文的开头。
	ciphertext = append(iv, ciphertext...)

	// 将密文编码为Base64字符串。
	encodedCiphertext := base64.StdEncoding.EncodeToString(ciphertext)

	return encodedCiphertext, nil
}

/**
 * @Description: 生成AES密钥，返回密钥字符串和错误信息
 * @author: jjq
 * @return string 密钥字符串
 * @return error 错误信息
 */
func GenerateAESKey() (string, error) {

	key := make([]byte, 16)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", key), nil
}
