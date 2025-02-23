<template>
  <div class="mytitle">
    <div class="config-icon">
      <t-tooltip content="点击查看帮助文档" placement="bottom">
        <HelpCircleIcon @click="showHelpDoc" style="margin-left: 20px" />
      </t-tooltip>
      <t-tooltip
        content="API配置，请勿随意修改"
        placement="bottom"
        theme="warning"
      >
        <Setting1Icon @click="showApiConfigDialog" style="margin-left: 20px" />
      </t-tooltip>
      <t-tooltip content="返回主页" placement="bottom">
        <HomeIcon
          @click="routeTo('/')"
          style="margin-left: 20px; margin-right: 20px"
        />
      </t-tooltip>
    </div>
    <t-typography-title>
      <t-tooltip
        content="点击返回“区块链复杂查询与审计模块”主页"
        placement="bottom"
      >
        <span @click="routeTo('/')" style="cursor: pointer"
          >区块链复杂查询与审计模块</span
        >
      </t-tooltip>

      <!-- <t-tag theme="primary" @click="routeTo('/')" style="cursor: pointer">
        V0.0.6
      </t-tag> -->
    </t-typography-title>
    <t-tooltip content="点击修改用户信息" placement="bottom">
      <div
        style="
          border-radius: 20px;
          background-color: #e7f5ff;
          padding: 5px;
          width: 200px;
          margin: 0 auto;
          font-size: 14px;
          cursor: pointer;
        "
        @click="checkAndSetUserId"
      >
        <div v-if="userId">用户: {{ userId }}</div>
        <div v-else>未设置用户账号</div>
      </div>
    </t-tooltip>
    <t-divider></t-divider>
    <ApiConfigDialog ref="ApiConfigDialogRef" />
    <UserIdCheckAndSetDialog ref="UserIdCheckAndSetDialogRef" />
  </div>
</template>

<script setup>
import { computed, ref } from "vue";
import { useUserStore } from "../stores/user";
import router from "../router";
import { Setting1Icon, HomeIcon, HelpCircleIcon } from "tdesign-icons-vue-next";
import ApiConfigDialog from "./ApiConfigDialog.vue";
import UserIdCheckAndSetDialog from "./UserIdCheckAndSetDialog.vue";

const userStore = useUserStore();
const userId = computed(() => userStore.userId);
const ApiConfigDialogRef = ref(null);

const routeTo = (path) => {
  router.push(path);
};

const showApiConfigDialog = () => {
  // ApiConfigDialogRef.value.visible = true;
  ApiConfigDialogRef.value.showApiConfigDialog();
};

const showHelpDoc = () => {
  window.open("https://www.yuque.com/jjq0425/ei0f4f/acgk6l3ax4gf98wh");
};

const UserIdCheckAndSetDialogRef = ref(null);
const checkAndSetUserId = () => {
  UserIdCheckAndSetDialogRef.value.checkAndShowUserIdSetDialog(true);
};
</script>

<style scoped>
.mytitle {
  text-align: center;
}

.config-icon {
  position: absolute;
  right: 20px;
  top: 20px;
  cursor: pointer;
}
</style>
