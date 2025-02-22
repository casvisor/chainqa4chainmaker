<template>
  <t-drawer
    v-model:visible="visible"
    :on-overlay-click="() => (visible = false)"
    placement="bottom"
    size="90%"
    @cancel="visible = false"
  >
    <template #footer>
      <div style="width: 100%; display: flex; justify-content: flex-end">
        <t-button theme="default" @click="visible = false"> 关闭 </t-button>
      </div>
    </template>
    <template #header>
      <div
        style="
          display: flex;
          justify-content: space-between;
          width: 100%;
          align-items: center;
        "
      >
        <div>
          日志解析：查询结果与查询条件&nbsp;&nbsp;&nbsp;&nbsp;
          <span style="color: #999; font-weight: 400; font-size: 12px"
            >查询记录ID:{{ query.queryId }}</span
          >
        </div>
        <div>
          <t-radio-group
            v-model="showType"
            variant="primary-filled"
            @change="onChange"
          >
            <t-radio-button value="result">查看查询结果</t-radio-button>
            <t-radio-button value="item">查看查询条件</t-radio-button>
          </t-radio-group>
        </div>
      </div>
    </template>

    <div style="padding: 20px; padding-top: 10px; overflow-y: auto">
      <query-result-inner
        ref="queryResultInnerRef"
        v-if="showType === 'result'"
      >
      </query-result-inner>
      <query-item-inner
        v-else-if="showType === 'item'"
        ref="queryItemInnerRef"
      ></query-item-inner>
    </div>
  </t-drawer>
</template>

<script setup>
import { ref, defineExpose, reactive, nextTick } from "vue";
import QueryResultInner from "./QueryResultInner.vue";
import QueryItemInner from "./QueryItemInner.vue";

const visible = ref(false);

const showType = ref("result"); // result, item

const query = reactive({
  queryResultJSON: "",
  queryItemJSON: "",
  queryId: "",
});

const queryResultInnerRef = ref(null);
const queryItemInnerRef = ref(null);
const showQueryDrawer = (queryResult_, queryItem_, queryId_) => {
  // 判断queryResult_是否是JSON字符串
  let queryResJSON;
  try {
    queryResJSON = JSON.stringify(JSON.parse(queryResult_));
  } catch (e) {
    // 非JSON字符串，说明是查询失败
    let failObj = {
      counts: -1,
      data: [],
      message: queryResult_,
    };
    queryResJSON = JSON.stringify(failObj);
  }

  let queryItemJSON = JSON.stringify(JSON.parse(queryItem_));
  console.log("queryResJSON", queryResJSON);
  query.queryResultJSON = queryResJSON;
  query.queryItemJSON = queryItemJSON;
  query.queryId = queryId_;

  showType.value = "result";
  nextTick(() => {
    queryResultInnerRef.value.showQueryResultInner(query.queryResultJSON);
  });

  visible.value = true;
};

const onChange = (e) => {
  if (e === "result") {
    showType.value = "result";
    nextTick(() => {
      queryResultInnerRef.value.showQueryResultInner(query.queryResultJSON);
    });
  } else {
    showType.value = "item";
    nextTick(() => {
      queryItemInnerRef.value.showQueryItemInner(query.queryItemJSON);
    });
  }
};

defineExpose({
  showQueryDrawer,
});
</script>

<style>
</style>