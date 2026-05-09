import { AppShell } from "@/components/app-shell"
import { Card, Field, SplitLayout } from "@/components/ui"

export default function AIWorkflowPage() {
  return (
    <AppShell title="AI Workflow Builder" description="Prompt-driven workflow generation with backend validation before save.">
      <SplitLayout
        left={
          <Card title="Prompt" description="Describe the approval routing you want.">
            <Field label="Prompt">
              <textarea
                className="textarea"
                defaultValue="Create a purchase request workflow. Requests above 5 million need finance approval. Requests above 25 million need director approval. Procurement should process after approval."
              />
            </Field>
            <div className="action-row">
              <button className="button button-solid" type="button">Generate</button>
              <button className="button button-outline" type="button">Save workflow</button>
            </div>
          </Card>
        }
        right={
          <div className="stack">
            <Card title="Generated JSON">
              <textarea className="textarea code-area" readOnly value={`{\n  "name": "Purchase Request Approval",\n  "type": "purchase_request",\n  "version": 1,\n  "steps": [\n    {\n      "key": "manager_approval",\n      "name": "Manager Approval",\n      "approver_role": "manager",\n      "scope": "requester_department",\n      "required": true\n    }\n  ]\n}`} />
            </Card>
            <Card title="Mermaid preview">
              <pre className="diagram">flowchart TD
  A[Submit Purchase Request] --&gt; S1[Manager Approval]
  S1 --&gt; Z[Completed]</pre>
            </Card>
            <Card title="Explanation and warnings">
              <p className="muted-copy">
                The builder explains the generated routing before an admin confirms the save.
              </p>
            </Card>
          </div>
        }
      />
    </AppShell>
  )
}
