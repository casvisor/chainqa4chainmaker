<script setup>
import { RouterLink, RouterView } from "vue-router";
import MyTitle from "./components/MyTitle.vue";
</script>

<template>
  <div>
    <MyTitle />
    <RouterView />
  </div>
</template>

<style scoped>
</style>
