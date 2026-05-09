"use client"

import { useState, type FormEvent } from "react"
import { useRouter } from "next/navigation"

import { Button } from "@/components/ui/button"
import { Textarea } from "@/components/ui/textarea"

type Props = {
  requestId: string
}

export function RequestActionPanel({ requestId }: Props) {
  const router = useRouter()
  const [error, setError] = useState<string | null>(null)
  const [pending, setPending] = useState(false)

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const submitter = (event.nativeEvent as SubmitEvent).submitter as HTMLButtonElement | null
    const action = submitter?.value ?? "approve"
    const formData = new FormData(event.currentTarget)
    const comment = String(formData.get("comment") ?? "")

    setPending(true)
    setError(null)
    try {
      const response = await fetch(`/api/purchase-requests/${requestId}/approval-actions`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ action, comment }),
      })
      const body = await response.json()
      if (!response.ok) {
        setError(body?.error?.message ?? "Approval action failed")
        return
      }
      event.currentTarget.reset()
      router.refresh()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Approval action failed")
    } finally {
      setPending(false)
    }
  }

  return (
    <form onSubmit={onSubmit} className="grid gap-4">
      <Textarea name="comment" placeholder="Approved. Proceed." />
      {error ? <p className="text-sm text-red-600">{error}</p> : null}
      <div className="flex flex-wrap gap-3">
        <Button type="submit" value="approve" disabled={pending}>
          {pending ? "Working..." : "Approve"}
        </Button>
        <Button type="submit" value="reject" disabled={pending} variant="outline">
          Reject
        </Button>
        <Button type="submit" value="request_revision" disabled={pending} variant="outline">
          Request revision
        </Button>
      </div>
    </form>
  )
}
