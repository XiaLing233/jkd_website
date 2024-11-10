<template>
    <div class="welcome-container">
        <div id = "welcome">
        <el-card shadow="never" id="wel-box">
            <h1>欢迎使用同济大学课程检索系统！</h1>
            <p>请点击右上方导航栏进行操作</p>
        </el-card>
    </div>
    </div>
</template>

<style scoped>
    #welcome h1 {
        text-align: center;
    }

    #welcome p {
        text-align: center;
        font-size: 25px;
    }

    .welcome-container {
        background-image: url('../assets/welcome_bg_1.jpg');
        /* background-image: url('../assets/welcome_bg.png'); */
        background-size: cover;
        background-position: center;
        height: calc(100vh - 130px);
        width: 100%;
    }

    #wel-box {
        float: right; 
        margin-top: 7%;
        margin-right: 7%;
        background-color: rgba(0, 0, 0, 0.4); 
        backdrop-filter: blur(20px);
        border-radius: 4px;
        box-shadow: 0 4px 36px 0 #0d5ead1a;
        padding: 24px;
        position:relative;
        color: white;
        border: 0px;
    }
</style>