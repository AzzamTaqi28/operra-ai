import { AppShell } from "@/components/app-shell"

const stats = [
  ["Open requests", "24"],
  ["Pending approvals", "8"],
  ["Active workflows", "3"],
  ["Audit events today", "116"],
]

export default function DashboardPage() {
  return (
    <AppShell title="Dashboard" description="Operational overview for purchase approvals, workflows, and audit trails.">
      <section className="card-grid">
        {stats.map(([label, value]) => (
          <article key={label} className="stat-card">
            <p>{label}</p>
            <strong>{value}</strong>
          </article>
        ))}
      </section>
    </AppShell>
  )
}
