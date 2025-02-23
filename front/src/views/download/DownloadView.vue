<template>
  <div class="download-view">
    <t-typography-title level="h5">
      <span style="color: #0052d9">下载文件</span>
    </t-typography-title>
    <t-form
      :rules="rules"
      :data="form"
      label-width="100px"
      @submit="handleSubmit"
    >
      <t-form-item label="IPFS地址" name="ipfsPos">
        <t-input v-model="form.ipfsPos" placeholder="请输入IPFS地址" />
      </t-form-item>
      <t-form-item label="AES密钥" name="aesKey">
        <t-input v-model="form.aesKey" placeholder="请输入AES密钥（可选）" />
      </t-form-item>
      <t-form-item>
        <t-button theme="primary" type="submit">下载</t-button>
      </t-form-item>
    </t-form>
    <DownloadIPFSModal ref="downloadModalRef" />
    <UserIdCheckAndSetDialog ref="UserIdCheckAndSetDialogRef" />
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from "vue";
import DownloadIPFSModal from "../upload/components/DownloadIPFSModal.vue";
import { MessagePlugin } from "tdesign-vue-next";
import UserIdCheckAndSetDialog from "@/components/UserIdCheckAndSetDialog.vue";

const form = reactive({
  ipfsPos: "",
  aesKey: "",
});

const rules = {
  ipfsPos: [{ required: true, message: "IPFS地址为必填项", type: "error" }],
};

const formRef = ref(null);
const downloadModalRef = ref(null);

const handleSubmit = ({ validateResult, firstError, e }) => {
  e.preventDefault(); // 阻止表单默认提交行为
  if (validateResult === true) {
    if (form.aesKey == "") {
      form.aesKey = null;
    }
    form.name = form.ipfsPos;
    downloadModalRef.value.row = form;
    downloadModalRef.value.visible = true;
  } else {
    MessagePlugin.error("请检查输入项");
  }
};

const UserIdCheckAndSetDialogRef = ref(null);
onMounted(() => {
  UserIdCheckAndSetDialogRef.value.checkAndShowUserIdSetDialog(false, false);
});
</script>

<style scoped>
.download-view {
  margin: 20px auto;
  width: 40%;
}

.input-container {
  margin: 20px 0;
}

.button-container {
  text-align: center;
}
</style>