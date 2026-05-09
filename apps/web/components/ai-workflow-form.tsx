"use client"

import { useState, type FormEvent } from "react"
import { useRouter } from "next/navigation"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Textarea } from "@/components/ui/textarea"
import type { WorkflowGenerationResult } from "@/lib/api"

export function AIWorkflowForm() {
  const router = useRouter()
  const [prompt, setPrompt] = useState(
    "Create a purchase request workflow. Requests above 5 million need finance approval. Requests above 25 million need director approval. Procurement should process after approval.",
  )
  const [generated, setGenerated] = useState<WorkflowGenerationResult | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [generatePending, setGeneratePending] = useState(false)
  const [savePending, setSavePending] = useState(false)

  async function generate(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setGeneratePending(true)
    setError(null)
    try {
      const response = await fetch("/api/workflows/generate-with-ai", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ prompt }),
      })
      const payload = await response.json()
      if (!response.ok) {
        setError(payload?.error?.message ?? "AI generation failed")
        return
      }
      setGenerated(payload.data as WorkflowGenerationResult)
    } catch (err) {
      setError(err instanceof Error ? err.message : "AI generation failed")
    } finally {
      setGeneratePending(false)
    }
  }

  async function saveWorkflow() {
    if (!generated) return
    setSavePending(true)
    setError(null)
    try {
      const response = await fetch("/api/workflows", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          name: generated.workflow_json.name,
          type: generated.workflow_json.type,
          config_json: generated.workflow_json,
        }),
      })
      const payload = await response.json()
      if (!response.ok) {
        setError(payload?.error?.message ?? "Failed to save workflow")
        return
      }
      router.push("/workflows")
      router.refresh()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to save workflow")
    } finally {
      setSavePending(false)
    }
  }

  return (
    <div className="split-layout">
      <div className="split-pane">
        <Card>
          <CardHeader>
            <CardTitle>Prompt</CardTitle>
            <CardDescription>Describe the approval routing you want.</CardDescription>
          </CardHeader>
          <CardContent className="stack">
            <form onSubmit={generate} className="stack">
              <label className="field">
                <span>Prompt</span>
                <Textarea value={prompt} onChange={(event) => setPrompt(event.target.value)} />
              </label>
              {error ? <p className="text-sm text-red-600">{error}</p> : null}
              <div className="action-row">
                <Button type="submit" disabled={generatePending}>
                  {generatePending ? "Generating..." : "Generate"}
                </Button>
                <Button type="button" variant="outline" disabled={!generated || savePending} onClick={saveWorkflow}>
                  {savePending ? "Saving..." : "Save workflow"}
                </Button>
              </div>
            </form>
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
                value={
                  generated
                    ? JSON.stringify(generated.workflow_json, null, 2)
                    : `{\n  "name": "Purchase Request Approval",\n  "type": "purchase_request",\n  "version": 1,\n  "steps": []\n}`
                }
              />
            </CardContent>
          </Card>
          <Card>
            <CardHeader>
              <CardTitle>Mermaid preview</CardTitle>
            </CardHeader>
            <CardContent>
              <pre className="diagram">{generated?.mermaid_diagram ?? "flowchart TD\n  A[Submit Purchase Request] --> S1[Manager Approval]\n  S1 --> Z[Completed]"}</pre>
            </CardContent>
          </Card>
          <Card>
            <CardHeader>
              <CardTitle>Explanation and warnings</CardTitle>
            </CardHeader>
            <CardContent className="stack">
              <p className="muted-copy">{generated?.explanation ?? "The builder explains the generated routing before an admin confirms the save."}</p>
              {generated?.warnings?.length ? (
                <div className="grid gap-2">
                  {generated.warnings.map((warning) => (
                    <p key={warning} className="text-sm text-amber-700">{warning}</p>
                  ))}
                </div>
              ) : null}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}
