import { proxyAuthedJson } from "@/lib/proxy"

export async function POST(request: Request, { params }: { params: Promise<{ id: string; versionId: string }> }) {
  const { id, versionId } = await params
  return proxyAuthedJson(request, `/api/v1/workflows/${id}/versions/${versionId}/activate`, {
    method: "POST",
    body: await request.text(),
  })
}
