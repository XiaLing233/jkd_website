<template>
    <div>
        <el-menu
          mode="horizontal"
          :ellipsis="false"
          :router="true"
          >
            <el-menu-item :style="`font-size: var(--el-font-size-medium)`" index="/Home"><el-icon class="el-icon--left"><House /></el-icon>页面介绍</el-menu-item>
            <el-menu-item :style="`font-size: var(--el-font-size-medium)`" index="/Search"><el-icon class="el-icon--left"><Search /></el-icon>查找课程</el-menu-item>
            <el-menu-item :style="`font-size: var(--el-font-size-medium)`" index="/Contact"><el-icon class="el-icon--left"><Link /></el-icon>友情链接</el-menu-item>
        </el-menu>
    </div>
</template>

<style scoped>

</style>