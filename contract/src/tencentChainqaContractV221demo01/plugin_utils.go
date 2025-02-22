package main

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"
	"time"
)

/**
 * StringToMD5：字符串转MD5
 * @param input string：输入字符串
 * @return string：MD5字符串
 */
func StringToMD5(input string) string {
	// 创建一个新的MD5哈希对象
	hasher := md5.New()
	// 写入要哈希的字符串
	hasher.Write([]byte(input))
	// 计算哈希值并返回字节数组
	bytes := hasher.Sum(nil)
	// 将字节数组转换为十六进制字符串
	md5Str := hex.EncodeToString(bytes)
	return md5Str
}

/**
 * StringTimestampToDateTime：字符串时间戳转日期时间字符串
 * @param timestamp：时间戳 （1735693850）
 * @return string：日期时间字符串（格式：YYYY_MM_DD_HH_mm_ss）
 * @return error：错误信息
 * 长安链给出的时间戳为例如：1735693850，需要转换为：2021_12_12_12_12_12。因为为了方便查询，上链的field只能包含字母、数字、下划线。如果直接用数字型时间戳，因为长安链比较的时候，是按照字符串比较的，所以会导致时间戳比较的时候，无法按照时间顺序来比较。关于时间戳，请具体参考：http://shijianchuo.wiicha.com/
 */
func StringTimestampToUnderlineTimeStamp(timestampStr string) (string, error) {
	// 将字符串转换为int64
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return "", err // 如果转换失败，返回错误
	}

	// 使用time.Unix将时间戳转换为time.Time对象
	tm := time.Unix(timestamp, 0)

	// 使用Format方法将time.Time对象格式化为字符串
	formattedTime := tm.Format("2006-01-02 15:04:05")

	// 替换日期和时间之间的空格为下划线
	formattedTimeWithUnderscores := strings.ReplaceAll(formattedTime, " ", "_")
	// 替换日期中的-为下划线
	formattedTimeWithUnderscores = strings.ReplaceAll(formattedTimeWithUnderscores, "-", "_")
	// 替换时间中的:为下划线
	formattedTimeWithUnderscores = strings.ReplaceAll(formattedTimeWithUnderscores, ":", "_")
	return formattedTimeWithUnderscores, nil
}
