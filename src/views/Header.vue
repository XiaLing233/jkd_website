<template>
    <el-header style="background-color: white; display: flex; align-items: center; justify-content: space-between">
        <div id="title-left" style="padding: 12px, 0; display: inline-flex; align-items: center; gap: 15px;">
            <a href="https://1.tongji.edu.cn"><img height=45px width=167px src='../assets/tongji.webp'></img></a>
            <a href="/"><img height="45px" width="auto" src="../assets/title.webp"></a>
            <div class="update-time-container">
                <span class="update-time-text">更新时间：{{ lastUpdateTime }}</span>
            </div>
        </div>
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
        },
        data() {
            return {
                lastUpdateTime: '加载中...'
            };
        },
        methods: {
            async fetchLastUpdate() {
                try {
                    const response = await fetch('/api/get_last_update');
                    const data = await response.json();
                    
                    if (data.status === 'OK' && data.data.fetchTime) {
                        this.lastUpdateTime = data.data.fetchTime;
                    } else {
                        this.lastUpdateTime = '未知';
                    }
                } catch (error) {
                    console.error('获取更新时间失败:', error);
                    this.lastUpdateTime = '获取失败';
                }
            }
        },
        mounted() {
            this.fetchLastUpdate();
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

.update-time-container {
    padding: 6px 12px;
    background-color: #f8f8f8;
    border: 1px solid rgba(60, 60, 60, 0.12);
    border-radius: 4px;
}

.update-time-text {
    font-size: 16px;
    color: black;
    white-space: nowrap;
}
</style>
