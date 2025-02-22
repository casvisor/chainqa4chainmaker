package main

import (
	"encoding/json"
)

func MatchEnvelopAndDecrpty(envelop string) (string, error) {
	// [{"File_prefix_path":"chainQA_test_0_007","AesKeyCipBase64":"d8tpaWWcRdY9Gr2wIqjONg0iAZjtFd6sr9nyRGIVu50VBjnsZcvBfWuBSf4OfWulehdACyqUZEi38HzFzR5gXaDGN9Q28b2kxKHGL173iRsMfsl97asy58WUFgngZ5L9XCh4oGVIrJJggqRVH2OUDceGcslubBOcwCZxiypiqpw="}]
	// 1. 把envelop转换为[]map[string]string
	var envelops []map[string]string
	err := json.Unmarshal([]byte(envelop), &envelops)
	if err != nil {
		return "", err
	}
	// 2. 遍历envelops，获取File_prefix_path和AesKeyCipBase64

	// filePrefixPath := ""
	aesKeyCipBase64 := ""
	for _, item := range envelops {
		// TODO: 根据File_prefix_path的不同值找，这里只有一个值
		// filePrefixPath = item["File_prefix_path"] // 获取文件前缀路径
		aesKeyCipBase64 = item["AesKeyCipBase64"] // 获取AES密钥的密文
	}
	// 3. 解密AES密钥
	aesKey, err := RSADecrypt(aesKeyCipBase64)
	if err != nil {
		return "", err
	}
	return aesKey, nil

}
