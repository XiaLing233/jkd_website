<template>
    
    <!-- 调试用 -->
    <!-- <div>
      <h2>Available Items:</h2>
      <ul>
        <li v-for="item in allItems" :key="item">{{ item }}</li>
      </ul>
    </div> -->

    <!-- <div>
      <h2>Available Items:</h2>
      <ul>
        <li v-for="item in allItems" :key="item">{{ item }}</li>
      </ul>
    </div> -->

    <div id="searchbody"> <!-- 搜索主体 -->
    <h1>课程检索</h1>
    <hr>

    <!-- 搜索条件 -->
    <div class="search-conditions-container">

      <!-- 搜索条件 -->
      <div v-for="(condition, index) in conditions" :key="index" class="search-condition">
        <!-- 显示连接符的下拉菜单（如果不是第一个条件） -->
        <div v-if="index !== 0">
          <span>布尔连接符:</span>  <!-- 连接符 -->
          <select v-model="condition.connector">
            <option value="AND">AND</option>
            <option value="OR">OR</option>
            <option value="NOT">NOT</option>
          </select>
        </div>

        <!-- 选择字段的下拉菜单 -->
        <select v-model="condition.selectedItem" @change="fetchFieldOptions(index)"> <!-- 当改变字段时，调用 fetchFieldOptions 方法 -->
          <option value="" disabled>请选择</option> <!-- 默认选项 -->
          <option v-for="item in availableItemsFor(index)" :key="item" :value="item">{{ item }}</option>  <!-- 可选字段 -->
        </select>

        <!-- 输入搜索词 -->
        <input
          type="text"
          v-model="condition.searchWord"
          placeholder="请输入搜索词"
          :list="'keywords-' + index"
        />
        <button @click="condition.searchWord = ''" style="margin-left: 5px;" class="search-button">清空关键词</button>
        <datalist :id="'keywords-' + index"> <!-- 根据条件的索引生成不同的 datalist -->
          <option v-for="option in condition.options" :key="option" :value="option">{{ option }}</option>
        </datalist>

        <button @click="removeCondition(index)" class="search-button">移除条件</button> <!-- 移除条件按钮 -->
      </div>
    </div>
    <div id="separator">
      <button @click="addCondition" :disabled="!canAddCondition" class="search-button">添加条件</button> <!-- 添加条件按钮 -->
      <button @click="submitSearch" class="search-button">搜索</button>   <!-- 搜索按钮 -->
    </div>

    <br>
    <br>

    
    <!-- 显示搜索结果 -->
    <!-- 数据库中表格的结构是：
     课程序号、课程名称、周学时、授课教师、教师工号、课程性质、授课语言、校区、听课专业、选课人数、起始周、结束周、排课信息、开课学院 -->
     <div v-if="searchResults == 'default'">
      <h1>请检索</h1>
    </div>

      <div v-else-if="searchResults.length > 0">
        <!-- 列选择复选框 -->
      <div class="tabs">
        <label><input type="checkbox" v-model="showCourseNumber" /> 课程序号</label>
        <label><input type="checkbox" v-model="showCourseName" /> 课程名称</label>
        <label><input type="checkbox" v-model="showTeacher" /> 授课教师（工号）</label>
        <label><input type="checkbox" v-model="showCourseType" /> 课程性质</label>
        <label><input type="checkbox" v-model="showCampus" /> 校区</label>
        <label><input type="checkbox" v-model="showMajor" /> 听课专业</label>
        <label><input type="checkbox" v-model="showSchedule" /> 排课信息</label>
        <label><input type="checkbox" v-model="showDepartment" /> 开课学院</label>
      </div>

      <table>
        <thead>
          <tr>
            <th v-if="showCourseNumber" style="width: 5%;">课程序号</th>
            <th v-if="showCourseName" style="width: 20%;">课程名称</th>
            <th v-if="showTeacher" style="width: 15%;">授课教师（工号）</th>
            <th v-if="showCourseType" style="width: 10%;">课程性质</th>
            <th v-if="showCampus" style="width: 7%;">校区</th>
            <th v-if="showMajor" style="width: 15%;">听课专业</th>
            <th v-if="showSchedule" style="width: 20%;">排课信息</th>
            <th v-if="showDepartment" style="width: 13%;">开课学院</th>

          </tr>
        </thead>
        <tbody>
          <tr v-for="(result, index) in paginatedResults" :key="index">
            <td v-if="showCourseNumber">{{ result[0] }}</td> <!-- 课程序号 -->
            <td v-if="showCourseName">{{ result[1] }}</td> <!-- 课程名称 -->
            <td v-if="showTeacher">{{ formatTeachers(result[3], result[4]) }}</td> <!-- 合并教师名称和工号 -->
            <td v-if="showCourseType">{{ result[5] }}</td> <!-- 课程性质 -->
            <td v-if="showCampus">{{ result[7] }}</td> <!-- 校区 -->
            <td v-if="showMajor">{{ result[8] }}</td> <!-- 听课专业 -->
            <td v-if="showSchedule">{{ result[12] }}</td> <!-- 排课信息 -->
            <td v-if="showDepartment">{{ result[13] }}</td> <!-- 开课学院 -->
          </tr>
        </tbody>
      </table>
      <!-- 分页控制 -->
      <div class="pagination-container">
        <div class="page-size-selector">
          <label for="pageSize">每页显示:</label>
          <select v-model="pageSize" id="pageSize">
            <option v-for="size in [10, 20, 50, 100, 1000]" :key="size" :value="size">{{ size }}</option>
          </select>
          行
        </div>
        
        <div class="pagination">
          <button @click="currentPage = 1" :disabled="currentPage === 1">首页</button>
          <button @click="currentPage--" :disabled="currentPage === 1">上一页</button>
          <span>第 {{ currentPage }} 页 / 共 {{ totalPages }} 页</span>
          <button @click="currentPage++" :disabled="currentPage === totalPages">下一页</button>
          <button @click="currentPage = totalPages" :disabled="currentPage === totalPages">末页</button>
        </div>
      </div>
    </div>

    <div v-else>
      <h1>{{searchResults}}</h1>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      conditions: [{ selectedItem: '', searchWord: '', connector: '' }], // 条件数组
      allItems: [], // 所有可检索的字段
      searchResults: ['default'], // 用于存储搜索结果
      showCourseNumber: true, // 是否显示课程序号
      showCourseName: true, // 是否显示课程名称
      showTeacher: true, // 是否显示授课教师
      showCourseType: true, // 是否显示课程性质
      showCampus: true, // 是否显示校区
      showMajor: false, // 是否显示听课专业
      showSchedule: true, // 是否显示排课信息
      showDepartment: true, // 是否显示开课学院
      currentPage: 1, // 当前页码
      pageSize: 10, // 每页显示的行数
    };
  },
  computed: {
    canAddCondition() {
      // 检查条件的数量是否小于所有字段数量
      return this.conditions.length < this.allItems.length;
    },

    filteredResults() {
      return this.searchResults.map(result => {
        return {
          showCourseNumber: this.showCourseNumber ? result[0] : null,
          showCourseName: this.showCourseName ? result[1] : null,
          showTeacher: this.showTeacher ? this.formatTeachers(result[3], result[4]) : null,
          showCourseType: this.showCourseType ? result[5] : null,
          showCampus: this.showCampus ? result[7] : null,
          showMajor: this.showMajor ? result[8] : null,
          showSchedule: this.showSchedule ? result[12] : null,
          showDepartment: this.showDepartment ? result[13] : null,
        };
      });
    },

    paginatedResults() {
      const start = (this.currentPage - 1) * this.pageSize;
      const end = start + this.pageSize;
      return this.searchResults.slice(start, end);
    },

    totalPages() {
      return Math.ceil(this.searchResults.length / this.pageSize);
    },

  },
  methods: {
    // 添加条件
    addCondition() {
      this.conditions.push({ selectedItem: '', searchWord: '', connector: '' });
    },

    // 移除条件
    removeCondition(index) {
      this.conditions.splice(index, 1);
    },

    // 获取当前条件可选的字段
    availableItemsFor(index) {
      // 获取除当前条件外所有已选择的字段
      const selectedItems = this.conditions
        .filter((cond, i) => i !== index)
        .map(cond => cond.selectedItem);
      // 返回未被选择的字段
      return this.allItems.filter(item => !selectedItems.includes(item));
    },

    // 获取字段的选项
    fetchFieldOptions(index) {
  const fieldName = this.conditions[index].selectedItem; // 获取当前条件的字段名

  // 请求字段的选项数据
  fetch('http://localhost:5000/get_field_options', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ field_name: fieldName }),
      })
        .then(response => response.json()) // 解析 JSON 数据
        .then(data => {
          // 更新当前条件的可选项
          this.conditions[index].options = data; // 直接赋值，Vue 3 支持响应式
        })
        .catch(error => {
          console.error('Error fetching field options:', error);
        });
    },


    // 提交搜索请求
    submitSearch() {
      const requestData = this.conditions.filter(cond => cond.selectedItem); // 过滤掉空条件
      // 发送 HTTP 请求到 Flask 后端
      fetch('http://localhost:5000/search', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestData),
      })
        .then(response => response.json())
        .then(data => {
          // 更新搜索结果
          this.searchResults = data;
          console.log('Search results:', data);
        })
        .catch(error => {
          console.error('Error:', error);
        });
    },

    // 格式化教师信息
    formatTeachers(teachers, ids) {
      const teacherList = teachers.split(','); // 假设教师用逗号分隔
      const idList = ids.split(','); // 假设工号用逗号分隔
      return teacherList.map((teacher, index) => `${teacher.trim()} (${idList[index].trim()})`).join(', '); // 格式化
    },

  },

  mounted() {
    // 获取所有可检索的字段
    fetch('http://localhost:5000/get_searchable_fields')
      .then(response => response.json()) // 解析 JSON 数据
      .then(data => { 
        this.allItems = data; // 更新所有字段
      })
      .catch(error => {
        console.error('Error fetching searchable fields:', error);
      });
  },
};
</script>

