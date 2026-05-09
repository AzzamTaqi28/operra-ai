import { AppShell } from "@/components/app-shell"

export default function AuditLogsPage() {
  return (
    <AppShell title="Audit Logs" description="Read-only event trail for meaningful mutations, approvals, exports, and AI generation.">
      <section className="panel">
        <h2>Event trail</h2>
        <p className="muted-copy">Owner, admin, finance, and auditor roles can filter and export logs here.</p>
      </section>
    </AppShell>
  )
}
