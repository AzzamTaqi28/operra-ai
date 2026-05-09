import { redirect } from "next/navigation"

import { AppShell } from "@/components/app-shell"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { RequestActionPanel } from "@/components/request-action-panel"
import { RequestCommentForm } from "@/components/request-comment-form"
import {
  apiGet,
  getToken,
  type ApiListResponse,
  type DepartmentListItem,
  type PurchaseRequestDetail,
  type PurchaseRequestComment,
  type UserListItem,
} from "@/lib/api"

function currency(amount: number, code = "IDR") {
  return new Intl.NumberFormat("id-ID", { style: "currency", currency: code, maximumFractionDigits: 0 }).format(amount)
}

export default async function RequestDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const token = await getToken()
  if (!token) redirect("/login")

  const { id } = await params

  const [detail, comments, users, departments] = await Promise.all([
    apiGet<{ data: PurchaseRequestDetail }>(`/api/v1/purchase-requests/${id}`, token),
    apiGet<{ data: PurchaseRequestComment[] }>(`/api/v1/purchase-requests/${id}/comments`, token).catch(() => ({ data: [] })),
    apiGet<ApiListResponse<UserListItem>>("/api/v1/users?page_size=100", token),
    apiGet<ApiListResponse<DepartmentListItem>>("/api/v1/departments?page_size=100", token),
  ])

  const item = detail.data
  const requesterName = users.data.find((user) => user.id === item.requester_id)?.name ?? item.requester_id.slice(0, 8)
  const departmentName = departments.data.find((department) => department.id === item.department_id)?.name ?? item.department_id.slice(0, 8)

  return (
    <AppShell title={`Purchase Request ${item.id.slice(0, 8)}`} description="Summary, attachments, approval state, comments, and audit history.">
      <div className="grid gap-4 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Request summary</CardTitle>
            <CardDescription>Current status and routing context.</CardDescription>
          </CardHeader>
          <CardContent className="grid gap-4">
            <div className="flex flex-wrap gap-2">
              <Badge>{item.status}</Badge>
              <Badge>{departmentName}</Badge>
              <Badge>{currency(item.estimated_amount, item.currency)}</Badge>
            </div>
            <p className="text-sm text-slate-600">Requester: {requesterName}. This request is locked to the workflow version active at submission time.</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Approval actions</CardTitle>
            <CardDescription>Approver controls stay compact and explicit.</CardDescription>
          </CardHeader>
          <CardContent className="grid gap-4">
            <RequestActionPanel requestId={item.id} />
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-4">
        <Card>
          <CardHeader>
            <CardTitle>Request fields</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 md:grid-cols-2">
              <div>
                <p className="text-sm text-slate-500">Title</p>
                <p className="font-medium">{item.title}</p>
              </div>
              <div>
                <p className="text-sm text-slate-500">Item name</p>
                <p className="font-medium">{item.item_name}</p>
              </div>
              <div className="md:col-span-2">
                <p className="text-sm text-slate-500">Description</p>
                <p className="font-medium">{item.description}</p>
              </div>
              <div>
                <p className="text-sm text-slate-500">Quantity</p>
                <p className="font-medium">{item.quantity}</p>
              </div>
              <div>
                <p className="text-sm text-slate-500">Urgency</p>
                <p className="font-medium">{item.urgency}</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Attachments</CardTitle>
          </CardHeader>
          <CardContent>
            {item.attachments.length === 0 ? (
              <p className="text-sm text-slate-500">No attachments uploaded yet.</p>
            ) : (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>File</TableHead>
                    <TableHead>Size</TableHead>
                    <TableHead>Type</TableHead>
                    <TableHead>Created At</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {item.attachments.map((attachment) => (
                    <TableRow key={attachment.id}>
                      <TableCell>{attachment.file_name}</TableCell>
                      <TableCell>{attachment.file_size}</TableCell>
                      <TableCell>{attachment.mime_type}</TableCell>
                      <TableCell>{new Date(attachment.created_at).toLocaleDateString("id-ID")}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Approval timeline</CardTitle>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Step key</TableHead>
                  <TableHead>Step name</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Role</TableHead>
                  <TableHead>Scope</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {item.approval_steps.map((step) => (
                  <TableRow key={step.id}>
                    <TableCell>{step.step_key}</TableCell>
                    <TableCell>{step.step_name}</TableCell>
                    <TableCell><Badge>{step.status}</Badge></TableCell>
                    <TableCell>{step.approver_role_key}</TableCell>
                    <TableCell>{step.scope}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Comments</CardTitle>
          </CardHeader>
          <CardContent className="grid gap-4">
            <RequestCommentForm requestId={item.id} />
            {comments.data.length === 0 ? (
              <p className="text-sm text-slate-500">No comments yet.</p>
            ) : (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Actor</TableHead>
                    <TableHead>Comment</TableHead>
                    <TableHead>Created At</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {comments.data.map((comment) => (
                  <TableRow key={comment.id}>
                      <TableCell>{users.data.find((user) => user.id === comment.actor_user_id)?.name ?? comment.actor_user_id.slice(0, 8)}</TableCell>
                      <TableCell>{comment.body}</TableCell>
                      <TableCell>{new Date(comment.created_at).toLocaleDateString("id-ID")}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            )}
          </CardContent>
        </Card>
      </div>
    </AppShell>
  )
}
