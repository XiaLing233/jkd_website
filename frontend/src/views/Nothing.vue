<template>
    <el-card style="align-items: center; display:flex; flex-direction: column; margin: 20px" shadow="never">
        <img src="../assets/notfound.webp" alt="notfound"/>
        <div>
        <h3>出错啦！该资源不存在！</h3>
        <p id="reason">
        </p>
        <div style="font-size: .8em;margin-top: 12px;">(错误代码：404)</div>
    </div>
</el-card>
</template>

<script>
export default {
    name: 'Nothing'
}
</script>
