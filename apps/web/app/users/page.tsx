import { AppShell } from "@/components/app-shell"

export default function UsersPage() {
  return (
    <AppShell title="Users" description="Tenant-scoped user administration and role assignment.">
      <section className="panel">
        <h2>User directory</h2>
        <p className="muted-copy">Owner/admin screens for user onboarding, department assignment, and role mapping.</p>
      </section>
    </AppShell>
  )
}
