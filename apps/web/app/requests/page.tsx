import { AppShell } from "@/components/app-shell"

const requests = [
  ["Laptop refresh", "in_review", "IDR 45,000,000"],
  ["Office supplies", "draft", "IDR 1,200,000"],
  ["Cloud licenses", "revision_requested", "IDR 12,500,000"],
]

export default function RequestsPage() {
  return (
    <AppShell title="Purchase Requests" description="Drafts, submissions, approval timelines, comments, and attachments.">
      <section className="panel">
        <h2>Recent requests</h2>
        <div className="table">
          {requests.map(([title, status, amount]) => (
            <div key={title} className="table-row">
              <span>{title}</span>
              <span>{status}</span>
              <span>{amount}</span>
            </div>
          ))}
        </div>
      </section>
    </AppShell>
  )
}
