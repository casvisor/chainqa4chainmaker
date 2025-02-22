<template>
  <div class="query-result">
    <div
      style="margin: 20px auto"
      v-if="dataStatus === 'loading'"
      class="load-card"
    >
      <t-loading
        style="margin-left: 50%"
        size="large"
        text="正在查询中..."
      ></t-loading>
    </div>

    <!-- 错误 -->
    <div v-else-if="dataStatus === 'error'" class="error-card">
      <t-alert
        theme="warning"
        title="查询出错了"
        :message="
          errMsg +
          ' 。若提示500，可能是未上传数字信封。请再次确认ipfs地址正确。'
        "
      >
        <template #operation>
          <span @click="handleRetry">点击重试</span>
        </template>
      </t-alert>
    </div>

    <!-- 成功 -->
    <div v-else-if="dataStatus === 'success'">
      <t-alert
        theme="success"
        title="查询成功"
        :message="successData.counts + '条数据'"
      >
        <template #operation>
          <span @click="handleRetry">点击重试</span>
        </template>
      </t-alert>

      <QueryResultTable ref="queryResultTableRef"></QueryResultTable>
    </div>

    <!-- 附加 -->
    <t-divider style="margin-top: 40px"> </t-divider>
    <div style="margin-top: 40px">
      <t-collapse :default-value="[]">
        <t-collapse-panel header="查询结果JSON" v-if="dataStatus !== 'loading'">
          <template #headerRightContent>
            <t-space size="small">
              <t-button
                size="small"
                @click="copyToClipBoard(JSON.stringify(resultJson))"
              >
                复制
              </t-button>
            </t-space>
          </template>
          <div>
            <t-typography-text style="word-break: break-all">
              {{ resultJson }}
            </t-typography-text>
          </div>
        </t-collapse-panel>
        <t-collapse-panel header="查询条件JSON">
          <template #headerRightContent>
            <t-space size="small">
              <t-button
                size="small"
                @click="copyToClipBoard(queryItem.toString())"
              >
                复制
              </t-button>
            </t-space>
          </template>
          <div>
            <t-typography-text style="word-break: break-all">{{
              queryItem
            }}</t-typography-text>
          </div>
        </t-collapse-panel>
        <t-collapse-panel header="查询涉及的IPFS">
          <t-typography-text>
            <span v-for="(f, idx) in filePoses" :key="f">
              <span style="font-weight: bold">【分片{{ idx + 1 }}】</span
              >{{ f }}&nbsp;&nbsp;<br />
            </span>
          </t-typography-text>
        </t-collapse-panel>
      </t-collapse>
    </div>
  </div>
</template>

<script setup>
import { ref, defineExpose, reactive, watch, nextTick } from "vue";
import { queryData } from "@/api/query";
import { MessagePlugin } from "tdesign-vue-next";
import QueryResultTable from "./QueryResultTable.vue";
import copyToClipBoard from "@/utils/copyToclipBoard";

const queryItem = ref(null);
const filePoses = ref(null);
const dataStatus = ref("loading"); //loading, success, error
const errMsg = ref("");
const resultJson = ref("");
const successData = reactive({
  counts: 0,
  data: [],
});

const queryAndShowResult = async (qItem, fPoses) => {
  dataStatus.value = "loading";
  queryItem.value = qItem;
  filePoses.value = fPoses;
  let resp = await queryData(qItem, fPoses).catch((err) => {
    dataStatus.value = "error";
    console.log(err);
    errMsg.value = "查询失败！原因 —— " + err?.response?.data?.msg;
    MessagePlugin.error("查询失败！");
  });
  // let resp = {
  //   code: 0,
  //   data: '{"counts":8,"data":[{"Cscore":"9","name":"大兰"},{"Cscore":"9","name":"大黑"},{"Cscore":"9","name":"大橙"},{"Cscore":"9","name":"大黄"},{"Cscore":"9","name":"小兰"},{"Cscore":"9","name":"小黑"},{"Cscore":"9","name":"小橙"},{"Cscore":"9","name":"小黄"}],"message":"查询成功"}',
  //   msg: "查询成功",
  // };
  resultJson.value = resp;
  if (resp.code === 0) {
    if (resp.msg.includes("出现错误")) {
      dataStatus.value = "error";
      MessagePlugin.warning(resp.data);
      errMsg.value = resp.data;
      return;
    } else {
      dataStatus.value = "success";
      // 处理查询结果，将结果存储到successData中
      successData.counts = JSON.parse(resp.data).counts; // 查询结果总数
      successData.data = JSON.parse(resp.data).data; // 查询结果数据
      showQueryResultTable(successData.counts, successData.data);
      MessagePlugin.success(JSON.parse(resp.data).message);
    }
  } else {
    dataStatus.value = "error";
    errMsg.value = "查询失败！原因未知";
    MessagePlugin.error("查询失败！");
  }
};

const queryResultTableRef = ref(null);
const showQueryResultTable = (counts, data) => {
  nextTick(() => {
    queryResultTableRef.value.ShowQueryResultTable(counts, data);
  });
};

const handleRetry = () => {
  queryAndShowResult(queryItem.value, filePoses.value);
};

defineExpose({
  queryAndShowResult,
});
</script>

<style>
</style>