"use client"

import { useState, type FormEvent } from "react"
import { useRouter } from "next/navigation"

import { Button } from "@/components/ui/button"
import { Textarea } from "@/components/ui/textarea"
import type { WorkflowListItem, WorkflowVersionSummary } from "@/lib/api"

type Props = {
  workflow: WorkflowListItem
  versions: WorkflowVersionSummary[]
}

export function WorkflowEditorForm({ workflow, versions }: Props) {
  const router = useRouter()
  const [jsonValue, setJsonValue] = useState(() =>
    JSON.stringify(workflow.active_version?.config_json ?? {}, null, 2),
  )
  const [error, setError] = useState<string | null>(null)
  const [pending, setPending] = useState(false)
  const [activeVersion, setActiveVersion] = useState<string | null>(null)

  async function saveVersion(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setPending(true)
    setError(null)
    try {
      const response = await fetch(`/api/workflows/${workflow.id}/versions`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ config_json: JSON.parse(jsonValue) }),
      })
      const payload = await response.json()
      if (!response.ok) {
        setError(payload?.error?.message ?? "Failed to save workflow version")
        return
      }
      router.refresh()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to save workflow version")
    } finally {
      setPending(false)
    }
  }

  async function activateVersion(versionId: string) {
    setActiveVersion(versionId)
    setError(null)
    try {
      const response = await fetch(`/api/workflows/${workflow.id}/versions/${versionId}/activate`, {
        method: "POST",
      })
      const payload = await response.json()
      if (!response.ok) {
        setError(payload?.error?.message ?? "Failed to activate workflow version")
        return
      }
      router.refresh()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to activate workflow version")
    } finally {
      setActiveVersion(null)
    }
  }

  return (
    <div className="grid gap-4">
      <form onSubmit={saveVersion} className="grid gap-4">
        <Textarea className="code-area" value={jsonValue} onChange={(event) => setJsonValue(event.target.value)} />
        {error ? <p className="text-sm text-red-600">{error}</p> : null}
        <div className="action-row">
          <Button type="submit" disabled={pending}>
            {pending ? "Saving..." : "Save new version"}
          </Button>
        </div>
      </form>

      <div className="grid gap-3">
        {versions.map((version) => (
          <div key={version.id} className="flex items-center justify-between gap-4 rounded-xl border border-slate-200 bg-white px-4 py-3">
            <div>
              <p className="font-medium">Version {version.version_number}</p>
              <p className="text-sm text-slate-500">{new Date(version.created_at).toLocaleString("id-ID")}</p>
            </div>
            <Button
              type="button"
              variant={workflow.active_version_id === version.id ? "outline" : "default"}
              onClick={() => activateVersion(version.id)}
              disabled={activeVersion === version.id}
            >
              {workflow.active_version_id === version.id ? "Active" : activeVersion === version.id ? "Activating..." : "Activate"}
            </Button>
          </div>
        ))}
      </div>
    </div>
  )
}
