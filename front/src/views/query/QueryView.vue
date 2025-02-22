<template>
  <div class="query-view">
    <t-typography-title level="h5">
      <span style="color: #0052d9">查询数据</span>
    </t-typography-title>
    <QueryForm @submitQuery="submitQuery" v-show="status === 'query'" />
    <QueryResult ref="queryResultRef" v-show="status === 'result'" />
    <UserIdCheckAndSetDialog ref="UserIdCheckAndSetDialogRef" />
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import QueryForm from "./components/QueryForm.vue";
import QueryResult from "./components/QueryResult.vue";
import UserIdCheckAndSetDialog from "@/components/UserIdCheckAndSetDialog.vue";

const status = ref("query"); //query or result

const queryResultRef = ref();
const submitQuery = (queryItem, filePoses) => {
  status.value = "result";
  console.log(queryItem, filePoses);
  console.log(queryResultRef.value);
  queryResultRef.value.queryAndShowResult(queryItem, filePoses);
};

const UserIdCheckAndSetDialogRef = ref(null);
onMounted(() => {
  UserIdCheckAndSetDialogRef.value.checkAndShowUserIdSetDialog();
});
</script>

<style scoped>
.query-view {
  margin: 20px auto;
  width: 60%;
}
</style>