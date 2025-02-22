<template>
  <div>
    <t-typography-title level="h3">查询结果</t-typography-title>
    <div style="color: #999; font-size: 12px; margin-bottom: 20px">
      说明：当前展示的是查询结果，点击右上角按钮可切换。
    </div>
    <div v-if="queryResult.counts < 0">
      <div style="color: red; font-size: 18px; margin-top: 30px">
        {{ queryResult.message }}
      </div>
      <div style="font-size: 14px; margin-top: 10px">
        在该条件下，查询失败。无法获取到有效数据。
      </div>
    </div>
    <div v-else>
      <div style="margin-top: 10px">
        {{ queryResult.message }}，共有{{ queryResult.counts }}条记录
      </div>
      <QueryResultTable ref="queryResultTableRef"></QueryResultTable>
      <t-divider style="margin-top: 40px"></t-divider>
      <t-collapse :default-value="[]">
        <t-collapse-panel header="查询结果JSON">
          <template #headerRightContent>
            <t-space size="small">
              <t-button size="small" @click="copyToClipBoard(queryResult)">
                复制
              </t-button>
            </t-space>
          </template>
          <div>
            <t-typography-text style="word-break: break-all">{{
              queryResult
            }}</t-typography-text>
          </div>
        </t-collapse-panel>
      </t-collapse>
    </div>
  </div>
</template>

<script setup>
import { ref, defineExpose, reactive, watch, nextTick } from "vue";
import QueryResultTable from "@/views/query/components/QueryResultTable.vue";
import copyToClipBoard from "@/utils/copyToclipBoard";

const queryResult = reactive({
  counts: 0,
  data: [],
  message: "",
});

const queryResultTableRef = ref(null);
const showQueryResultInner = (queryResult_) => {
  let queryResObj = JSON.parse(queryResult_);
  queryResult.counts = queryResObj.counts;
  queryResult.data = queryResObj.data;
  queryResult.message = queryResObj.message;
  nextTick(() => {
    queryResultTableRef.value.ShowQueryResultTable(
      queryResult.counts,
      queryResult.data
    );
  });
};

defineExpose({
  showQueryResultInner,
});
</script>

<style>
</style>