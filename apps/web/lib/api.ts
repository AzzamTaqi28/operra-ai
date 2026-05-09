import { cookies } from "next/headers"

const API_BASE_URL = process.env.API_URL ?? "http://localhost:8080"

export async function getToken() {
  const store = await cookies()
  return store.get("operra_token")?.value ?? null
}

export async function apiGet<T>(path: string, token: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...init,
    cache: "no-store",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
      ...(init?.headers ?? {}),
    },
  })

  const data = await response.json().catch(() => null)
  if (!response.ok) {
    const message = data?.error?.message ?? `Request failed: ${response.status}`
    throw new Error(message)
  }

  return data as T
}

export type ApiListResponse<T> = {
  data: T[]
  pagination: {
    page: number
    page_size: number
    total: number
  }
}

export type PurchaseRequest = {
  id: string
  title: string
  item_name: string
  description: string
  quantity: number
  estimated_amount: number
  currency: string
  urgency: string
  department_id: string
  requester_id: string
  status: string
  current_step_instance_id?: string | null
  submitted_at?: string | null
  completed_at?: string | null
  created_at: string
  updated_at: string
}

export type PurchaseRequestDetail = PurchaseRequest & {
  attachments: Array<{
    id: string
    file_name: string
    file_size: number
    mime_type: string
    created_at: string
  }>
  approval_steps: Array<{
    id: string
    step_name: string
    step_key: string
    status: string
    sequence_number: number
    approver_role_key: string
    scope: string
    created_at: string
    updated_at: string
  }>
  approval_actions: Array<{
    id: string
    action: string
    comment?: string | null
    created_at: string
  }>
}

export type ApprovalQueueItem = {
  request_id: string
  step_id: string
  step_name: string
  role_key: string
  scope: string
  status: string
}

export type UserListItem = {
  id: string
  name: string
  email: string
  department_id?: string | null
  status: string
}

export type DepartmentListItem = {
  id: string
  name: string
  code: string
  status?: string
}
