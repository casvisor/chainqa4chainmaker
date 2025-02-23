<template>
  <div class="log-view">
    <t-typography-title level="h5">
      <span style="color: #0052d9">审计日志</span>
    </t-typography-title>
    <div class="query-method">
      <t-radio-group
        v-model="queryMethod"
        variant="outline"
        :readonly="logSearchLoading"
      >
        <t-radio-button value="userId">按用户账号查询</t-radio-button>
        <t-radio-button value="timeRange">按时间范围查询</t-radio-button>
      </t-radio-group>
    </div>
    <div class="query-form">
      <t-form :data="form" ref="formRef" label-width="100px" :rules="rules">
        <t-form-item
          v-if="queryMethod === 'userId'"
          label="用户账号"
          name="userId"
        >
          <t-input
            v-model="form.userId"
            placeholder="请输入用户账号"
            :readonly="logSearchLoading"
          />
        </t-form-item>
        <t-form-item
          v-if="queryMethod === 'timeRange'"
          label="时间范围"
          name="timeRange"
        >
          <t-date-range-picker
            v-model="form.timeRange"
            placeholder="请选择时间范围"
            enable-time-picker
            :presets="presets"
            :readonly="logSearchLoading"
          />
        </t-form-item>
        <t-button
          theme="primary"
          @click="handleSubmit"
          style="margin-top: 30px"
          :loading="logSearchLoading"
        >
          查询日志
        </t-button>
      </t-form>
    </div>

    <!-- loading -->
    <t-loading :loading="logSearchLoading" text="查询中...">
      <div style="margin-top: 20px; color: grey">
        共有{{ resultCounts }}条日志记录
      </div>
      <div style="margin-top: 20px; color: grey; font-size: 12px">
        说明：由于查询条件和查询结果的内容较多，需点击“查看”打开弹窗展示
      </div>
      <t-table
        :data="tableData"
        :columns="columns"
        :pagination="false"
        style="margin-top: 20px"
      >
        <template #QueryStatus="{ row }">
          <t-tag
            v-if="row.QueryStatus === -1"
            variant="light-outline"
            theme="danger"
            >查询失败</t-tag
          >
          <span v-else>{{ row.QueryStatus }}条</span>
        </template>
        <template #Timestamp="{ row }">
          {{ dayjs(row.Timestamp * 1000).format("YYYY-MM-DD HH:mm:ss") }}
        </template>
        <template #QueryItemAndResult="{ row }">
          <t-button
            @click="
              ShowQueryItemAndResultDrawer(
                row.QueryResult,
                row.QueryItem,
                row.QueryId
              )
            "
            theme="primary"
            variant="text"
            >查看</t-button
          >
        </template>
      </t-table>
      <t-button @click="downloadLog" theme="default" variant="outline"
        >下载日志记录（Excel）</t-button
      >
      <query-item-and-result-drawer
        ref="QueryItemAndResultDrawerRef"
      ></query-item-and-result-drawer>
    </t-loading>
    <UserIdCheckAndSetDialog ref="UserIdCheckAndSetDialogRef" />
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from "vue";
import { MessagePlugin } from "tdesign-vue-next";
import dayjs from "dayjs";
import { logByTimeRange, logByUid } from "@/api/log";
import * as XLSX from "xlsx";
import QueryItemAndResultDrawer from "./components/QueryItemAndResultDrawer.vue";
import UserIdCheckAndSetDialog from "@/components/UserIdCheckAndSetDialog.vue";

const queryMethod = ref("userId"); // 用于切换查询方式
const form = reactive({
  userId: "",
  timeRange: [],
});

const presets = ref({
  最近7天: [dayjs().subtract(6, "day").toDate(), dayjs().toDate()],
  最近3天: [dayjs().subtract(2, "day").toDate(), dayjs().toDate()],
  今天: [dayjs().toDate(), dayjs().toDate()],
});

const rules = {
  userId: [{ required: true, message: "请输入用户账号", trigger: "blur" }],
  timeRange: [{ required: true, message: "请选择时间范围", trigger: "change" }],
};

