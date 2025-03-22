<template>
    <div id="searchbody" style="margin: 20px"> <!-- 搜索主体 -->
    <!-- 搜索条件 -->
     <el-card shadow="never" style="width: 100%; margin-bottom: 10px;">
      <template #header>
        <div style="font-weight: bold;">筛选条件</div>
      </template>
      <el-collapse v-model="activeNames">
      <el-collapse-item name="term">
        <template #title>
          <el-icon class="el-icon--left"><Calendar /></el-icon>学年学期
        </template>
        <el-checkbox
          v-model="checkAll"
          :indeterminate="isIndeterminate"
          @change="handleCheckAllChange"
          >
            全选
        </el-checkbox>

        <el-checkbox-group 
          v-model="selectedTerms"
          @change="handleTermChange"
          >
          <el-checkbox
            v-for="term in terms"
            :key="term"
            :value="term"
            :label="term"
            >
              {{ term }}
          </el-checkbox>
        </el-checkbox-group>
      </el-collapse-item>

      <el-collapse-item name="condition">
        <template #title>
          <el-icon class="el-icon--left"><Operation /></el-icon>逻辑条件
        </template>
        <div style="display: flex; flex-wrap: wrap;">
          <div v-for="(condition, index) in conditions" :key="index">
            <!-- 显示连接符的下拉菜单（如果不是第一个条件） -->
              <el-card shadow="hover" style="width: 550px">
                <el-select 
                  placeholder="连接符" 
                  clearable 
                  v-model="condition.connector"
                  style="width: 100px"
                  >
                    <el-option v-if="index !== 0" value="AND">AND</el-option>
                    <el-option v-if="index !== 0" value="OR">OR</el-option>
                    <el-option value="NOT">NOT</el-option>
                  </el-select>
              

                <!-- 选择字段的下拉菜单 -->
                <el-select 
                  v-model="condition.selectedItem" 
                  placeholder="字段选择" 
                  clearable 
                  :loading="loading_field"
                  style="width: 120px"
                  > <!-- 当改变字段时，调用 fetchFieldOptions 方法 -->
                  <el-option v-for="item in allItems" :key="item" :value="item">{{ item }}</el-option>  <!-- 可选字段 -->
                </el-select>

                <!-- 输入搜索词 -->
                <el-select
                  v-model="condition.searchWord"
                  @focus="fetchFieldOptions(index)"
                  placeholder="请输入搜索词"
                  filterable
                  allow-create
                  clearable
                  style="width: 200px;"
                  :loading="loading_field"
                >
                  <el-option 
                    v-for="option in condition.options" 
                    :key="option" 
                    :value="option">{{ option }}
                  </el-option>
                </el-select>

                <el-button type="primary" @click="removeCondition(index)">移除条件</el-button> <!-- 移除条件按钮 -->
              </el-card>
              
          </div>
          </div>
        
          <div style="margin: 15px 0 0 0">
            <!-- <button @click="addCondition" :disabled="!canAddCondition" class="search-button">添加条件</button> 添加条件按钮 -->
            <el-button type="primary" plain @click="addCondition">添加条件</el-button> <!-- 添加条件按钮 -->
            <el-button type="primary" @click="submitSearch"><el-icon class="el-icon--left"><Search /></el-icon>搜索</el-button>   <!-- 搜索按钮 -->
          </div>
      </el-collapse-item>
      </el-collapse>
     </el-card>
      
     <el-card shadow="never" style="width: 100%;">
      <template #header>
        <div style="font-weight: bold;">教学任务查询</div>
      </template>
 <!-- 显示搜索结果 -->
    <!-- 数据库中表格的结构是：
     课程序号、课程名称、周学时、授课教师、教师工号、课程性质、授课语言、校区、听课专业、选课人数、起始周、结束周、排课信息、开课学院 -->
      <div v-if="searchResults.length > 0">
        <!-- 列选择复选框 -->
      <div style="margin-bottom: 10px; display: flex; align-items: center; justify-items: center;">
        <div style="margin-right: 10px; padding-bottom: 2px; display: inline-flex"><el-icon class="el-icon--left"><Filter /></el-icon><el-text>显示列：</el-text></div>
        <el-checkbox v-model="showTerm" label="学期" />
        <el-checkbox v-model="showCourseNumber" label="课程序号" />
        <el-checkbox v-model="showCourseName" label="课程名称" /> 
        <el-checkbox v-model="showTeacher" label="授课教师（工号）" />
        <el-checkbox v-model="showCourseType" label="课程性质" />
        <el-checkbox v-model="showCampus" label="校区" />
        <el-checkbox v-model="showMajor" label="听课专业" />
        <el-checkbox v-model="showSchedule" label="排课信息" />
        <el-checkbox v-model="showDepartment" label="开课学院" />
        <el-checkbox v-model="showStudentNumber" label="选课人数" />
      </div>

      <el-table
        :data="paginatedData" 
        stripe 
        :header-cell-style= "{'text-align': 'center'}"
        :cell-style= "{'text-align': 'center'}"
        style="width: 100%;">
        <el-table-column sortable v-if="showTerm" prop="学期" label="学期"></el-table-column>
        <el-table-column sortable v-if="showCourseNumber" prop="课程序号" label="课程序号"></el-table-column>
        <el-table-column v-if="showCourseName" prop="课程名称" label="课程名称"></el-table-column>
         <!-- 使用 formatter 函数来格式化教师和工号的显示 -->
        <el-table-column
          v-if="showTeacher"
          :formatter="formatTeachers"
          label="授课教师（工号）"
        ></el-table-column>
        <el-table-column v-if="showCourseType" prop="课程性质" label="课程性质"></el-table-column>
        <el-table-column v-if="showCampus" prop="校区" sortable label="校区"></el-table-column>
        <el-table-column v-if="showMajor" prop="听课专业" label="听课专业"></el-table-column>
        <el-table-column v-if="showSchedule" prop="排课信息" label="排课信息"></el-table-column>
        <el-table-column v-if="showDepartment" prop="开课学院" label="开课学院"></el-table-column>
        <el-table-column v-if="showStudentNumber" prop="选课人数" label="选课人数"></el-table-column>
      </el-table>

      <!-- 分页控制 -->
       <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :pageSizes="[20, 50, 100, 500, 1000]"
        layout="total, sizes, prev, pager, next, jumper"
        :total="searchResults.length"
        style=" float: right; margin: 10px 60px 10px 0px;"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
       />
    </div>

    <div v-else-if="searchResults.length === 0">
      <img src="../assets/notfound.webp" alt="No data" style="display: block; margin: 0 auto; width: 25%; height: auto;" />
    </div>
     </el-card>


    <!-- 错误提示框 -->
      <el-dialog
        v-model="isError"
        title="错误提示"
        width="500"
        draggable
        >
        <div v-html="msg"></div>

        <template #footer>
          <div>
            <el-button type="primary" @click="isError = false">确定</el-button>
          </div>
        </template> 
      </el-dialog>  
  </div>

  <div
    v-loading.fullscreen.lock="isLoading"
    element-loading-text="Loading"
  >

  </div>
