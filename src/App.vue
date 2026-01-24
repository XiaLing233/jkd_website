<template>
    <div class="app-layout">
        <div class="app-header">
            <Header />
        </div>
        
        <div class="app-main">
            <router-view />
        </div>
    </div>
</template>

<style>
    * {
        font-family: Inter, 'Helvetica Neue', Helvetica, 'PingFang SC',
  'Hiragino Sans GB', 'Microsoft YaHei', '微软雅黑', Arial, sans-serif;
    }

    html, body {
        margin: 0;
        padding: 0;
        height: 100%;
        overflow: hidden;
    }

    #app {
        height: 100%;
    }

    .app-layout {
        display: flex;
        flex-direction: column;
        height: 100vh;
        overflow: hidden;
    }

    .app-header {
        flex-shrink: 0;
        z-index: 1000;
    }

    .app-main {
        flex: 1;
        overflow-y: auto;
        overflow-x: hidden;
    }
</style>

<script>
    import Header from './views/Header.vue';
    export default {
        components: {
            Header,
        },
    };
</script>
