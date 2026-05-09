import { AppShell } from "@/components/app-shell"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"

export default function SettingsPage() {
  return (
    <AppShell title="Settings" description="Organization-level configuration and deployment details.">
      <Card>
        <CardHeader>
          <CardTitle>Organization settings</CardTitle>
        </CardHeader>
        <CardContent className="stack">
          <div className="form-grid">
            <label className="field">
              <span>Organization name</span>
              <Input defaultValue="Demo Company" />
            </label>
            <label className="field">
              <span>Organization slug</span>
              <Input defaultValue="demo-company" />
            </label>
            <label className="field">
              <span>Timezone</span>
              <Input defaultValue="Asia/Jakarta" />
            </label>
            <label className="field">
              <span>Default currency</span>
              <Input defaultValue="IDR" />
            </label>
          </div>
          <div className="action-row">
            <Button type="button">Save settings</Button>
          </div>
        </CardContent>
      </Card>
    </AppShell>
  )
}
