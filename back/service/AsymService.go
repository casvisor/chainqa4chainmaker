package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
)

// ====================== 链下调用函数 ======================

/**
 * RSA加密字节数组,返回字节数组
 * @param originalBytes 需要加密的字节数组
 * @param pkString 公钥字符串
 * @return string 加密后的字符串的数字信封数组。id和aesSecretKey
 * @return error 错误信息
 * @Note: 公钥字符串为提取的公钥信息，不是公钥文件
 * @调用时间：调用时，传入公钥字符串，返回加密后的字符串.用于自行加密数据。形成数字信封
 */
type AesEnvelope struct {
	// File_prefix_path string // 文件前缀路径(一般是合约名)
	AesKeyCipBase64 string // AES密钥(被RSA加密)
}

func RSAEncryptAndReturnEnvelop(file_prefix_path, originalText, pkString string) (string, error) {
	// 1、读取公钥文件,解析出公钥对象
	publicKey, err := ReadParsePublicKeyFromString(pkString)
	if err != nil {
		return "", err
	}
	// 2、RSA加密,参数是随机数、公钥对象、需要加密的字节
	// Reader是一个全局共享的密码安全的强大的伪随机生成器
	cipherBytes, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(originalText))
	if err != nil {
		return "", err
	}

	// 转为AesEnvelope格式
	envelope := AesEnvelope{
		// File_prefix_path: file_prefix_path,                               // 文件前缀路径(一般是合约名)
		AesKeyCipBase64: base64.StdEncoding.EncodeToString(cipherBytes), // aeskey 密文（被RSA加密） 然后base64编码
	}
	// fmt.Println("aesKeyCipBase64:", base64.StdEncoding.EncodeToString(cipherBytes))
	// envelop转为一个数组
	envelopeArray := []AesEnvelope{envelope}
	// envelop转为json
	envelopArrayJson, err := json.Marshal(envelopeArray)
	if err != nil {

		return "", err
	}
	return string(envelopArrayJson), nil
}

func ReadParsePublicKeyFromString(publicKeyString string) (*rsa.PublicKey, error) {
	publicKeyBytes, _ := base64.StdEncoding.DecodeString(publicKeyString)
	// 2、解码公钥字节,生成加密块对象
	block, _ := pem.Decode(publicKeyBytes)
	if block == nil {
		return nil, errors.New("公钥信息错误！")
	}
	// 3、解析DER编码的公钥,生成公钥接口
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 4、公钥接口转型成公钥对象
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	return publicKey, nil
}
