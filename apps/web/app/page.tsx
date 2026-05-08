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
          This frontend is the initial Next.js shell for the Operra monorepo.
          It will host request, approval, workflow, audit, and export screens.
        </p>

        <div className="panel">
          <h2>Current implementation focus</h2>
          <ul>
            {highlights.map((item) => (
              <li key={item}>{item}</li>
            ))}
          </ul>
        </div>
      </section>
    </main>
  )
}
