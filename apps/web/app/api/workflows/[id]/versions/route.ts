import { proxyAuthedJson } from "@/lib/proxy"

export async function POST(request: Request, { params }: { params: Promise<{ id: string }> }) {
  const { id } = await params
  return proxyAuthedJson(request, `/api/v1/workflows/${id}/versions`, {
    method: "POST",
    body: await request.text(),
  })
}
