<template>
    <div id="searchbody" style="margin: 20px;"> <!-- 搜索主体 -->
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
          <el-icon class="el-icon--left"><Operation /></el-icon>检索条件
        </template>
        
        <div class="search-tip" style="margin-bottom: 15px; padding: 10px; background-color: #f0f9ff; border-radius: 4px; font-size: 13px; color: #606266; display: flex; align-items: center;">
          <el-icon class="el-icon--left"><InfoFilled /></el-icon>
          <span>提示：<strong>条件组之间</strong>是 <el-tag size="small" type="primary">AND</el-tag> 关系，<strong>同一组内</strong>的条件是 <el-tag size="small" type="success">OR</el-tag> 关系。</span>
        </div>

        <div v-for="(group, groupIndex) in conditionGroups" :key="groupIndex" style="margin-bottom: 15px;">
          <el-card shadow="hover">
            <template #header>
              <div style="display: flex; justify-content: space-between; align-items: center;">
                <span style="font-weight: bold;">
                  <el-tag v-if="groupIndex > 0" type="primary" style="margin-right: 8px;">AND</el-tag>
                  条件组 {{ groupIndex + 1 }}
                  <el-tag type="success" size="small" style="margin-left: 8px;">组内 OR</el-tag>
                </span>
                <el-button type="danger" size="small" @click="removeGroup(groupIndex)" :disabled="conditionGroups.length === 1">
                  删除组
                </el-button>
              </div>
            </template>
            
            <div v-for="(condition, condIndex) in group.conditions" :key="condIndex" style="display: flex; align-items: center; margin-bottom: 10px; flex-wrap: wrap; gap: 8px;">
              <!-- OR 标签（组内第二个条件开始显示） -->
              <el-tag v-if="condIndex > 0" type="success" style="margin-right: 4px;">OR</el-tag>
              
              <!-- 选择字段 -->
              <el-select 
                v-model="condition.field" 
                placeholder="选择字段" 
                clearable 
                style="width: 120px"
              >
                <el-option v-for="item in allItems" :key="item.field" :value="item.field">{{ item.field }}</el-option>
              </el-select>

              <!-- 匹配类型 -->
              <el-select 
                v-model="condition.matchType" 
                style="width: 100px"
              >
                <el-option value="contains" label="包含"></el-option>
                <el-option value="not_contains" label="不包含"></el-option>
              </el-select>

              <!-- 搜索词 - 根据字段类型动态切换 -->
              <!-- 选项式字段 -->
              <el-select
                v-if="getFieldType(condition.field) === 'select'"
                v-model="condition.value"
                @focus="fetchFieldOptions(groupIndex, condIndex)"
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
              <!-- 填写式字段 -->
              <el-input
                v-else
                v-model="condition.value"
                placeholder="请输入搜索词"
                clearable
                style="width: 200px;"
              />

              <!-- 删除条件按钮 -->
              <el-button 
                type="danger" 
                :icon="DeleteIcon" 
                circle 
                size="small"
                @click="removeCondition(groupIndex, condIndex)"
                :disabled="group.conditions.length === 1"
              />
            </div>
            
            <!-- 添加条件按钮（组内） -->
            <el-button type="success" plain size="small" @click="addConditionToGroup(groupIndex)">
              <el-icon class="el-icon--left"><component :is="PlusIcon" /></el-icon>添加 OR 条件
            </el-button>
          </el-card>
        </div>
        
        <div style="margin: 15px 0 0 0">
          <el-button type="primary" plain @click="addGroup">
            <el-icon class="el-icon--left"><component :is="PlusIcon" /></el-icon>添加 AND 条件组
          </el-button>
          <el-button type="primary" @click="submitSearch">
            <el-icon class="el-icon--left"><Search /></el-icon>搜索
          </el-button>
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
        <el-checkbox v-model="showCampus" label="校区" />
        <el-checkbox v-model="showDepartment" label="开课学院" />
        <el-checkbox v-model="showPeriod" label="总学时" />
        <el-checkbox v-model="showCourseType" label="课程性质" />
        <el-checkbox v-model="showAssessment" label="考核方式" />
        <el-checkbox v-model="showTeacher" label="授课教师" />
        <el-checkbox v-model="showSchedule" label="排课信息" />
        <el-checkbox v-model="showMajor" label="听课专业" />
        <el-checkbox v-model="showCapacity" label="额定人数" />
        <el-checkbox v-model="showStudentNumber" label="选课人数" />
      </div>

      <div style="overflow-x: auto;">
      <el-table
        :data="paginatedData" 
        stripe 
        :header-cell-style= "{'text-align': 'center'}"
        :cell-style= "{'text-align': 'center'}"
        style="width: 100%;">
        <el-table-column v-if="showTerm" prop="term" label="学期" width="180"></el-table-column>
        <el-table-column v-if="showCourseNumber" prop="courseCode" label="课程序号" width="120"></el-table-column>
        <el-table-column v-if="showCourseName" prop="courseName" label="课程名称" min-width="150"></el-table-column>
        <el-table-column v-if="showCampus" prop="campus" label="校区" width="120"></el-table-column>
        <el-table-column v-if="showDepartment" prop="faculty" label="开课学院" min-width="150"></el-table-column>
        <el-table-column v-if="showPeriod" prop="totalHours" label="总学时" width="90"></el-table-column>
        <el-table-column v-if="showCourseType" prop="courseType" label="课程性质" width="120"></el-table-column>
        <el-table-column v-if="showAssessment" prop="assessmentMode" label="考核方式" width="100"></el-table-column>
        <el-table-column v-if="showTeacher" prop="teachers" label="授课教师" min-width="180"></el-table-column>
        <el-table-column v-if="showSchedule" label="排课信息" min-width="250">
          <template #default="scope">
            <div style="white-space: pre-line; text-align: center;">{{ scope.row.schedule }}</div>
          </template>
        </el-table-column>
        <el-table-column v-if="showMajor" prop="majors" label="听课专业" min-width="250"></el-table-column>
        <el-table-column v-if="showCapacity" prop="capacity" label="额定人数" width="100"></el-table-column>
        <el-table-column v-if="showStudentNumber" prop="enrolled" label="选课人数" width="100"></el-table-column>
      </el-table>
      </div>

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
    element-loading-svg-view-box="0 0 1024 1024"
    element-loading-svg='<path fill="currentColor" d="M512 64a32 32 0 0 1 32 32v192a32 32 0 0 1-64 0V96a32 32 0 0 1 32-32m0 640a32 32 0 0 1 32 32v192a32 32 0 1 1-64 0V736a32 32 0 0 1 32-32m448-192a32 32 0 0 1-32 32H736a32 32 0 1 1 0-64h192a32 32 0 0 1 32 32m-640 0a32 32 0 0 1-32 32H96a32 32 0 0 1 0-64h192a32 32 0 0 1 32 32M195.2 195.2a32 32 0 0 1 45.248 0L376.32 331.008a32 32 0 0 1-45.248 45.248L195.2 240.448a32 32 0 0 1 0-45.248m452.544 452.544a32 32 0 0 1 45.248 0L828.8 783.552a32 32 0 0 1-45.248 45.248L647.744 692.992a32 32 0 0 1 0-45.248M828.8 195.264a32 32 0 0 1 0 45.184L692.992 376.32a32 32 0 0 1-45.248-45.248l135.808-135.808a32 32 0 0 1 45.248 0m-452.544 452.48a32 32 0 0 1 0 45.248L240.448 828.8a32 32 0 0 1-45.248-45.248l135.808-135.808a32 32 0 0 1 45.248 0"></path>'
  >

  </div>
