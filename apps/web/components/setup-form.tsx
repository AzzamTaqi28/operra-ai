"use client"

import { useState, type FormEvent } from "react"
import { useRouter } from "next/navigation"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"

export function SetupForm() {
  const router = useRouter()
  const [error, setError] = useState<string | null>(null)
  const [pending, setPending] = useState(false)

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const formData = new FormData(event.currentTarget)
    const payload = {
      organization_name: String(formData.get("organization_name") ?? ""),
      organization_slug: String(formData.get("organization_slug") ?? ""),
      owner_name: String(formData.get("owner_name") ?? ""),
      owner_email: String(formData.get("owner_email") ?? ""),
      password: String(formData.get("password") ?? ""),
    }

    setPending(true)
    setError(null)
    try {
      const response = await fetch("/api/auth/register-organization", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      })
      const body = await response.json()
      if (!response.ok) {
        setError(body?.error?.message ?? "Registration failed")
        return
      }
      router.push("/dashboard")
      router.refresh()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Registration failed")
    } finally {
      setPending(false)
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Register organization</CardTitle>
        <CardDescription>Create the first tenant and initial owner account.</CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={onSubmit} className="stack">
          <div className="form-grid">
            <label className="field">
              <span>Organization name</span>
              <Input name="organization_name" placeholder="Demo Company" required />
            </label>
            <label className="field">
              <span>Organization slug</span>
              <Input name="organization_slug" placeholder="demo-company" required />
            </label>
            <label className="field">
              <span>Owner name</span>
              <Input name="owner_name" placeholder="Taqi" required />
            </label>
            <label className="field">
              <span>Owner email</span>
              <Input name="owner_email" type="email" placeholder="taqi@example.com" required />
            </label>
            <label className="field">
              <span>Password</span>
              <Input name="password" type="password" placeholder="Choose a password" required />
            </label>
          </div>
          {error ? <p className="text-sm text-red-600">{error}</p> : null}
          <div className="action-row">
            <Button type="submit" disabled={pending}>
              {pending ? "Creating..." : "Create organization"}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  )
}
