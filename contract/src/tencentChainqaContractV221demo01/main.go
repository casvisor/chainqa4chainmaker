package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"chainmaker/pb/protogo"
	"chainmaker/shim"
)

type ChainQA struct {
}

func (f *ChainQA) InitContract(stub shim.CMStubInterface) protogo.Response {
	return shim.Success([]byte("Init Success"))
}

func (f *ChainQA) InvokeContract(stub shim.CMStubInterface) protogo.Response {
	method := string(stub.GetArgs()["method"])
	switch method {
	case "getPk":
		// 获取公钥（后期改为隐私合约）
		return f.getPK(stub)
	case "updateDataDigtalEnvelop":
		// 上传数字信封
		return f.updateDataDigtalEnvelop(stub)
	case "getAesKey":
		// 解密数字信封，来获取AES密钥（后期改为隐私合约）
		return f.getAesKey(stub)
	case "updateQueryLog":
		// 上传查询日志
		return f.updateQueryLog(stub)
	case "getAllQueryLogByUid":
		// 根据用户ID获取查询日志
		return f.getAllQueryLogByUid(stub)
	case "getAllQueryLogByTimestamp":
		// 根据时间戳获取查询日志
		return f.getAllQueryLogByTimestamp(stub)
	default:
		return shim.Error("invalid method")
	}
}

// ====================== 合约部分 ======================
// getPK：获取公钥
func (f *ChainQA) getPK(stub shim.CMStubInterface) protogo.Response {
	// 获取公钥
	pkString, err := GetPublicKeyString()
	if err != nil {
		return shim.Error(fmt.Sprintf("[chainqa getPK CONTRACT]获取公钥失败：: %s", err))
	}
	return shim.Success([]byte(pkString))
}

// updateDataDigtalEnvelop：更新数字信封
// 数字信封
type DataDigtalEnvelop struct {
	Uid       string //用户ID
	TimeStamp string //时间戳,日志中仍然按照原始类型的时间戳存储,例如：1735693850
	Pos       string //文件位置
	Envelop   string //数字信封(JSON)
}

/**
 * 更新数字信封
 * @param uId 用户ID
 * @param pos 文件位置
 * @param envelop 数字信封（json字符串） 格式为:[{"AesKeyCipBase64":"密文"}] 【注意是数组哦，后期可能会扩展】
 * 数字信封是一个json数组，包含了一个“AES密钥”的密文。“AES密钥”用于解密IPFS文件，使用RSA加密“AES密钥”，将其存储在数字信封中
 */
func (f *ChainQA) updateDataDigtalEnvelop(stub shim.CMStubInterface) protogo.Response {
	// ----------------- 获取参数 -----------------
	params := stub.GetArgs()
	// 获取参数
	uIdStr := string(params["uId"]) //用户ID
	// 初始化时间戳
	timestampNumberStr, err := stub.GetTxTimeStamp() //时间戳
	if err != nil {
		return shim.Error("[chainqa updateDataDigtalEnvelop CONTRACT]时间戳获取失败")
	}
	// 文件位置
	posStr := string(params["pos"]) //文件位置
	// 数字信封
	envelopStr := string(params["envelop"]) //数字信封

	// 判断前三个是否为空
	if uIdStr == "" || timestampNumberStr == "" || posStr == "" || envelopStr == "" {
		return shim.Error("[chainqa updateDataDigtalEnvelop CONTRACT]参数不能为空")
	}

	// ----------------- 序列化数字信封 -----------------
	// 初始化一个数字信封对象
	dataDigtalEnvelop := DataDigtalEnvelop{
		Uid:       uIdStr,
		TimeStamp: timestampNumberStr, // 时间戳
		Pos:       posStr,
		Envelop:   envelopStr,
	}
	// 将数字信封对象序列化为JSON
	dataDigtalEnvelopBytes, err := json.Marshal(dataDigtalEnvelop)
	if err != nil {
		return shim.Error(fmt.Sprintf("[chainqa updateDataDigtalEnvelop CONTRACT]序列化失败：: %s", err))
	}

	// ----------------- 写入区块链 -----------------
	// 将数字信封对象写入区块链
	err = stub.PutStateByte("chain_data_digtal_envelop", posStr, dataDigtalEnvelopBytes) //key为pos，value为数字信封对象. 因为pos是唯一的，所以不会重复。
	if err != nil {
		return shim.Error(fmt.Sprintf("[chainqa updateDataDigtalEnvelop CONTRACT]写入区块链失败：: %s", err))
	}
	return shim.Success([]byte("上传数字信封成功"))
}

