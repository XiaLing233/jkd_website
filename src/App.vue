<template>
    <div id="app" style="width: 100%; height: 100%">
        <Header />
        
        <el-main>
            <router-view />
        </el-main>
        
        <Footer />
    </div>

</template>

<style scoped>
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
