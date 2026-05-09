import { AppShell } from "@/components/app-shell"
import { Card, Table } from "@/components/ui"

const departments = [
  ["Finance", "FIN", "8"],
  ["Engineering", "ENG", "14"],
  ["Operations", "OPS", "6"],
]

export default function DepartmentsPage() {
  return (
    <AppShell title="Departments" description="Department setup used for requester_department approval scopes.">
      <Card title="Department management">
        <Table headers={["Name", "Code", "User Count"]} rows={departments.map((row) => row.map((cell) => cell))} />
      </Card>
    </AppShell>
  )
}