/**
 * 解密数字信封获取AES密钥
 * @param pos 文件前缀路径(一般是合约名)
 */
func (f *ChainQA) getAesKey(stub shim.CMStubInterface) protogo.Response {
	// ----------------- 获取参数 -----------------
	params := stub.GetArgs()
	pos := string(params["pos"]) //文件前缀路径(一般是合约名)
	if pos == "" {
		return shim.Error("[chainqa getPK CONTRACT]pos参数不能为空")
	}

	// ----------------- 根据pos找出数字信封 -----------------
	// 根据pos找出数字信封
	dataDigtalEnvelopBytes, err := stub.GetStateByte("chain_data_digtal_envelop", pos)
	if err != nil {
		return shim.Error(fmt.Sprintf("[chainqa chainQuery CONTRACT]获取数字信封失败：: %s", err))
	}
	if dataDigtalEnvelopBytes == nil || len(dataDigtalEnvelopBytes) == 0 {
		return shim.Error(fmt.Sprintf("[chainqa chainQuery CONTRACT]数字信封不存在"))
	}
	dataDigtalEnvelop := DataDigtalEnvelop{}
	err = json.Unmarshal(dataDigtalEnvelopBytes, &dataDigtalEnvelop)
	if err != nil {
		return shim.Error(fmt.Sprintf("[chainqa chainQuery CONTRACT]反序列化失败：: %s", err))
	}

	// ----------解密数字信封----------
	aesKey, err := MatchEnvelopAndDecrpty(dataDigtalEnvelop.Envelop)
	if err != nil {
		return shim.Error(fmt.Sprintf("[chainqa chainQuery CONTRACT]解密数字信封失败：: %s", err))
	}
	return shim.Success([]byte(aesKey))
}

// [日志部分]：其实不如让区块链管理系统来做这个事情，这里只是一个简单的示例

// getQueryLog：获取查询日志
type QueryLog struct {
	QueryId     string //查询ID
	Uid         string //用户ID
	Timestamp   string //时间戳，日志中仍然按照原始类型的时间戳存储，例如：1735693850
	QueryItem   string //查询项
	QueryStatus int    //查询状态
	QueryResult string //查询结果
}

/**
 * 上传查询日志
 * @param uId 用户ID
 * @param queryItem 查询项
 * @param queryStatus 查询状态
 * @param queryResult 查询结果
 */
func (f *ChainQA) updateQueryLog(stub shim.CMStubInterface) protogo.Response {
	// ----------------- 获取参数 -----------------
	params := stub.GetArgs()
	// 获取参数

	uId := string(params["uId"])                 //用户ID
	queryItem := string(params["queryItem"])     //查询项
	queryStatus := string(params["queryStatus"]) //查询状态
	// 将queryStatus转为int
	queryStatusInt, err := strconv.Atoi(queryStatus)
	queryResult := string(params["queryResult"]) //查询结果
	// 初始化时间戳
	timestampNumberStr, err := stub.GetTxTimeStamp() //时间戳
	if err != nil {
		return shim.Error("[chainqa updateQueryLog CONTRACT]时间戳获取失败")
	}
	// 初始化查询ID
	queryId := StringToMD5(uId + timestampNumberStr) //查询ID
	timestampStrUnderline, err := StringTimestampToUnderlineTimeStamp(timestampNumberStr)
	if err != nil {
		return shim.Error("[chainqa updateQueryLog CONTRACT]时间戳转换失败")
	}

	// ----------------- 序列化查询日志 -----------------
	// 初始化查询日志对象
	QueryLog := QueryLog{
		QueryId:     queryId,
		Uid:         uId,
		Timestamp:   timestampNumberStr, // 日志中仍然按照原始时间戳存储
		QueryItem:   queryItem,
		QueryStatus: queryStatusInt,
		QueryResult: queryResult,
	}
	// 将查询日志对象序列化为JSON
	QueryLogBytes, err := json.Marshal(QueryLog)
	if err != nil {
		return shim.Error(fmt.Sprintf("[chainqa updateQueryLog CONTRACT]序列化失败：: %s", err))
	}
	// 字符串转byte
	QueryIdBytes := []byte(queryId)

	// ----------------- 写入区块链 -----------------
	// 将查询日志对象写入区块链
	err = stub.PutStateByte("chain_query_log", queryId, QueryLogBytes) //key为queryId， 因为queryId是唯一的，所以不会重复。
	if err != nil {
		return shim.Error(fmt.Sprintf("[chainqa updateQueryLog CONTRACT]写入区块链失败：: %s", err))
	}

	// 辅助查询
	// 按时间查询
	//* 因为范围查询时，必须精确指定field，所以如果按时间查询，只能以时间为field。时间戳仅精确到秒，所以如果一秒内有多次并发，会覆盖！
	err = stub.PutStateByte("logHash_timestamp", timestampStrUnderline, QueryIdBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("[chainqa updateQueryLog CONTRACT]按时间写入区块链失败：: %s", err))
	}
	// 按用户查询，存储用户ID和时间戳
	err = stub.PutStateByte("logHash_uId", uId+"__"+timestampStrUnderline, QueryIdBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("[chainqa updateQueryLog CONTRACT]按用户写入区块链失败：: %s", err))
	}
	return shim.Success([]byte("上传查询日志成功"))

}

