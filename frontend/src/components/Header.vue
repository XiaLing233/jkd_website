<template>
  <el-header class="bg-white flex items-center justify-between px-5">
    <div class="py-3 flex items-center gap-4 min-w-0">
      <a href="https://1.tongji.edu.cn" class="shrink-0">
        <img height="45" width="167" src="@/assets/tongji.webp" alt="tongji" />
      </a>
      <a href="/" class="shrink-0">
        <img height="45" src="@/assets/title.webp" alt="title" />
      </a>
      <div class="px-3 py-1.5 bg-gray-100 border border-gray-200 rounded text-base text-black whitespace-nowrap">
        更新时间：{{ lastUpdateTime }}
      </div>
    </div>
    <div class="py-3 shrink-0">
      <Menu />
    </div>
  </el-header>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { fetchLastUpdate } from '@/api/client'
import Menu from './Menu.vue'

const lastUpdateTime = ref('加载中...')

onMounted(async () => {
  try {
    const resp = await fetchLastUpdate()
    lastUpdateTime.value = resp.data?.fetchTime ?? '未知'
  } catch {
    lastUpdateTime.value = '获取失败'
  }
})
</script>
