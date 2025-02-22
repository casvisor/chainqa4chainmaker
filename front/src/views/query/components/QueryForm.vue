<template>
  <t-form
    :rules="rules"
    :data="form"
    ref="formRef"
    label-width="120px"
    @submit="handleSubmit"
    colon
    labelAlign="right"
    resetType="initial"
  >
    <t-form-item
      label="联表类型"
      name="queryConcatType"
      class="query-item"
      :help="
        form.queryConcatType == 'single'
          ? '单表查询也可以输入多个IPFS地址，但需要保证数据集格式一致'
          : '联表查询需要输入联表条件'
      "
    >
      <t-radio-group
        v-model="form.queryConcatType"
        @change="onQueryConcatTypeChange"
        :readonly="readOnly"
      >
        <t-radio-button value="single">单表查询</t-radio-button>
        <t-radio-button value="multi">联表查询</t-radio-button>
      </t-radio-group>
    </t-form-item>

    <t-form-item label="分片与展示" class="query-item" requiredMark>
      <t-tabs
        v-model="queryIPFSaddrAndShowColTab"
        :addable="form.queryConcatType === 'multi' && readOnly == false"
        @add="onAddQueryIPFSAddressAndShowColumn"
        @remove="onRemoveQueryIPFSAddressAndShowColumn"
        style="margin-left: 0px; border: 1px solid var(--td-component-stroke)"
      >
        <t-tab-panel
          v-for="(queryCol, index) in form.queryIPFSAdressAndShowColumn"
          :key="queryCol.id"
          :value="queryCol.id"
          :label="`数据集${form.queryConcatType == 'multi' ? index + 1 : ''}`"
          :destroy-on-hide="false"
          :removable="
            form.queryConcatType === 'multi' &&
            form.queryIPFSAdressAndShowColumn.length > 1 &&
            readOnly == false
          "
        >
          <div style="margin: 20px; min-width: 670px">
            <t-form-item
              requiredMark
              class="query-item"
              label="IPFS分片地址"
              :name="`queryIPFSAdressAndShowColumn[${index}].ipfsAddress`"
              help="按回车确认每个分片。为保证查询效率，最多输入3个IPFS地址，同一数据集的各分片格式应该一致，整个分片合起来为一个数据集。"
              :rules="[
                {
                  validator: (value) => {
                    if (value.length > 3) {
                      return {
                        result: false,
                        message: '最多输入3个IPFS地址',
                        type: 'error',
                      };
                    } else if (new Set(value).size !== value.length) {
                      return {
                        result: false,
                        message: 'IPFS地址重复',
                        type: 'error',
                      };
                    } else if (value.length === 0) {
                      return {
                        result: false,
                        message: 'IPFS地址为必填项',
                        type: 'error',
                      };
                    } else {
                      return {
                        result: true,
                      };
                    }
                  },
                },
              ]"
            >
              <t-tag-input
                v-model="queryCol.ipfsAddress"
                :placeholder="`请输入数据集${
                  form.queryConcatType == 'multi' ? index + 1 : ''
                }的IPFS地址，按回车确认每个分片`"
                :max="3"
                :readonly="readOnly"
              >
                <template #tag="{ value }">
                  <span>
                    <t-avatar
                      size="18px"
                      shape="round"
                      style="background-color: #0052d9; color: white"
                    >
                      <!-- 找到是第几个 IPFS地址 -->
                      {{ queryCol.ipfsAddress.indexOf(value) + 1 }}
                    </t-avatar>
                    <span>&nbsp;{{ value }}</span>
                  </span>
                </template>
              </t-tag-input>
            </t-form-item>
            <t-form-item
              :label="`数据集${
                form.queryConcatType === 'single' ? '' : index + 1
              }展示范围`"
              :name="`queryIPFSAdressAndShowColumn[${index}].queryShowColumnRange`"
              :label-width="130"
              help="仅标识查询结果所展示的范围，不影响查询结果"
              style="margin-bottom: 20px"
            >
              <t-radio-group
                v-model="
                  form.queryIPFSAdressAndShowColumn[index].queryShowColumnRange
                "
                :readonly="readOnly"
              >
                <t-radio value="*">展示本数据集全部列</t-radio>
                <t-radio value="select">选择部分列展示</t-radio>
                <t-radio value="non"> 仅统计条数，不展示结果 </t-radio>
              </t-radio-group>
            </t-form-item>
            <t-form-item
              v-if="
                form.queryIPFSAdressAndShowColumn[index].queryShowColumnRange ==
                'select'
              "
              :label="`数据集${
                form.queryConcatType === 'single' ? '' : index + 1
              }展示列`"
              :name="`queryIPFSAdressAndShowColumn[${index}].queryShowColumn`"
              help="请输入本数据集中返回的列，可输入多列，回车确认。若输入不存在的列，则会忽略"
              :label-width="120"
            >
              <t-tag-input
                v-model="
                  form.queryIPFSAdressAndShowColumn[index].queryShowColumn
                "
                placeholder="请输入展示列"
                :clearable="!readOnly"
                :readonly="readOnly"
              />
            </t-form-item>
          </div>
        </t-tab-panel>
      </t-tabs>
    </t-form-item>
    <t-form-item
      label="查询条件组"
      class="query-item"
      help="不设置条件组为默认查询所有数据。条件组之间为OR关系，条件组内为AND关系。"
    >
      <t-alert
        v-if="form.queryIPFSAdressAndShowColumn[0]?.ipfsAddress?.length === 0"
        type="info"
        showIcon
        :message="`请先输入IPFS地址${
          form.queryConcatType == 'multi' ? '(至少设置数据集1的IPFS地址)' : ''
        }后设置，输入IPFS地址后请勿修改IPFS地址`"
        style="margin-bottom: 20px"
      ></t-alert>
      <t-button
        variant="outline"
        @click="onAddQueryConditionGroup"
        v-else-if="form.queryConditionGroups.length == 0 && readOnly == false"
      >
        新增条件组
      </t-button>
      <div
        v-else-if="form.queryConditionGroups.length == 0 && readOnly == true"
      >
        未设置条件
      </div>
      <t-tabs
        v-else
        v-model="queryConditionGroupTab"
        style="margin-left: 0px; border: 1px solid var(--td-component-stroke)"
        :addable="readOnly == false"
        @add="onAddQueryConditionGroup"
        @remove="onRemoveQueryConditionGroup"
      >
        <t-tab-panel
          v-for="(item, groupIndex) in form.queryConditionGroups"
          :key="item.tabIndex"
          :value="item.tabIndex"
          :label="`条件组${groupIndex + 1}`"
          :removable="readOnly == false"
        >
          <div style="margin: 20px; min-width: 670px">
            <t-table
              ref="tableRef"
              rowKey="id"
              :columns="queryConditionsColumns"
              :data="form.queryConditionGroups[groupIndex].queryConditions"
            >
              <template #queryColumn="{ row, rowIndex }">
                <t-form-item
                  labelWidth="0"
                  :name="`queryConditionGroups[${groupIndex}]`"
                  style="margin-bottom: 20px"
                  :rules="[
                    {
                      validator: (value) => {
                        if (value.queryConditions[rowIndex].queryColumn == '') {
                          return {
                            result: false,
                            message: '请输入查询列',
                            type: 'error',
                          };
                        } else {
                          return {
                            result: true,
                          };
                        }
                      },
                    },
                  ]"
                >
                  <t-input
                    v-model="row.queryColumn"
                    :placeholder="`请输入查询列`"
                    :clearable="!readOnly"
                    :style="{ width: '100%' }"
                    :readonly="readOnly"
                  />
                </t-form-item>
              </template>
              <!-- 联表查询会用到 -->
              <template #queryFile="{ row, rowIndex }">
                <t-form-item
                  v-if="form.queryConcatType === 'multi'"
                  labelWidth="0"
                  :name="`queryConditionGroups[${groupIndex}]`"
                  style="margin-bottom: 20px"
                  :rules="[
                    {
                      validator: (value) => {
                        if (value.queryConditions[rowIndex].queryFile == '') {
                          return {
                            result: false,
                            message: '请选择列所在数据集',
                            type: 'error',
                          };
                        } else {
                          return {
                            result: true,
                          };
                        }
                      },
                    },
                  ]"
                >
                  <t-select
                    v-model="row.queryFile"
                    :placeholder="`选择查询列所在数据集`"
                    :clearable="!readOnly"
                    :style="{ width: '100%' }"
                    :readonly="readOnly"
                  >
                    <t-option
                      v-for="(
                        item, index
                      ) in form.queryIPFSAdressAndShowColumn.filter(
                        (item) => item.ipfsAddress.length > 0
                      )"
                      :key="index"
                      :label="`数据集${index + 1}`"
                      :value="`${form.queryIPFSAdressAndShowColumn[index]?.ipfsAddress[0]}`"
                    />
                  </t-select>
                </t-form-item>
                <span v-else> 数据集 </span>
              </template>
              <template #queryOperator="{ row, rowIndex }">
                <t-form-item
                  labelWidth="0"
                  :name="`queryConditionGroups[${groupIndex}]`"
                  style="margin-bottom: 20px"
                  :rules="[
                    {
                      validator: (value) => {
                        if (
                          value.queryConditions[rowIndex].queryOperator == '' ||
                          value.queryConditions[rowIndex].queryOperator.length <
                            2
                        ) {
                          return {
                            result: false,
                            message: '请选择查询操作',
                            type: 'error',
                          };
                        } else {
                          return {
                            result: true,
                          };
                        }
                      },
                    },
                  ]"
                >
                  <t-cascader
                    v-model="
                      form.queryConditionGroups[groupIndex].queryConditions[
                        rowIndex
                      ].queryOperator
                    "
                    valueType="full"
                    :options="queryOperatorOptions"
                    :clearable="!readOnly"
                    :readonly="readOnly"
                  >
                    <template #valueDisplay="{ value, selectedOptions }">
                      <div v-if="selectedOptions?.length === 2">
                        <span>{{
                          selectedOptions[1]?.label?.split(" ")[0]
                        }}</span>
                        <span>({{ selectedOptions[0]?.value }})</span>
                      </div>
                    </template>
                  </t-cascader>
                </t-form-item>
              </template>
              <template #baseValue="{ row, rowIndex }">
                <t-form-item
                  labelWidth="0"
                  :name="`queryConditionGroups[${groupIndex}]`"
                  style="margin-bottom: 20px"
                  :rules="[
                    {
                      validator: (value) => {
                        if (value.queryConditions[rowIndex].baseValue == '') {
                          return {
                            result: false,
                            message: '请输入基准值(或正则/前后缀)',
                            type: 'error',
                          };
                        } else {
                          return {
                            result: true,
                          };
                        }
                      },
                    },
                  ]"
                >
                  <t-input
                    v-model="row.baseValue"
                    :placeholder="`请输入基准值(或正则/前后缀)`"
                    :clearable="!readOnly"
                    :style="{ width: '100%' }"
                  />
                </t-form-item>
              </template>

              <template #action="{ row, rowIndex }">
                <t-button
                  variant="text"
                  @click="onRemoveQueryCondition(groupIndex, rowIndex)"
                  theme="danger"
                  :disabled="
                    form.queryConditionGroups[groupIndex].queryConditions
                      .length <= 1 || readOnly
                  "
                >
                  删除
                </t-button>
              </template>
            </t-table>

            <t-button
              variant="outline"
              @click="onAddQueryCondition(groupIndex)"
              style="margin-top: 20px"
              v-show="!readOnly"
            >
              新增条件
            </t-button>
          </div>
        </t-tab-panel>
      </t-tabs>
    </t-form-item>
    <t-form-item
      label="数据集联表条件"
      v-if="form.queryConcatType === 'multi'"
      class="query-item"
      help="联表查询时，需要输入联表条件，暂不支持联合条件"
      requiredMark
    >
      <t-alert
        v-if="
          form.queryIPFSAdressAndShowColumn.length < 2 ||
          form.queryIPFSAdressAndShowColumn[0]?.ipfsAddress?.length == 0 ||
          form.queryIPFSAdressAndShowColumn[1]?.ipfsAddress?.length == 0
        "
        type="info"
        showIcon
        :message="`请先输入IPFS地址${
          form.queryConcatType == 'multi'
            ? '(至少设置数据集1、数据集2的IPFS地址)'
            : ''
        }后设置，输入IPFS地址后请勿修改IPFS地址`"
        style="margin-bottom: 20px"
      ></t-alert>
      <!-- 联表条件，table -->
      <div v-else style="width: 100%">
        <t-table
          rowKey="id"
          ref="jointTableRef"
          :columns="jointConditionsColumns"
          :data="form.jointConditions"
        >
          <template #pos1="{ row, rowIndex }">
            <t-form-item
              labelWidth="0"
              :name="`jointConditions[${rowIndex}].pos1`"
              style="margin-bottom: 20px"
              :rules="[
                {
                  validator: (value) => {
                    if (value == '' || value == null) {
                      return {
                        result: false,
                        message: '请选择列1数据集',
                        type: 'error',
                      };
                    } else if (value == row?.pos2) {
                      return {
                        result: false,
                        message: '列1和列2数据集不能相同',
                        type: 'error',
                      };
                    } else {
                      return {
                        result: true,
                      };
                    }
                  },
                },
              ]"
            >
              <t-select
                v-model="row.pos1"
                :placeholder="`选择列1数据集`"
                :clearable="!readOnly"
                :style="{ width: '100%' }"
                :readonly="readOnly"
              >
                <t-option
                  v-for="(
                    item, index
                  ) in form.queryIPFSAdressAndShowColumn.filter(
                    (item) => item.ipfsAddress.length > 0
                  )"
                  :key="index"
                  :label="`数据集${index + 1}`"
                  :value="`${form.queryIPFSAdressAndShowColumn[index]?.ipfsAddress[0]}`"
                  :disabled="
                    row?.pos2 ==
                    form.queryIPFSAdressAndShowColumn[index]?.ipfsAddress[0]
                  "
                />
              </t-select>
            </t-form-item>
          </template>
          <template #field1="{ row, rowIndex }">
            <t-form-item
              labelWidth="0"
              :name="`jointConditions[${rowIndex}].field1`"
              style="margin-bottom: 20px"
              :rules="[
                {
                  validator: (value) => {
                    if (value == '' || value == null) {
                      return {
                        result: false,
                        message: '请输入列1名',
                        type: 'error',
                      };
                    } else {
                      return {
                        result: true,
                      };
                    }
                  },
                },
              ]"
            >
              <t-input
                v-model="row.field1"
                :placeholder="`请输入列1名`"
                :clearable="!readOnly"
                :style="{ width: '100%' }"
                :readonly="readOnly"
              />
            </t-form-item>
          </template>
          <template #pos2="{ row, rowIndex }">
            <t-form-item
              labelWidth="0"
              :name="`jointConditions[${rowIndex}].pos2`"
              style="margin-bottom: 20px"
              :rules="[
                {
                  validator: (value) => {
                    if (value == '' || value == null) {
                      return {
                        result: false,
                        message: '请选择列2数据集',
                        type: 'error',
                      };
                    } else if (value == row?.pos1) {
                      return {
                        result: false,
                        message: '列1和列2数据集不能相同',
                        type: 'error',
                      };
                    } else {
                      return {
                        result: true,
                      };
                    }
                  },
                },
              ]"
            >
              <t-select
                v-model="row.pos2"
                :placeholder="`选择列2数据集`"
                :clearable="!readOnly"
                :style="{ width: '100%' }"
                :readonly="readOnly"
              >
                <t-option
                  v-for="(
                    item, index
                  ) in form.queryIPFSAdressAndShowColumn.filter(
                    (item) => item.ipfsAddress.length > 0
                  )"
                  :key="index"
                  :label="`数据集${index + 1}`"
                  :value="`${form.queryIPFSAdressAndShowColumn[index]?.ipfsAddress[0]}`"
                  :disabled="
                    row.pos1 ==
                    form.queryIPFSAdressAndShowColumn[index]?.ipfsAddress[0]
                  "
                />
              </t-select>
            </t-form-item>
          </template>
          <template #field2="{ row, rowIndex }">
            <t-form-item
              labelWidth="0"
              :name="`jointConditions[${rowIndex}].field2`"
              style="margin-bottom: 20px"
              :rules="[
                {
                  validator: (value) => {
                    if (value == '' || value == null) {
                      return {
                        result: false,
                        message: '请输入列2名',
                        type: 'error',
                      };
                    } else {
                      return {
                        result: true,
                      };
                    }
                  },
                },
              ]"
            >
              <t-input
                v-model="row.field2"
                :placeholder="`请输入列2名`"
                :clearable="!readOnly"
                :style="{ width: '100%' }"
                :readonly="readOnly"
              />
            </t-form-item>
          </template>

          <template #jointType="{ row, rowIndex }">
            <t-form-item
              labelWidth="0"
              :name="`jointConditions[${rowIndex}].jointType`"
              style="margin-bottom: 20px"
              :rules="[
                {
                  validator: (value) => {
                    if (value == '' || value == null) {
                      return {
                        result: false,
                        message: '请选择联表类型',
                        type: 'error',
                      };
                    } else {
                      return {
                        result: true,
                      };
                    }
                  },
                },
              ]"
            >
              <t-select
                v-model="row.jointType"
                :placeholder="`选择联表类型`"
                :clearable="!readOnly"
                :style="{ width: '100%' }"
                :readonly="readOnly"
              >
                <t-option
                  v-for="item in jointTypeOptions"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                />
              </t-select>
            </t-form-item>
          </template>
          <template #compareAndType="{ row, rowIndex }">
            <t-form-item
              labelWidth="0"
              :name="`jointConditions[${rowIndex}].compareAndType`"
              style="margin-bottom: 20px"
              :rules="[
                {
                  validator: (value) => {
                    if (value == '' || value == null) {
                      return {
                        result: false,
                        message: '请选择联表条件',
                        type: 'error',
                      };
                    } else {
                      return {
                        result: true,
                      };
                    }
                  },
                },
              ]"
            >
              <t-cascader
                v-model="row.compareAndType"
                :options="compareAndTypeOptions"
                :clearable="!readOnly"
                :style="{ width: '100%' }"
                valueType="full"
                :readonly="readOnly"
              >
                <template #valueDisplay="{ value, selectedOptions }">
                  <div v-if="selectedOptions?.length === 2">
                    <span>{{ selectedOptions[1]?.label?.split(" ")[0] }}</span>
                    <span>({{ selectedOptions[0]?.value }})</span>
                  </div>
                </template>
              </t-cascader>
            </t-form-item>
          </template>
          <template #action="{ row, rowIndex }">
            <t-button
              variant="text"
              @click="onRemoveJointCondition(rowIndex)"
              theme="danger"
              :disabled="readOnly"
            >
              删除
            </t-button>
          </template>
        </t-table>
        <div
          style="
            display: flex;
            justify-content: flex-start;
            align-items: center;
            margin-top: 20px;
          "
        >
          <div>
            <t-button
              variant="outline"
              @click="onAddJointCondition"
              style="margin-top: 0px"
              :disabled="
                form.queryIPFSAdressAndShowColumn.length - 1 ==
                  form.jointConditions.length || readOnly
              "
            >
              新增联表条件
            </t-button>
          </div>
          <div>
            <t-typography-text
              :style="{
                color:
                  form.queryIPFSAdressAndShowColumn.length - 1 ==
                  form.jointConditions.length
                    ? 'black'
                    : 'red',
              }"
            >
              （应有{{
                form.queryIPFSAdressAndShowColumn.length - 1
              }}个条件，现有{{ form.jointConditions.length }}个条件）
            </t-typography-text>
          </div>
        </div>
      </div>
    </t-form-item>
    <t-form-item class="query-btn" v-show="readOnly == false">
      <t-button theme="primary" type="submit">查询</t-button>
      <t-button
        theme="default"
        variant="base"
        @click="ShowImportDialog()"
        style="margin-left: 10px"
        >导入条件</t-button
      >
      <import-dialog ref="importDialogRef" @importDataFinish="importNow" />
      <t-button
        theme="default"
        variant="base"
        type="reset"
        style="margin-left: 10px"
        >重置</t-button
      >
    </t-form-item>
    <t-form-item class="query-btn" v-show="readOnly == false">
      <t-typography-text strong>
        查询条件JSON生成操作（可用于分享）：
      </t-typography-text>
      <t-button
        theme="primary"
        @click="onGenerateQueryItemJson"
        style="margin-left: 10px"
        variant="text"
        >生成JSON查询项</t-button
      >
      <!-- <t-button
        theme="primary"
        @click="onGenerateQueryItemStr"
        style="margin-left: 10px"
        variant="text"
        >生成字符串查询项</t-button
      > -->
    </t-form-item>
  </t-form>
</template>

<script setup>
import { reactive, ref, defineEmits, defineExpose } from "vue";
import { MessagePlugin } from "tdesign-vue-next";
import ImportDialog from "./ImportDialog.vue";
const formRef = ref(null);
const form = reactive({
  queryConcatType: "single", // 1: 单表查询 2: 多表查询
  //   ipfsAddress: [],
  queryIPFSAdressAndShowColumn: [
    {
      ipfsAddress: [],
      id: 0,
      queryShowColumn: [],
      queryShowColumnRange: "*", // 默认查询全部列
    },
  ],
  queryConditionGroups: [],
  jointConditions: [], // 联表条件
});

const queryIPFSaddrAndShowColTab = ref(0);
const queryConditionGroupTab = ref(1);

const rules = {
  queryConcatType: [
    { required: true, message: "联表类型为必填项", type: "error" },
  ],
  ipfsAddress: [{ required: true, message: "IPFS地址为必填项", type: "error" }],
  queryShowColumnRange: [
    { required: true, message: "展示范围为必填项", type: "error" },
  ],
  queryShowColumn: [
    { required: true, message: "展示列为必填项", type: "error" },
  ],
};

// ========== 表单提交 ============

const emit = defineEmits(["submitQuery"]);

import checkCanTopologicalSortOfEdges from "@/utils/checkCanTopologicalSortOfEdges";
import { read } from "xlsx";
const checkBeforeSubmit = () => {
  // 检查是否满足联表条件
  if (form.queryConcatType === "multi") {
    if (form.queryIPFSAdressAndShowColumn.length < 2) {
      MessagePlugin.warning("联表查询需至少有两个数据集，否则请使用单表查询");
      return false;
    }

    if (
      form.jointConditions.length !==
      form.queryIPFSAdressAndShowColumn.length - 1
    ) {
      MessagePlugin.warning("联表条件个数不足，请添加联表条件");
      return false;
    }

    if (!checkCanTopologicalSortOfEdges(form.jointConditions)) {
      MessagePlugin.warning("联表条件可能存在环，请检查");
      return false;
    }
    return true;
  } else {
    return true;
  }
};

const handleSubmit = ({ validateResult, firstError, e }) => {
  e.preventDefault(); // 阻止表单默认提交行为
  let flag = true;
  if (validateResult === true) {
    flag = checkBeforeSubmit();
  }
  if (validateResult === true && flag) {
    let queryItem = composeTheQueryItems();
    // 将filePos转为一维数组
    let filePoses = new Set();
    form.queryIPFSAdressAndShowColumn.forEach((item) => {
      item.ipfsAddress.forEach((ipfsAddr) => {
        filePoses.add(ipfsAddr);
      });
    });
    // 把set转数组
    filePoses = Array.from(filePoses);
    // console.log(filePoses);
    emit("submitQuery", queryItem, filePoses);
  } else if (flag == true) {
    MessagePlugin.error("请检查输入项");
  }
};

const onGenerateQueryItemJson = async () => {
  // 先触发验证
  const validateResult = await formRef.value.validate(); // 获取验证结果
  if (validateResult === true) {
    let queryItem = composeTheQueryItems();

    let clipboard = {
      writeText: (text) => {
        let copyInput = document.createElement("input");
        copyInput.value = text;
        document.body.appendChild(copyInput);
        copyInput.select();
        document.execCommand("copy");
        document.body.removeChild(copyInput);
      },
    };
    if (clipboard) {
      await clipboard.writeText(queryItem);
      MessagePlugin.success("查询条件（JSON）已复制到剪贴板");
    }
  } else {
    MessagePlugin.error("请检查输入项");
  }
};

const onGenerateQueryItemStr = async () => {
  // 先触发验证
  const validateResult = await formRef.value.validate(); // 获取验证结果
  if (validateResult === true) {
    let queryItem = composeTheQueryItems();
    const escapedString = queryItem.replace(/"/g, '\\"');

    let clipboard = {
      writeText: (text) => {
        let copyInput = document.createElement("input");
        copyInput.value = text;
        document.body.appendChild(copyInput);
        copyInput.select();
        document.execCommand("copy");
        document.body.removeChild(copyInput);
      },
    };
    if (clipboard) {
      await clipboard.writeText(escapedString);
      MessagePlugin.success(
        "查询条件（字符串）已复制到剪贴板，适合粘贴到postman等请求中"
      );
    }
    // navigator.clipboard.writeText(escapedString).then(
    //   function () {
    //     MessagePlugin.success(
    //       "查询条件（字符串）已复制到剪贴板，适合粘贴到postman等请求中"
    //     );
    //   },
    //   function () {
    //     MessagePlugin.error("复制失败，请手动复制");
    //   }
    // );
  } else {
    MessagePlugin.error("请检查输入项");
  }
};

// ============ 视觉组件 ============

const onQueryConcatTypeChange = (value) => {
  // 重置表单
  formRef.value.reset();
  form.queryConditionGroups = []; // 重置条件组
  form.queryIPFSAdressAndShowColumn = [];
  onAddQueryIPFSAddressAndShowColumn();
  form.queryConcatType = value;
  form.jointConditions = []; // 重置联表条件
};

var queryIPFSAddressAndShowColumnTabCounts = 1;
const onAddQueryIPFSAddressAndShowColumn = () => {
  queryIPFSAddressAndShowColumnTabCounts++;
  form.queryIPFSAdressAndShowColumn.push({
    ipfsAddress: [],
    id: queryIPFSAddressAndShowColumnTabCounts,
    queryShowColumn: [],
    queryShowColumnRange: "*",
  });
  if (form.queryIPFSAdressAndShowColumn.length == 1) {
    queryIPFSaddrAndShowColTab.value = queryIPFSAddressAndShowColumnTabCounts;
  }
};
const onRemoveQueryIPFSAddressAndShowColumn = (option) => {
  var needChange = false;
  if (
    queryIPFSaddrAndShowColTab.value ==
    form.queryIPFSAdressAndShowColumn[option.index].id
  ) {
    needChange = true;
  }
  form.queryIPFSAdressAndShowColumn = form.queryIPFSAdressAndShowColumn.filter(
    (item, index) => index !== option.index
  );
  if (needChange) {
    queryIPFSaddrAndShowColTab.value = form.queryIPFSAdressAndShowColumn[0].id;
  }
};

var queryConditionGroupsTabCounts = 0;
const onAddQueryConditionGroup = () => {
  queryConditionGroupsTabCounts++;
  form.queryConditionGroups.push({
    tabIndex: queryConditionGroupsTabCounts,
    queryConditions: [],
  });
  onAddQueryCondition(form.queryConditionGroups.length - 1);
  // 如果是第一个条件组，则默认选中
  if (form.queryConditionGroups.length == 1) {
    queryConditionGroupTab.value = queryConditionGroupsTabCounts;
  }
};

const onRemoveQueryConditionGroup = (option) => {
  var needChange = false;
  if (
    queryConditionGroupTab.value ==
    form.queryConditionGroups[option.index].tabIndex
  ) {
    needChange = true;
  }
  form.queryConditionGroups = form.queryConditionGroups.filter(
    (item, index) => index !== option.index
  );
  if (needChange) {
    queryConditionGroupTab.value = form.queryConditionGroups[0].tabIndex;
  }
};

// ============ 组装提交 ============
const composeTheQueryItems = () => {
  var queryItems = {};
  queryItems.queryConcatType = form.queryConcatType; // 查询连接类型
  queryItems.filePos = composeIPFSAddress(); // IPFS地址
  queryItems.returnField = composeReturnField(); // 返回展示的字段
  queryItems.queryConditions = composeQueryConditionGroups(); // 查询条件
  if (form.queryConcatType === "multi") {
    queryItems.jointConditions = composeJointConditions(); // 联表条件
  }
  let queryItemsJSON = JSON.stringify(queryItems);
  return queryItemsJSON;
};

const composeIPFSAddress = () => {
  var filePos = [];
  for (var i = 0; i < form.queryIPFSAdressAndShowColumn.length; i++) {
    filePos.push(form.queryIPFSAdressAndShowColumn[i].ipfsAddress);
  }
  return filePos;
};

const composeReturnField = () => {
  var returnField = [];
  for (var i = 0; i < form.queryIPFSAdressAndShowColumn.length; i++) {
    if (form.queryIPFSAdressAndShowColumn[i].queryShowColumnRange == "*") {
      returnField.push(
        form.queryIPFSAdressAndShowColumn[i].ipfsAddress[0] + "_*"
      );
    } else if (
      form.queryIPFSAdressAndShowColumn[i].queryShowColumnRange == "select"
    ) {
      for (
        var j = 0;
        j < form.queryIPFSAdressAndShowColumn[i].queryShowColumn.length;
        j++
      ) {
        returnField.push(
          form.queryIPFSAdressAndShowColumn[i].ipfsAddress[0] +
            "_" +
            form.queryIPFSAdressAndShowColumn[i].queryShowColumn[j]
        );
      }
    } else {
      // 仅统计条数 不操作
    }
  }
  return returnField;
};

const composeQueryConditionGroups = () => {
  var queryConditionGroups = [];
  for (var i = 0; i < form.queryConditionGroups.length; i++) {
    var queryConditionGroup = composeQueryConditionGroupTable(
      form.queryConditionGroups[i]
    ); // 某个条件组
    queryConditionGroups.push(queryConditionGroup); // 添加到条件组列表
  }
  return queryConditionGroups;
};

const composeQueryConditionGroupTable = (queryConditionGroup) => {
  var queryConditionGroupTable = [];
  for (var i = 0; i < queryConditionGroup.queryConditions.length; i++) {
    var queryConditionGroupRow = composeQueryConditionGroupRow(
      queryConditionGroup.queryConditions[i]
    );
    queryConditionGroupTable.push(queryConditionGroupRow);
  }
  return queryConditionGroupTable;
};

const composeQueryConditionGroupRow = (queryConditionGroupRowPASS) => {
  var queryConditionGroupRow = {};
  queryConditionGroupRow.field = queryConditionGroupRowPASS.queryColumn;
  queryConditionGroupRow.pos = queryConditionGroupRowPASS.queryFile;

  queryConditionGroupRow.compare = queryConditionGroupRowPASS.queryOperator[1];
  queryConditionGroupRow.val = queryConditionGroupRowPASS.baseValue;
  queryConditionGroupRow.type = queryConditionGroupRowPASS.queryOperator[0];
  return queryConditionGroupRow;
};

const composeJointConditions = () => {
  var jointConditions = [];
  for (var i = 0; i < form.jointConditions.length; i++) {
    var jointCondition = composeJointCondition(form.jointConditions[i]);
    jointConditions.push(jointCondition);
  }
  return jointConditions;
};

const composeJointCondition = (jointCondition) => {
  var jointConditionRow = {};
  jointConditionRow.pos1 = jointCondition.pos1;
  jointConditionRow.field1 = jointCondition.field1;
  jointConditionRow.pos2 = jointCondition.pos2;
  jointConditionRow.field2 = jointCondition.field2;
  jointConditionRow.compare = jointCondition.compareAndType[1];
  jointConditionRow.type = jointCondition.compareAndType[0];
  jointConditionRow.jointType = jointCondition.jointType;
  return jointConditionRow;
};

// ============ 表格 ============

const onAddQueryCondition = (groupIndex) => {
  //   console.log(form.queryConditionGroups[groupIndex].queryConditions);
  if (form.queryConcatType === "single") {
    form.queryConditionGroups[groupIndex].queryConditions.push({
      queryFile: form.queryIPFSAdressAndShowColumn[0].ipfsAddress[0],
      queryColumn: "",
      queryOperator: "",
      baseValue: "",
    });
  } else {
    form.queryConditionGroups[groupIndex].queryConditions.push({
      queryFile: "",
      queryColumn: "",
      queryOperator: "",
      baseValue: "",
    });
  }
};

const onRemoveQueryCondition = (groupIndex, index) => {
  form.queryConditionGroups[groupIndex].queryConditions.splice(index, 1);
};

const queryConditionsColumns = [
  {
    colKey: "serial-number",
    width: 50,
    title: "序号",
  },
  {
    colKey: "queryColumn",
    title: "列名",
    width: 200,
    cell: "queryColumn",
  },
  {
    colKey: "queryFile",
    title: "列所在数据集",
    width: 200,
    cell: "queryFile",
  },
  {
    colKey: "queryOperator",
    title: "操作符和类型",
    width: 200,
    cell: "queryOperator",
  },
  {
    colKey: "baseValue",
    title: "基准值",
    width: 200,
    ellipsis: true,
    cell: "baseValue",
  },
  {
    colKey: "action",
    title: "操作",
    fixed: "right",
    cell: "action",
  },
];

const queryOperatorOptions = [
  {
    label: "整数(int)",
    value: "int",
    children: [
      {
        label: "等于 =",
        value: "eq",
      },
      {
        label: "不等于 !=",
        value: "ne",
      },
      {
        label: "大于 >",
        value: "gt",
      },
      {
        label: "大于等于 >=",
        value: "ge",
      },
      {
        label: "小于 <",
        value: "lt",
      },
      {
        label: "小于等于 <=",
        value: "le",
      },
    ],
  },
  {
    label: "浮点数(float)",
    value: "float",
    children: [
      {
        label: "等于 =",
        value: "eq",
      },
      {
        label: "不等于 !=",
        value: "ne",
      },
      {
        label: "大于 >",
        value: "gt",
      },
      {
        label: "大于等于 >=",
        value: "ge",
      },
      {
        label: "小于 <",
        value: "lt",
      },
      {
        label: "小于等于 <=",
        value: "le",
      },
    ],
  },
  {
    label: "字符串(string)",
    value: "string",
    children: [
      {
        label: "等于 =",
        value: "eq",
      },
      {
        label: "不等于 !=",
        value: "ne",
      },
      {
        label: "大于 >",
        value: "gt",
      },
      {
        label: "大于等于 >=",
        value: "ge",
      },
      {
        label: "小于 <",
        value: "lt",
      },
      {
        label: "小于等于 <=",
        value: "le",
      },
      {
        label: "正则匹配",
        value: "regexp",
      },
      {
        label: "包含",
        value: "contain",
      },
      {
        label: "前缀",
        value: "prefix",
      },
      {
        label: "后缀",
        value: "suffix",
      },
    ],
  },
];

// ---------------- 联表条件 ----------------
const jointConditionsColumns = [
  {
    colKey: "serial-number",
    width: 50,
    title: "序号",
  },
  {
    colKey: "pos1",
    title: "列1数据集",
    width: 160,
    cell: "pos1",
  },
  {
    colKey: "field1",
    title: "列1名",
    width: 200,
    cell: "field1",
  },
  {
    colKey: "compareAndType",
    title: "操作符及类型",
    width: 200,
    cell: "compareAndType",
  },
  {
    colKey: "pos2",
    title: "列2数据集",
    width: 160,
    cell: "pos2",
  },
  {
    colKey: "field2",
    title: "列2名",
    width: 200,
    cell: "field2",
  },
  {
    colKey: "jointType",
    title: "联表类型",
    width: 160,
    cell: "jointType",
  },
  {
    colKey: "action",
    title: "操作",
    fixed: "right",
    cell: "action",
  },
];

const onAddJointCondition = () => {
  form.jointConditions.push({
    pos1: "",
    field1: "",
    compareAndType: "",
    pos2: "",
    field2: "",
    jointType: "INNER",
  });
};

const onRemoveJointCondition = (index) => {
  form.jointConditions.splice(index, 1);
};

const compareAndTypeOptions = [
  {
    label: "整数(int)",
    value: "int",
    children: [
      {
        label: "等于 =",
        value: "eq",
      },
      {
        label: "不等于 !=",
        value: "ne",
      },
      {
        label: "大于 >",
        value: "gt",
      },
      {
        label: "大于等于 >=",
        value: "ge",
      },
      {
        label: "小于 <",
        value: "lt",
      },
      {
        label: "小于等于 <=",
        value: "le",
      },
    ],
  },
  {
    label: "浮点数(float)",
    value: "float",
    children: [
      {
        label: "等于 =",
        value: "eq",
      },
      {
        label: "不等于 !=",
        value: "ne",
      },
      {
        label: "大于 >",
        value: "gt",
      },
      {
        label: "大于等于 >=",
        value: "ge",
      },
      {
        label: "小于 <",
        value: "lt",
      },
      {
        label: "小于等于 <=",
        value: "le",
      },
    ],
  },
  {
    label: "字符串(string)",
    value: "string",
    children: [
      {
        label: "等于 =",
        value: "eq",
      },
      {
        label: "不等于 !=",
        value: "ne",
      },
      {
        label: "大于 >",
        value: "gt",
      },
      {
        label: "大于等于 >=",
        value: "ge",
      },
      {
        label: "小于 <",
        value: "lt",
      },
      {
        label: "小于等于 <=",
        value: "le",
      },
      {
        label: "正则匹配",
        value: "regexp",
      },
      {
        label: "包含",
        value: "contain",
      },
      {
        label: "前缀",
        value: "prefix",
      },
      {
        label: "后缀",
        value: "suffix",
      },
    ],
  },
];

const jointTypeOptions = [
  {
    label: "内联（inner）",
    value: "INNER",
  },
  // {
  //   label: "左联",
  //   value: "left",
  // },
  // {
  //   label: "右联",
  //   value: "right",
  // },
];

// ========== 导入数据 ============
const importDialogRef = ref(null);
const ShowImportDialog = () => {
  importDialogRef.value.showImportDialog();
};

const readOnly = ref(false);
const importNow = (data_, readOnly_ = false) => {
  let data = JSON.parse(data_);
  form.queryConcatType = data.queryConcatType;
  form.queryIPFSAdressAndShowColumn = data.queryIPFSAdressAndShowColumn;
  form.queryConditionGroups = data.queryConditionGroups;
  form.jointConditions = data.jointConditions;
  queryIPFSaddrAndShowColTab.value = 0;
  queryConditionGroupsTabCounts = form.queryConditionGroups.length;
  if (form.queryConditionGroups.length > 0) {
    queryConditionGroupTab.value = 0;
  }
  readOnly.value = readOnly_;
};

defineExpose({
  importNow,
});
</script>

<style scoped>
.query-item {
  margin-bottom: 30px;
}

.query-btn {
  display: flex;
  justify-content: flex-end;
}
</style>
