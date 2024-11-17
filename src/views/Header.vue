<template>
        <el-affix><el-header style="background-color: white; display: flex; align-items: center; justify-content: space-between">
            <div id="title-left" style="padding: 12px, 0; display: inline-flex;">
                <a href="https://1.tongji.edu.cn"><img height=45px width=167px src='../assets/tongji.webp'></img></a>
                <a href="/"><img height="45px" width="auto" src="../assets/title.webp"></a>
            </div>

            <div id="title-right" style="padding: 12px, 0;">
                <Menu />
            </div>

        </el-header></el-affix>
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

</style>