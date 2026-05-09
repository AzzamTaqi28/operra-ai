import { redirect } from "next/navigation"

import { AppShell } from "@/components/app-shell"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { apiGet, getToken, type CurrentUser } from "@/lib/api"

export default async function SettingsPage() {
  const token = await getToken()
  if (!token) redirect("/login")

  const me = await apiGet<{ data: CurrentUser }>("/api/v1/auth/me", token)

  return (
    <AppShell title="Settings" description="Organization-level configuration and deployment details.">
      <Card>
        <CardHeader>
          <CardTitle>Organization settings</CardTitle>
          <CardDescription>Live session context from the signed-in user.</CardDescription>
        </CardHeader>
        <CardContent className="stack">
          <div className="form-grid">
            <label className="field">
              <span>Your name</span>
              <Input defaultValue={me.data.name} />
            </label>
            <label className="field">
              <span>Your email</span>
              <Input defaultValue={me.data.email} />
            </label>
            <label className="field">
              <span>Organization ID</span>
              <Input defaultValue={me.data.organization_id} />
            </label>
            <label className="field">
              <span>Roles</span>
              <Input defaultValue={me.data.roles.join(", ")} />
            </label>
          </div>
          <div className="action-row">
            <Button type="button">Save settings</Button>
          </div>
        </CardContent>
      </Card>

      <Card className="mt-4">
        <CardHeader>
          <CardTitle>Session context</CardTitle>
        </CardHeader>
        <CardContent className="flex flex-wrap gap-2">
          <Badge>{me.data.status}</Badge>
          {me.data.roles.map((role) => (
            <Badge key={role}>{role}</Badge>
          ))}
        </CardContent>
      </Card>
    </AppShell>
  )
}
