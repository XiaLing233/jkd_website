<template>
    <div id="contactbody">
        <h1>友情链接</h1>
        <hr>
        <p>
        如果有改进的建议，或者发现了 bug，欢迎联系我。
        </p>

        <ul class="blank_link">
                <li><a id="mail"  href="mailto:wangxilin@tongji.edu.cn">通讯地址：wangxilin@tongji.edu.cn</a></li>
                <li><a id="link"  href="https://github.com/XiaLing233">Github</a></li>
        </ul>

        <p>
            资料来源：
        </p>

        <ul class="blank_link">
                    <li><a href="https://1.tongji.edu.cn">1系统</a></li>
                </ul>

        <p>
            采用的技术：
        </p>

        <ul class="blank_link">
                    <li><a href="https://vuejs.org/">Vue.js</a></li>
                    <li><a href="https://flask.palletsprojects.com/">Flask</a></li>
                    <li><a href="https://www.python.org/">Python</a></li>
                    <li><a href="https://www.mysql.com/">MySQL</a></li>
                    <li><a href="https://nginx.org/">Nginx</a></li>
                    <li><a href="https://gunicorn.org/">gunicorn</a></li>
                    <li><a href="https://aws.amazon.com/">Aws</a></li>
        </ul>
        

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

    #mail {
        text-decoration: none;
        color: black;
    }

    #contactbody {
        font-family: 'Times New Roman', "新宋体";
        /* margin: 20px; */
        /* max-width: 1000px; */
        min-width: 800px;
    }

    #link {
        text-decoration: none;
        color: black;
    }

    #contactbody ul {
        list-style-type: none;
        padding: 0;
        /* background-color: cornflowerblue; */
        max-width: 1000px;
        /* font-size: 20px; */
        display: block;

    }

    #contactbody ul li {
        padding-top: 10px;
        padding-bottom: 15px;
    }

    #contactbody ul li a:hover {
        background-color: #e0e0e0; /* 悬停时背景颜色变深 */
    }

    #contactbody ul li a {
        text-decoration: none;
        background-color: #f0f0f0;
        color: black;
        /* color: white; */
        /* font-weight: bold; */
        padding-left: 20px;
        padding-right: 20px;
        padding-top: 10px;
        padding-bottom: 10px;
    }

</style>