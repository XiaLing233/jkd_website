<template>
    <div id="app">
        <Header />
        <br>
        <div id="main">
            <router-view />
        </div>
        <br>
        <Footer />
    </div>

</template>

<style scoped>
    * {
        font-family: 'Times New Roman', "宋体", serif;
    }
    #app {
        background-color: whitesmoke;
    }
    
    #main {
        padding-top: 10px;
        padding-bottom: 10px;
        /* margin-left: 80px; */
        /* margin-right: 80px; */
        /* max-width: 1000px; */
        margin: 0 auto;
        background-color:white;
        min-height: 800px;
        justify-content: center;
        /* align-items: center; */
        display: flex;
    }

</style>

<script>
    import Footer from './views/Footer.vue';
    import Header
     from './views/Header.vue';
    export default {
        components: {
            // Menu,
            Footer,
            Header,
        },
    };
</script>
