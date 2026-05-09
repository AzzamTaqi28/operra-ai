import { AppShell } from "@/components/app-shell"

export default function DepartmentsPage() {
  return (
    <AppShell title="Departments" description="Department setup used for requester_department approval scopes.">
      <section className="panel">
        <h2>Department management</h2>
        <p className="muted-copy">Create and maintain the organizational structure used by workflow approvers.</p>
      </section>
    </AppShell>
  )
}
