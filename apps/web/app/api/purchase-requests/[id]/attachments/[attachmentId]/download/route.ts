import { NextResponse } from "next/server"

import { API_BASE_URL, getSessionToken } from "@/lib/proxy"

export async function GET(_request: Request, { params }: { params: Promise<{ id: string; attachmentId: string }> }) {
  const token = await getSessionToken()
  if (!token) {
    return NextResponse.json({ error: { code: "UNAUTHORIZED", message: "unauthorized" } }, { status: 401 })
  }

  const { id, attachmentId } = await params
  const upstream = await fetch(`${API_BASE_URL}/api/v1/purchase-requests/${id}/attachments/${attachmentId}/download`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
    cache: "no-store",
  })

  if (!upstream.ok) {
    const payload = await upstream.json().catch(() => null)
    return NextResponse.json(payload ?? {}, { status: upstream.status })
  }

  const contentType = upstream.headers.get("content-type") ?? "application/octet-stream"
  const contentDisposition = upstream.headers.get("content-disposition") ?? `attachment; filename="${attachmentId}"`
  const data = await upstream.arrayBuffer()
  return new NextResponse(data, {
    status: upstream.status,
    headers: {
      "Content-Type": contentType,
      "Content-Disposition": contentDisposition,
    },
  })
}
