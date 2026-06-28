<template>
  <el-card id="introbody" shadow="never" class="m-5">
    <h1 class="text-center font-bold">页面介绍</h1>
    <el-divider />

    <h2 class="text-xl hover:underline">功能</h2>
    <p>根据提供的字段进行课程查找，是 1 系统的增强版。</p>

    <h2 class="text-xl hover:underline mt-4">缘起</h2>
    <p>因为 1 系统的查询系统灵活度太低，无法满足需求，因此自己编写一个查询网站。</p>

    <h2 class="text-xl hover:underline mt-4">涉及的主要技术</h2>
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

    <h3 class="text-lg mt-4 hover:underline">前端</h3>
    <p>使用 <code class="bg-gray-200 px-0.5 font-mono">Vue 3</code> + <code class="bg-gray-200 px-0.5 font-mono">TypeScript</code> + <code class="bg-gray-200 px-0.5 font-mono">Vite</code>，配合 <code class="bg-gray-200 px-0.5 font-mono">Element Plus</code> 和 <code class="bg-gray-200 px-0.5 font-mono">Tailwind CSS</code>。</p>

    <h3 class="text-lg mt-4 hover:underline">后端</h3>
    <p>采用 <code class="bg-gray-200 px-0.5 font-mono">Go</code> + <code class="bg-gray-200 px-0.5 font-mono">Gin</code>，跨库聚合多学期数据。</p>

    <h3 class="text-lg mt-4 hover:underline">数据库</h3>
    <p>采用 <code class="bg-gray-200 px-0.5 font-mono">MariaDB</code>，和 <a href="https://xk.xialing.icu" target="_blank">xk.xialing.icu</a> 共用同一个数据库。爬虫的具体实现请参见<a href="https://github.com/XiaLing233/tongji-course-scheduler" target="_blank">仓库</a>。</p>
  </el-card>
</template>

<script setup lang="ts">
</script>
