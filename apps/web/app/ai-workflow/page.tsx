import { AppShell } from "@/components/app-shell"
import { AIWorkflowForm } from "@/components/ai-workflow-form"

export default function AIWorkflowPage() {
  return (
    <AppShell title="AI Workflow Builder" description="Prompt-driven workflow generation with backend validation before save.">
      <AIWorkflowForm />
    </AppShell>
  )
}
