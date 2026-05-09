import { AppShell } from "@/components/app-shell"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Textarea } from "@/components/ui/textarea"

export default function AIWorkflowPage() {
  return (
    <AppShell title="AI Workflow Builder" description="Prompt-driven workflow generation with backend validation before save.">
      <div className="split-layout">
        <div className="split-pane">
          <Card>
            <CardHeader>
              <CardTitle>Prompt</CardTitle>
              <CardDescription>Describe the approval routing you want.</CardDescription>
            </CardHeader>
            <CardContent className="stack">
              <label className="field">
                <span>Prompt</span>
                <Textarea defaultValue="Create a purchase request workflow. Requests above 5 million need finance approval. Requests above 25 million need director approval. Procurement should process after approval." />
              </label>
              <div className="action-row">
                <Button type="button">Generate</Button>
                <Button type="button" variant="outline">Save workflow</Button>
              </div>
            </CardContent>
          </Card>
        </div>
        <div className="split-pane">
          <div className="stack">
            <Card>
              <CardHeader>
                <CardTitle>Generated JSON</CardTitle>
              </CardHeader>
              <CardContent>
                <Textarea
                  className="code-area"
                  readOnly
                  value={`{\n  "name": "Purchase Request Approval",\n  "type": "purchase_request",\n  "version": 1,\n  "steps": [\n    {\n      "key": "manager_approval",\n      "name": "Manager Approval",\n      "approver_role": "manager",\n      "scope": "requester_department",\n      "required": true\n    }\n  ]\n}`}
                />
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle>Mermaid preview</CardTitle>
              </CardHeader>
              <CardContent>
                <pre className="diagram">flowchart TD
  A[Submit Purchase Request] --&gt; S1[Manager Approval]
  S1 --&gt; Z[Completed]</pre>
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle>Explanation and warnings</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="muted-copy">
                  The builder explains the generated routing before an admin confirms the save.
                </p>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </AppShell>
  )
}
