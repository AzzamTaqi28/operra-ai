import { redirect } from "next/navigation"

import { AppShell } from "@/components/app-shell"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { apiGet, getToken, type ApiListResponse, type DepartmentListItem, type UserListItem } from "@/lib/api"
import { UserEditForm } from "@/components/user-edit-form"

export default async function UserDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const token = await getToken()
  if (!token) redirect("/login")

  const { id } = await params

  const [userRes, departments] = await Promise.all([
    apiGet<{ data: UserListItem }>(`/api/v1/users/${id}`, token),
    apiGet<ApiListResponse<DepartmentListItem>>("/api/v1/departments?page_size=100", token),
  ])

  return (
    <AppShell title="Edit User" description="Update tenant-scoped user profile information.">
      <Card>
        <CardHeader>
          <CardTitle>User profile</CardTitle>
          <CardDescription>Edit the basic user fields. Role management remains separate.</CardDescription>
        </CardHeader>
        <CardContent>
          <UserEditForm user={userRes.data} departmentOptions={departments.data} />
        </CardContent>
      </Card>
    </AppShell>
  )
}
