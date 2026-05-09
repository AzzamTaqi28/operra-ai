import { proxyAuthedJson } from "@/lib/proxy"

export async function PATCH(request: Request, { params }: { params: Promise<{ id: string }> }) {
  const { id } = await params
  return proxyAuthedJson(request, `/api/v1/departments/${id}`, {
    method: "PATCH",
    body: await request.text(),
  })
}
