/**
 * 根据queryItemJSON生成queryItem（适用于填表的）对象JSON
 * @param {*} queryItemJSON  
 */
const backComposeQueryItem = (queryItemJSON) => {
    try {


        let queryItemObj = JSON.parse(queryItemJSON);

        let form = {
            queryConcatType: "single", // 1: 单表查询 2: 多表查询
            //   ipfsAddress: [],
            queryIPFSAdressAndShowColumn: [

            ],
            queryConditionGroups: [],
            jointConditions: [], // 联表条件
        };

        // 1. 联表条件
        if (queryItemObj.queryConcatType == "single") {
            form.queryConcatType = "single";

        } else {
            form.queryConcatType = "multi";
        }

        // 2. 查询条件
        for (let i = 0; i < queryItemObj.queryConditions.length; i++) {
            let queryConditionGroup = queryItemObj.queryConditions[i];
            let group = {
                tabIndex: i,
                queryConditions: []
            }
            let groupArr = [];
            for (let j = 0; j < queryConditionGroup.length; j++) {
                let queryCondition = queryConditionGroup[j];
                let condition = {
                    queryFile: queryCondition.pos,
                    queryColumn: queryCondition.field,
                    baseValue: queryCondition.val,
                    // colType: queryCondition.type,//废弃
                    queryOperator: [queryCondition.type, queryCondition.compare]
                };
                groupArr.push(condition);
            }
            group.queryConditions = groupArr;
            form.queryConditionGroups.push(group);
        }

        // 3. 联表条件
        if (queryItemObj.queryConcatType === "multi") {
            for (let i = 0; i < queryItemObj.jointConditions.length; i++) {
                let jointCondition = queryItemObj.jointConditions[i];
                let condition = {
                    pos1: jointCondition.pos1,
                    pos2: jointCondition.pos2,
                    field1: jointCondition.field1,
                    field2: jointCondition.field2,
                    compareAndType: [jointCondition.type, jointCondition.compare],
                    jointType: jointCondition.jointType
                }
                form.jointConditions.push(condition);
            }
        } else {
            form.jointConditions = []; // 单表查询没有联表条件
        }

        // 4. 分片与查询列
        let showColCounts = 0;
        for (let i = 0; i < queryItemObj.filePos.length; i++) {
            let filePoses = queryItemObj.filePos[i];// 一个数据集的多个分片
            let queryIPFSAdressAndShowColumnSingle = {
                ipfsAddress: [],
                id: i,
                queryShowColumn: [],
                queryShowColumnRange: "non"
            };//一个数据集的ipfs地址和查询列
            for (let j = 0; j < filePoses.length; j++) {
                let filePos = filePoses[j];
                queryIPFSAdressAndShowColumnSingle.ipfsAddress.push(filePos);

                if (j == 0) {
                    // 只取第一个分片的查询列，此为代表列
                    // 在returnFields中搜索以filePos开头的列名，加入到queryShowColumn中

                    for (let k = 0; k < queryItemObj.returnField.length; k++) {
                        let returnField = queryItemObj.returnField[k];

                        if (returnField.indexOf(filePos) == 0) {
                            let colName = returnField.split("_").slice(1).join("_");

                            if (colName === "*") {
                                queryIPFSAdressAndShowColumnSingle.queryShowColumnRange = "*";
                                queryIPFSAdressAndShowColumnSingle.queryShowColumn = [];
                                showColCounts++;
                                break;
                            } else {
                                queryIPFSAdressAndShowColumnSingle.queryShowColumn.push(colName);
                                queryIPFSAdressAndShowColumnSingle.queryShowColumnRange = "select";
                                showColCounts++;
                            }

                        }
                    }
                    // 找完了都找不到，那么是non

                }
            }
            form.queryIPFSAdressAndShowColumn.push(queryIPFSAdressAndShowColumnSingle);

        }

        if (showColCounts != queryItemObj.returnField.length) {
            return "生成失败，请检查展示列是否正确。请注意：展示列的所属的分片地址必须是该数据集的第一个分片地址。";
        }



        return JSON.stringify(form);
    } catch (e) {
        return "生成失败，请检查字段是否正确：" + e;
    }


}

export default backComposeQueryItem