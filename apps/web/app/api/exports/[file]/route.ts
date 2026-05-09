import { NextResponse } from "next/server"

import { API_BASE_URL, getSessionToken } from "@/lib/proxy"

const allowed = new Map([
  ["purchase-requests.csv", "/api/v1/exports/purchase-requests.csv"],
  ["approval-history.csv", "/api/v1/exports/approval-history.csv"],
  ["audit-logs.csv", "/api/v1/exports/audit-logs.csv"],
])

export async function GET(request: Request, { params }: { params: Promise<{ file: string }> }) {
  const token = await getSessionToken()
  if (!token) {
    return NextResponse.json({ error: { code: "UNAUTHORIZED", message: "unauthorized" } }, { status: 401 })
  }

  const { file } = await params
  const upstreamPath = allowed.get(file)
  if (!upstreamPath) {
    return NextResponse.json({ error: { code: "NOT_FOUND", message: "resource not found" } }, { status: 404 })
  }

  const requestUrl = new URL(request.url)
  const search = requestUrl.searchParams.toString()
  const upstream = await fetch(`${API_BASE_URL}${upstreamPath}${search ? `?${search}` : ""}`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
    cache: "no-store",
  })

  const contentType = upstream.headers.get("content-type") ?? "text/csv; charset=utf-8"
  const content = await upstream.text()
  return new NextResponse(content, {
    status: upstream.status,
    headers: {
      "Content-Type": contentType,
      "Content-Disposition": `attachment; filename="${file.replace(/"/g, "")}"`,
    },
  })
}
