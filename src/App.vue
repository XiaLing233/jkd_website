<template>
    <div id="app">
        <header>
            <div id="title">
                <a href="https://1.tongji.edu.cn"><img height=45px width=167px src="./assets/tongji.png"></img></a>
                <h1>课程检索系统</h1>
            </div>

            <div id="menuholder">
                <ul id="menu">
                    <li><router-link to="/Home">页面介绍</router-link></li>
                    <li><router-link to="/Search">查找课程</router-link></li>
                    <li><router-link to="/Contact">联系方式</router-link></li>
            </ul>
            </div>

        </header>

        <br>
        <main>
            <router-view />
        </main>

        <footer>
            <p>{{ new Date().getFullYear() }} — <strong>夏凌</strong></p>
            <p><strong>计算机科学导论</strong>网站大作业</p>
            <p>更新时间：2024/10/13</p>
        </footer>
        </div>

</template>

<style>
    * {
        font-family: 'Times New Roman', "宋体", serif;
    }
    #app {
                background-color: whitesmoke;
    }

    #menu {
        list-style-type: none;
        margin: 0;
        padding: 0;
        background-color: cornflowerblue;
        max-width: 1000px;
        /* margin: 0 auto; */
        font-size: 20px;
        display: flex;
        align-items: center;
        /* height: 100%; */
        /* margin: 0; */
    }

    .classitems {
        display: inline;
        text-align: center;
        padding-left: 20px;
        padding-right: 20px;
    }

    footer {
        background-color: whitesmoke;
        text-align: center;
        position: relative;
        bottom: 0;
        margin-top: 40px;
        margin-bottom: 10px;
        font-size: 14px;
    }
    footer p {
        margin: 0;
        padding-bottom: 5px;
    }

    header {
        background-color: white;
        justify-content: space-between;
        align-items: center;
        /* padding: 10px 20px; */
        display: flex;
    }

    #title {
        /* background-color: white; */
        text-align: center;
        display: flex; 
        align-items: center;
    }

    #title h1 {
        font-family: "楷体";
        /* font-weight: 500; */
    }
    
    #menu li {
        /* display: inline; */
        /* text-align: center; */
        /* padding-top: 30px; */
        padding-left: 20px;
        padding-right: 20px;
        height: 100%;
        justify-content: center;
        align-items: center;
        /* margin: 0; */
        /* margin-top: 0; */
        /* margin-bottom: 0; */
        /* padding: 0; */
    }

    #menu li:hover {
        background-color:dodgerblue;
        height: 100%;
        text-align: center;
        display: inline;
        /* align-items: center; */
    }

    a {
        text-decoration: none;
        color: white;
        font-weight: bold;
    }

    main {
        /* background-image: url("./assets/body-bg.png");
        background-size: cover; */
        /* border: 1px solid black; */
        padding-top: 10px;
        padding-bottom: 10px;
        margin-left: 80px;
        margin-right: 80px;
        max-width: 1000px;
        margin: 0 auto;
        background-color:white;
    }

</style>