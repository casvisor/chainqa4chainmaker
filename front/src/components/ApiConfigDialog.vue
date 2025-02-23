<template>
  <t-dialog
    :visible.sync="visible"
    header="API配置设置"
    confirm-btn="保存"
    @confirm="saveConfig"
    @close="closeDialog"
    :destory-on-close="true"
  >
    <div class="config-item">
      <t-input v-model="api" label="后端API地址：" />
    </div>
    <div class="config-item">
      <t-input v-model="ipfsServiceUrl" label="IPFS服务器地址：" />
    </div>
    <div class="config-item">
      <t-input v-model="chainServiceUrl" label="区块链INVOKE SDK地址：" />
    </div>
    <div class="config-item">
      <t-input v-model="contractName" label="合约名称：" />
    </div>
    <div class="config-item" style="font-size: 12px; color: #999">
      版本号：V0.0.8.2
    </div>
    <t-alert theme="warning">
      <template #message>请勿随意更改！</template>
    </t-alert>
  </t-dialog>
</template>

<script setup>
import { MessagePlugin } from "tdesign-vue-next";
import { ref, onMounted } from "vue";
import { useApiStore } from "../stores/api";

const apiStore = useApiStore();
const visible = ref(false);
const api = ref(apiStore.api);
const ipfsServiceUrl = ref(apiStore.ipfsServiceUrl);
const chainServiceUrl = ref(apiStore.chainServiceUrl);
const contractName = ref(apiStore.contractName);

const saveConfig = () => {
  apiStore.setApi(api.value);
  apiStore.setIpfsServiceUrl(ipfsServiceUrl.value);
  apiStore.setChainServiceUrl(chainServiceUrl.value);
  apiStore.setContractName(contractName.value);
  MessagePlugin.success("保存成功");
};

const showApiConfigDialog = () => {
  visible.value = true;
  api.value = apiStore.api;
  ipfsServiceUrl.value = apiStore.ipfsServiceUrl;
  chainServiceUrl.value = apiStore.chainServiceUrl;
  contractName.value = apiStore.contractName;
};
const closeDialog = () => {
  visible.value = false;
};

defineExpose({
  visible,
  showApiConfigDialog,
});
</script>

<style scoped>
.config-item {
  margin-bottom: 20px;
}
</style>
