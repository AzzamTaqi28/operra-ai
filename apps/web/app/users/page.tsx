import { AppShell } from "@/components/app-shell"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

const users = [
  ["Siti", "siti@operra.ai", "Finance", "owner, admin", "active"],
  ["Ari", "ari@operra.ai", "Engineering", "requester, manager", "active"],
  ["Mira", "mira@operra.ai", "Finance", "finance", "active"],
]

export default function UsersPage() {
  return (
    <AppShell title="Users" description="Tenant-scoped user administration and role assignment.">
      <div className="toolbar">
        <div className="toolbar-filters">
          <Badge>Department</Badge>
          <Badge>Role</Badge>
          <Badge>Status</Badge>
        </div>
        <Button type="button">Create user</Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>User directory</CardTitle>
          <CardDescription>Organization-scoped users and role assignments.</CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Name</TableHead>
                <TableHead>Email</TableHead>
                <TableHead>Department</TableHead>
                <TableHead>Roles</TableHead>
                <TableHead>Status</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {users.map((row) => (
                <TableRow key={row[1]}>
                  {row.map((cell) => (
                    <TableCell key={`${row[1]}-${cell}`}>{cell}</TableCell>
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
