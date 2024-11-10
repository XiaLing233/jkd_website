<template>
        <Header />
        
        <el-main style="padding: 0px; margin: 0px">
            <router-view />
        </el-main>
        
        <Footer />
</template>

<style>
    * {
        font-family: Inter, 'Helvetica Neue', Helvetica, 'PingFang SC',
  'Hiragino Sans GB', 'Microsoft YaHei', '微软雅黑', Arial, sans-serif;
    }

</style>

<script>
    import Footer from './views/Footer.vue';
    import Header
     from './views/Header.vue';
    export default {
        components: {
            // Menu,
            Footer,
            Header,
        },
    };
</script>
