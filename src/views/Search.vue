<template>
    <div id="searchbody"> <!-- 搜索主体 -->
    <h1>课程检索</h1>
    <hr>
    <!-- 搜索条件 -->
    <div class="search-conditions-container">

      <!-- 选择学期 -->
      <div class="search-terms">
      <!-- 全选/取消全选 checkbox -->

      <!-- 下拉式选择器 -->
      <div>
      <button @click="toggleDropdown" class="search-button" style="margin-bottom: 10px;">{{ dropdownOpen ? '收起学期' : '展开学期' }}</button>

      <div v-if="dropdownOpen">
        <div class="tabs select-all">
        <input type="checkbox" v-model="selectAll" @change="toggleSelectAll" id="selectAll" />
        <label for="selectAll">全选/取消全选</label>
       </div>

       <div class="terms-container">
        <div v-for="term in terms" :key="term" class="tabs term">
          <input type="checkbox" 
                 :value="term" 
                 v-model="selectedTerms" 
                 :id="term" />
          <label :for="term">{{ term }}</label>
          </div>
        </div>
    </div>
  </div>
</div>

    <br>

    <button @click="toggleDropdown2" class="search-button" style="margin-bottom: 10px;">{{ dropdownOpen2 ? '收起条件' : '展开条件' }}</button>

    <!-- 搜索条件 -->

      <div v-if="dropdownOpen2" id="query-container">
        <div class="search-query-container">
          <div v-for="(condition, index) in conditions" :key="index" class="search-condition">
            <!-- 显示连接符的下拉菜单（如果不是第一个条件） -->
            <div v-if="index !== 0">
              <span>连接:</span>  <!-- 连接符 -->
              <select v-model="condition.connector">
                <option value="AND">AND</option>
                <option value="OR">OR</option>
                <option value="NOT">NOT</option>
              </select>
            </div>

            <!-- 选择字段的下拉菜单 -->
            <select v-model="condition.selectedItem" @change="fetchFieldOptions(index)"> <!-- 当改变字段时，调用 fetchFieldOptions 方法 -->
              <option value="" disabled>请选择</option> <!-- 默认选项 -->
              <!-- <option v-for="item in availableItemsFor(index)" :key="item" :value="item">{{ item }}</option>  可选字段 -->
              <option v-for="item in allItems" :key="item" :value="item">{{ item }}</option>  <!-- 可选字段 -->
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
            <!-- <button @click="addCondition" :disabled="!canAddCondition" class="search-button">添加条件</button> 添加条件按钮 -->
            <button @click="addCondition" class="search-button">添加条件</button> <!-- 添加条件按钮 -->
            <button @click="submitSearch" class="search-button">搜索</button>   <!-- 搜索按钮 -->
          </div>
      </div>
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
        <label><input type="checkbox" v-model="showTerm" /> 学期</label>
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
            <th v-if="showTerm" style="width: 10%;">学期</th>
            <th v-if="showCourseNumber" style="width: 5%;">课程序号</th>
            <th v-if="showCourseName" style="width: 20%;">课程名称</th>
            <th v-if="showTeacher" style="width: 15%;">授课教师（工号）</th>
            <th v-if="showCourseType" style="width: 10%;">课程性质</th>
            <th v-if="showCampus" style="width: 7%;">校区</th>
            <th v-if="showMajor" style="width: 15%;">听课专业</th>
            <th v-if="showSchedule" style="width: 20%;">排课信息</th>
            <th v-if="showDepartment" style="width: 15%;">开课学院</th>

          </tr>
        </thead>
        <tbody>
          <tr v-for="(result, index) in paginatedResults" :key="index">
            <td v-if="showTerm">{{ result[0] }}</td> <!-- 学期 -->
            <td v-if="showCourseNumber">{{ result[1] }}</td> <!-- 课程序号 -->
            <td v-if="showCourseName">{{ result[2] }}</td> <!-- 课程名称 -->
            <td v-if="showTeacher">{{ formatTeachers(result[4], result[5]) }}</td> <!-- 合并教师名称和工号 -->
            <td v-if="showCourseType">{{ result[7] }}</td> <!-- 课程性质 -->
            <td v-if="showCampus">{{ result[9] }}</td> <!-- 校区 -->
            <td v-if="showMajor">{{ result[10] }}</td> <!-- 听课专业 -->
            <td v-if="showSchedule">{{ result[14] }}</td> <!-- 排课信息 -->
            <td v-if="showDepartment">{{ result[15] }}</td> <!-- 开课学院 -->
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

    <div v-else-if="searchResults.length === 0">
      <img src="../assets/notfound.png" alt="No data" style="display: block; margin: 0 auto; width: 25%; height: auto;" />
    </div>

    <div v-else>
      <h1>{{searchResults['message']}}</h1>
      <img src="../assets/error.png" alt="Error" style="display: block; margin: 0 auto; width: 20%; height: auto;" />
    </div>

    <!-- 加载动画 -->
    <div v-if="isLoading" class="loading-overlay">
      <img src="../assets/loading.gif" alt="Loading..." />
      <p>Loading</p>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      conditions: [{ selectedItem: '', searchWord: '', connector: '' }], // 条件数组
      allItems: [], // 所有可检索的字段
      terms: [], // 所有学期
      selectedTerms: [],
      dropdownOpen: false,
      dropdownOpen2: true,
      selectAll: false,
      searchResults: ['default'], // 用于存储搜索结果
      showTerm : true, // 是否显示学期
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
      isLoading: false, // 是否正在加载数据
    };
  },
  computed: {
    // canAddCondition() {
    //   // 检查条件的数量是否小于所有字段数量
    //   return this.conditions.length < this.allItems.length;
    // },

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

    // // 获取当前条件可选的字段
    // availableItemsFor(index) {
    //   // 获取除当前条件外所有已选择的字段
    //   const selectedItems = this.conditions
    //     .filter((cond, i) => i !== index)
    //     .map(cond => cond.selectedItem);
    //   // 返回未被选择的字段
    //   return this.allItems.filter(item => !selectedItems.includes(item));
    // },

    // 获取字段的选项
    fetchFieldOptions(index) {
      const fieldName = this.conditions[index].selectedItem; // 获取当前条件的字段名

      // 请求字段的选项数据
      fetch('/api/get_field_options', {
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

    updateConditionsWithSelectedTerms() {
    // 获取所有学期
    const allTerms = this.terms; // 假设 this.terms 存储了所有可用的学期
    const excludedTerms = allTerms.filter(term => !this.selectedTerms.includes(term)); // 筛选未被选中的学期

    // 如果没有未选择的学期，返回空数组
    if (excludedTerms.length === 0) {
      return [];
    }

    // 创建一个新条件数组
    const newConditions = excludedTerms.map(term => ({
      selectedItem: '学期',
      searchWord: term,
      connector: 'NOT' // 使用 NOT 连接
    }));

    return newConditions; // 返回新条件数组
  },
  
    // 提交搜索请求
    submitSearch() {
      this.isLoading = true; // 开始加载数据
      const excludedTermConditions = this.updateConditionsWithSelectedTerms();

      if (this.selectedTerms.length === 0) {
        alert('请选择至少一个学期');
        this.isLoading = false; // 停止加载数据
        return;
      }

      // 创建 conditions 的副本，并添加未选择的学期条件
      const requestData = [
        ...this.conditions.map(cond => ({ ...cond })), // 复制原始条件
        ...excludedTermConditions // 添加未选择的学期条件
      ].filter(cond => cond.selectedItem); // 过滤掉空条件

      console.log('Request data:', requestData);
        // 发送 HTTP 请求到 Flask 后端
        fetch('/api/search', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(requestData),
        })
          .then(response => response.json())
          .then(data => {
            // 更新搜索结果
            this.isLoading = false; // 停止加载数据
            this.searchResults = data;
            console.log('Search results:', data);
          })
          .then(() => {
            this.currentPage = 1; // 重置页码
          })
          .catch(error => {
            this.isLoading = false; // 停止加载数据
            console.error('Error:', error);
          });
      },      

    // 格式化教师信息
    formatTeachers(teachers, ids) {
      if (!teachers) return ''; // 如果 teachers 为空，返回空字符串
      const teacherList = teachers.split(','); // 假设教师用逗号分隔
      const idList = ids.split(','); // 假设工号用逗号分隔
      return teacherList.map((teacher, index) => `${teacher.trim()} (${idList[index].trim()})`).join(', '); // 格式化
    },

    fetchTerms() {
      fetch('/api/get_terms')
        .then(response => response.json())
        .then(data => {
          this.terms = data;
          // 默认选择最近的学期
          if (this.terms.length > 0) {
            this.selectedTerms.push(this.terms[this.terms.length - 1]); // 假设最近的学期在数组最后
          }
        })
        .catch(error => {
          console.error('Error fetching terms:', error);
        });
    },
    toggleDropdown() {
      this.dropdownOpen = !this.dropdownOpen;
    },
    toggleSelectAll() {
      if (this.selectAll) {
        this.selectedTerms = [...this.terms]; // 全选
      } else {
        this.selectedTerms = []; // 取消全选
      }
    },

    toggleDropdown2() {
      this.dropdownOpen2 = !this.dropdownOpen2;
    },

  },

  mounted() {
    // 获取所有可检索的字段
    fetch('/api/get_searchable_fields')
      .then(response => response.json()) // 解析 JSON 数据
      .then(data => { 
        this.allItems = data; // 更新所有字段
      })
      .catch(error => {
        console.error('Error fetching searchable fields:', error);
      });

    this.fetchTerms();
  },
};
</script>

<style scoped>
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
  font-family: "新宋体", Arial, sans-serif; /* 字体 */
  width: 100%; /* 使内容占满整个宽度 */
  margin: 0; /* 外边距 */
  padding: 20px; /* 内边距 */
}

#searchbody h1 {
  text-align: center; /* 居中 */
  font-weight: bold; /* 加粗 */
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
  display: block; /* 使用 Flexbox 布局 */
  flex-wrap: wrap; /* 允许换行 */
  gap: 15px; /* 每个条件之间的间距 */
  margin-bottom: 15px; /* 下方的间距 */
}

.search-query-container {
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
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

#query-container {
  display: block;
}

/* 全选按钮样式 */
.select-all {
  margin-bottom: 10px; /* 与下面选项的间距 */
}

/* 学期选项容器样式 */
.terms-container {
  display: flex; /* 使用 flexbox 布局 */
  flex-wrap: wrap; /* 允许换行 */
}

/* 学期选项样式 */
.term {
  flex: 1 1 auto; /* 自适应宽度，允许换行 */
  min-width: 150px; /* 设置最小宽度 */
  display: flex; /* 使用 flexbox 来对齐 checkbox 和 label */
  align-items: center; /* 垂直居中对齐 */
  margin-right: 10px; /* 右边间距 */
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

.loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(255, 255, 255, 0.8);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color:#209ce0;
  font-weight: bold;
  font-size: 24px;
  font-family: 'Courier New', Courier, monospace;
}


</style>
