import type {
  APIResponse,
  PaginatedData,
  CalendarInfo,
  SearchField,
  FieldOption,
  SearchRequest,
  SearchResult,
  LastUpdate,
} from './types'

const BASE = '/api'

async function request<T>(url: string, options?: RequestInit): Promise<APIResponse<T>> {
  const resp = await fetch(`${BASE}${url}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  })
  if (!resp.ok) {
    throw new Error(`HTTP ${resp.status}`)
  }
  return resp.json()
}

// GET /api/search-fields
export function fetchSearchFields(): Promise<APIResponse<SearchField[]>> {
  return request('/search-fields')
}

// POST /api/search-fields/options
export function fetchFieldOptions(
  field: string,
  calendarIds: number[],
): Promise<APIResponse<FieldOption[]>> {
  return request('/search-fields/options', {
    method: 'POST',
    body: JSON.stringify({ field, calendar_ids: calendarIds }),
  })
}

// GET /api/calendars
export function fetchCalendars(): Promise<APIResponse<CalendarInfo[]>> {
  return request('/calendars')
}

// POST /api/courses/search
export function fetchSearch(req: SearchRequest): Promise<APIResponse<PaginatedData<SearchResult>>> {
  return request('/courses/search', {
    method: 'POST',
    body: JSON.stringify(req),
  })
}

// GET /api/last-update
export function fetchLastUpdate(): Promise<APIResponse<LastUpdate>> {
  return request('/last-update')
}
