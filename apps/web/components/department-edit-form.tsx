"use client"

import { useState, type FormEvent } from "react"
import { useRouter } from "next/navigation"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"

type Props = {
  department: {
    id: string
    name: string
    code: string
  }
}

export function DepartmentEditForm({ department }: Props) {
  const router = useRouter()
  const [error, setError] = useState<string | null>(null)
  const [pending, setPending] = useState(false)

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const formData = new FormData(event.currentTarget)
    setPending(true)
    setError(null)
    try {
      const response = await fetch(`/api/departments/${department.id}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          name: String(formData.get("name") ?? ""),
          code: String(formData.get("code") ?? ""),
        }),
      })
      const payload = await response.json()
      if (!response.ok) {
        setError(payload?.error?.message ?? "Failed to update department")
        return
      }
      router.refresh()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update department")
    } finally {
      setPending(false)
    }
  }

  return (
    <form onSubmit={onSubmit} className="stack">
      <div className="form-grid">
        <label className="field">
          <span>Name</span>
          <Input name="name" defaultValue={department.name} required />
        </label>
        <label className="field">
          <span>Code</span>
          <Input name="code" defaultValue={department.code} />
        </label>
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
