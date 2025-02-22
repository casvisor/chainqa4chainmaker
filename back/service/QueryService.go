package service

import (
	"errors"
	"fmt"

	// 查询
	"encoding/json"

	"regexp"
	"strconv"
	"strings"
	// "chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

// ----------------查询代码（QUERY MODULE）-------------------
// 查询条件结构
// 查询项（一次查询传入的数据）
type QueryItem struct {
	QueryConcatType string             `json:"queryConcatType"` // 查询条件组合类型（AND/OR）
	QueryConditions [][]QueryCondition `json:"queryConditions"` // 查询条件
	FilePos         [][]string         `json:"filePos"`         // 文件在IPFS中的位置(教师文件)
	ReturnField     []string           `json:"returnField"`     // 返回的列名
	JointConditions []JointCondition   `json:"jointConditions"` // 联表查询条件
}

// QueryCondition 定义细化查询条件结构
type QueryCondition struct {
	Field   string `json:"field"`   // 字段名（列名）
	Val     string `json:"val"`     // 基准值
	Pos     string `json:"pos"`     // 位置（多联表的时候来标记）
	Compare string `json:"compare"` // 比较符
	Type    string `json:"type"`    // 类型（string/int/float）
}

// JointCondition 定义联表查询条件结构
type JointCondition struct {
	Pos1      string `json:"pos1"`    // 位置1
	Field1    string `json:"field1"`  // 字段1
	Pos2      string `json:"pos2"`    // 位置2
	Field2    string `json:"field2"`  // 字段2
	Compare   string `json:"compare"` // 比较符
	Type      string `json:"type"`
	JointType string `json:"jointType"` // 联表类型（INNER/OUTER）
}

// 返回结构（实际返回的是这个结构形成的JSON字符串）
type QueryResult struct {
	// 统计结果
	Counts int `json:"counts"` // 若为-1，则表示查询错误；若为0，则表示查询无结果；若为其他值，则表示查询结果数量

	// 返回的data，即查询结果。
	// 采用JSON格式，key为列名，value为该列的值，均以字符串标识。
	Data []interface{} `json:"data"` // 当Counts为0或-1时，为空；当Counts为其他值时，非空

	// 返回的信息，string格式
	Message string `json:"message"` // 主要是：错误信息。当Counts为-1时，非空

}

/**
 * @Description: 解析表头
 * @author: jjq
 * @param tableStr 表字符串(从IPFS读入，为教师上传的文件的内容)
 * @return map[string]int 表头字段名和索引的映射关系
 * @usage: tableHeaderMap["id"]=0。
 */
func parseTableHeader(fix_prefix_pos, tableStr string) map[string]int {
	// 读取文件，按行读取，第一行是表头
	lines := strings.Split(tableStr, "\n")
	// 获取表头，按空格拆分
	headers := strings.Fields(lines[0])
	// tableHeaderMap 用于存储表头字段名和索引序号（第x列，从0开始）的映射关系
	tableHeaderMap := make(map[string]int)
	for i, header := range headers {
		tableHeaderMap[fix_prefix_pos+"_"+header] = i
		// 列名为header ，索引为i
	}

	return tableHeaderMap
	// usage: tableHeaderMap["id"]=0。
}

/**
 * @description: 比较两个map是否相等，即比较两个map的键值对是否完全相同（包括键和值）
 * @param a map[string]int 第一个map
 * @param b map[string]int 第二个map
 * @return bool 如果两个map相等，则返回true，否则返回false
 */
func compareMaps(a, b map[string]int) bool {
	// 如果两个 map 的长度不同，它们肯定不相等
	if len(a) != len(b) {
		return false
	}

	// 遍历 map a，检查每个键值对是否在 map b 中
	for key, value := range a {
		if bValue, exists := b[key]; !exists || bValue != value {
			return false
		}
	}

	// 所有键值对都相等，两个 map 相等
	return true
}

/**
 * matchesConditions 检查行是否匹配查询条件数组
 * @author: jjq
 * @param {string} row 行字符串，即一行数据，将按空格拆分
 * @param tableHeaderMap 表头字段名和索引的映射关系
 * @param conditions 查询条件数组（数组内条件按AND逻辑关联）
 * @return bool 本行是否满足查询条件组匹配要求
 */
func matchesConditions(row string, tableHeaderMap map[string]int, conditions []QueryCondition) (bool, error) {

	// 遍历每个conditions。该condition为同一组conditionGroup的，其为AND关系，必须同时满足。
	// 将row按空格拆分
	cells := strings.Fields(row)
	if len(cells) == 0 {
		return false, nil
	}
	for _, condition := range conditions {
		// 检查字段是否存在
		_, ok := tableHeaderMap[condition.Pos+"_"+condition.Field]
		if !ok {
			errorMsg := fmt.Sprintf("查询条件中 %s 列不存在", condition.Field)
			return false, errors.New(errorMsg)
		}

		// 获取条件基准值
		condVal := condition.Val
		// 获取单元格值：根据表头获取到比较列的索引（从0开始），然后根据索引获取到单元格值

		cellValue := cells[tableHeaderMap[condition.Pos+"_"+condition.Field]]

		// 获取类型
		cellType := condition.Type
		// 获取比较符
		compare := condition.Compare

		// 根据cellType调用对应的函数（cellType可以为int，float，string）
		switch cellType {
		case "int":
			flag, err := compareInt(cellValue, condVal, compare)
			if err != nil {
				return false, err
			} else if !flag && err == nil {
				return false, nil
			} else {
				// 继续
			}
		case "float":
			flag, err := compareFloat(cellValue, condVal, compare)
			if err != nil {
				return false, err
			} else if !flag && err == nil {
				return false, nil
			} else {
				// 继续
			}
		case "string":
			flag, err := compareString(cellValue, condVal, compare)
			if err != nil {
				return false, err
			} else if !flag && err == nil {
				return false, nil
			} else {
				// 继续
			}
		default:
			errorMsg := fmt.Sprintf("查询条件中暂时不支持 %s 列的类型: %s", condition.Field, condition.Type)
			return false, errors.New(errorMsg)
		}

	}
	// 所有条件都通过，才为true！
	return true, nil
}

func compareString(cellValue string, condVal string, compare string) (bool, error) {
	switch compare {
	case "eq":
		return cellValue == condVal, nil
	case "ne":
		return cellValue != condVal, nil
	case "lt":
		return strings.Compare(cellValue, condVal) < 0, nil
	case "gt":
		return strings.Compare(cellValue, condVal) > 0, nil
	case "le":
		return strings.Compare(cellValue, condVal) <= 0, nil
	case "ge":
		return strings.Compare(cellValue, condVal) >= 0, nil
	case "regexp":
		re, err := regexp.Compile(condVal)
		if err != nil {
			errorMsg := fmt.Sprintf("查询条件传入的正则表达式错误: %s", condVal)
			return false, errors.New(errorMsg)
		}
		return re.MatchString(cellValue), nil
	case "contain":
		return strings.Contains(cellValue, condVal), nil
	case "suffix":
		// cellValue是否以condVal结尾
		return strings.HasSuffix(cellValue, condVal), nil
	case "prefix":
		return strings.HasPrefix(cellValue, condVal), nil
	default:
		errorMsg := fmt.Sprintf("查询条件传入的运算比较符 %s 暂不被string类型支持", compare)
		return false, errors.New(errorMsg)
	}

}

/**
 * 比较浮点数
 * @author: jjq
 * @param cellValue 单元格值
 * @param condVal 查询条件基准值
 * @param compare 比较符
 * @return bool 比较结果
 * @return error 错误信息
 */
func compareFloat(cellValue string, condVal string, compare string) (bool, error) {
	// 将cellValue转换为float64类型
	cellValueFloat, err := strconv.ParseFloat(cellValue, 64)
	if err != nil {
		errorMsg := fmt.Sprintf("数据 %v 无法转换为浮点数: %s，该列可能非数字列，建议使用string比较。请修改查询条件的type", cellValue, err)
		return false, errors.New(errorMsg)
	}
	condValFloat, err := strconv.ParseFloat(condVal, 64)
	if err != nil {
		errorMsg := fmt.Sprintf("查询条件中基准值 %v 无法转换为浮点数: %s，请修改查询条件的type", condVal, err)
		return false, errors.New(errorMsg)
	}
	// return compareNum[float64](cellValueFloat, condValFloat, compare)
	switch compare {
	case "gt":
		return cellValueFloat > condValFloat, nil
	case "lt":
		return cellValueFloat < condValFloat, nil
	case "ge":
		return cellValueFloat >= condValFloat, nil
	case "le":
		return cellValueFloat <= condValFloat, nil
	case "eq":
		return cellValueFloat == condValFloat, nil
	case "ne":
		return cellValueFloat != condValFloat, nil

	default:
		errorMsg := fmt.Sprintf("查询条件传入的运算比较符 %s 暂不被float类型支持", compare)
		return false, errors.New(errorMsg)
	}
}

/**
 * 比较整数
 * @author: jjq
 * @param cellValue 单元格值
 * @param condVal 查询条件基准值
 * @param compare 比较符
 * @return bool 比较结果
 * @return error 错误信息
 */
func compareInt(cellValue string, condVal string, compare string) (bool, error) {
	// 将cellValue转换为float64类型
	cellValueInt, err := strconv.ParseInt(cellValue, 10, 64)
	if err != nil {
		errorMsg := fmt.Sprintf("数据 %v 无法转换为整数: %s，该列可能非数字列，建议使用string比较。请修改查询条件的type", cellValue, err)
		return false, errors.New(errorMsg)
	}
	condValInt, err := strconv.ParseInt(condVal, 10, 64)
	if err != nil {
		errorMsg := fmt.Sprintf("查询条件中基准值 %v 无法转换为整数: %s 。请修改查询条件的type", condVal, err)
		return false, errors.New(errorMsg)
	}
	// return compareNum[int64](cellValueInt, condValInt, compare)
	switch compare {
	case "gt":
		return cellValueInt > condValInt, nil
	case "lt":
		return cellValueInt < condValInt, nil
	case "ge":
		return cellValueInt >= condValInt, nil
	case "le":
		return cellValueInt <= condValInt, nil
	case "eq":
		return cellValueInt == condValInt, nil
	case "ne":
		return cellValueInt != condValInt, nil

	default:
		errorMsg := fmt.Sprintf("查询条件传入的运算比较符 %s 暂不被int类型支持", compare)
		return false, errors.New(errorMsg)
	}
}

/**
 * 【废弃】：因为泛型编译连接时间过久，因此本函数作废，不引入constraints.Ordered泛型包
 * 比较数字：compareInt和compareFloat将数据转换为int64和float64，再传入本函数进行比较。泛型进行类型匹配。
 * @param cellValue 单元格值
 * @param condVal 查询条件基准值
 * @param compare 比较符
 * @return bool 比较结果
 * @return error 错误信息

func compareNum[T constraints.Ordered](cellValue T, condVal T, compare string) (bool, error) {

	switch compare {
	case "gt":
		return cellValue > condVal, nil
	case "lt":
		return cellValue < condVal, nil
	case "ge":
		return cellValue >= condVal, nil
	case "le":
		return cellValue <= condVal, nil
	case "eq":
		return cellValue == condVal, nil
	case "ne":
		return cellValue != condVal, nil

	default:
		errorMsg := fmt.Sprintf("查询条件传入的运算比较符 %s 暂不被int/float类型支持", compare)
		return false, errors.New(errorMsg)
	}
}

*/

/**
 * AggregateSliceINDataSet 聚合同一个数据集多个分片（多个分片文件）的数据，返回表头和表数据字符串
 * @param filePosesSingleDataSet 文件位置（CID）数组，该数组内的文件位置为同一个数据集的多个分片（格式应该一样）
 * @param filePosAndDataMap 文件位置和文件内容的map（原始文件内容字符串）
 * @return map[string]int 表头map
 * @return string 表数据字符串
 * @return error 错误信息
 */
func AggregateSliceINDataSet(filePosesSingleDataSet []string, filePosAndDataMap map[string]string) (map[string]int, string, error) {
	tableHeaderMap := make(map[string]int)
	tableStr := ""
	for idx, filePos := range filePosesSingleDataSet {
		if idx == 0 {
			tableHeaderMap = parseTableHeader(filePos, filePosAndDataMap[filePos])
			tableStr += filePosAndDataMap[filePos]
			tableStr += "\n"
		} else {
			// 先比较是否表头一致
			if !compareMaps(tableHeaderMap, parseTableHeader(filePosesSingleDataSet[0], filePosAndDataMap[filePos])) {
				return nil, "", errors.New(filePos + "与" + filePosesSingleDataSet[0] + "表头不一致")
			}
			// 删掉表头
			tableStr += filePosAndDataMap[filePos][strings.Index(filePosAndDataMap[filePos], "\n")+1:]
			// 如果不是最后一个数据集，那么加上换行符
			if idx != len(filePosesSingleDataSet)-1 {
				tableStr += "\n"
			}
		}
	}
	return tableHeaderMap, tableStr, nil
}

// =====================联表操作=====================

/**
 * checkRowPairJoinConditionSatisfied 检查两行的联表条件是否满足，返回bool值
 * @param jointCondition 联表条件
 * @param line1Arr 表1的行数据数组
 * @param line2Arr 表2的行数据数组
 * @param field1Index 表1的字段索引
 * @param field2Index 表2的字段索引
 * @return bool 比较结果
 * @return error 错误信息
 */
func checkRowPairJoinConditionSatisfied(jointCondition JointCondition, line1Arr []string, line2Arr []string, field1Index int, field2Index int) (bool, error) {
	// 比较两个值是否相等
	switch jointCondition.Type {
	case "string":
		return compareInt(line1Arr[field1Index], line2Arr[field2Index], jointCondition.Compare)
	case "int":
		return compareInt(line1Arr[field1Index], line2Arr[field2Index], jointCondition.Compare)
	case "float":
		return compareFloat(line1Arr[field1Index], line2Arr[field2Index], jointCondition.Compare)
	default:
		return false, errors.New("联表条件类型错误，不支持类型：" + jointCondition.Type)
	}
}

/**
 * JointTwoTableInner 联表操作，返回表头和表数据字符串（INNER连接）
 * @param jointCondition 联表条件
 * @param tableHeaderMap1 表1表头map
 * @param tableHeaderMap2 表2表头map
 * @param tableStr1 表1数据字符串
 * @param tableStr2 表2数据字符串
 * @return string 表数据字符串
 * @return map[string]int 表头map
 * @return error 错误信息
 */
func JointTwoTableInner(jointCondition JointCondition, tableHeaderMap1 map[string]int, tableHeaderMap2 map[string]int, tableStr1 string, tableStr2 string) (string, map[string]int, error) {
	tableStrReturn := ""
	tableHeaderMapReturn := make(map[string]int)
	// 解析表1和表2
	field1Index, ok1 := tableHeaderMap1[jointCondition.Pos1+"_"+jointCondition.Field1]
	if !ok1 {
		return "", nil, errors.New("表" + jointCondition.Pos1 + "中不存在联表字段" + jointCondition.Field1)
	}
	field2Index, ok2 := tableHeaderMap2[jointCondition.Pos2+"_"+jointCondition.Field2]
	if !ok2 {
		return "", nil, errors.New("表" + jointCondition.Pos2 + "中不存在联表字段" + jointCondition.Field2)
	}

	// 逐个比较
	lines1 := strings.Split(tableStr1, "\n") // 表1的行
	lines2 := strings.Split(tableStr2, "\n") // 表2的行
	// 先把表头line合并
	tableStrReturn = lines1[0] + " " + lines2[0] + "\n"
	// 每行按空格拆分
	for idx1, line1 := range lines1 {
		if idx1 == 0 || line1 == "" {
			continue
		}
		line1Arr := strings.Split(line1, " ")
		for idx2, line2 := range lines2 {
			if idx2 == 0 || line2 == "" {
				continue
			}
			line2Arr := strings.Split(line2, " ")
			// 比较两个值是否相等
			flag, err := checkRowPairJoinConditionSatisfied(jointCondition, line1Arr, line2Arr, field1Index, field2Index)
			if err != nil {
				return "联表条件比较错误：" + err.Error(), nil, errors.New("联表条件比较错误：" + err.Error())
			} else if flag {
				// 如果满足联表条件，那么加入到tableStrReturn中
				tableStrReturnNewLine := line1 + " " + line2
				// 删除每行首位空格
				tableStrReturnNewLine = strings.Trim(tableStrReturnNewLine, " ")
				tableStrReturn += tableStrReturnNewLine + "\n"
			}

		}
	}

	// 删除最后一个\n
	if tableStrReturn[len(tableStrReturn)-1] == '\n' {
		tableStrReturn = tableStrReturn[:len(tableStrReturn)-1]
	}

	// 整合表头
	// 对于每一个tableHeaderMap2的value，加上len(tableHeaderMap1)，然后加入到tableHeaderMapReturn中.作为新的表头Header
	for key, value := range tableHeaderMap2 {
		tableHeaderMapReturn[key] = value + len(tableHeaderMap1)
	}
	for key, value := range tableHeaderMap1 {
		tableHeaderMapReturn[key] = value
	}

	return tableStrReturn, tableHeaderMapReturn, nil

}

/**
 * JointTwoTable 两表联表操作，返回表头和表数据字符串
 * @param jointCondition 联表条件
 * @param tableHeaderMap1 表1表头map
 * @param tableHeaderMap2 表2表头map
 * @param tableStr1 表1数据字符串
 * @param tableStr2 表2数据字符串
 * @return string 表数据字符串
 * @return map[string]int 表头map
 * @return error 错误信息
 */
func JointTwoTable(jointCondition JointCondition, tableHeaderMap1 map[string]int, tableHeaderMap2 map[string]int, tableStr1 string, tableStr2 string) (string, map[string]int, error) {

	switch jointCondition.JointType {
	case "INNER":
		// 内连接
		return JointTwoTableInner(jointCondition, tableHeaderMap1, tableHeaderMap2, tableStr1, tableStr2)
	default:
		return "不支持的联表类型", nil, errors.New("不支持的联表类型:" + jointCondition.JointType)
	}
}

func JointTables(jointConditions []JointCondition, tableHeaderMap_map map[string]map[string]int, tableStrMap map[string]string) (string, map[string]int, error) {
	tableHeaderMap := make(map[string]int) // 联表后的表头
	tableStr := ""                         // 最终的表头和表数据字符串
	posHasJoint := make([]string, 0)       // 已联表的数据集

	for idx, jointCondition := range jointConditions {
		if idx == 0 {
			// 对第一个进行联表
			newTableStr, newTableHeader, err := JointTwoTable(jointCondition, tableHeaderMap_map[jointCondition.Pos1], tableHeaderMap_map[jointCondition.Pos2], tableStrMap[jointCondition.Pos1], tableStrMap[jointCondition.Pos2]) // 联表
			if err != nil {
				return "", nil, err
			}
			posHasJoint = append(posHasJoint, jointCondition.Pos1, jointCondition.Pos2)
			tableStr = newTableStr
			tableHeaderMap = newTableHeader
		} else {
			// 对后续的进行联表
			// 判断是否已经联表
			if strIsInSlice(posHasJoint, jointCondition.Pos1) && strIsInSlice(posHasJoint, jointCondition.Pos2) {
				continue
			} else if strIsInSlice(posHasJoint, jointCondition.Pos1) && !strIsInSlice(posHasJoint, jointCondition.Pos2) {
				newTableStr, newTableHeader, err := JointTwoTable(jointCondition, tableHeaderMap, tableHeaderMap_map[jointCondition.Pos2], tableStr, tableStrMap[jointCondition.Pos2]) // 联表
				if err != nil {
					return "", nil, err
				}
				posHasJoint = append(posHasJoint, jointCondition.Pos2)
				tableStr = newTableStr
				tableHeaderMap = newTableHeader
			} else if !strIsInSlice(posHasJoint, jointCondition.Pos1) && strIsInSlice(posHasJoint, jointCondition.Pos2) {
				newTableStr, newTableHeader, err := JointTwoTable(jointCondition, tableHeaderMap_map[jointCondition.Pos1], tableHeaderMap, tableStrMap[jointCondition.Pos1], tableStr) // 联表
				if err != nil {
					return "", nil, err
				}
				posHasJoint = append(posHasJoint, jointCondition.Pos1)
				tableStr = newTableStr
				tableHeaderMap = newTableHeader
			} else {
				return "", nil, errors.New("联表条件中存在未联表的数据集")
			}

		}
	}
	if len(posHasJoint) != len(jointConditions)+1 {
		return "", nil, errors.New("联表条件中存在未联表的数据集")
	}
	return tableStr, tableHeaderMap, nil
}

/**
 * TopologicalSortOfEdges 对联表条件进行拓扑排序，返回排序后的联表条件数组（注：拓扑排序后，位于数组后面的联表条件中要么pos1已经被连接，要么pos2已经被连接）
 * @param edges 联表条件数组
 * @return []JointCondition 排序后的联表条件数组
 * @return error 错误信息
 */
func TopologicalSortOfEdges(edges []JointCondition) ([]JointCondition, error) {

	adjList := make(map[string][]string)
	inDegree := make(map[string]int)
	allPos := make(map[string]bool)

	// Initialize graph
	for _, edge := range edges {
		for _, pos := range []string{edge.Pos1, edge.Pos2} {
			if _, exists := allPos[pos]; !exists {
				adjList[pos] = []string{}
				allPos[pos] = true
				inDegree[pos] = 0
			}
		}
	}

	// Build graph
	for _, edge := range edges {
		adjList[edge.Pos1] = append(adjList[edge.Pos1], edge.Pos2)
		adjList[edge.Pos2] = append(adjList[edge.Pos2], edge.Pos1)
		inDegree[edge.Pos1]++
		inDegree[edge.Pos2]++
	}

	// Initialize queue with nodes having 0 in-degree
	var queue []string
	for pos, degree := range inDegree {
		if degree == 1 { // Adjusted for undirected graph
			queue = append(queue, pos)
		}
	}

	// Process nodes in the queue
	var sortedEdges []JointCondition
	visited := make(map[string]bool)
	for len(queue) > 0 {
		pos := queue[0]
		queue = queue[1:]

		// Find the edge that should be added to the sorted list
		for _, edge := range edges {
			if (edge.Pos1 == pos || edge.Pos2 == pos) && !visited[edge.Pos1+edge.Pos2] && !visited[edge.Pos2+edge.Pos1] {
				sortedEdges = append(sortedEdges, edge)
				visited[edge.Pos1+edge.Pos2] = true
				visited[edge.Pos2+edge.Pos1] = true
				break
			}
		}

		// Decrease in-degree of adjacent nodes
		for _, adj := range adjList[pos] {
			inDegree[adj]--
			if inDegree[adj] == 1 { // Adjusted for undirected graph
				queue = append(queue, adj)
			}
		}
	}

	// Check for cycles in the graph
	if len(sortedEdges) != len(edges) {

		return nil, fmt.Errorf("联表条件中存在环，无法进行拓扑排序")
	}

	// Reverse the sortedEdges to get the final order
	for i, j := 0, len(sortedEdges)-1; i < j; i, j = i+1, j-1 {
		sortedEdges[i], sortedEdges[j] = sortedEdges[j], sortedEdges[i]
	}

	return sortedEdges, nil
}

// =====================辅助函数=====================
/**
 * utils：strIsInSlice 检查字符串str是否在字符串slice数组中
 * @author: jjq
 * @param slice 字符串slice数组
 * @param str 字符串
 * @return bool 是否在数组中
 */
func strIsInSlice(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

/**
 * QueryModule 查询模块，返回查询结果（JSON字符串）和查询结果数量
 * @param queryConditions 查询条件数组
 * @param returnField 返回的列名
 * @param tableStr 表数据字符串
 * @param tableHeaderMap 表头map
 * @param isMulti 是否联表查询
 * @return string 查询结果（JSON字符串）
 * @return int 查询结果数量
 * @return error 错误信息
 */
func QueryModule(queryConditions [][]QueryCondition, returnField []string, tableStr string, tableHeaderMap map[string]int, isMulti bool) (string, int, error) {
	queryResultData := QueryResult{} // 初始化返回结果

	needReturnAllSlices := make([]string, 0) // 需要全部返回的分片数据(含*)
	for _, returnFieldSingle := range returnField {
		if strings.Join(strings.Split(returnFieldSingle, "_")[1:], "_") == "*" {
			// 如果有*，那么返回该数据集所有字段，则将此pos加入到needReturnAllSlices中
			needReturnAllSlices = append(needReturnAllSlices, strings.Split(returnFieldSingle, "_")[0]) // 返回所有字段的数据集
		}
	}

	// 逐行判断是否满足查询条件
	lines := strings.Split(tableStr, "\n")
	countLinesSatisfy := 0 // countLinesSatisfy: 满足条件的行数，计数

	queryResultData.Data = nil
	for _, line := range lines[1:] {
		yesForConditionFlag := false // 本行是否满足查询条件

		// 逐个conditionGroup判断
		// 只要满足其中一个conditionGroup条件，就认为本行满足查询条件【数组之间为OR关系】
		// 如果没有查询条件，那么默认为全查询
		if queryConditions == nil || len(queryConditions) == 0 {
			yesForConditionFlag = true
		} else {
			for _, conditionGroup := range queryConditions {
				condFlag, err := matchesConditions(line, tableHeaderMap, conditionGroup)
				if err != nil {
					// fmt.Println("matchesConditions error: ", err)
					return errorQueryResult(err.Error()), -1, err
				}
				// 满足条件即可返回，不同数组的条件为OR关系，同组条件为AND关系
				if condFlag {
					yesForConditionFlag = true
					break
				}
			}
		}

		if yesForConditionFlag == true && len(line) > 0 {
			// 将一行的数据变为JSON格式，key为tableHeaderMap的key，value为行数据
			rowData := make(map[string]string)

			for key, index := range tableHeaderMap {

				// 遍历每个列，判断只返回returnField中的字段。如果为*（在needReturnAllSlices数组中），那么该数据集的所有列恒为真，也需要返回该字段
				if strIsInSlice(returnField, key) || (len(returnField) > 0 && strIsInSlice(needReturnAllSlices, strings.Split(key, "_")[0])) {
					KeyWithoutPos := strings.Join(strings.Split(key, "_")[1:], "_")
					// line只需要读取前len(line)-1个字符，因为最后一个是换行。读取时去掉前缀
					// 去掉了反而会吞掉最后一个字符，不知道为什么。所以不去掉。
					if isMulti {
						// 联表查询需要考虑到前缀问题
						// // * 同一个列名需要代表相同的值！无论是否联表。否则会造成覆盖

						// rowData[KeyWithoutPos] = strings.Split(line[:len(line)], " ")[index]

						// 有前缀
						rowData[key] = strings.Split(line[:len(line)], " ")[index]
					} else {
						rowData[KeyWithoutPos] = strings.Split(line[:len(line)], " ")[index]
					}
				}
			}

			queryResultData.Data = append(queryResultData.Data, rowData) // 添加到数组中
			countLinesSatisfy++

		}

	}
	queryResultData.Counts = countLinesSatisfy
	if isMulti {
		queryResultData.Message = "联表查询成功"
	} else {
		queryResultData.Message = "查询成功"
	}
	returnData, _ := json.Marshal(queryResultData)
	return string(returnData), queryResultData.Counts, nil
}

/**
 * 返回查询失败的结果
 * @author: jjq
 * @param errMsg 错误信息
 * @return string 查询结果（携带错误信息）
 */
func errorQueryResult(errMsg string) string {
	queryResultData := QueryResult{}
	queryResultData.Data = nil
	queryResultData.Message = errMsg
	queryResultData.Counts = -1
	returnData, _ := json.Marshal(queryResultData)
	return string(returnData)
}

/**
 * 进行查询流程、返回查询结果
 * @author: jjq
 * @param queryItem 查询条件
 * @return string 查询结果, int 查询结果数量, error 错误信息
 * @Description: 当遇到查询错误时（包括：IPFS 文件不存在、IPFS 文件解析错误、列名不存在、类型不匹配（比如 string 类型列标明为 int）、正则表达式不正确、运算符不正确等等情况），直接返回errorQueryResult(errMsg)函数，不继续执行
 */
func ReturnQuerySingle(queryItemData QueryItem, filePosAndDataMap map[string]string) (string, int, error) {
	// queryResultData := QueryResult{} // 初始化返回结果

	filePos2DimArr := queryItemData.FilePos // 文件位置（CID）:此为二维数组，每个元素为同一个数据集
	returnField := queryItemData.ReturnField
	queryConditions := queryItemData.QueryConditions

	// 将FilePos转为一维数组，元素为每个子元素的第一个元素
	filePosArr := make([]string, 0)
	for _, filePos := range filePos2DimArr {
		filePosArr = append(filePosArr, filePos[0])
	}
	// 把filePosAndDataMap中的数据读取出来，拼接在一起。

	tableHeaderMap, tableStr, err := AggregateSliceINDataSet(filePos2DimArr[0], filePosAndDataMap) // 聚合同一个数据集多个分片（多个分片文件）的数据，返回表头和表数据字符串
	if err != nil {
		return errorQueryResult(err.Error()), -1, err
	}

	// fmt.Println("表头索引", tableHeaderMap) // 结果——map[id:0 name:1 score:2]

	// -----------------查询部分-----------------
	return QueryModule(queryConditions, returnField, tableStr, tableHeaderMap, false)
}

// =====================联表查询=====================

func ReturnQueryMulti(queryItemData QueryItem, filePosAndDataMap map[string]string) (string, int, error) {

	filePos2DimArr := queryItemData.FilePos // 文件位置（CID）:此为二维数组，每个元素为同一个数据集
	returnField := queryItemData.ReturnField
	queryConditions := queryItemData.QueryConditions
	jointConditions := queryItemData.JointConditions // 联表条件

	if len(filePos2DimArr) != len(jointConditions)+1 {
		return "联表条件数量与数据集数量不匹配（联表条件=数据集数量-1）", -1, errors.New("联表条件数量与数据集数量不匹配")
	}
	// 拓扑排序联表条件
	jointConditions, err := TopologicalSortOfEdges(jointConditions)
	if err != nil {
		return err.Error(), -1, err
	}

	tableStrMap := make(map[string]string)
	tableHeaderMap_map := make(map[string]map[string]int)
	filePosArr := make([]string, 0)

	// 合并数据集，并返回合并后的数据集
	for _, filePosesSingleDataSet := range filePos2DimArr {
		// 逐个解析每个数据集，filePosesSingleDataSet代表一个数据集（多个分片）的数组
		tableStr := ""
		// 解析表头字符串
		tableHeaderMap, tableStr, err := AggregateSliceINDataSet(filePosesSingleDataSet, filePosAndDataMap) // 聚合同一个数据集多个分片（多个分片文件）的数据，返回表头和表数据字符串
		if err != nil {
			return errorQueryResult(err.Error()), -1, err
		}

		// 赋值，以第0个数据集的主索引
		tableStrMap[filePosesSingleDataSet[0]] = tableStr
		tableHeaderMap_map[filePosesSingleDataSet[0]] = tableHeaderMap
		filePosArr = append(filePosArr, filePosesSingleDataSet[0])

	}
	tableStr, tableHeaderMap, err := JointTables(jointConditions, tableHeaderMap_map, tableStrMap)
	if err != nil {
		return errorQueryResult(err.Error()), -1, err
	}

	// fmt.Println("tableStr:", tableStr)
	// fmt.Println("tableHeaderMap:", tableHeaderMap)

	// -----------------查询部分-----------------
	return QueryModule(queryConditions, returnField, tableStr, tableHeaderMap, true)

}

func GetQueryResult(queryItem string, filePosAndDataMap map[string]string) (string, int) {

	// 解析查询条件，将queryItem转为JSON
	var queryItemData QueryItem
	err := json.Unmarshal([]byte(queryItem), &queryItemData)
	if err != nil {
		// fmt.Println("解析查询条件失败:", err)
		return errorQueryResult("解析查询条件失败:" + err.Error()), -1
	}

	if queryItemData.QueryConcatType == "single" {
		returnStr, resultCounts, err := ReturnQuerySingle(queryItemData, filePosAndDataMap)
		if err != nil {
			return "查询失败: " + err.Error(), -1
		} else {
			return returnStr, resultCounts
		}
	} else if queryItemData.QueryConcatType == "multi" {
		returnStr, resultCounts, err := ReturnQueryMulti(queryItemData, filePosAndDataMap)
		// return "联表查询暂未开放", -1

		if err != nil {
			return "查询失败: " + err.Error(), -1
		} else {
			return returnStr, resultCounts
		}
	} else {
		return "查询条件中的QueryConcatType字段错误", -1
	}

}