/**
 * 根据用户ID获取查询日志
 * @param uId 用户ID
 */
// * 建议uId不要以"__"结尾，否则会引起查询到别人的查询日志
func (f *ChainQA) getAllQueryLogByUid(stub shim.CMStubInterface) protogo.Response {
	// ----------------- 获取参数 -----------------
	params := stub.GetArgs()
	// 获取参数
	uId := string(params["uId"]) //用户ID
	if uId == "" {
		return shim.Error("[chainqa getAllQueryLogByUid CONTRACT]uId参数不能为空")
	}

	// ----------------- 根据uId找出查询日志 -----------------
	// 定义一个数组，用于存储queryId
	queryIdArray := []string{}
	// 前缀为uId+"__"的查询日志。建议uId不要以"__"结尾
	rsKv, err := stub.NewIteratorPrefixWithKeyField("logHash_uId", uId+"__")
	if err != nil {
		return shim.Error(fmt.Sprintf("[chainqa getAllQueryLogByUid CONTRACT]获取查询日志失败：: %s", err))
	}
	for rsKv.HasNext() {
		kvkey, kvfeild, kvvalue, err := rsKv.Next()
		if err != nil || kvkey == "" || kvfeild == "" || kvvalue == nil {
			continue
		}
		queryIdArray = append(queryIdArray, string(kvvalue))
	}

	// 定义一个数组，用于存储查询日志
	queryLogArray := []QueryLog{}
	counts := 0
	for _, queryId := range queryIdArray {
		queryLogBytes, err := stub.GetStateByte("chain_query_log", queryId)
		if err != nil {
			continue
		}
		var queryLog QueryLog
		err = json.Unmarshal(queryLogBytes, &queryLog)
		if err != nil {
			continue
		}
		queryLogArray = append(queryLogArray, queryLog)
		counts++
	}
	type QueryLogArray struct {
		Count         int
		QueryLogArray []QueryLog
	}
	QueryLogSearchResult := QueryLogArray{
		Count:         counts,
		QueryLogArray: queryLogArray,
	}
	QueryLogSearchResultBytes, err := json.Marshal(QueryLogSearchResult)
	if err != nil {
		return shim.Error(fmt.Sprintf("[chainqa getAllQueryLogByUid CONTRACT]序列化失败：: %s", err))
	}
	return shim.Success(QueryLogSearchResultBytes)

}

/**
 * 根据时间戳获取查询日志
 * @param startTime 开始时间戳
 * @param endTime 结束时间戳（若为Now，则获取当前时间）
 */