<style scoped>
/* 页面整体样式 */
#searchbody {
  font-family: "新宋体", Arial, sans-serif; /* 字体 */
  margin: 0;
  padding: 20px; /* 页面内边距 */
  min-width: 1200px; /* 最小宽度 */
}

/* 标题样式 */
#searchbody h1 {
  text-align: center; /* 居中 */
  font-weight: bold; /* 加粗 */
}

/* 选项卡样式 */
.tabs {
  display: flex; /* 使用 Flexbox 布局 */
  gap: 15px; /* 每个选项之间的间距 */
  margin-bottom: 10px; /* 下方的间距 */
}

.tabs label {
  display: flex; /* 使选项卡内部的复选框与文本对齐 */
  align-items: center; /* 垂直居中 */
  background-color: #f0f0f0; /* 选项卡背景颜色 */
  border-radius: 5px; /* 圆角 */
  padding: 5px 10px; /* 内边距 */
  cursor: pointer; /* 鼠标悬停时变成手指 */
  transition: background-color 0.3s; /* 背景颜色渐变效果 */
}

.tabs label:hover {
  background-color: #e0e0e0; /* 悬停时背景颜色变深 */
}

/* 搜索条件样式 */
#searchbody {
  width: 100%; /* 使内容占满整个宽度 */
}

