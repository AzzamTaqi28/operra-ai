import Link from "next/link"
import { redirect } from "next/navigation"

import { AppShell } from "@/components/app-shell"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
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
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Request ID</TableHead>
                <TableHead>Title</TableHead>
                <TableHead>Requester</TableHead>
                <TableHead>Department</TableHead>
                <TableHead>Amount</TableHead>
                <TableHead>Current Step</TableHead>
                <TableHead>Scope</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Action</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {pendingApprovals.data.map((item) => {
                const request = requestById.get(item.request_id)
                return (
                  <TableRow key={item.step_id}>
                    <TableCell>{item.request_id.slice(0, 8)}</TableCell>
                    <TableCell>{request?.title ?? "Unknown request"}</TableCell>
                    <TableCell>{request ? (userById.get(request.requester_id) ?? request.requester_id.slice(0, 8)) : "—"}</TableCell>
                    <TableCell>{request ? (departmentById.get(request.department_id) ?? request.department_id.slice(0, 8)) : "—"}</TableCell>
                    <TableCell>{request ? currency(request.estimated_amount, request.currency) : "—"}</TableCell>
                    <TableCell>{item.step_name}</TableCell>
                    <TableCell>{item.scope}</TableCell>
                    <TableCell><Badge>{item.status}</Badge></TableCell>
                    <TableCell>
                      <Button asChild variant="outline" size="sm">
                        <Link href={`/requests/${item.request_id}`}>Open</Link>
                      </Button>
                    </TableCell>
                  </TableRow>
                )
              })}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </AppShell>
  )
}
