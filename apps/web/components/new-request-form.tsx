"use client"

import { useState, type FormEvent } from "react"
import { useRouter } from "next/navigation"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Textarea } from "@/components/ui/textarea"

type DepartmentOption = {
  id: string
  name: string
  code: string
}

type Props = {
  departments: DepartmentOption[]
}

const inputClassName =
  "flex h-11 w-full rounded-md border border-slate-200 bg-white px-3 py-2 text-sm ring-offset-white placeholder:text-slate-400 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--accent)] focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"

export function NewRequestForm({ departments }: Props) {
  const router = useRouter()
  const [error, setError] = useState<string | null>(null)
  const [pending, setPending] = useState(false)

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const submitter = (event.nativeEvent as SubmitEvent).submitter as HTMLButtonElement | null
    const mode = submitter?.value === "submit" ? "submit" : "draft"
    const formData = new FormData(event.currentTarget)

    const payload = {
      department_id: String(formData.get("department_id") ?? ""),
      title: String(formData.get("title") ?? ""),
      item_name: String(formData.get("item_name") ?? ""),
      description: String(formData.get("description") ?? ""),
      quantity: Number(formData.get("quantity") ?? 0),
      estimated_amount: Number(formData.get("estimated_amount") ?? 0),
      currency: String(formData.get("currency") ?? "IDR"),
      urgency: String(formData.get("urgency") ?? "normal"),
      expected_date: String(formData.get("expected_date") ?? ""),
      vendor_name: String(formData.get("vendor_name") ?? ""),
      notes: String(formData.get("notes") ?? ""),
    }

    setPending(true)
    setError(null)
    try {
      const createResponse = await fetch("/api/purchase-requests", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      })
      const createBody = await createResponse.json()
      if (!createResponse.ok) {
        setError(createBody?.error?.message ?? "Failed to create request")
        return
      }

      const requestId = createBody?.data?.id as string | undefined
      if (!requestId) {
        setError("Request created without an identifier")
        return
      }

      if (mode === "submit") {
        const submitResponse = await fetch(`/api/purchase-requests/${requestId}/submit`, {
          method: "POST",
        })
        const submitBody = await submitResponse.json()
        if (!submitResponse.ok) {
          setError(submitBody?.error?.message ?? "Failed to submit request")
          return
        }
      }

      router.push(`/requests/${requestId}`)
      router.refresh()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to save request")
    } finally {
      setPending(false)
    }
  }

  return (
    <form onSubmit={onSubmit} className="stack">
      <div className="form-grid">
        <label className="field">
          <span>Department</span>
          <select name="department_id" className={inputClassName} required defaultValue="">
            <option value="" disabled>
              Select department
            </option>
            {departments.map((department) => (
              <option key={department.id} value={department.id}>
                {department.name} ({department.code || department.id.slice(0, 6)})
              </option>
            ))}
          </select>
        </label>
        <label className="field">
          <span>Title</span>
          <Input name="title" placeholder="Buy laptops for new engineers" required />
        </label>
        <label className="field">
          <span>Item name</span>
          <Input name="item_name" placeholder="Laptop" required />
        </label>
        <label className="field">
          <span>Quantity</span>
          <Input name="quantity" type="number" min="1" step="1" defaultValue="1" required />
        </label>
        <label className="field">
          <span>Estimated amount</span>
          <Input name="estimated_amount" type="number" min="0" step="1" placeholder="45000000" required />
        </label>
        <label className="field">
          <span>Currency</span>
          <Input name="currency" defaultValue="IDR" required />
        </label>
        <label className="field">
          <span>Urgency</span>
          <Input name="urgency" defaultValue="high" />
        </label>
        <label className="field">
          <span>Expected date</span>
          <Input name="expected_date" type="date" />
        </label>
        <label className="field">
          <span>Vendor name optional</span>
          <Input name="vendor_name" placeholder="Vendor ABC" />
        </label>
        <label className="field md:col-span-2">
          <span>Description</span>
          <Textarea name="description" placeholder="Describe the need and business justification" required />
        </label>
        <label className="field md:col-span-2">
          <span>Notes optional</span>
          <Textarea name="notes" placeholder="Prefer business warranty" />
        </label>
      </div>

      {error ? <p className="text-sm text-red-600">{error}</p> : null}

      <div className="action-row">
        <Button type="submit" value="draft" disabled={pending} variant="outline">
          {pending ? "Saving..." : "Save draft"}
        </Button>
        <Button type="submit" value="submit" disabled={pending}>
          {pending ? "Submitting..." : "Submit request"}
        </Button>
      </div>
    </form>
  )
}
