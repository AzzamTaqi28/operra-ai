"use client"

import { useMemo, useState } from "react"
import Link from "next/link"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"

export type ApprovalRow = {
  request_id: string
  title: string
  requester: string
  department: string
  amount: string
  step_name: string
  scope: string
  status: string
}

type Props = {
  rows: ApprovalRow[]
}

export function ApprovalsTable({ rows }: Props) {
  const [query, setQuery] = useState("")

  const filteredRows = useMemo(() => {
    const q = query.trim().toLowerCase()
    if (!q) return rows
    return rows.filter((row) =>
      [row.request_id, row.title, row.requester, row.department, row.amount, row.step_name, row.scope, row.status]
        .join(" ")
        .toLowerCase()
        .includes(q),
    )
  }, [query, rows])

  return (
    <div className="grid gap-4">
      <div className="toolbar">
        <div className="toolbar-filters">
          <Input
            value={query}
            onChange={(event) => setQuery(event.target.value)}
            placeholder="Search requests, people, roles, or status"
            className="w-[min(100%,24rem)]"
          />
        </div>
        <Badge>{filteredRows.length} visible</Badge>
      </div>

      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Request ID</TableHead>
            <TableHead>Title</TableHead>
            <TableHead>Requester</TableHead>
            <TableHead>Department</TableHead>
            <TableHead>Amount</TableHead>
            <TableHead>Current Step</TableHead>
            <TableHead>Scope</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Action</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {filteredRows.map((row) => (
            <TableRow key={row.request_id + row.step_name}>
              <TableCell>{row.request_id.slice(0, 8)}</TableCell>
              <TableCell>{row.title}</TableCell>
              <TableCell>{row.requester}</TableCell>
              <TableCell>{row.department}</TableCell>
              <TableCell>{row.amount}</TableCell>
              <TableCell>{row.step_name}</TableCell>
              <TableCell>{row.scope}</TableCell>
              <TableCell>
                <Badge>{row.status}</Badge>
              </TableCell>
              <TableCell>
                <Button asChild variant="outline" size="sm">
                  <Link href={`/requests/${row.request_id}`}>Open</Link>
                </Button>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}