#searchbody .search-condition {
  display: flex; /* 每个条件使用 Flexbox */
  align-items: center; /* 垂直居中 */
  gap: 10px; /* 每个元素之间的间距 */
  border: 1px solid #ddd; /* 边框 */
  padding: 10px; /* 内边距 */
  border-radius: 5px; /* 圆角 */
  background-color: #f9f9f9; /* 背景颜色 */
  margin-bottom: 10px; /* 每个条件之间的垂直间距 */
  flex-wrap: nowrap; /* 不换行，保持在一行内 */
}

#searchbody .search-conditions-container {
  display: flex; /* 使用 Flexbox 布局 */
  flex-wrap: wrap; /* 允许换行 */
  gap: 15px; /* 每个条件之间的间距 */
  margin-bottom: 15px; /* 下方的间距 */
}



#searchbody .search-condition:first-child select {
  margin-right: 10px; /* 为布尔连接符留出空间 */
}

#separator {
  margin-top: 15px; /* 添加与搜索条件的间距 */
}

.search-button {
  margin-left: 10px; /* 按钮之间的间距 */
  background-color: #8a9a9d; /* 调整按钮背景色为浅灰色 */
  color: #ffffff; /* 按钮文字颜色 */
  border: none; /* 移除默认边框 */
  border-radius: 5px; /* 圆角 */
  padding: 8px 12px; /* 内边距 */
  cursor: pointer; /* 鼠标样式 */
  transition: background-color 0.3s; /* 添加过渡效果 */
}

.search-button:hover {
  background-color: #7a8a8d; /* 悬停时的背景色 */
}


/* 搜索结果表格样式 */
table {
  width: 100%; /* 表格宽度 */
  border-collapse: collapse; /* 合并边框 */
  margin-top: 20px; /* 顶部间距 */
}

th,
td {
  padding: 10px; /* 内边距 */
  border: 1px solid #dee2e6; /* 边框 */
  text-align: left; /* 左对齐 */
}

th {
  background-color: #e9ecef; /* 表头背景色 */
  font-weight: bold; /* 加粗文本 */
}

/* 数据提示信息样式 */
.no-data {
  text-align: center; /* 居中 */
  margin-top: 20px; /* 顶部间距 */
  font-size: 1.2em

}

/* 分页按钮样式 */
.pagination-container {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
  margin-top: 20px;
}

.page-size-selector {
  display: flex;
  align-items: center;
  gap: 5px;
}

.pagination {
  display: flex;
  align-items: center;
  gap: 10px;
}

.pagination button {
  padding: 5px 10px;
  background-color: #f0f0f0;
  /* border: 1px solid #ccc; */
  border: none;
  cursor: pointer;
  transition: background-color 0.3s;
  border-radius: 5px; /* 圆角 */
  font-size: 14px;
}

.pagination button:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

.pagination button:hover:not(:disabled) {
  background-color: #e0e0e0; /* 悬停时背景颜色变深 */
}

</style>
