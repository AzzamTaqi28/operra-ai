"use client"

import { useState, type FormEvent } from "react"
import { useRouter } from "next/navigation"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"

type Props = {
  requestId: string
}

export function AttachmentUploadForm({ requestId }: Props) {
  const router = useRouter()
  const [error, setError] = useState<string | null>(null)
  const [pending, setPending] = useState(false)

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const formData = new FormData(event.currentTarget)
    const file = formData.get("file")
    if (!(file instanceof File) || file.size === 0) {
      setError("Please choose a file")
      return
    }

    const upload = new FormData()
    upload.set("file", file)

    setPending(true)
    setError(null)
    try {
      const response = await fetch(`/api/purchase-requests/${requestId}/attachments`, {
        method: "POST",
        body: upload,
      })
      const payload = await response.json()
      if (!response.ok) {
        setError(payload?.error?.message ?? "Failed to upload attachment")
        return
      }
      event.currentTarget.reset()
      router.refresh()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to upload attachment")
    } finally {
      setPending(false)
    }
  }

  return (
    <form onSubmit={onSubmit} className="grid gap-4">
      <Input name="file" type="file" />
      {error ? <p className="text-sm text-red-600">{error}</p> : null}
      <div className="flex justify-end">
        <Button type="submit" disabled={pending} variant="outline">
          {pending ? "Uploading..." : "Upload attachment"}
        </Button>
      </div>
    </form>
  )
}
