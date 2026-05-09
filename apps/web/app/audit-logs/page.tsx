import { AppShell } from "@/components/app-shell"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

const logs = [
  ["2026-05-08 10:11", "Siti", "request.submitted", "purchase_request", "PR-1024"],
  ["2026-05-08 10:25", "Mira", "approval.approved", "approval_step_instance", "asi-022"],
  ["2026-05-08 10:31", "Siti", "csv.exported", "export", "-"],
]

export default function AuditLogsPage() {
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
              {logs.map((row) => (
                <TableRow key={`${row[0]}-${row[4]}`}>
                  {row.map((cell) => (
                    <TableCell key={`${row[0]}-${cell}`}>{cell}</TableCell>
                  ))}
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </AppShell>
  )
}
