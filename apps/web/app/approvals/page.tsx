import { AppShell } from "@/components/app-shell"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

const approvals = [
  ["PR-1024", "Laptop refresh", "Siti", "Finance", "IDR 45,000,000", "Director Approval", "1h ago"],
  ["PR-1019", "Security training", "Budi", "Operations", "IDR 6,500,000", "Manager Approval", "3h ago"],
]

export default function ApprovalsPage() {
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
                <TableHead>Waiting Since</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {approvals.map((row) => (
                <TableRow key={row[0]}>
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
