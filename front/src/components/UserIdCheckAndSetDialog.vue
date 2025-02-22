<template>
  <t-dialog
    header=""
    :visible="visible"
    :confirm-btn="null"
    :cancel-btn="null"
    :width="400"
    :closeOnEscKeydown="false"
    :closeOnOverlayClick="false"
    :closeBtn="false"
  >
    <t-typography-title level="h5">请先设置用户ID</t-typography-title>

    <t-input
      v-model="userId"
      placeholder="请输入用户ID"
      style="margin-top: 10px"
    ></t-input>

    <div>
      <t-button theme="primary" @click="setAndClose" style="margin-top: 20px"
        >设置并关闭</t-button
      >
    </div>
  </t-dialog>
</template>

<script setup>
import { ref, defineExpose } from "vue";
import { useUserStore } from "@/stores/user";
import { MessagePlugin } from "tdesign-vue-next";
import router from "@/router";
const userStore = useUserStore();

const visible = ref(false);
const userId = ref("");

const setAndClose = () => {
  if (userId.value) {
    userStore.setUserId(userId.value);
    visible.value = false;
    MessagePlugin.success("设置成功");
  } else {
    MessagePlugin.warning("用户ID不能为空");
  }
};

const checkAndShowUserIdSetDialog = () => {
  // 检查router中的userId是否为空
  if (router.currentRoute.value.query.userId) {
    userStore.setUserId(router.currentRoute.value.query.userId);
    return;
  } else if (userStore.userId != null && userStore.userId != "") {
    return;
  }
  MessagePlugin.warning("此操作需先设置用户ID");
  visible.value = true;
};

defineExpose({
  checkAndShowUserIdSetDialog,
});
</script>

<style>
</style>