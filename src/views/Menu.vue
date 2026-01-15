<template>
    <div>
        <el-menu
          mode="horizontal"
          :ellipsis="false"
          :router="true"
          class="header-menu"
          >
            <el-menu-item class="menu-item" index="/Home"><el-icon class="el-icon--left menu-icon"><House /></el-icon><span class="menu-text">页面介绍</span></el-menu-item>
            <el-menu-item class="menu-item" index="/Search"><el-icon class="el-icon--left menu-icon"><Search /></el-icon><span class="menu-text">查找课程</span></el-menu-item>
            <el-menu-item class="menu-item" index="/Contact"><el-icon class="el-icon--left menu-icon"><Link /></el-icon><span class="menu-text">友情链接</span></el-menu-item>
        </el-menu>
    </div>
</template>
