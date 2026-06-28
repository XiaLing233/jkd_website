// ─── API 统一响应 ───

export interface APIResponse<T = unknown> {
  code: number
  msg: string
  data: T
}

export interface PaginatedData<T> {
  items: T[]
  page: number
  pageSize: number
  totalCount: number
}

// ─── 日历 ───

export interface CalendarInfo {
  calendarId: number
  calendarName: string
}

// ─── 搜索字段 ───

export interface SearchField {
  id: string
  label: string
  searchType: 'select' | 'input'
}

export interface FieldOption {
  value: string
  label: string
}

// ─── 搜索 ───

export interface SearchCondition {
  field: string
  matchType: 'contains' | 'not_contains'
  value: string
}

export interface SearchGroup {
  conditions: SearchCondition[]
}

export interface SearchRequest {
  groups: SearchGroup[]
  calendar_ids: number[]
  page: number
  page_size: number
}

export interface SearchResult {
  calendarId: number
  term: string
  courseCode: string
  courseName: string
  campus: string
  faculty: string
  majors: string
  totalHours: number | null
  courseType: string
  assessmentMode: string
  capacity: number | null
  enrolled: number | null
  teachers: string
  schedule: string
}

// ─── 最后更新 ───

export interface LastUpdate {
  fetchTime: string
  message: string
}
