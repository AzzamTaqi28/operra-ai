import { cookies } from "next/headers"
import { NextResponse } from "next/server"

export const API_BASE_URL = process.env.API_URL ?? "http://localhost:8080"

export async function getSessionToken() {
  const store = await cookies()
  return store.get("operra_token")?.value ?? null
}

export async function proxyAuthedJson(request: Request, path: string, init: RequestInit = {}) {
  const token = await getSessionToken()
  if (!token) {
    return NextResponse.json(
      { error: { code: "UNAUTHORIZED", message: "unauthorized" } },
      { status: 401 },
    )
  }

  const upstream = await fetch(`${API_BASE_URL}${path}`, {
    ...init,
    cache: "no-store",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
      ...(init.headers ?? {}),
    },
  })

  const payload = await upstream.json().catch(() => null)
  return NextResponse.json(payload ?? {}, { status: upstream.status })
}
