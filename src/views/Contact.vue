<template>
    <div id="contactbody">
        <h1>联系方式</h1>
        <hr>
        <p>
        如果有改进的建议，或者发现了 bug，欢迎联系我。
        </p>
        <div>
            <ul>
                <li><a id="mail"  href="mailto:wangxilin@tongji.edu.cn">通讯地址：wangxilin@tongji.edu.cn</a></li>
                <li><a id="link"  href="https://github.com/XiaLing233">Github</a></li>
            </ul>
            
        </div>
        

    </div>
</template>

<script>
export default {
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
</style>