const logSearchLoading = ref(false); // 查询loading
const formRef = ref(null);
const handleSubmit = async () => {
  const validateResult = await formRef.value.validate(); // 获取验证结果
  if (validateResult === true) {
    // 验证通过
    logSearchLoading.value = true;
    if (queryMethod.value === "userId") {
      await searchByUserId();
    } else {
      await searchByTimeRange();
    }
    logSearchLoading.value = false;
  } else {
    // 验证失败
    MessagePlugin.error("请检查输入项");
  }
};

const searchByUserId = async () => {
  // 根据用户账号查询
  var uId = form.userId;
  var result = await logByUid(uId).catch((err) => {
    MessagePlugin.error("查询失败");
    logSearchLoading.value = false;
  });

  if (result.code === 0) {
    tableData.splice(
      0,
      tableData.length,
      ...JSON.parse(result.data).QueryLogArray
    );
    resultCounts.value = JSON.parse(result.data).Count;
    MessagePlugin.success("查询成功");
  } else {
    MessagePlugin.error("查询失败");
  }
};

const searchByTimeRange = async () => {
  // 根据时间范围查询
  // 将时间均转为秒级时间戳
  var startTime = dayjs(form.timeRange[0]).unix().toString();
  var endTime = dayjs(form.timeRange[1]).unix().toString();
  console.log(startTime, endTime);
  var result = await logByTimeRange(startTime, endTime).catch((err) => {
    MessagePlugin.error("查询失败");
    logSearchLoading.value = false;
  });

  if (result.code === 0) {
    tableData.splice(
      0,
      tableData.length,
      ...JSON.parse(result.data).QueryLogArray
    );
    resultCounts.value = JSON.parse(result.data).Count;

    MessagePlugin.success("查询成功");
  } else {
    MessagePlugin.error("查询失败");
  }
};

// 表格数据
const tableData = reactive([]);
const resultCounts = ref(0);
const columns = [
  { colKey: "serial-number", title: "序号", width: 80 },
  { colKey: "QueryId", title: "查询ID" },
  { colKey: "Uid", title: "用户账号" },
  { colKey: "Timestamp", title: "查询时间", cell: "Timestamp", width: 150 },
  // { colKey: "QueryItem", title: "查询项", width: 300 },
  { colKey: "QueryStatus", title: "查询状态", cell: "QueryStatus" },
  // { colKey: "QueryResult", title: "查询结果", width: 300, cell: "QueryResult" },
  {
    colKey: "QueryItemAndResult",
    title: "查询条件与结果",
    cell: "QueryItemAndResult",
  },
];

const downloadLog = () => {
  // 创建一个工作簿
  const wb = XLSX.utils.book_new();
  let tableDataCopy = JSON.parse(JSON.stringify(tableData));
  tableDataCopy.forEach((item) => {
    item.Timestamp = dayjs(item.Timestamp * 1000).format("YYYY-MM-DD HH:mm:ss");
  });
  // 在第一行插入表头
  tableDataCopy.unshift({
    QueryId: "查询ID",
    Uid: "用户账号",
    Timestamp: "查询时间",
    QueryItem: "查询项",
    QueryStatus: "查询状态",
    QueryResult: "查询结果",
  });
  // 行末尾
  tableDataCopy.push({
    QueryId: "",
  });
  tableDataCopy.push({
    QueryId: "共有" + resultCounts.value + "条日志记录",
  });
  tableDataCopy.push({
    QueryId:
      "备注：查询状态为-1代表查询错误，为正数代表查询匹配的条数（0标识无结果）",
  });
  const ws = XLSX.utils.json_to_sheet(tableDataCopy);

  // 将工作表添加到工作簿
  XLSX.utils.book_append_sheet(wb, ws, "Sheet1");
  // 将工作簿写入文件
  XLSX.writeFile(wb, "log.xlsx");
  MessagePlugin.success("下载成功！");
};

const QueryItemAndResultDrawerRef = ref(null);
const ShowQueryItemAndResultDrawer = (queryResult, queryItem, queryId) => {
  QueryItemAndResultDrawerRef.value.showQueryDrawer(
    queryResult,
    queryItem,
    queryId
  );
};

const UserIdCheckAndSetDialogRef = ref(null);
onMounted(() => {
  UserIdCheckAndSetDialogRef.value.checkAndShowUserIdSetDialog(false, false);
});
</script>

<style scoped>
.log-view {
  margin: 20px auto;
  width: 60%;
}

.query-method {
  margin-bottom: 20px;
}

.query-form {
  margin-top: 20px;
}
</style>