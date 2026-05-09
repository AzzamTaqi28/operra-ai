import { AppShell } from "@/components/app-shell"

export default function ExportsPage() {
  return (
    <AppShell title="Exports" description="CSV exports for purchase requests, approval history, and audit logs.">
      <section className="panel">
        <h2>Available exports</h2>
        <p className="muted-copy">Export endpoints are audited and respect tenant permissions.</p>
      </section>
    </AppShell>
  )
}
