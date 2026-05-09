import { AppShell } from "@/components/app-shell"

const workflows = [
  ["Purchase Request Approval", "active", "v3"],
  ["Travel Request Approval", "draft", "v1"],
]

export default function WorkflowsPage() {
  return (
    <AppShell title="Workflows" description="JSON-first workflow definitions with validation, Mermaid preview, and versioning.">
      <section className="panel">
        <h2>Workflow versions</h2>
        <div className="table">
          {workflows.map(([name, status, version]) => (
            <div key={name} className="table-row">
              <span>{name}</span>
              <span>{status}</span>
              <span>{version}</span>
            </div>
          ))}
        </div>
      </section>
    </AppShell>
  )
}