</template>

<script>
import { ElNotification } from 'element-plus';

export default {
  data() {
    return {
      // 搜索条件
      conditions: [{ selectedItem: '', searchWord: '', connector: '' }], // 条件数组
      allItems: [], // 所有可检索的字段
      terms: [], // 所有学期
      selectedTerms: [],
      searchResults: [], // 用于存储搜索结果

      // 表格列的显示控制
      showTerm : true, // 是否显示学期
      showCourseNumber: true, // 是否显示课程序号
      showCourseName: true, // 是否显示课程名称
      showTeacher: true, // 是否显示授课教师
      showCourseType: true, // 是否显示课程性质
      showCampus: true, // 是否显示校区
      showMajor: false, // 是否显示听课专业
      showSchedule: true, // 是否显示排课信息
      showDepartment: true, // 是否显示开课学院
      showStudentNumber: true, // 是否显示选课人数

      // 页码控制
      currentPage: 1, // 当前页码
      pageSize: 20, // 每页显示的行数
      
      // 加载状态
      isLoading: false, // 是否正在加载数据
      isError: false,
      isWarning: false,
      msg: '',

      // 以下是新 UI
      loading_field: false, // 下拉菜单是否正在加载数据
      isIndeterminate: true,
      checkAll: false,
      activeNames: ['term', 'condition'],

      //
    };
  },

  computed: {
    paginatedData() {
      const start = (this.currentPage - 1) * this.pageSize;
      return this.searchResults.slice(start, start + this.pageSize);
    },
  },

  methods: {
    // 添加条件
    addCondition() {
      this.conditions.push({ selectedItem: '', searchWord: '', connector: '' });
    },

    // 移除条件
    removeCondition(index) {
      if (this.conditions.length === 1) {
        this.isError = true;
        this.msg = '至少需要一个条件';
        return;
      }

      this.conditions.splice(index, 1);
      
      // 如果只剩下一个条件，把 connector 设置为空
      if (this.conditions.length === 1) {
        this.conditions[0].connector = '';
      }
    },

    // 获取字段的选项
    fetchFieldOptions(index) {
      const fieldName = this.conditions[index].selectedItem; // 获取当前条件的字段名
      const selectTerm = this.selectedTerms; // 获取当前选择的学期
      this.loading_field = true; // 开始加载数据
      // console.log(this.loading_field)

      // 请求字段的选项数据
      fetch('/api/get_field_options', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ field_name: fieldName, select_term: selectTerm}),
      })
        .then(response => response.json()) // 解析 JSON 数据
        .then(data => {
          // 更新当前条件的可选项
          console.log('Field options:', data['data']);
          this.conditions[index].options = data['data']; // 直接赋值，Vue 3 支持响应式
          data['status'] === 'ERROR' ? this.isError = true : this.isError = false;
          data['status'] === 'WARNING' ? this.isWarning = true : this.isWarning = false;
          // console.log('status:', data['status']);
          // console.log('isWarning:', this.isWarning);
          this.msg = data['message'];
          this.loading_field = false; // 停止加载数据
        })
        .catch(error => {
          console.error('Error fetching field options:', error);
          this.isError = true;
          this.msg= '后端服务器未响应，请稍后再试';
          // 把 options 设置为空数组
          this.conditions[index].options = [];
          this.loading_field = false; // 停止加载数据
        });


      // console.log(this.loading_field)
    },

    updateConditionsWithSelectedTerms() {
    // map 所有已经选择的学期，创建新的条件对象
    const newConditions = this.selectedTerms.map(term => ({
      selectedItem: '学期',
      searchWord: term,
      connector: 'OR' // 使用 OR 连接
    }));

    return newConditions; // 返回新条件数组
  },
  
    // 提交搜索请求
    submitSearch() {
      this.isLoading = true; // 开始加载数据
      const excludedTermConditions = this.updateConditionsWithSelectedTerms();

      // 复制一份 conditions，确保原始数据不被修改
      let requestData = this.conditions.map(cond => ({ ...cond }));

      // 重新排序，确保相同的字段在一起，而且固定第一个条件的位置不变
      requestData = [
        requestData[0], // 第一个条件
        ...requestData.slice(1).sort((a, b) => a.selectedItem.localeCompare(b.selectedItem)),
      ];

      if (this.selectedTerms.length === 0) {
        this.isError = true;
        this.msg = '请选择至少一个学期';
        this.isLoading = false; // 停止加载数据
        return;
      }

      // 添加学期条件并过滤空条件
      requestData = [...requestData, ...excludedTermConditions].filter(cond => cond.selectedItem || cond.searchWord);

      // console.log('Request data:', requestData);
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
            this.searchResults = data['data'];
            this.msg = data['message'];
            console.log('Search results:', data['data']);
            data['status'] === 'ERROR' ? this.isError = true : this.isError = false;
          })
          .then(() => {
            this.currentPage = 1; // 重置页码
          })
          .catch(error => {
            this.isLoading = false; // 停止加载数据
            this.isError = true;
            this.msg= '后端服务器未响应，请稍后再试';
            console.error('Error:', error);
          });
      },      

    // 格式化教师信息
    formatTeachers(row) {
      const teachers = row['授课教师'] ? row['授课教师'].split(',') : [];
      const ids = row['教师工号'] ? row['教师工号'].split(',') : [];

      return teachers.map((teacher, index) => `${teacher.trim()} (${ids[index]?.trim() || ''})`).join(', ');
    },

    fetchTerms() {
    this.isLoading = true; // 开始加载数据
    fetch('/api/get_terms')
      .then(response => response.json())
      .then(data => {
        // 将学期数据转换为对象数组，避免使用 label 作为值
        this.terms = data['data'];
        this.msg = data['message'];

        // 默认选择最近的学期
        if (this.terms.length > 0) {
            this.selectedTerms.push(this.terms[this.terms.length - 1]); // 假设最近的学期在数组最后
        }

        this.isLoading = false; // 停止加载数据
        data['status'] === 'ERROR' ? this.isError = true : this.isError = false;
      })
      .catch(error => {
        console.error('Error fetching terms:', error);
        this.msg= '后端服务器未响应，请稍后再试';
        this.isError = true;
        this.isLoading = false; // 停止加载数据
      });
  },

    handleCheckAllChange(val) {
      this.selectedTerms = val ? this.terms : [];
      this.isIndeterminate = false;
    },

    handleTermChange(val) {
      const checkedCount = val.length;
      this.checkAll = checkedCount === this.terms.length;
      this.isIndeterminate = checkedCount > 0 && checkedCount < this.terms.length;
    },

    handleSizeChange(val) {
      this.pageSize = val;
    },

    handleCurrentChange(val) {
      this.currentPage = val;
    },
  },

  watch: {
    isWarning(newVal) {
      if (newVal) {
        ElNotification({
          title: '提示',
          message: this.msg,
          type: 'warning',
          duration: 0,
          onClose: () => {
            this.isWarning = false;
          },
      })
      }
    }
  },

  mounted() {
    // 获取所有可检索的字段
    this.isLoading = true; // 开始加载数据
    fetch('/api/get_searchable_fields')
      .then(response => response.json()) // 解析 JSON 数据
      .then(data => { 
        this.allItems = data['data']; // 更新所有字段
        this.isLoading = false; // 停止加载数据
        this.msg = data['message'];

        data['status'] === 'ERROR' ? this.isError = true : this.isError = false;
      })
      .catch(error => {
        console.error('Error fetching searchable fields:', error);
        this.msg= '后端服务器未响应，请稍后再试';
        this.isError = true;
        this.isLoading = false; // 停止加载数据
      });

    this.fetchTerms();
  },
};
</script>

<style scoped>
</style>
