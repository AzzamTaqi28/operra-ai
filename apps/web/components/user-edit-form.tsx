"use client"

import { useState, type FormEvent } from "react"
import { useRouter } from "next/navigation"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"

type Props = {
  user: {
    id: string
    name: string
    email: string
    department_id?: string | null
    status: string
    roles?: string[]
  }
  departmentOptions: Array<{ id: string; name: string }>
}

export function UserEditForm({ user, departmentOptions }: Props) {
  const router = useRouter()
  const [error, setError] = useState<string | null>(null)
  const [pending, setPending] = useState(false)

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const formData = new FormData(event.currentTarget)
    setPending(true)
    setError(null)
    try {
      const response = await fetch(`/api/users/${user.id}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          name: String(formData.get("name") ?? ""),
          email: String(formData.get("email") ?? ""),
          department_id: String(formData.get("department_id") ?? ""),
          status: String(formData.get("status") ?? ""),
        }),
      })
      const payload = await response.json()
      if (!response.ok) {
        setError(payload?.error?.message ?? "Failed to update user")
        return
      }
      router.refresh()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update user")
    } finally {
      setPending(false)
    }
  }

  return (
    <form onSubmit={onSubmit} className="stack">
      <div className="form-grid">
        <label className="field">
          <span>Name</span>
          <Input name="name" defaultValue={user.name} required />
        </label>
        <label className="field">
          <span>Email</span>
          <Input name="email" type="email" defaultValue={user.email} required />
        </label>
        <label className="field">
          <span>Department</span>
          <select name="department_id" className="h-11 w-full rounded-md border border-slate-200 bg-white px-3 py-2 text-sm" defaultValue={user.department_id ?? ""}>
            <option value="">None</option>
            {departmentOptions.map((department) => (
              <option key={department.id} value={department.id}>
                {department.name}
              </option>
            ))}
          </select>
        </label>
        <label className="field">
          <span>Status</span>
          <Input name="status" defaultValue={user.status} required />
        </label>
      </div>
      <div className="field">
        <span>Roles</span>
        <p className="text-sm text-slate-500">{user.roles?.length ? user.roles.join(", ") : "No roles assigned."}</p>
        <p className="text-xs text-slate-500">Role assignment is managed separately in v0.1.</p>
      </div>
      {error ? <p className="text-sm text-red-600">{error}</p> : null}
      <div className="action-row">
        <Button type="submit" disabled={pending}>
          {pending ? "Saving..." : "Save changes"}
        </Button>
      </div>
    </form>
  )
}
