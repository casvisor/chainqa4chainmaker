package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"

	"os"
)

// ====================== 链中调用函数 ======================

/**
 * 获取公钥字符串
 * @return string 公钥字符串
 * @return error 错误信息
 * @Note: 用于获取公钥字符串，传给区块链
 * @调用时间：调用时，返回公钥字符串。用于加密数字信封
 */
func GetPublicKeyString() (string, error) {
	// 作弊方式
	return SECRET_GETPKBASE64(), nil

	// 文件方式
	// // 1、读取公钥文件
	// file, err := os.Open("./" + file_prefix_path + "/chainQA_certs/chainQA_private.pem")
	// if err != nil {
	// 	return "", err
	// }
	// defer file.Close()
	// // 2、读取公钥文件内容
	// derPublicKeyFile, err := ioutil.ReadAll(file)
	// if err != nil {
	// 	return "", err
	// }
	// if derPublicKeyFile == nil {
	// 	return "", errors.New("公钥文件内容为空")
	// }
	// // 3、返回公钥字符串
	// return base64.StdEncoding.EncodeToString(derPublicKeyFile), nil
}

/**
 * RSA解密字符串,返回字符串
 * @param cipherlText 需要解密的字符串，即aesKeyCipBase64
 * @param file_prefix_path 文件前缀路径(一般是合约名)
 * @return string 解密后的字符串
 * @return error 错误信息
 * @调用时间: 智能合约调用，解密解出密文（即IPFS文件的AES密钥）
 */

// RSA解密字节数组,返回字节数组
func RSADecrypt(cipherlText string) (string, error) {
	cipherBytes, _ := base64.StdEncoding.DecodeString(cipherlText)
	if cipherBytes == nil {
		return "", errors.New("解密失败，密文为空")
	}

	// 1、读取私钥文件，解析出私钥对象
	privateKey, err := ReadParsePrivaterKey()
	if err != nil {
		return "", err
	}

	// 2、ras解密,参数是随机数、私钥对象、需要解密的字节
	originalBytes, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherBytes)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return string(originalBytes), nil
}

// 读取私钥文件,解析出私钥对象
func ReadParsePrivaterKey() (*rsa.PrivateKey, error) {

	// 1、读取私钥文件,获取私钥字节
	privateKeyBytes, err := base64.StdEncoding.DecodeString(SECRET_GETSKBASE64())
	// privateKeyBytes, err := ioutil.ReadFile("./" + file_prefix_path + "/chainQA_certs/chainQA_private.pem")
	if err != nil {
		return nil, err
	}
	// 2、对私钥文件进行编码,生成加密块对象
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return nil, errors.New("私钥信息错误！")
	}
	// 3、解析DER编码的私钥,生成私钥对象
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// ====================== 链下调用函数 ======================
/**
 * 生成RSA加密字符串（生成私钥文件和公钥文件）
 * @return error 错误信息
 * @Note: 生成的私钥文件和公钥文件在当前目录下的chainQA_certs文件夹下
 * @调用时间：部署合约前。[部署合约前！！！]
 */
func BeforeCommitContract_GenerateRSAKey() error {
	// 1、RSA生成私钥文件的核心步骤
	// 1)、生成RSA密钥对
	// 密钥长度,默认值为1024位
	bits := 1024
	privateKer, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	// 2)、将私钥对象转换成DER编码形式
	derPrivateKer := x509.MarshalPKCS1PrivateKey(privateKer)
	// 3)、创建私钥pem文件
	file, err := os.Create("./chainQA_certs/chainQA_private.pem")
	if err != nil {
		return err
	}
	// 4)、对密钥信息进行编码,写入到私钥文件中
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derPrivateKer,
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	// 等待文件写入
	defer file.Close()

	// 2、RSA生成公钥文件的核心步骤
	// 1)、生成公钥对象
	publicKey := &privateKer.PublicKey
	// 2)、将公钥对象序列化为DER编码格式
	derPublicKey, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	// 3)、创建公钥pem文件
	file, err = os.Create("./chainQA_certs/chainQA_public.pem")
	if err != nil {
		return err
	}
	// 4)、对公钥信息进行编码,写入到公钥文件中
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPublicKey,
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}
