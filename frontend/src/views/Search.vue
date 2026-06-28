<template>
  <div id="searchbody" style="margin: 20px">
    <el-card shadow="never" style="width: 100%; margin-bottom: 10px">
      <template #header><div style="font-weight: bold">筛选条件</div></template>
      <el-collapse v-model="activeNames">
        <!-- 学期选择 -->
        <el-collapse-item name="term">
          <template #title><el-icon class="el-icon--left"><Calendar /></el-icon>学年学期</template>
          <el-checkbox v-model="checkAll" :indeterminate="isIndeterminate" @change="handleCheckAll">全选</el-checkbox>
          <el-checkbox-group v-model="selectedTerms" @change="handleTermChange">
            <el-checkbox v-for="t in terms" :key="t.calendarId" :value="t.calendarId" :label="t.calendarId">
              {{ t.calendarName }}
            </el-checkbox>
          </el-checkbox-group>
        </el-collapse-item>

        <!-- 检索条件 -->
        <el-collapse-item name="condition">
          <template #title><el-icon class="el-icon--left"><Operation /></el-icon>检索条件</template>
          <div style="margin-bottom:15px;padding:10px;background:#f0f9ff;border-radius:4px;font-size:13px;color:#606266">
            <el-icon class="el-icon--left"><InfoFilled /></el-icon>
            提示：<strong>条件组之间</strong>是 <el-tag size="small" type="primary">AND</el-tag> 关系，<strong>同一组内</strong>是 <el-tag size="small" type="success">OR</el-tag> 关系。
          </div>

          <div v-for="(group, gi) in conditionGroups" :key="gi" style="margin-bottom:15px">
            <el-card shadow="hover">
              <template #header>
                <div style="display:flex;justify-content:space-between;align-items:center">
                  <span><el-tag v-if="gi>0" type="primary">AND</el-tag> 条件组 {{ gi+1 }} <el-tag type="success" size="small">组内 OR</el-tag></span>
                  <el-button type="danger" size="small" @click="removeGroup(gi)" :disabled="conditionGroups.length===1">删除组</el-button>
                </div>
              </template>

              <div v-for="(cond, ci) in group.conditions" :key="ci" style="display:flex;align-items:center;margin-bottom:10px;flex-wrap:wrap;gap:8px">
                <el-tag v-if="ci>0" type="success">OR</el-tag>
                <el-select v-model="cond.field" placeholder="选择字段" clearable style="width:120px">
                  <el-option v-for="f in allItems" :key="f.id" :value="f.id" :label="f.label" />
                </el-select>
                <el-select v-model="cond.matchType" style="width:100px">
                  <el-option value="contains" label="包含" />
                  <el-option value="not_contains" label="不包含" />
                </el-select>
                <el-select v-if="getFieldType(cond.field)==='select'"
                  v-model="cond.value" @focus="fetchOpts(gi,ci)" filterable allow-create clearable
                  style="width:200px" :loading="loadingField">
                  <el-option v-for="o in cond.options" :key="o" :value="o" :label="o" />
                </el-select>
                <el-input v-else v-model="cond.value" placeholder="请输入搜索词" clearable style="width:200px" />
                <el-button type="danger" :icon="Delete" circle size="small" @click="removeCond(gi,ci)" :disabled="group.conditions.length===1" />
              </div>
              <el-button type="success" plain size="small" @click="addCond(gi)"><el-icon><Plus /></el-icon>添加 OR 条件</el-button>
            </el-card>
          </div>

          <el-button type="primary" plain @click="addGroup"><el-icon><Plus /></el-icon>添加 AND 条件组</el-button>
          <el-button type="primary" @click="doSearchFromStart"><el-icon><Search /></el-icon>搜索</el-button>
        </el-collapse-item>
      </el-collapse>
    </el-card>

    <!-- 结果 -->
    <el-card shadow="never" style="width:100%">
      <template #header><div style="font-weight:bold">教学任务查询</div></template>
      <div v-if="searchResults.length>0">
        <div class="col-checks" style="margin-bottom:10px;display:flex;align-items:center;overflow-x:auto;white-space:nowrap">
          <el-text style="margin-right:10px"><el-icon><Filter /></el-icon>显示列：</el-text>
          <el-checkbox v-model="showTerm" label="学期" />
          <el-checkbox v-model="showCourseCode" label="课程序号" />
          <el-checkbox v-model="showCourseName" label="课程名称" />
          <el-checkbox v-model="showCampus" label="校区" />
          <el-checkbox v-model="showFaculty" label="开课学院" />
          <el-checkbox v-model="showPeriod" label="总学时" />
          <el-checkbox v-model="showCourseType" label="课程性质" />
          <el-checkbox v-model="showAssessment" label="考核方式" />
          <el-checkbox v-model="showTeacher" label="授课教师" />
          <el-checkbox v-model="showSchedule" label="排课信息" />
          <el-checkbox v-model="showMajor" label="听课专业" />
          <el-checkbox v-model="showCapacity" label="额定人数" />
          <el-checkbox v-model="showEnrolled" label="选课人数" />
        </div>
        <div style="overflow-x:auto">
          <el-table :data="searchResults" stripe :header-cell-style="{'text-align':'center'}" :cell-style="{'text-align':'center'}" style="width:100%">
            <el-table-column v-if="showTerm" prop="term" label="学期" width="180" />
            <el-table-column v-if="showCourseCode" prop="courseCode" label="课程序号" width="120" />
            <el-table-column v-if="showCourseName" prop="courseName" label="课程名称" min-width="150" />
            <el-table-column v-if="showCampus" prop="campus" label="校区" width="120" />
            <el-table-column v-if="showFaculty" prop="faculty" label="开课学院" min-width="150" />
            <el-table-column v-if="showPeriod" prop="totalHours" label="总学时" width="90" />
            <el-table-column v-if="showCourseType" prop="courseType" label="课程性质" width="120" />
            <el-table-column v-if="showAssessment" prop="assessmentMode" label="考核方式" width="100" />
            <el-table-column v-if="showTeacher" prop="teachers" label="授课教师" min-width="180" />
            <el-table-column v-if="showSchedule" label="排课信息" min-width="250">
              <template #default="scope"><div style="white-space:pre-line;text-align:center">{{ scope.row.schedule }}</div></template>
            </el-table-column>
            <el-table-column v-if="showMajor" prop="majors" label="听课专业" min-width="250" />
            <el-table-column v-if="showCapacity" prop="capacity" label="额定人数" width="100" />
            <el-table-column v-if="showEnrolled" prop="enrolled" label="选课人数" width="100" />
          </el-table>
        </div>
        <el-pagination
          v-model:current-page="currentPage" v-model:page-size="pageSize"
          :page-sizes="[20,50,100]"
          layout="total,sizes,prev,pager,next,jumper"
          :total="totalCount"
          style="float:right;margin:10px 60px 10px 0"
          @size-change="doSearchFromStart"
          @current-change="doSearch"
        />
      </div>
      <div v-else-if="!loading">
        <img src="/assets/notfound.webp" alt="No data" style="display:block;margin:0 auto;width:25%;height:auto" />
      </div>
    </el-card>

    <el-dialog v-model="isError" title="错误提示" width="500" draggable>
      <div v-html="msg" />
      <template #footer><el-button type="primary" @click="isError=false">确定</el-button></template>
    </el-dialog>
  </div>

  <div v-loading.fullscreen.lock="loading" element-loading-text="Loading" />
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, markRaw } from 'vue'
import { ElNotification } from 'element-plus'
import { Delete, Plus } from '@element-plus/icons-vue'
import { fetchSearchFields, fetchFieldOptions, fetchCalendars, fetchSearch } from '@/api/client'
import type { SearchField, SearchResult, CalendarInfo, FieldOption } from '@/api/types'

