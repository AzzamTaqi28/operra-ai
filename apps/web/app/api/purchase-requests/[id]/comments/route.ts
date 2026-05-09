import { NextResponse } from "next/server"

import { API_BASE_URL, getSessionToken } from "@/lib/proxy"

export async function POST(request: Request, { params }: { params: Promise<{ id: string }> }) {
  const token = await getSessionToken()
  if (!token) {
    return NextResponse.json({ error: { code: "UNAUTHORIZED", message: "unauthorized" } }, { status: 401 })
  }

  const body = await request.json().catch(() => null)
  if (!body || typeof body.body !== "string") {
    return NextResponse.json({ error: { code: "VALIDATION_ERROR", message: "invalid request body" } }, { status: 400 })
  }

  const { id } = await params
  const upstream = await fetch(`${API_BASE_URL}/api/v1/purchase-requests/${id}/comments`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(body),
    cache: "no-store",
  })

  const payload = await upstream.json().catch(() => null)
  return NextResponse.json(payload ?? {}, { status: upstream.status })
}
