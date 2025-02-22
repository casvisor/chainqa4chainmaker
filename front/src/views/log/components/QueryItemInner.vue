<template>
  <div>
    <t-typography-title level="h3">查询条件</t-typography-title>

    <div>
      <div style="color: #999; font-size: 12px; margin-bottom: 20px">
        说明：当前展示的是查询条件，点击右上角按钮可切换。
        <span style="font-weight: bold">本页内容只读，不能进行修改。</span>
      </div>
      <QueryForm ref="queryFormRef" />
      <t-divider style="margin-top: 40px"></t-divider>
      <t-collapse :default-value="[]">
        <t-collapse-panel header="查询条件JSON">
          <template #headerRightContent>
            <t-space size="small">
              <t-button size="small" @click="copyToClipBoard(queryItemJson)">
                复制
              </t-button>
            </t-space>
          </template>
          <div>
            <t-typography-text style="word-break: break-all">{{
              queryItemJson
            }}</t-typography-text>
          </div>
        </t-collapse-panel>
      </t-collapse>
    </div>
  </div>
</template>

<script setup>
import { ref, defineExpose, reactive, watch, nextTick } from "vue";
import QueryForm from "@/views/query/components/QueryForm.vue";
import backComposeQueryItem from "@/utils/backComposeQueryItem";
import copyToClipBoard from "@/utils/copyToclipBoard";
const queryFormRef = ref();

const queryItemJson = ref("");
const showQueryItemInner = (queryItem_) => {
  queryItemJson.value = queryItem_;
  let queryItemObj = backComposeQueryItem(queryItem_);
  nextTick(() => {
    queryFormRef.value.importNow(queryItemObj, true);
  });
};

defineExpose({
  showQueryItemInner,
});
</script>

<style>
</style>