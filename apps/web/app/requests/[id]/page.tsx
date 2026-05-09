import { AppShell } from "@/components/app-shell"
import { Card, Chip, Field, SplitLayout, Table } from "@/components/ui"

const approvalTimeline = [
  ["manager_approval", "Manager Approval", "approved", "Ari", "2026-05-07 08:12"],
  ["finance_approval", "Finance Approval", "approved", "Mira", "2026-05-07 14:55"],
  ["director_approval", "Director Approval", "pending", "-", "Waiting"],
]

const comments = [
  ["Mira", "Please confirm vendor warranty terms before approval."],
  ["Siti", "Updated quotation attached in the latest revision."],
]

export default function RequestDetailPage() {
  return (
    <AppShell title="Purchase Request PR-1024" description="Summary, attachments, approval state, comments, and audit history.">
      <SplitLayout
        left={
          <Card title="Request summary" description="Current status and routing context.">
            <div className="summary-grid">
              <Chip>in_review</Chip>
              <Chip>Director Approval</Chip>
              <Chip>IDR 45,000,000</Chip>
              <Chip>Engineering</Chip>
            </div>
            <p className="muted-copy">
              Requester: Siti. This request is locked to the workflow version active at submission time.
            </p>
          </Card>
        }
        right={
          <Card title="Approval actions" description="Approver controls stay compact and explicit.">
            <div className="action-row">
              <button className="button button-solid" type="button">Approve</button>
              <button className="button button-outline" type="button">Reject</button>
              <button className="button button-outline" type="button">Request revision</button>
            </div>
            <Field label="Comment"><textarea className="textarea" placeholder="Approved. Proceed." /></Field>
          </Card>
        }
      />

      <div className="stack">
        <Card title="Request fields">
          <div className="detail-grid">
            <Field label="Title"><input className="input" value="Laptop refresh" readOnly /></Field>
            <Field label="Item name"><input className="input" value="Laptop" readOnly /></Field>
            <Field label="Description"><textarea className="textarea" value="Need devices for new engineers" readOnly /></Field>
            <Field label="Attachments"><div className="file-list">quotation.pdf · signed PO draft.pdf</div></Field>
          </div>
        </Card>

        <Card title="Approval timeline">
          <Table headers={["Step key", "Step name", "Status", "Actor", "Updated at"]} rows={approvalTimeline.map((row) => row.map((cell) => cell))} />
        </Card>

        <Card title="Comments">
          <Table headers={["Author", "Comment"]} rows={comments.map((row) => row.map((cell) => cell))} />
        </Card>

        <Card title="Audit history">
          <p className="muted-copy">Audit trail is visible to authorized roles and backed by the API audit log endpoint.</p>
        </Card>
      </div>
    </AppShell>
  )
}
