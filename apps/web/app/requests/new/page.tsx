import { AppShell } from "@/components/app-shell"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Textarea } from "@/components/ui/textarea"

export default function NewRequestPage() {
  return (
    <AppShell title="Create Purchase Request" description="Save a draft or submit directly into the approval workflow.">
      <Card>
        <CardHeader>
          <CardTitle>Request form</CardTitle>
          <CardDescription>Attachment upload, comments, and submission are part of the same operational flow.</CardDescription>
        </CardHeader>
        <CardContent className="stack">
          <div className="form-grid">
            <label className="field">
              <span>Title</span>
              <Input placeholder="Buy laptops for new engineers" />
            </label>
            <label className="field">
              <span>Department</span>
              <Input placeholder="Engineering" />
            </label>
            <label className="field">
              <span>Item name</span>
              <Input placeholder="Laptop" />
            </label>
            <label className="field">
              <span>Quantity</span>
              <Input type="number" placeholder="3" />
            </label>
            <label className="field">
              <span>Estimated amount</span>
              <Input type="number" placeholder="45000000" />
            </label>
            <label className="field">
              <span>Currency</span>
              <Input placeholder="IDR" />
            </label>
            <label className="field">
              <span>Urgency</span>
              <Input placeholder="high" />
            </label>
            <label className="field">
              <span>Expected date</span>
              <Input type="date" />
            </label>
            <label className="field">
              <span>Vendor name optional</span>
              <Input placeholder="Vendor ABC" />
            </label>
            <label className="field">
              <span>Notes optional</span>
              <Textarea placeholder="Prefer business warranty" />
            </label>
          </div>

          <div className="action-row">
            <Button type="button" variant="outline">Save draft</Button>
            <Button type="button">Submit request</Button>
          </div>
        </CardContent>
      </Card>
    </AppShell>
  )
}
