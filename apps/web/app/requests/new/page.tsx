import { AppShell } from "@/components/app-shell"
import { Card, Field } from "@/components/ui"

export default function NewRequestPage() {
  return (
    <AppShell title="Create Purchase Request" description="Save a draft or submit directly into the approval workflow.">
      <Card title="Request form" description="Attachment upload, comments, and submission are part of the same operational flow.">
        <div className="form-grid">
          <Field label="Title"><input className="input" placeholder="Buy laptops for new engineers" /></Field>
          <Field label="Department"><input className="input" placeholder="Engineering" /></Field>
          <Field label="Item name"><input className="input" placeholder="Laptop" /></Field>
          <Field label="Quantity"><input className="input" type="number" placeholder="3" /></Field>
          <Field label="Estimated amount"><input className="input" type="number" placeholder="45000000" /></Field>
          <Field label="Currency"><input className="input" placeholder="IDR" /></Field>
          <Field label="Urgency"><input className="input" placeholder="high" /></Field>
          <Field label="Expected date"><input className="input" type="date" /></Field>
          <Field label="Vendor name optional"><input className="input" placeholder="Vendor ABC" /></Field>
          <Field label="Notes optional"><textarea className="textarea" placeholder="Prefer business warranty" /></Field>
        </div>

        <div className="action-row">
          <button className="button button-outline" type="button">Save draft</button>
          <button className="button button-solid" type="button">Submit request</button>
        </div>
      </Card>
    </AppShell>
  )
}
