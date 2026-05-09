import Link from "next/link"
import { redirect } from "next/navigation"

import { AppShell } from "@/components/app-shell"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { apiGet, getToken, type ApiListResponse, type DepartmentListItem, type PurchaseRequest, type UserListItem } from "@/lib/api"

function currency(amount: number, code = "IDR") {
  return new Intl.NumberFormat("id-ID", { style: "currency", currency: code, maximumFractionDigits: 0 }).format(amount)
}

export default async function RequestsPage() {
  const token = await getToken()
  if (!token) redirect("/login")

  const [requests, users, departments] = await Promise.all([
    apiGet<ApiListResponse<PurchaseRequest>>("/api/v1/purchase-requests?page_size=20", token),
    apiGet<ApiListResponse<UserListItem>>("/api/v1/users?page_size=100", token),
    apiGet<ApiListResponse<DepartmentListItem>>("/api/v1/departments?page_size=100", token),
  ])

  const requesterById = new Map(users.data.map((user) => [user.id, user.name]))
  const departmentById = new Map(departments.data.map((department) => [department.id, department.name]))

  return (
    <AppShell title="Purchase Requests" description="Drafts, submissions, approval timelines, comments, and attachments.">
      <div className="toolbar">
        <div className="toolbar-filters">
          <Badge>Status</Badge>
          <Badge>Department</Badge>
          <Badge>Date range</Badge>
          <Badge>Search</Badge>
        </div>
        <div className="toolbar-actions">
          <Button asChild>
            <Link href="/requests/new">Create request</Link>
          </Button>
          <Button asChild variant="outline">
            <Link href="/exports">Export CSV</Link>
          </Button>
        </div>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Request list</CardTitle>
          <CardDescription>Request ID, title, requester, department, amount, status, and current step.</CardDescription>
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
              {requests.data.map((item) => (
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
    </AppShell>
  )
}
