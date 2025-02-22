<template>
  <t-card>
    <div v-if="fileDataList.length > 0">
      <t-typography-title level="h5">转换结果</t-typography-title>
      <t-collapse>
        <t-collapse-panel
          v-for="fileData in fileDataList"
          :key="fileData.name"
          :header="fileData.name"
        >
          <template #headerRightContent>
            <t-button
              theme="primary"
              variant="text"
              shape="round"
              @click="downloadTxtFile(fileData.result, fileData.name)"
            >
              下载明文
            </t-button>
            <t-tooltip :content="fileData.aesKey" placement="top">
              <t-button theme="primary" variant="text" shape="round">
                查看密钥
              </t-button>
            </t-tooltip>
          </template>
          <pre style="max-height: 400px; overflow-y: auto">{{
            fileData.result
          }}</pre>
        </t-collapse-panel>
      </t-collapse>
    </div>
    <div v-else>
      <t-typography-title level="h5">没有转换结果</t-typography-title>
    </div>

    <div style="margin-top: 30px">
      <t-typography-text theme="secondary">
        确认转换结果后，将提交至后端，完成数字信封生成，并将转换后文本上载到IPFS中
      </t-typography-text>
    </div>

    <div style="margin-top: 20px; text-align: right">
      <t-button
        @click="lastStep"
        theme="primary"
        variant="outline"
        shape="round"
        >上一步
      </t-button>
      <t-button
        theme="primary"
        variant="outline"
        shape="round"
        @click="confirmResult"
        style="margin-left: 20px"
      >
        确认转换结果
      </t-button>
    </div>
  </t-card>
</template>

<script setup>
import { ref, defineExpose, defineEmits, reactive } from "vue";
import { useRoute } from "vue-router";
// import { getAesKey } from "@/api/upload";
import { MessagePlugin } from "tdesign-vue-next";

const route = useRoute();
const fileDataList = reactive([]);

const PassData = (newVal) => {
  fileDataList.splice(0, fileDataList.length, ...JSON.parse(newVal));
  //   getAesKey(fileDataList.length).then((res) => {
  //     for (let i = 0; i < fileDataList.length; i++) {
  //       console.log(res.data[i]);
  //       fileDataList[i].aesKey = res.data[i];
  //     }
  //   });

  fileDataList.forEach((fileData) => {
    fileData.aesKey = generateSecureAesKey();
  });
};

const generateSecureAesKey = (length = 16) => {
  const array = new Uint8Array(length);
  window.crypto.getRandomValues(array);
  return Array.from(array, (byte) => ("0" + byte.toString(16)).slice(-2)).join(
    ""
  );
};

const emit = defineEmits(["lastStep", "confirmResult"]);
const lastStep = () => {
  emit("lastStep");
};
const confirmResult = () => {
  if (fileDataList.length === 0) {
    MessagePlugin.warning("请先上传文件", 3000);
    return;
  } else {
    fileDataList.forEach((fileData) => {
      if (!fileData.result || !fileData.aesKey) {
        MessagePlugin.warning("请先完成所有文件的转换并生成密钥", 3000);
        return;
      }
    });
  }
  emit("confirmResult", JSON.stringify(fileDataList));
};

// 下载txt文件
const downloadTxtFile = (content, filename) => {
  // 下载包含content的txt文件
  const element = document.createElement("a");
  const file = new Blob([content], { type: "text/plain" });
  element.href = URL.createObjectURL(file);
  (element.download = filename.split(".").slice(0, -1).join(".")) + ".txt";
  document.body.appendChild(element);
  element.click();
  document.body.removeChild(element);
};

defineExpose({
  PassData,
});
</script>

<style scoped>
.file-data {
  margin-top: 20px;
}
</style>
