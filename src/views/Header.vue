<template>
        <div id="header">
            <div id="title">
                <a href="https://1.tongji.edu.cn"><img height=45px width=167px src='../assets/tongji.png'></img></a>
                <a href="/"><h1>课程检索系统</h1></a>
            </div>
        <Menu />
        </div>
</template>

<script>
    import Menu from './Menu.vue';

    export default {
        components: {
            Menu,
        }
    }

</script>

<style scoped>
    #header {
        background-color: white;
        justify-content: space-between;
        align-items: center;
        /* padding: 10px 20px; */
        display: flex;
    }

    #title a {
        text-decoration: none;
        color: black;
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
</style>