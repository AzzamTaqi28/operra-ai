import { AppShell } from "@/components/app-shell"

export default function AIWorkflowPage() {
  return (
    <AppShell title="AI Workflow Builder" description="Prompt-driven workflow generation with backend validation before save.">
      <section className="panel">
        <h2>Prompt</h2>
        <p className="muted-copy">
          Describe the approval routing you want, then review the generated JSON, Mermaid diagram, and validation warnings before saving.
        </p>
      </section>
    </AppShell>
  )
}