//* 本方法使用了范围查询。但范围查询时，必须精确指定field，所以如果按时间查询，只能以时间为field。时间戳仅精确到秒，所以如果一秒内有多次并发，会覆盖！
func (f *ChainQA) getAllQueryLogByTimestamp(stub shim.CMStubInterface) protogo.Response {
	// ----------------- 获取参数 -----------------
	params := stub.GetArgs()
	// 获取参数
	startTimeStamp := string(params["startTime"]) //时间戳
	endTimeStamp := string(params["endTime"])     //时间戳
	if startTimeStamp == "" {
		return shim.Error(fmt.Sprintf("[chainqa getAllQueryLogByTimestamp CONTRACT]开始时间不能为空"))
	}

	if endTimeStamp == "" || endTimeStamp == "Now" {
		// 获取当前时间
		var err error
		endTimeStamp, err = stub.GetTxTimeStamp()
		if err != nil {
			return shim.Error(fmt.Sprintf("[chainqa getAllQueryLogByTimestamp CONTRACT]时间戳获取失败：: %s", err))
		}
	}

	// 时间戳转换
	startTime, err := StringTimestampToUnderlineTimeStamp(startTimeStamp)
	if err != nil {
		return shim.Error(fmt.Sprintf("[chainqa getAllQueryLogByTimestamp CONTRACT]时间戳转换失败：: %s", err))
	}
	endTime, err := StringTimestampToUnderlineTimeStamp(endTimeStamp)
	if err != nil {
		return shim.Error(fmt.Sprintf("[chainqa getAllQueryLogByTimestamp CONTRACT]时间戳转换失败：: %s", err))
	}

	// ----------------- 根据时间戳找出查询日志 -----------------
	// 定义一个数组，用于存储queryId
	queryIdArray := []string{}
	rsKv, err := stub.NewIteratorWithField("logHash_timestamp", startTime, endTime)
	if err != nil {
		return shim.Error(fmt.Sprintf("[chainqa getAllQueryLogByTimestamp CONTRACT]获取查询日志失败：: %s", err))
	}
	for rsKv.HasNext() {
		kvkey, kvfeild, kvvalue, err := rsKv.Next()
		if err != nil || kvkey == "" || kvfeild == "" || kvvalue == nil {
			continue
		}
		queryIdArray = append(queryIdArray, string(kvvalue))

	}
	// 定义一个数组，用于存储查询日志
	queryLogArray := []QueryLog{}
	counts := 0
	for _, queryId := range queryIdArray {
		queryLogBytes, err := stub.GetStateByte("chain_query_log", queryId)
		if err != nil {
			continue
		}
		var queryLog QueryLog
		err = json.Unmarshal(queryLogBytes, &queryLog)
		if err != nil {
			continue
		}
		queryLogArray = append(queryLogArray, queryLog)
		counts++
	}
	type QueryLogArray struct {
		Count         int
		QueryLogArray []QueryLog
	}
	QueryLogSearchResult := QueryLogArray{
		Count:         counts,
		QueryLogArray: queryLogArray,
	}
	QueryLogSearchResultBytes, err := json.Marshal(QueryLogSearchResult)
	if err != nil {
		return shim.Error(fmt.Sprintf("[chainqa getAllQueryLogByTimestamp CONTRACT]序列化失败：: %s", err))
	}
	return shim.Success(QueryLogSearchResultBytes)

}

// ====================== 废弃 ======================
// func (f *ChainQA) save(stub shim.CMStubInterface) protogo.Response {
// 	params := stub.GetArgs()
// 	// 获取参数
// 	fileHash := string(params["file_hash"])
// 	fileName := string(params["file_name"])
// 	timeStr := string(params["time"])
// 	time, err := strconv.Atoi(timeStr)
// 	if err != nil {
// 		msg := "time is [" + timeStr + "] not int"
// 		stub.Log(msg)
// 		return shim.Error(msg)
// 	}
// 	// 构建结构体
// 	fact := NewFact(fileHash, fileName, int32(time))
// 	// 序列化
// 	factBytes, _ := json.Marshal(fact)
// 	// 发送事件
// 	stub.EmitEvent("topic_vx", []string{fact.FileHash, fact.FileName})
// 	// 存储数据
// 	err = stub.PutStateByte("fact_bytes", fact.FileHash, factBytes)
// 	if err != nil {
// 		return shim.Error("fail to save fact bytes")
// 	}
// 	// 记录日志
// 	stub.Log("[save] FileHash=" + fact.FileHash)
// 	stub.Log("[save] FileName=" + fact.FileName)
// 	// 返回结果
// 	return shim.Success([]byte(fact.FileName + fact.FileHash))

// }

// func (f *ChainQA) findByFileHash(stub shim.CMStubInterface) protogo.Response {
// 	// 获取参数
// 	FileHash := string(stub.GetArgs()["file_hash"])
// 	// 查询结果
// 	result, err := stub.GetStateByte("fact_bytes", FileHash)
// 	if err != nil {
// 		return shim.Error("failed to call get_state")
// 	}
// 	// 反序列化
// 	var fact Fact
// 	_ = json.Unmarshal(result, &fact)
// 	// 记录日志
// 	stub.Log("[find_by_file_hash] FileHash=" + fact.FileHash)
// 	stub.Log("[find_by_file_hash] FileName=" + fact.FileName)
// 	// 返回结果
// 	return shim.Success(result)
// }

func main() {
	err := shim.Start(new(ChainQA))
	if err != nil {
		log.Fatal(err)
	}
}
