import { redirect } from "next/navigation"
import Link from "next/link"

import { AppShell } from "@/components/app-shell"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { apiGet, getToken, type ApiListResponse, type DepartmentListItem, type UserListItem } from "@/lib/api"

export default async function UsersPage() {
  const token = await getToken()
  if (!token) redirect("/login")

  const [users, departments] = await Promise.all([
    apiGet<ApiListResponse<UserListItem>>("/api/v1/users?page_size=100", token),
    apiGet<ApiListResponse<DepartmentListItem>>("/api/v1/departments?page_size=100", token),
  ])

  const departmentById = new Map(departments.data.map((department) => [department.id, department.name]))

  return (
    <AppShell title="Users" description="Tenant-scoped user administration and role assignment.">
      <div className="toolbar">
        <div className="toolbar-filters">
          <Badge>Department</Badge>
          <Badge>Role</Badge>
          <Badge>Status</Badge>
        </div>
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
                <TableHead>Action</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {users.data.map((user) => (
                <TableRow key={user.id}>
                  <TableCell>{user.name}</TableCell>
                  <TableCell>{user.email}</TableCell>
                  <TableCell>{user.department_id ? departmentById.get(user.department_id) ?? user.department_id.slice(0, 8) : "—"}</TableCell>
                  <TableCell>{(user.roles ?? []).length > 0 ? (user.roles ?? []).join(", ") : "—"}</TableCell>
                  <TableCell><Badge>{user.status}</Badge></TableCell>
                  <TableCell>
                    <Button asChild variant="outline" size="sm">
                      <Link href={`/users/${user.id}`}>Edit</Link>
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </AppShell>
  )
}