const DeleteIcon = markRaw(Delete)
const PlusIcon = markRaw(Plus)
const loading = ref(false)
const loadingField = ref(false)
const isError = ref(false)
const msg = ref('')

// 学期
const terms = ref<CalendarInfo[]>([])
const selectedTerms = ref<number[]>([])
const checkAll = ref(false)
const isIndeterminate = ref(true)

// 字段
const allItems = ref<SearchField[]>([])
const getFieldType = (id: string) => allItems.value.find((f: SearchField) => f.id === id)?.searchType ?? 'input'

// 条件
type MatchType = 'contains' | 'not_contains'
interface Condition { field: string; matchType: MatchType; value: string; options: string[] }
interface Group { conditions: Condition[] }
const conditionGroups = ref<Group[]>([{ conditions: [{ field: '', matchType: 'contains', value: '', options: [] }] }])

const addGroup = () => conditionGroups.value.push({ conditions: [{ field: '', matchType: 'contains', value: '', options: [] }] })
const removeGroup = (i: number) => conditionGroups.value.splice(i, 1)
const addCond = (gi: number) => conditionGroups.value[gi]!.conditions.push({ field: '', matchType: 'contains', value: '', options: [] })
const removeCond = (gi: number, ci: number) => conditionGroups.value[gi]!.conditions.splice(ci, 1)

