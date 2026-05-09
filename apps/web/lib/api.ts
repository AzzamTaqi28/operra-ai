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

export type PurchaseRequestComment = {
  id: string
  actor_user_id: string
  body: string
  created_at: string
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
  roles?: string[]
}

export type DepartmentListItem = {
  id: string
  name: string
  code: string
  status?: string
}

export type ApprovalActionResult = {
  purchase_request_id: string
  approval_step_instance_id: string
  action: string
  status: string
  current_step_instance_id?: string | null
  completed_at?: string | null
}

export type WorkflowGenerationResult = {
  workflow_json: {
    name: string
    type: string
    version: number
    steps: unknown[]
    [key: string]: unknown
  }
  explanation: string
  mermaid_diagram: string
  validation: {
    valid: boolean
    errors: unknown[]
  }
  warnings: string[]
  provider: string
  model?: string
}

export type WorkflowCreateResponse = {
  workflow: WorkflowListItem
  version: WorkflowVersionSummary
}

export type WorkflowVersionSummary = {
  id: string
  workflow_id: string
  version_number: number
  config_json: unknown
  mermaid_diagram?: string
  explanation?: string
  created_by: string
  created_at: string
}

export type WorkflowListItem = {
  id: string
  organization_id: string
  name: string
  type: string
  status: string
  active_version_id?: string | null
  created_by: string
  created_at: string
  updated_at: string
  active_version?: {
    id: string
    organization_id: string
    workflow_id: string
    version_number: number
    config_json: unknown
    mermaid_diagram?: string
    explanation?: string
    created_by: string
    created_at: string
  } | null
}

export type AuditLogItem = {
  id: string
  organization_id: string
  actor_user_id?: string | null
  action: string
  entity_type: string
  entity_id?: string | null
  old_value?: unknown
  new_value?: unknown
  ip_address?: string | null
  user_agent?: string | null
  created_at: string
}

export type CurrentUser = {
  id: string
  organization_id: string
  department_id?: string | null
  name: string
  email: string
  roles: string[]
  status: string
}