</template>

<script>
import { ElNotification } from 'element-plus';
import { Delete, Plus, InfoFilled } from '@element-plus/icons-vue';
import { markRaw } from 'vue';

export default {
  data() {
    return {
      // 图标
      DeleteIcon: markRaw(Delete),
      PlusIcon: markRaw(Plus),
      
      // 新的条件组结构
      conditionGroups: [
        {
          conditions: [
            { field: '', matchType: 'contains', value: '', options: [] }
          ]
        }
      ],
      allItems: [], // 所有可检索的字段
      terms: [], // 所有学期
      selectedTerms: [],
      searchResults: [], // 用于存储搜索结果

      // 表格列的显示控制
      showTerm : true, // 是否显示学期
      showCourseNumber: true, // 是否显示课程序号
      showCourseName: true, // 是否显示课程名称
      showCampus: true, // 是否显示校区
      showDepartment: true, // 是否显示开课学院
      showPeriod: true, // 是否显示总学时
      showCourseType: true, // 是否显示课程性质
      showAssessment: true, // 是否显示考核方式
      showTeacher: false, // 是否显示授课教师
      showSchedule: true, // 是否显示排课信息
      showMajor: false, // 是否显示听课专业
      showCapacity: true, // 是否显示额定人数
      showStudentNumber: true, // 是否显示选课人数

      // 页码控制
      currentPage: 1, // 当前页码
      pageSize: 20, // 每页显示的行数
      
      // 加载状态
      isLoading: false, // 是否正在加载数据
      isError: false,
      isWarning: false,
      msg: '',

      // UI 状态
      loading_field: false, // 下拉菜单是否正在加载数据
      isIndeterminate: true,
      checkAll: false,
      activeNames: ['term', 'condition'],
    };
  },

  computed: {
    paginatedData() {
      const start = (this.currentPage - 1) * this.pageSize;
      return this.searchResults.slice(start, start + this.pageSize);
    },
  },

  methods: {
    // 添加新的条件组
    addGroup() {
      this.conditionGroups.push({
        conditions: [
          { field: '', matchType: 'contains', value: '', options: [] }
        ]
      });
    },

    // 删除条件组
    removeGroup(groupIndex) {
      if (this.conditionGroups.length === 1) {
        this.isError = true;
        this.msg = '至少需要一个条件组';
        return;
      }
      this.conditionGroups.splice(groupIndex, 1);
    },

    // 在指定组内添加条件
    addConditionToGroup(groupIndex) {
      this.conditionGroups[groupIndex].conditions.push({
        field: '',
        matchType: 'contains',
        value: '',
        options: []
      });
    },

    // 删除指定组内的条件
    removeCondition(groupIndex, condIndex) {
      if (this.conditionGroups[groupIndex].conditions.length === 1) {
        this.isError = true;
        this.msg = '每个条件组至少需要一个条件';
        return;
      }
      this.conditionGroups[groupIndex].conditions.splice(condIndex, 1);
    },

    // 获取字段类型
    getFieldType(fieldName) {
      const item = this.allItems.find(item => item.field === fieldName);
      return item ? item.type : 'input';
    },

    // 获取字段的选项（仅对 select 类型字段）
    fetchFieldOptions(groupIndex, condIndex) {
      const condition = this.conditionGroups[groupIndex].conditions[condIndex];
      const fieldName = condition.field;
      
      // 如果是 input 类型字段，不需要获取选项
      if (this.getFieldType(fieldName) === 'input') {
        return;
      }
      
      const selectTerm = this.selectedTerms;
      this.loading_field = true;

      fetch('/api/get_field_options', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ field_name: fieldName, select_term: selectTerm }),
      })
        .then(response => response.json())
        .then(data => {
          console.log('Field options:', data['data']);
          condition.options = data['data'];
          data['status'] === 'ERROR' ? this.isError = true : this.isError = false;
          data['status'] === 'WARNING' ? this.isWarning = true : this.isWarning = false;
          this.msg = data['message'];
          this.loading_field = false;
        })
        .catch(error => {
          console.error('Error fetching field options:', error);
          this.isError = true;
          this.msg = '后端服务器未响应，请稍后再试';
          condition.options = [];
          this.loading_field = false;
        });
    },

    // 提交搜索请求
    submitSearch() {
      this.isLoading = true;

      if (this.selectedTerms.length === 0) {
        this.isError = true;
        this.msg = '请选择至少一个学期';
        this.isLoading = false;
        return;
      }

      // 检查是否有有效的条件
      const hasValidCondition = this.conditionGroups.some(group =>
        group.conditions.some(cond => cond.field && cond.value)
      );

      // 如果没有有效条件且选择了超过2个学期，则报错
      if (!hasValidCondition && this.selectedTerms.length > 2) {
        this.isError = true;
        this.msg = '至少选择1个检索条件！不允许在不选择条件的情况下查看超过两个学期的课程！';
        this.isLoading = false;
        return;
      }

      // 构建新的请求数据格式
      const requestData = {
        groups: this.conditionGroups
          .map(group => ({
            conditions: group.conditions
              .filter(cond => cond.field && cond.value)
              .map(cond => ({
                field: cond.field,
                matchType: cond.matchType,
                value: cond.value
              }))
          }))
          .filter(group => group.conditions.length > 0),
        terms: this.selectedTerms
      };

      console.log('Request data:', requestData);

      fetch('/api/search', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestData),
      })
        .then(response => response.json())
        .then(data => {
          this.isLoading = false;
          this.searchResults = data['data'];
          this.msg = data['message'];
          console.log('Search results:', data['data']);
          data['status'] === 'ERROR' ? this.isError = true : this.isError = false;
        })
        .then(() => {
          this.currentPage = 1;
        })
        .catch(error => {
          this.isLoading = false;
          this.isError = true;
          this.msg = '后端服务器未响应，请稍后再试';
          console.error('Error:', error);
        });
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
        // data['data'] 现在是对象数组 [{field: '课程序号', type: 'input'}, ...]
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

<!-- 全局样式 - fullscreen loading 挂载在 body 上，scoped 样式无法覆盖 -->
<style>
/* 旧版 Loading 样式 - 覆盖 Element Plus 的 CSS 变量 */
:root {
  --el-loading-fullscreen-spinner-size: 42px;
}

/* 旧版 Loading 样式 - 与 app.c0acfd8887604e1d87ce.css 保持一致 */
.el-loading-mask.is-fullscreen .el-loading-spinner .el-loading-text {
  color: #409eff;
  font-size: 14px;
  margin: 3px 0;
}

.el-loading-mask.is-fullscreen .el-loading-spinner .circular {
  height: 42px !important;
  width: 42px !important;
}

.el-loading-mask.is-fullscreen .el-loading-spinner svg path {
  fill: #409eff;
}

/* el-table 样式自定义 */
.el-table thead {
  color: #4c5c70 !important;
}

.el-table tbody {
  color: #2b3b4e !important;
}
</style>
