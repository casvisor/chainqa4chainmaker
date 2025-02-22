<template>
  <t-dialog
    v-model:visible="visible"
    header="导入查询条件"
    width="50%"
    :destroy-on-close="true"
    :footer="false"
  >
    <t-form
      :data="importData"
      :rules="importDataRules"
      ref="importDataRef"
      label-align="top"
      :colon="true"
    >
      <t-form-item label="导入查询条件JSON" name="importJson">
        <t-textarea
          v-model="importData.importJson"
          placeholder="请输入导入的查询条件"
          rows="10"
          :autosize="{ minRows: 8, maxRows: 15 }"
        ></t-textarea>
      </t-form-item>
    </t-form>
    <t-alert
      theme="info"
      message="导入后，原有的查询条件将会被覆盖"
      style="margin-top: 20px"
    ></t-alert>
    <div style="margin-top: 10px; text-align: right">
      <t-button theme="default" @click="visible = false">关闭</t-button>
      <t-button theme="primary" @click="importDataNow" style="margin-left: 20px"
        >导入</t-button
      >
    </div>
  </t-dialog>
</template>

<script setup>
import { reactive, defineExpose, ref, defineEmits } from "vue";
import { MessagePlugin } from "tdesign-vue-next";
import backComposeQueryItem from "@/utils/backComposeQueryItem";

const visible = ref(false);
const importData = reactive({
  importJson: "",
});
const importDataRef = ref(null);

const emit = defineEmits(["importDataFinish"]);

const importDataRules = {
  importJson: [
    { required: true, message: "请输入导入条件（JSON格式）", trigger: "blur" },
    {
      validator: (value) => {
        try {
          JSON.parse(value);
          // 检查是否含有必要字段
          const data = JSON.parse(value);
          if (
            !data.hasOwnProperty("queryConcatType") ||
            !data.hasOwnProperty("filePos") ||
            !data.hasOwnProperty("returnField") ||
            !data.hasOwnProperty("queryConditions")
          ) {
            return {
              result: false,
              message:
                "检查是否包含必要字段：queryConcatType, filePos, returnField, queryConditions",
              type: "error",
            };
          } else if (
            data.queryConcatType == "multi" &&
            !data.hasOwnProperty("jointConditions")
          ) {
            return {
              result: false,
              message: "当为联表查询时，还需包含必要字段：jointConditions",
              type: "error",
            };
          } else {
            return {
              result: true,
            };
          }
        } catch (error) {
          return {
            result: false,
            message: "请输入正确的JSON格式",
            type: "error",
          };
        }
      },
      trigger: "blur",
    },
  ],
};
const showImportDialog = () => {
  importData.importJson = "";
  visible.value = true;
};

const importDataNow = async () => {
  const validateResult = await importDataRef.value.validate();
  if (validateResult === true) {
    // 导入查询条件
    const backComposeQueryItemRes = backComposeQueryItem(importData.importJson);
    if (backComposeQueryItemRes.indexOf("生成失败") == 0) {
      MessagePlugin.warning(backComposeQueryItemRes);
      return;
    }
    emit("importDataFinish", backComposeQueryItemRes);
    MessagePlugin.success("导入查询条件成功");
    visible.value = false;
  } else {
    MessagePlugin.error("请输入正确的查询条件JSON");
  }
};

defineExpose({
  showImportDialog,
});
</script>

<style>
</style>