<template>
<div id="introbody">
    <div id="heading">
        <h1 id='页面介绍'>页面介绍</h1>
        </div>
    <div id="separator">
        <hr>
    </div>
    <h2 id='功能'>功能</h2>
    <p>根据提供的字段进行课程查找，是 1 系统的增强版。</p>
    <h2 id='缘起'>缘起</h2>
    <p>因为 1 系统的查询系统灵活度太低，无法满足需求，因此自己编写一个查询网站。</p>
    <h2 id='涉及的技术'>涉及的技术</h2>
    <figure class='table-figure'><table>
    <thead>
    <tr>
        <th>前端</th>
        <th>后端</th>
        <th>数据库</th>
    </tr>
</thead>
    <tbody>
        <tr>
            <td>vue</td>
            <td>flask</td>
            <td>MySQL</td>
        </tr>
    </tbody>
    </table></figure>
    <h3 id='前端'>前端</h3>
    <p>使用 <code>vue</code> 框架，刚上手有些难度。</p>
    <p>主要采用的技术有：</p>
    <ul class="intro">
    <li><code>vue-routers</code> 进行页面跳转</li>
    <li><code>CSS</code> 进行页面美化</li>
    <li>一些方法，用来实现前端的内容筛选</li>
    <li>根据表格的不同字段，向后端发送数据请求</li>

    </ul>
    <h3 id='后端'>后端</h3>
    <p>采用 <code>flask</code></p>
    <p>主要的功能是：</p>
    <ul class="intro">
    <li>接收前端的请求，向数据库发送 <code>SQL</code> 语句</li>
    <li>返回查找好的列表内容给前端</li>

    </ul>
    <h3 id='数据库'>数据库</h3>
    <p>采用了 <code>MySQL</code>，共一个数据库，储存了一些表格，每个表格是一学期的所有课程。</p>
    <ul class="intro">
    <li>对于从 1 系统导出的 <code>excel</code> 表格进行处理，需要先导出为 <code>.csv</code> 逗号分隔文件，再通过 <code>pandas</code> 库进行数据的优化整理。</li>

    </ul>
    <h2 id='可以改进的地方'>可以改进的地方</h2>
    <ul class="intro">
    <li>使用更好看的图标</li>
    <li>根据返回的有效内容，生成一个课程时间表</li>
    </ul>
</div>
</template>

<script>
export default {
};
</script>

<style> 
    #introbody {
        font-family: 'Times New Roman', "新宋体";
        margin: 20px;
    }

    #introbody h2:hover, h3:hover {
        text-decoration: underline;
    }

    #introbody h1 {
        text-align: center;
        font-weight: bold;
    }

    #introbody h2 {
        font-size: 20px;
    }
    #introbody h3 {
        font-size: 18px;
    }
    .table-figure {
        margin: 0;
        padding: 0;
    }
    table {
        border-collapse: collapse;
        width: 100%;
    }
    th, td {
        border: 1px solid black;
        text-align: center;
    }
    li {
        margin: 10px;   
    }
    code {
        background-color: rgb(236, 236, 236);
        padding: 2px;
        font-family: consolas;
    }
</style>
