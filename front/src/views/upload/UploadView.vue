<template>
  <div class="upload-view">
    <t-steps :current="nowStep">
      <t-step-item
        title="上传文件"
        content="上传原始文件，转换为特定格式"
        :value="0"
        @click="changeToStep(0)"
      />
      <t-step-item
        title="转换确认"
        content="确认转换结果"
        :value="1"
        @click="changeToStep(1)"
      />
      <t-step-item
        title="上载生成"
        content="上载加密文本至IPFS并生成数字信封"
        :value="2"
      />
    </t-steps>
    <!-- <component :is="`step$
  components: { uploadInitFile },{nowStep.value}`" /> -->
    <div class="card-container">
      <upload-init-file v-show="nowStep == 0" @initFileOK="initFileOK" />
      <display-file-data
        v-show="nowStep == 1"
        ref="displayFileDataRef"
        @lastStep="lastStep"
        @confirmResult="confirmResult"
      />
      <uploadToIPFS v-show="nowStep == 2" ref="uploadToIPFSRef" />
    </div>
    <UserIdCheckAndSetDialog ref="UserIdCheckAndSetDialogRef" />
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import uploadInitFile from "./components/uploadInitFile.vue";
import displayFileData from "./components/DisplayFileData.vue";
import uploadToIPFS from "./components/UploadToIPFS.vue";
import { MessagePlugin } from "tdesign-vue-next";
import UserIdCheckAndSetDialog from "@/components/UserIdCheckAndSetDialog.vue";

const nowStep = ref(0);
const displayFileDataRef = ref(null);
const uploadToIPFSRef = ref(null);

const lastStep = () => {
  if (nowStep.value > 0) {
    nowStep.value--;
  }
};
const nextStep = () => {
  if (nowStep.value < 3) {
    nowStep.value++;
  }
};

const changeToStep = (step) => {
  if (nowStep.value == 2) {
    MessagePlugin.info("上载中不允许返回上一步");
    return;
  }
  nowStep.value = step;
};

// 接收子组件传递的jsonSuccessFileDataList
const initFileOK = (jsonSuccessFileDataList) => {
  nextStep();
  // 将jsonSuccessFileDataList传递到下一步
  displayFileDataRef.value.PassData(jsonSuccessFileDataList);
};

// 第二步确认了转换结果
const confirmResult = (jsonSuccessFileDataList) => {
  nextStep();
  // 将jsonSuccessFileDataList传递到下一步
  uploadToIPFSRef.value.PassData(jsonSuccessFileDataList);
};

const UserIdCheckAndSetDialogRef = ref(null);
onMounted(() => {
  UserIdCheckAndSetDialogRef.value.checkAndShowUserIdSetDialog();
});
</script>

<style>
.upload-view {
  margin: 20px auto;
  width: 60%;
}
.btn-group {
  display: flex;
  justify-content: center;
  margin-top: 30px;
}

.card-container {
  margin-top: 20px;
}
</style>