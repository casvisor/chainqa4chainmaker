<template>
  <t-dialog
    :visible.sync="visible"
    header="选择操作"
    :confirm-btn="null"
    @close="closeDialog"
  >
    <div style="margin-bottom: 20px">
      <div>您将要下载分片：{{ row?.ipfsPos }}</div>
      <t-tag theme="success" variant="light-outline" v-if="row?.aesKey != null"
        >已传入密钥，可进行解密</t-tag
      >
      <t-tag theme="warning" variant="light-outline" v-if="row?.aesKey == null"
        >未传入密钥，只能下载密文</t-tag
      >
    </div>
    <div class="modal-content">
      <t-button theme="primary" @click="downloadIPFS">下载密文</t-button>
      <t-button
        theme="primary"
        @click="tryDecrypt"
        style="margin-left: 20px"
        :disabled="row == null || row?.aesKey == null"
      >
        解密下载
      </t-button>
    </div>
    <t-alert
      style="margin-top: 20px"
      theme="info"
      message="下载密文后不能自主解密，均需通过本接口解密，因为AES的密钥填充、IV值等情况略有不同；上述按钮所有操作均是从IPFS下载后进行，未涉及到区块链数字信封操作，可验证用户是否上传成功。"
    />

    <t-alert
      v-if="row == null || row?.aesKey == null"
      style="margin-top: 20px"
      theme="warning"
      message="当前未传入AES密钥，无法进行解密操作"
    />
    <!-- <div>{{ row }}</div> -->
  </t-dialog>
</template>

<script setup>
import { ref } from "vue";
import { downloadIPFSFile, tryDecryptFile } from "@/api/download";
import { MessagePlugin } from "tdesign-vue-next";

const visible = ref(false);
const row = ref(null);

const downloadIPFS = async () => {
  var response = await downloadIPFSFile(row.value.ipfsPos).catch((err) => {
    MessagePlugin.warning("下载IPFS文件失败：" + err?.response?.data?.data);
  });
  if (response.code == 0) {
    var blob = new Blob([response.data], { type: "application/octet-stream" });
    var link = document.createElement("a");
    link.href = window.URL.createObjectURL(blob);
    link.download =
      row.value.name.split(".").slice(0, -1).join(".") + "(密文)" + ".txt";
    link.click();
    MessagePlugin.success("下载IPFS文件成功，请注意该文件为密文文件");
  } else {
    MessagePlugin.error("下载IPFS文件失败");
  }
  visible.value = false;
};

const tryDecrypt = async () => {
  var response = await tryDecryptFile(
    row.value.ipfsPos,
    row.value.aesKey
  ).catch((err) => {
    MessagePlugin.warning("尝试解密失败：" + err?.response?.data?.data);
  });
  if (response.code == 0) {
    var blob = new Blob([response.data], { type: "application/octet-stream" });
    var link = document.createElement("a");
    link.href = window.URL.createObjectURL(blob);
    link.download =
      row.value.name.split(".").slice(0, -1).join(".") + "(明文)" + ".txt";
    link.click();
    MessagePlugin.success("尝试解密成功，请注意该文件为明文文件");
  } else {
    MessagePlugin.error("尝试解密失败");
  }
  visible.value = false;
};

const closeDialog = () => {
  visible.value = false;
};

defineExpose({
  visible,
  row,
});
</script>

<style scoped>
.modal-content {
  display: flex;
  justify-content: center;
  gap: 20px;
}
</style>
