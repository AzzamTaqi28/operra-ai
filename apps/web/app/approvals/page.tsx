import Link from "next/link"
import { redirect } from "next/navigation"

import { AppShell } from "@/components/app-shell"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { ApprovalsTable } from "@/components/approvals-table"
import { apiGet, getToken, type ApiListResponse, type ApprovalQueueItem, type DepartmentListItem, type PurchaseRequest, type UserListItem } from "@/lib/api"

function currency(amount: number, code = "IDR") {
  return new Intl.NumberFormat("id-ID", { style: "currency", currency: code, maximumFractionDigits: 0 }).format(amount)
}

export default async function ApprovalsPage() {
  const token = await getToken()
  if (!token) redirect("/login")

  const [pendingApprovals, requests, users, departments] = await Promise.all([
    apiGet<ApiListResponse<ApprovalQueueItem>>("/api/v1/approvals/pending?page_size=100", token),
    apiGet<ApiListResponse<PurchaseRequest>>("/api/v1/purchase-requests?page_size=100", token),
    apiGet<ApiListResponse<UserListItem>>("/api/v1/users?page_size=100", token),
    apiGet<ApiListResponse<DepartmentListItem>>("/api/v1/departments?page_size=100", token),
  ])

  const requestById = new Map(requests.data.map((request) => [request.id, request]))
  const userById = new Map(users.data.map((user) => [user.id, user.name]))
  const departmentById = new Map(departments.data.map((department) => [department.id, department.name]))

  return (
    <AppShell title="Pending Approvals" description="Current user's actionable queue, scoped by role and organization.">
      <Card>
        <CardHeader>
          <CardTitle>Approval queue</CardTitle>
          <CardDescription>Requests waiting on your role-specific action.</CardDescription>
        </CardHeader>
        <CardContent>
          <ApprovalsTable
            rows={pendingApprovals.data.map((item) => {
              const request = requestById.get(item.request_id)
              return {
                request_id: item.request_id,
                title: request?.title ?? "Unknown request",
                requester: request ? (userById.get(request.requester_id) ?? request.requester_id.slice(0, 8)) : "—",
                department: request ? (departmentById.get(request.department_id) ?? request.department_id.slice(0, 8)) : "—",
                amount: request ? currency(request.estimated_amount, request.currency) : "—",
                step_name: item.step_name,
                scope: item.scope,
                status: item.status,
              }
            })}
          />
        </CardContent>
      </Card>
    </AppShell>
  )
}
