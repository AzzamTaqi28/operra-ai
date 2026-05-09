import { AppShell } from "@/components/app-shell"
import { Card, Field } from "@/components/ui"

export default function SettingsPage() {
  return (
    <AppShell title="Settings" description="Organization-level configuration and deployment details.">
      <Card title="Organization settings">
        <div className="form-grid">
          <Field label="Organization name"><input className="input" defaultValue="Demo Company" /></Field>
          <Field label="Organization slug"><input className="input" defaultValue="demo-company" /></Field>
          <Field label="Timezone"><input className="input" defaultValue="Asia/Jakarta" /></Field>
          <Field label="Default currency"><input className="input" defaultValue="IDR" /></Field>
        </div>
        <div className="action-row">
          <button className="button button-solid" type="button">Save settings</button>
        </div>
      </Card>
    </AppShell>
  )
}
