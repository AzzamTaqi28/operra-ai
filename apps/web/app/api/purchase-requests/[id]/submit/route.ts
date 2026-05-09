import { NextResponse } from "next/server"

import { API_BASE_URL, getSessionToken } from "@/lib/proxy"

export async function POST(_request: Request, { params }: { params: Promise<{ id: string }> }) {
  const token = await getSessionToken()
  if (!token) {
    return NextResponse.json({ error: { code: "UNAUTHORIZED", message: "unauthorized" } }, { status: 401 })
  }

  const { id } = await params
  const upstream = await fetch(`${API_BASE_URL}/api/v1/purchase-requests/${id}/submit`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    cache: "no-store",
  })

  const payload = await upstream.json().catch(() => null)
  return NextResponse.json(payload ?? {}, { status: upstream.status })
}
