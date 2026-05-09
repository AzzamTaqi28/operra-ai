import Link from "next/link"
import { redirect } from "next/navigation"

import { AppShell } from "@/components/app-shell"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import {
  apiGet,
  getToken,
  type ApiListResponse,
  type ApprovalQueueItem,
  type DepartmentListItem,
  type PurchaseRequest,
  type UserListItem,
} from "@/lib/api"

function currency(amount: number, code = "IDR") {
  return new Intl.NumberFormat("id-ID", { style: "currency", currency: code, maximumFractionDigits: 0 }).format(amount)
}

function daysToDuration(hours: number) {
  if (!Number.isFinite(hours) || hours < 0) return "—"
  const totalMinutes = Math.round(hours * 60)
  const d = Math.floor(totalMinutes / (60 * 24))
  const h = Math.floor((totalMinutes % (60 * 24)) / 60)
  const m = totalMinutes % 60
  if (d > 0) return `${d}d ${h}h`
  if (h > 0) return `${h}h ${m}m`
  return `${m}m`
}

function averageApprovalTime(items: PurchaseRequest[]) {
  const durations = items
    .filter((item) => item.completed_at && item.submitted_at)
    .map((item) => {
      const submitted = new Date(item.submitted_at as string).getTime()
      const completed = new Date(item.completed_at as string).getTime()
      return (completed - submitted) / (1000 * 60 * 60)
    })
    .filter((value) => value > 0)

  if (durations.length === 0) return "—"
  const average = durations.reduce((sum, value) => sum + value, 0) / durations.length
  return daysToDuration(average)
}

export default async function DashboardPage() {
  const token = await getToken()
  if (!token) redirect("/login")

  const [
    allRequests,
    approvedRequests,
    rejectedRequests,
    completedRequests,
    pendingApprovals,
    users,
    departments,
  ] = await Promise.all([
    apiGet<ApiListResponse<PurchaseRequest>>("/api/v1/purchase-requests?page_size=100", token),
    apiGet<ApiListResponse<PurchaseRequest>>("/api/v1/purchase-requests?status=approved&page_size=100", token),
    apiGet<ApiListResponse<PurchaseRequest>>("/api/v1/purchase-requests?status=rejected&page_size=100", token),
    apiGet<ApiListResponse<PurchaseRequest>>("/api/v1/purchase-requests?status=completed&page_size=100", token),
    apiGet<ApiListResponse<ApprovalQueueItem>>("/api/v1/approvals/pending?page_size=100", token),
    apiGet<ApiListResponse<UserListItem>>("/api/v1/users?page_size=100", token),
    apiGet<ApiListResponse<DepartmentListItem>>("/api/v1/departments?page_size=100", token),
  ])

  const recentRequests = allRequests.data.slice(0, 5)
  const avgApproval = averageApprovalTime(completedRequests.data)
  const approvedThisMonth = approvedRequests.pagination.total + completedRequests.pagination.total
  const requesterById = new Map(users.data.map((user) => [user.id, user.name]))
  const departmentById = new Map(departments.data.map((department) => [department.id, department.name]))

  return (
    <AppShell title="Dashboard" description="Operational overview for purchase approvals, workflows, and audit trails.">
      <section className="card-grid">
        <Card>
          <CardHeader>
            <CardDescription>Total purchase requests</CardDescription>
            <CardTitle>{allRequests.pagination.total}</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader>
            <CardDescription>Pending approvals</CardDescription>
            <CardTitle>{pendingApprovals.pagination.total}</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader>
            <CardDescription>Approved this month</CardDescription>
            <CardTitle>{approvedThisMonth}</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader>
            <CardDescription>Rejected this month</CardDescription>
            <CardTitle>{rejectedRequests.pagination.total}</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader>
            <CardDescription>Average approval time</CardDescription>
            <CardTitle>{avgApproval}</CardTitle>
          </CardHeader>
        </Card>
      </section>

      <div className="stack">
        <Card>
          <CardHeader>
            <CardTitle>Recent requests</CardTitle>
            <CardDescription>Status should be obvious at a glance.</CardDescription>
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
                  <TableHead>Status</TableHead>
                  <TableHead>Current Step</TableHead>
                  <TableHead>Created At</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {recentRequests.map((item) => (
                  <TableRow key={item.id}>
                    <TableCell>
                      <Link href={`/requests/${item.id}`} className="font-medium text-[var(--accent)] hover:underline">
                        {item.id.slice(0, 8)}
                      </Link>
                    </TableCell>
                    <TableCell>{item.title}</TableCell>
                    <TableCell>{requesterById.get(item.requester_id) ?? item.requester_id.slice(0, 8)}</TableCell>
                    <TableCell>{departmentById.get(item.department_id) ?? item.department_id.slice(0, 8)}</TableCell>
                    <TableCell>{currency(item.estimated_amount, item.currency)}</TableCell>
                    <TableCell>
                      <Badge>{item.status}</Badge>
                    </TableCell>
                    <TableCell>{item.current_step_instance_id ? item.current_step_instance_id.slice(0, 8) : "Draft"}</TableCell>
                    <TableCell>{new Date(item.created_at).toLocaleDateString("id-ID")}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>My pending approvals</CardTitle>
            <CardDescription>Current actionable items for the signed-in user.</CardDescription>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Request ID</TableHead>
                  <TableHead>Current Step</TableHead>
                  <TableHead>Waiting Since</TableHead>
                  <TableHead>Action</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {pendingApprovals.data.map((item) => (
                  <TableRow key={item.step_id}>
                    <TableCell>{item.request_id.slice(0, 8)}</TableCell>
                    <TableCell>{item.step_name}</TableCell>
                    <TableCell>Recently assigned</TableCell>
                    <TableCell>
                      <Button asChild variant="outline" size="sm">
                        <Link href={`/requests/${item.request_id}`}>Open</Link>
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </div>
    </AppShell>
  )
}
