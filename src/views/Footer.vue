<template>
        <el-footer
        style="text-align: center; margin-bottom: 5px; margin-top: 5px"
        >
            <p>{{ new Date().getFullYear() }} — <strong>夏凌</strong></p>
            <p><strong>计算机科学导论</strong>网站大作业</p>
            <p>更新时间：{{ lastUpdateTime }}</p>
        </el-footer>
</template>

<script>
export default {
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
};
</script>

<style scoped>
    p {
        margin: 0;
    }
</style>