<template>
    <div id="contactbody" style="margin: 20px">
        <el-card shadow="never">
            <h1>友情链接</h1>
        <el-divider />
        <el-card shadow="never" class="link_set" style="margin-bottom: 20px;">
            <template #header>
                <div>如果有改进的建议，或者发现了 bug，欢迎联系我。</div>
            </template>

            <ul>
                <li><el-link id="mail" type="primary" href="mailto:wangxilin@tongji.edu.cn"><el-icon class="el-icon--left"><Message /></el-icon>通讯地址：wangxilin@tongji.edu.cn</el-link></li>
                <li><el-link id="link" type="primary" href="https://github.com/XiaLing233">Github</el-link></li>
            </ul>


        </el-card>

        <el-card shadow="never" class="link_set" style="margin-bottom: 20px;">
            <template #header>
                <div>资料来源：</div>
            </template>
            <ul>
                <li>
                    <el-link type="primary" href="https://1.tongji.edu.cn">1系统</el-link>
                </li>
            </ul>

        </el-card>

        <el-card shadow="never" class="link_set">
            <template #header>
                <div> 采用的技术：</div>
            </template>
            <ul>
                <li><el-link type="primary" href="https://vuejs.org/">Vue.js</el-link></li>
                <li><el-link type="primary" href="https://flask.palletsprojects.com/">Flask</el-link></li>
                <li><el-link type="primary" href="https://www.python.org/">Python</el-link></li>
                <li><el-link type="primary" href="https://www.mysql.com/">MySQL</el-link></li>
                <li><el-link type="primary" href="https://nginx.org/">Nginx</el-link></li>
                <li><el-link type="primary" href="https://gunicorn.org/">gunicorn</el-link></li>
                <li><el-link type="primary" href="https://aws.amazon.com/">Aws</el-link></li>
                <li><el-link type="primary" href="https://element-plus.org/">Element Plus</el-link></li>
            </ul>
        </el-card>
        </el-card>
    </div>
</template>

<script>
export default {
    mounted() {
        // 把所有 blank_link 类的 ul 下的 a 标签的 target 属性设置为 _blank
        const links = document.querySelectorAll('.blank_link a');
        for (let i = 0; i < links.length; i++) {
            links[i].setAttribute('target', '_blank');
        }
    }
};
</script>

<style scoped> 
    h1 {
        text-align: center;
        font-weight: bold;
    }

    li {
        list-style-type: none;
        padding: 2px;
    }
</style>