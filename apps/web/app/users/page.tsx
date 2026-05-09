import { AppShell } from "@/components/app-shell"
import { Card, Chip, Table } from "@/components/ui"

const users = [
  ["Siti", "siti@operra.ai", "Finance", "owner, admin", "active"],
  ["Ari", "ari@operra.ai", "Engineering", "requester, manager", "active"],
  ["Mira", "mira@operra.ai", "Finance", "finance", "active"],
]

export default function UsersPage() {
  return (
    <AppShell title="Users" description="Tenant-scoped user administration and role assignment.">
      <div className="toolbar">
        <div className="toolbar-filters">
          <Chip>Department</Chip>
          <Chip>Role</Chip>
          <Chip>Status</Chip>
        </div>
        <button className="button button-solid" type="button">Create user</button>
      </div>

      <Card title="User directory">
        <Table headers={["Name", "Email", "Department", "Roles", "Status"]} rows={users.map((row) => row.map((cell) => cell))} />
      </Card>
    </AppShell>
  )
}
