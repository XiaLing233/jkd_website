<template>
  <el-card id="introbody" shadow="never" class="m-5">
    <div class="my-6">
      <h1 class="text-center font-bold text-3xl">页面介绍</h1>
    </div>
    <el-divider />

    <h2 class="text-2xl font-bold hover:underline mb-2">功能</h2>
    <p class="mb-3">根据提供的字段进行课程查找，是 1 系统的增强版。</p>

    <h2 class="text-2xl font-bold hover:underline mt-6 mb-2">缘起</h2>
    <p class="mb-3">因为 1 系统的查询系统灵活度太低，无法满足需求，因此自己编写一个查询网站。</p>

    <h2 class="text-2xl font-bold hover:underline mt-6 mb-2">涉及的主要技术</h2>
    <table class="max-w-100 w-full border-collapse my-2">
      <thead>
        <tr>
          <th class="border border-black text-center">前端</th>
          <th class="border border-black text-center">后端</th>
          <th class="border border-black text-center">数据库</th>
          <th class="border border-black text-center">部署</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td class="border border-black text-center">Vue 3 + TS</td>
          <td class="border border-black text-center">Go + Gin</td>
          <td class="border border-black text-center">MariaDB</td>
          <td class="border border-black text-center">Docker</td>
        </tr>
      </tbody>
    </table>

    <h3 class="text-lg font-bold mt-6 hover:underline mb-2">前端</h3>
    <p class="mb-3">使用 <code class="bg-gray-200 px-1 py-0.5 rounded font-mono text-sm">Vue 3</code> + <code class="bg-gray-200 px-1 py-0.5 rounded font-mono text-sm">TypeScript</code> + <code class="bg-gray-200 px-1 py-0.5 rounded font-mono text-sm">Vite</code>，配合 <code class="bg-gray-200 px-1 py-0.5 rounded font-mono text-sm">Element Plus</code> 和 <code class="bg-gray-200 px-1 py-0.5 rounded font-mono text-sm">Tailwind CSS</code>。</p>

    <h3 class="text-lg font-bold mt-6 hover:underline mb-2">后端</h3>
    <p class="mb-3">采用 <code class="bg-gray-200 px-1 py-0.5 rounded font-mono text-sm">Go</code> + <code class="bg-gray-200 px-1 py-0.5 rounded font-mono text-sm">Gin</code>，跨库聚合多学期数据。</p>

    <h3 class="text-lg font-bold mt-6 hover:underline mb-2">数据库</h3>
    <p class="mb-3">采用 <code class="bg-gray-200 px-1 py-0.5 rounded font-mono text-sm">MariaDB</code>，和 <a class="text-blue-500 hover:underline" href="https://xk.xialing.icu" target="_blank">xk.xialing.icu</a> 共用同一个数据库。爬虫的具体实现请参见<a class="text-blue-500 hover:underline" href="https://github.com/XiaLing233/tongji-course-scheduler" target="_blank">仓库</a>。</p>
  </el-card>
</template>

<script setup lang="ts">
</script>
