<template>
    <el-header style="background-color: white; display: flex; align-items: center; justify-content: space-between">
        <div id="title-left" style="padding: 12px, 0; display: inline-flex;">
            <a href="https://1.tongji.edu.cn"><img height=45px width=167px src='../assets/tongji.webp'></img></a>
            <a href="/"><img height="45px" width="auto" src="../assets/title.webp"></a>            </div>
        <div id="title-right" class="title-right">
            <Menu />
        </div>
    </el-header>
</template>

<script>
    import Menu from './Menu.vue';

    export default {
        components: {
            Menu,
        }
    }

</script>

<style scoped>
.header-container {
    background-color: white;
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
    padding: 0 20px;
    box-sizing: border-box;
}

.title-left {
    padding: 12px 0;
    display: flex;
    align-items: center;
    gap: 10px;
    min-width: 0;
}

.title-right {
    padding: 12px 0;
    flex-shrink: 0;
}

.logo-link {
    flex-shrink: 0;
}

.logo-title {
    min-width: 0;
}

.logo-img {
    height: 45px;
    width: auto;
    max-width: 150px;
}

.title-img {
    height: 45px;
    width: auto;
}
</style>
