import { AppShell } from "@/components/app-shell"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

const departments = [
  ["Finance", "FIN", "8"],
  ["Engineering", "ENG", "14"],
  ["Operations", "OPS", "6"],
]

export default function DepartmentsPage() {
  return (
    <AppShell title="Departments" description="Department setup used for requester_department approval scopes.">
      <Card>
        <CardHeader>
          <CardTitle>Department management</CardTitle>
          <CardDescription>Departments are used in request routing and access rules.</CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Name</TableHead>
                <TableHead>Code</TableHead>
                <TableHead>User Count</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {departments.map((row) => (
                <TableRow key={row[1]}>
                  {row.map((cell) => (
                    <TableCell key={`${row[1]}-${cell}`}>{cell}</TableCell>
                  ))}
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </AppShell>
  )
}
