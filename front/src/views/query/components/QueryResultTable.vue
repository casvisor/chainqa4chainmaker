<template>
  <div>
    <t-table :data="getData()" :columns="getCol()"></t-table>
    <div>
      <t-button type="primary" @click="downloadResult"
        >下载结果（Excel）</t-button
      >
    </div>
  </div>
</template>

<script setup>
import { ref, defineExpose, reactive, watch } from "vue";
import * as XLSX from "xlsx";

const successData = reactive({
  counts: 0,
  data: [],
});

// 检测到successData.data变化时，重新计算列，使用watch
watch(
  () => successData.data,
  () => {
    getCol();
    getData();
  }
);
const getCol = () => {
  if (successData.data && successData.data.length > 0) {
    let colarr = Object.keys(successData.data[0]).map((key) => ({
      title: key,
      colKey: key,
    }));
    // 在最前面增加序号列
    colarr.unshift({
      title: "序号",
      colKey: "serial-number",
      width: 80,
      align: "center",
    });
    return colarr;
  }
  return [];
};
const getData = () => {
  if (successData.data && successData.data.length > 0) {
    return successData.data;
  }
  return [];
};

const downloadResult = () => {
  // 创建一个工作簿
  const wb = XLSX.utils.book_new();
  const ws = XLSX.utils.json_to_sheet(successData.data);
  // 将工作表添加到工作簿
  XLSX.utils.book_append_sheet(wb, ws, "Sheet1");
  // 将工作簿写入文件
  XLSX.writeFile(wb, "result.xlsx");
  MessagePlugin.success("下载成功！");
};

const ShowQueryResultTable = (counts_, data_) => {
  successData.counts = counts_;
  successData.data = data_;
};

defineExpose({
  ShowQueryResultTable,
});
</script>

<style>
</style>