const fetchOpts = (gi: number, ci: number) => {
  const cond = conditionGroups.value[gi]?.conditions[ci]
  if (!cond || !cond.field || getFieldType(cond.field) !== 'select') return
  loadingField.value = true
  fetchFieldOptions(cond.field, selectedTerms.value)
    .then(r => { cond.options = (r.data || []).map((o: any) => o.value); loadingField.value = false })
    .catch(() => { isError.value = true; msg.value = '获取选项失败'; loadingField.value = false })
}

// 搜索
const searchResults = ref<SearchResult[]>([])
const totalCount = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

const doSearch = () => {
  if (selectedTerms.value.length === 0) { isError.value = true; msg.value = '请选择至少一个学期'; return }
  const hasCond = conditionGroups.value.some(g => g.conditions.some(c => c.field && c.value))
  if (!hasCond && selectedTerms.value.length > 2) { isError.value = true; msg.value = '至少选择1个检索条件！'; return }
  loading.value = true
  fetchSearch({
    groups: conditionGroups.value
      .map(g => ({ conditions: g.conditions.filter((c: Condition) => c.field && c.value).map((c: Condition) => ({ field: c.field, matchType: c.matchType, value: c.value })) }))
      .filter(g => g.conditions.length > 0),
    calendar_ids: selectedTerms.value,
    page: currentPage.value,
    page_size: pageSize.value,
  })
    .then(r => { searchResults.value = r.data.items; totalCount.value = r.data.totalCount; loading.value = false })
    .catch(() => { loading.value = false; isError.value = true; msg.value = '搜索请求失败' })
}

const doSearchFromStart = () => { currentPage.value = 1; doSearch() }

// 学期
const handleCheckAll = (val: boolean) => { selectedTerms.value = val ? terms.value.map((t: CalendarInfo) => t.calendarId) : []; isIndeterminate.value = false }
const handleTermChange = (val: number[]) => { checkAll.value = val.length === terms.value.length; isIndeterminate.value = val.length > 0 && val.length < terms.value.length }

// 列显示
const showTerm = ref(true); const showCourseCode = ref(true); const showCourseName = ref(true)
const showCampus = ref(true); const showFaculty = ref(true); const showPeriod = ref(true)
const showCourseType = ref(true); const showAssessment = ref(true); const showTeacher = ref(false)
const showSchedule = ref(true); const showMajor = ref(false); const showCapacity = ref(true); const showEnrolled = ref(true)

const activeNames = ref(['term', 'condition'])

onMounted(async () => {
  loading.value = true
  try {
    const [fields, cals] = await Promise.all([fetchSearchFields(), fetchCalendars()])
    allItems.value = fields.data || []
    terms.value = cals.data || []
    if (terms.value.length > 0) {
      selectedTerms.value.push(terms.value[0]!.calendarId)
      doSearch()
    }
  } catch { isError.value = true; msg.value = '初始化失败' }
  loading.value = false
})
</script>

<style scoped>
.col-checks > * { flex-shrink: 0 }
.col-checks::-webkit-scrollbar { height: 6px }
.col-checks::-webkit-scrollbar-thumb { background: #c0c4cc; border-radius: 3px }
</style>

<style>
:root { --el-loading-fullscreen-spinner-size: 42px }
.el-loading-mask.is-fullscreen .el-loading-spinner .el-loading-text { color:#409eff;font-size:14px;margin:3px 0 }
.el-loading-mask.is-fullscreen .el-loading-spinner .circular { height:42px!important;width:42px!important }
.el-loading-mask.is-fullscreen .el-loading-spinner svg path { fill:#409eff }
.el-table thead { color:#4c5c70!important }
.el-table tbody { color:#2b3b4e!important }
</style>
