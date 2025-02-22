<template>
  <t-card>
    <div class="upload-container-loading" v-if="uploadStatus == 'loading'">
      <div style="margin: 10px auto; width: 100px">
        <t-loading size="large" />
      </div>
      <div>
        <t-typography-title level="h6">
          正在上传至IPFS，共{{ fileDataList.length }}个任务，正在处理第
          {{ uploadingNowIdx + 1 }}个任务，成功
          {{ successFileDataList.length }}
          个，失败{{ failFileDataList.length }}个
        </t-typography-title>
      </div>
    </div>

    <div v-if="successFileDataList.length >= 1" class="success-table">
      <div style="color: #2ba471; font-weight: bold; font-size: 20px">
        成功列表
      </div>
      <t-table
        :data="successFileDataList"
        row-key="name"
        :columns="successColumns"
      >
        <template #operation="{ row }">
          <t-link
            theme="primary"
            hover="color"
            @click="handleDownloadModal(row)"
          >
            下载
          </t-link>
          <t-link
            theme="primary"
            hover="color"
            style="margin-left: 20px"
            @click="handleCopyUploadDataConfig(row)"
          >
            复制信息
          </t-link>
          <DownloadIPFSModal ref="downloadIPFSModalRef"></DownloadIPFSModal>
        </template>
      </t-table>
    </div>

    <div
      v-if="failFileDataList.length >= 1"
      class="fail-table"
      style="margin-top: 30px"
    >
      <div style="color: #ff922b; font-weight: bold; font-size: 20px">
        失败列表
      </div>
      <t-table :data="failFileDataList" row-key="name" :columns="failColumns">
      </t-table>
    </div>

    <div style="margin-top: 40px" v-else-if="uploadStatus == 'finish'">
      <div>
        <t-typography-title level="h6">
          任务结束，共{{ fileDataList.length }}个任务，成功
          {{ successFileDataList.length }}
          个，失败{{ failFileDataList.length }}个
        </t-typography-title>
      </div>
    </div>

    <t-progress
      theme="line"
      :percentage="(uploadingNowIdx / fileDataList.length) * 100"
      :status="loadingConfig.status"
    />
    <t-typography-text theme="secondary">
      上传至IPFS的文件将会被加密，加密密钥的数字信封将会被保存在区块链上
    </t-typography-text>
  </t-card>
</template>

<script setup>
import { ref, reactive, defineExpose } from "vue";
import { uploadFile } from "@/api/upload";
import { downloadIPFSFile, tryDecryptFile } from "@/api/download";
import { MessagePlugin } from "tdesign-vue-next";
import DownloadIPFSModal from "./DownloadIPFSModal.vue";

const fileDataList = reactive([]);
const uploadStatus = ref("loading"); // loading状态/finish状态
const loadingConfig = reactive({
  status: "active",
});
const uploadingNowIdx = ref(0); // 正在上传的文件索引
const failFileDataList = reactive([]);
const successFileDataList = reactive([]);

const PassData = async (jsonSuccessFileDataList) => {
  uploadStatus.value = "loading";
  loadingConfig.status = "active";
  uploadingNowIdx.value = 0;
  successFileDataList.splice(0, successFileDataList.length);
  failFileDataList.splice(0, failFileDataList.length);
  fileDataList.splice(
    0,
    fileDataList.length,
    ...JSON.parse(jsonSuccessFileDataList)
  );

  for (let i = 0; i < fileDataList.length; i++) {
    uploadingNowIdx.value = i; // 正在上传的文件索引
    let response = await uploadFile(
      fileDataList[i].aesKey,
      fileDataList[i].result,
      fileDataList[i].name
    ).catch((err) => {
      fileDataList[i].err = err?.response?.msg;
      // failFileDataList.splice(0, 0, fileDataList[i]);
    });

    // let response = {
    //   code: 0,
    //   data: {
    //     fileName: "7人测试数据（语文成绩单）.xlsx",
    //     pos: "QmX798NJ5UrTwkRKNwwvg7goD4JwiXY8SA8zs8qoxLJcQP",
    //   },
    //   msg: "上传文件7人测试数据（语文成绩单）.xlsx到IPFS成功",
    // };
    // key:27454d43ee3599afbb9a720be96e74d8
    if (response.code == 0) {
      fileDataList[i].ipfsPos = response.data.pos;
      successFileDataList.splice(0, 0, fileDataList[i]);
    } else {
      fileDataList[i].err = response.msg;
      failFileDataList.splice(0, 0, fileDataList[i]);
    }
    // 等待3s
    await new Promise((resolve) => {
      setTimeout(() => {
        resolve();
      }, 3000);
    });
  }

  if (successFileDataList.length == fileDataList.length) {
    uploadingNowIdx.value = fileDataList.length; // 正在上传的文件索引
    loadingConfig.status = "success";
  } else {
    uploadingNowIdx.value = fileDataList.length; // 正在上传的文件索引
    loadingConfig.status = "warning";
  }
  uploadStatus.value = "finish";
};

// ==================== 表格配置 ====================
const successColumns = [
  {
    title: "文件名",
    colKey: "name",
  },
  {
    title: "IPFS位置",
    colKey: "ipfsPos",
  },
  {
    title: "加密密钥",
    colKey: "aesKey",
  },
  {
    colKey: "operation",
    title: "操作",
    width: "180",
    foot: "-",
    fixed: "right",
  },
];

const failColumns = [
  {
    title: "文件名",
    colKey: "name",
  },
  {
    title: "错误信息",
    colKey: "err",
  },
];

// ==================== 操作 ====================
const downloadIPFSModalRef = ref(null);
const handleDownloadModal = (row) => {
  downloadIPFSModalRef.value.row = row;
  downloadIPFSModalRef.value.visible = true;
};

const handleCopyUploadDataConfig = async (row) => {
  const text = `文件名：${row.name}\nIPFS位置：${row.ipfsPos}\n加密密钥：${row.aesKey}`;
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
    await clipboard.writeText(text);
    MessagePlugin.success("复制成功");
  }
};

defineExpose({
  PassData,
});
</script>

<style scoped>
.upload-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 200px;
}

.success-table {
  margin-top: 20px;
}
</style>
