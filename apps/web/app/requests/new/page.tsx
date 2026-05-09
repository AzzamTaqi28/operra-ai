import { redirect } from "next/navigation"

import { AppShell } from "@/components/app-shell"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { NewRequestForm } from "@/components/new-request-form"
import { apiGet, getToken, type ApiListResponse, type DepartmentListItem } from "@/lib/api"

export default async function NewRequestPage() {
  const token = await getToken()
  if (!token) redirect("/login")

  const departments = await apiGet<ApiListResponse<DepartmentListItem>>("/api/v1/departments?page_size=100", token)

  return (
    <AppShell title="Create Purchase Request" description="Save a draft or submit directly into the approval workflow.">
      <Card>
        <CardHeader>
          <CardTitle>Request form</CardTitle>
          <CardDescription>Attachment upload, comments, and submission are part of the same operational flow.</CardDescription>
        </CardHeader>
        <CardContent className="stack">
          <NewRequestForm departments={departments.data} />
        </CardContent>
      </Card>
    </AppShell>
  )
}
