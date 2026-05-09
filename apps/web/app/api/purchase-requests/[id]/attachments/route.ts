import { NextResponse } from "next/server"

import { API_BASE_URL, getSessionToken } from "@/lib/proxy"

export async function POST(request: Request, { params }: { params: Promise<{ id: string }> }) {
  const token = await getSessionToken()
  if (!token) {
    return NextResponse.json({ error: { code: "UNAUTHORIZED", message: "unauthorized" } }, { status: 401 })
  }

  const formData = await request.formData().catch(() => null)
  const file = formData?.get("file")
  if (!(file instanceof File)) {
    return NextResponse.json({ error: { code: "VALIDATION_ERROR", message: "file is required" } }, { status: 400 })
  }

  const upload = new FormData()
  upload.set("file", file)

  const { id } = await params
  const upstream = await fetch(`${API_BASE_URL}/api/v1/purchase-requests/${id}/attachments`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${token}`,
    },
    body: upload,
    cache: "no-store",
  })

  const payload = await upstream.json().catch(() => null)
  return NextResponse.json(payload ?? {}, { status: upstream.status })
}
