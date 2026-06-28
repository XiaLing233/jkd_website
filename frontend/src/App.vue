<template>
  <div class="flex flex-col h-screen overflow-hidden">
    <div class="shrink-0 z-1000">
      <Header />
    </div>
    <div class="flex-1 overflow-y-auto overflow-x-hidden">
      <router-view />
    </div>
  </div>
</template>

<style>
* {
  font-family: Inter, 'Helvetica Neue', Helvetica, 'PingFang SC',
    'Hiragino Sans GB', 'Microsoft YaHei', '微软雅黑', Arial, sans-serif;
}
</style>

<script setup lang="ts">
import Header from '@/components/Header.vue'
</script>
