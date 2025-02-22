<template>
  <div class="home">
    <t-typography-title level="h5">
      <span style="color: #0052d9">#1 设置用户ID</span>
    </t-typography-title>
    <div class="input-container input-user-id">
      <t-input v-model="uId" placeholder="请输入用户Id" />
      <t-button @click="saveUserId">设置ID</t-button>
    </div>

    <t-typography-title level="h5" style="color: #0052d9">
      <span style="color: #0052d9">#2 选择操作</span>
    </t-typography-title>
    <div class="choose">
      <t-card
        class="option-card"
        hover-shadow
        @click="routeTo('upload')"
        subtitle="上传源文件以供查询使用"
      >
        <template #title>上传信息</template>
      </t-card>
      <t-card
        class="option-card"
        hover-shadow
        @click="routeTo('query')"
        subtitle="对他人上传的数据源进行查询"
      >
        <template #title>查询信息</template>
      </t-card>
    </div>
    <div class="choose">
      <t-card
        class="option-card"
        hover-shadow
        @click="routeTo('download')"
        subtitle="从IPFS上下载文件，若已知AES密钥可解密"
      >
        <template #title>下载文件</template>
      </t-card>
      <t-card
        class="option-card"
        hover-shadow
        @click="routeTo('log')"
        subtitle="查看查询操作的审计日志"
      >
        <template #title>审计日志</template>
      </t-card>
    </div>
    <t-guide v-model="currentGuide" :steps="GuideSteps" :show-overlay="false" />
  </div>
</template>


<script setup>
import { ref } from "vue";
import { useUserStore } from "../stores/user";
import router from "../router";
import { MessagePlugin } from "tdesign-vue-next";

const userStore = useUserStore();
const uId = ref("");

const saveUserId = () => {
  userStore.setUserId(uId.value);
};

const routeTo = (name_) => {
  // 判断是否设置了用户ID
  if (name_ != "download" && name_ != "log" && !userStore.userId) {
    MessagePlugin.error("此操作需先设置用户ID");
    currentGuide.value = 0;
    return;
  }
  router.push({ name: name_ });
};

// guide
const currentGuide = ref(-1);
const GuideSteps = [
  {
    element: ".input-user-id",
    title: "设置用户ID",
    body: "查询操作和上传操作均需先设置用户ID",
    placement: "bottom-right",
  },
];
</script>

<style scoped>
.home {
  margin: 20px auto;
  width: 30%;
}

.input-container {
  display: flex;
  justify-content: center;
  margin: 20px auto;
  width: 100%;
}

.choose {
  display: flex;
  justify-content: flex-start;
  align-items: flex-start;
  text-align: left;
  margin-top: 20px;
  gap: 20px;
}

.option-card {
  width: 200px;
  cursor: pointer;
  text-align: center;
}
</style>