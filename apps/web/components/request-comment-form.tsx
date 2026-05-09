"use client"

import { useState, type FormEvent } from "react"
import { useRouter } from "next/navigation"

import { Button } from "@/components/ui/button"
import { Textarea } from "@/components/ui/textarea"

type Props = {
  requestId: string
}

export function RequestCommentForm({ requestId }: Props) {
  const router = useRouter()
  const [error, setError] = useState<string | null>(null)
  const [pending, setPending] = useState(false)

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const formData = new FormData(event.currentTarget)
    const body = String(formData.get("body") ?? "")

    if (!body.trim()) {
      setError("Comment cannot be empty")
      return
    }

    setPending(true)
    setError(null)
    try {
      const response = await fetch(`/api/purchase-requests/${requestId}/comments`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ body }),
      })
      const payload = await response.json()
      if (!response.ok) {
        setError(payload?.error?.message ?? "Failed to add comment")
        return
      }
      event.currentTarget.reset()
      router.refresh()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to add comment")
    } finally {
      setPending(false)
    }
  }

  return (
    <form onSubmit={onSubmit} className="grid gap-4">
      <Textarea name="body" placeholder="Add a comment or context for reviewers." />
      {error ? <p className="text-sm text-red-600">{error}</p> : null}
      <div className="flex justify-end">
        <Button type="submit" disabled={pending} variant="outline">
          {pending ? "Posting..." : "Add comment"}
        </Button>
      </div>
    </form>
  )
}
