import { redirect } from "next/navigation"
import Link from "next/link"

import { AppShell } from "@/components/app-shell"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { apiGet, getToken, type ApiListResponse, type DepartmentListItem, type UserListItem } from "@/lib/api"

export default async function DepartmentsPage() {
  const token = await getToken()
  if (!token) redirect("/login")

  const [departments, users] = await Promise.all([
    apiGet<ApiListResponse<DepartmentListItem>>("/api/v1/departments?page_size=100", token),
    apiGet<ApiListResponse<UserListItem>>("/api/v1/users?page_size=100", token),
  ])

  const countByDepartmentId = new Map<string, number>()
  for (const user of users.data) {
    if (!user.department_id) continue
    countByDepartmentId.set(user.department_id, (countByDepartmentId.get(user.department_id) ?? 0) + 1)
  }

  return (
    <AppShell title="Departments" description="Department setup used for requester_department approval scopes.">
      <Card>
        <CardHeader>
          <CardTitle>Department management</CardTitle>
          <CardDescription>Departments are used in request routing and access rules.</CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Name</TableHead>
                <TableHead>Code</TableHead>
                <TableHead>User Count</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Action</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {departments.data.map((department) => (
                <TableRow key={department.id}>
                  <TableCell>{department.name}</TableCell>
                  <TableCell>{department.code || "—"}</TableCell>
                  <TableCell>{countByDepartmentId.get(department.id) ?? 0}</TableCell>
                  <TableCell><Badge>{department.status ?? "active"}</Badge></TableCell>
                  <TableCell>
                    <Button asChild variant="outline" size="sm">
                      <Link href={`/departments/${department.id}`}>Edit</Link>
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
