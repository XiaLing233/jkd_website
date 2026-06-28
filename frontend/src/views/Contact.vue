<template>
  <div id="contactbody" style="margin: 20px">
    <el-card shadow="never">
      <div class="my-6">
        <h1 class="text-center font-bold text-3xl">友情链接</h1>
      </div>
      <el-divider />

      <el-card shadow="never" class="link_set" style="margin-bottom: 20px">
        <template #header>联系方式</template>
        <ul class="list-none">
          <li class="py-0.5">
            <el-link type="primary" href="mailto:wangxilin@tongji.edu.cn">
              <el-icon class="el-icon--left"><Message /></el-icon>wangxilin@tongji.edu.cn
            </el-link>
          </li>
          <li>
            <el-link type="primary" href="https://github.com/XiaLing233" target="_blank">
              GitHub
            </el-link>
          </li>
        </ul>
      </el-card>

      <el-card shadow="never" class="link_set" style="margin-bottom: 20px">
        <template #header>数据来源</template>
        <ul class="list-none">
          <li class="py-0.5">
            <el-link type="primary" href="https://1.tongji.edu.cn" target="_blank">
              同济大学 1 系统
            </el-link>
          </li>
        </ul>
      </el-card>

      <el-card shadow="never" class="link_set">
        <template #header>采用的技术：</template>
        <ul class="list-none">
          <li class="py-0.5"><el-link type="primary" href="https://vuejs.org/" target="_blank">Vue 3 + TypeScript + Vite</el-link></li>
          <li><el-link type="primary" href="https://go.dev/" target="_blank">Go + Gin</el-link></li>
          <li><el-link type="primary" href="https://mariadb.org/" target="_blank">MariaDB</el-link></li>
          <li><el-link type="primary" href="https://nginx.org/" target="_blank">Nginx</el-link></li>
          <li><el-link type="primary" href="https://www.docker.com/" target="_blank">Docker Compose</el-link></li>
          <li><el-link type="primary" href="https://element-plus.org/" target="_blank">Element Plus</el-link></li>
        </ul>
      </el-card>
    </el-card>
  </div>
</template>

<script setup lang="ts">
</script>

