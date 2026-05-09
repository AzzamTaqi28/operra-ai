import Link from "next/link"

const highlights = [
  "Multi-tenant organization scoping",
  "Role-based approvals",
  "Deterministic workflow execution",
  "MinIO-compatible attachment storage",
]

export default function HomePage() {
  return (
    <main className="page-shell">
      <section className="hero">
        <p className="eyebrow">Operra v0.1</p>
        <h1>Purchase approvals with a system of record.</h1>
        <p className="lede">
          The frontend shell now includes the core product routes for requests,
          workflows, approvals, audit logs, exports, and AI workflow generation.
        </p>

        <div className="panel">
          <h2>Current implementation focus</h2>
          <ul>
            {highlights.map((item) => (
              <li key={item}>{item}</li>
            ))}
          </ul>
          <div style={{ display: "flex", gap: 12, marginTop: 20, flexWrap: "wrap" }}>
            <Link className="nav-link" href="/dashboard">Open dashboard</Link>
            <Link className="nav-link" href="/login">Login route</Link>
          </div>
        </div>
      </section>
    </main>
  )
}
