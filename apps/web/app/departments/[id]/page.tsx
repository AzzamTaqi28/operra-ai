import { redirect } from "next/navigation"

import { AppShell } from "@/components/app-shell"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { apiGet, getToken, type DepartmentListItem } from "@/lib/api"
import { DepartmentEditForm } from "@/components/department-edit-form"

export default async function DepartmentDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const token = await getToken()
  if (!token) redirect("/login")

  const { id } = await params
  const department = await apiGet<{ data: DepartmentListItem }>(`/api/v1/departments/${id}`, token)

  return (
    <AppShell title="Edit Department" description="Update department code and name.">
      <Card>
        <CardHeader>
          <CardTitle>Department profile</CardTitle>
          <CardDescription>Edit the department metadata used in routing and access rules.</CardDescription>
        </CardHeader>
        <CardContent>
          <DepartmentEditForm department={department.data} />
        </CardContent>
      </Card>
    </AppShell>
  )
}
