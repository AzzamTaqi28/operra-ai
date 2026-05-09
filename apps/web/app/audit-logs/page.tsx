import { redirect } from "next/navigation"

import { AppShell } from "@/components/app-shell"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { apiGet, getToken, type ApiListResponse, type AuditLogItem, type UserListItem } from "@/lib/api"

export default async function AuditLogsPage() {
  const token = await getToken()
  if (!token) redirect("/login")

  const [logs, users] = await Promise.all([
    apiGet<ApiListResponse<AuditLogItem>>("/api/v1/audit-logs?page_size=100", token),
    apiGet<ApiListResponse<UserListItem>>("/api/v1/users?page_size=100", token),
  ])

  const userById = new Map(users.data.map((user) => [user.id, user.name]))

  return (
    <AppShell title="Audit Logs" description="Read-only event trail for meaningful mutations, approvals, exports, and AI generation.">
      <div className="toolbar">
        <div className="toolbar-filters">
          <Badge>Action</Badge>
          <Badge>Entity type</Badge>
          <Badge>Actor</Badge>
          <Badge>Date range</Badge>
        </div>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Event trail</CardTitle>
          <CardDescription>Mutations, approvals, exports, and AI generation events.</CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Timestamp</TableHead>
                <TableHead>Actor</TableHead>
                <TableHead>Action</TableHead>
                <TableHead>Entity Type</TableHead>
                <TableHead>Entity ID</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {logs.data.map((log) => (
                <TableRow key={log.id}>
                  <TableCell>{new Date(log.created_at).toLocaleString("id-ID")}</TableCell>
                  <TableCell>{log.actor_user_id ? userById.get(log.actor_user_id) ?? log.actor_user_id.slice(0, 8) : "—"}</TableCell>
                  <TableCell>{log.action}</TableCell>
                  <TableCell>{log.entity_type}</TableCell>
                  <TableCell>{log.entity_id ?? "—"}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </AppShell>
  )
}
