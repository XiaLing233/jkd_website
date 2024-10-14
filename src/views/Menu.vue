<template>
    <ul id="menu">
        <li><router-link to="/Home">页面介绍</router-link></li>
        <li><router-link to="/Search">查找课程</router-link></li>
        <li><router-link to="/Contact">联系方式</router-link></li>
    </ul>
</template>

<style>
    #menu li {
        /* padding-left: 20px; */
        /* padding-right: 20px; */
        padding-top: 10px;
        padding-bottom: 10px;
        justify-content: center;
        align-items: center;
        list-style: none;
        margin: 0;
    }

    #menu li a:hover {
        background-color:dodgerblue;
        display: inline;
        align-items: center;
    }

    #menu li a {
        text-decoration: none;
        color: white;
        font-weight: bold;
        padding-left: 20px;
        padding-right: 20px;
        padding-top: 10px;
        padding-bottom: 10px;
    }

    #menu {
        list-style-type: none;
        padding: 0;
        background-color: cornflowerblue;
        max-width: 1000px;
        font-size: 20px;
        display: flex;
    }
</style>