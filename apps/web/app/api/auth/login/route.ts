import { NextResponse } from "next/server"

const API_BASE_URL = process.env.API_URL ?? "http://localhost:8080"

export async function POST(request: Request) {
  const body = await request.json().catch(() => null)
  if (!body || typeof body.email !== "string" || typeof body.password !== "string") {
    return NextResponse.json({ error: { code: "VALIDATION_ERROR", message: "invalid request body" } }, { status: 400 })
  }

  const upstream = await fetch(`${API_BASE_URL}/api/v1/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
    cache: "no-store",
  })

  const payload = await upstream.json().catch(() => null)
  if (!upstream.ok) {
    return NextResponse.json(payload ?? { error: { code: "INTERNAL_ERROR", message: "login failed" } }, { status: upstream.status })
  }

  const response = NextResponse.json(payload)
  response.cookies.set("operra_token", payload?.data?.token ?? "", {
    httpOnly: true,
    sameSite: "lax",
    secure: process.env.NODE_ENV === "production",
    path: "/",
  })
  return response
}
