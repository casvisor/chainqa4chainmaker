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
    <t-typography-title level="h5">请设置用户账号</t-typography-title>

    <t-input
      v-model="userId"
      placeholder="请输入用户账号"
      style="margin-top: 10px"
    ></t-input>
    <p style="color: gray; font-size: 10px">
      *用户账号只能包含数字/字母/下划线，不能包含其他符号。
    </p>

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
    MessagePlugin.warning("用户账号不能为空");
  }
};

const checkAndShowUserIdSetDialog = (forceShow = false, required = true) => {
  // 检查router中的userId是否为空
  if (router.currentRoute.value.query.userId && !forceShow) {
    userStore.setUserId(router.currentRoute.value.query.userId);
    return;
  } else if (userStore.userId != null && userStore.userId != "" && !forceShow) {
    return;
  }
  // 不需要且不强制显示，直接返回
  if (!required && !forceShow) {
    return;
  }
  MessagePlugin.warning(forceShow ? "请设置用户账号" : "请先设置用户账号");
  visible.value = true;
};

defineExpose({
  checkAndShowUserIdSetDialog,
});
</script>

<style>
</style>