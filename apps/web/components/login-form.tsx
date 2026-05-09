"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"

export function LoginForm() {
  const router = useRouter()
  const [error, setError] = useState<string | null>(null)
  const [pending, setPending] = useState(false)

  async function onSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const form = event.currentTarget
    const formData = new FormData(form)
    const payload = {
      email: String(formData.get("email") ?? ""),
      password: String(formData.get("password") ?? ""),
    }

    setPending(true)
    setError(null)
    try {
      const response = await fetch("/api/auth/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      })
      const body = await response.json()
      if (!response.ok) {
        setError(body?.error?.message ?? "Login failed")
        return
      }
      router.push("/dashboard")
      router.refresh()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Login failed")
    } finally {
      setPending(false)
    }
  }

  return (
    <Card className="border-slate-200/80 shadow-xl shadow-slate-900/5">
      <CardHeader>
        <CardTitle>Login</CardTitle>
        <CardDescription>Sign in to the tenant-scoped approval workspace.</CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={onSubmit} className="grid gap-4">
          <label className="grid gap-2 text-sm font-medium text-slate-900">
            Email
            <Input name="email" type="email" placeholder="taqi@example.com" autoComplete="email" required />
          </label>
          <label className="grid gap-2 text-sm font-medium text-slate-900">
            Password
            <Input name="password" type="password" placeholder="••••••••" autoComplete="current-password" required />
          </label>
          {error ? <p className="text-sm text-red-600">{error}</p> : null}
          <div className="flex gap-3">
            <Button type="submit" disabled={pending}>
              {pending ? "Signing in..." : "Login"}
            </Button>
            <Button type="button" variant="outline" onClick={() => router.push("/setup")}>
              First-time setup
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  )
}
