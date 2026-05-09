import type { CSSProperties, ReactNode } from "react"

export function PageGrid({ children }: Readonly<{ children: ReactNode }>) {
  return <div className="stack">{children}</div>
}

export function Card({ title, description, children }: Readonly<{ title?: string; description?: string; children: ReactNode }>) {
  return (
    <section className="panel">
      {title ? <h2>{title}</h2> : null}
      {description ? <p className="muted-copy">{description}</p> : null}
      {children}
    </section>
  )
}

export function MetricCard({ label, value, hint }: Readonly<{ label: string; value: string; hint?: string }>) {
  return (
    <article className="stat-card">
      <p>{label}</p>
      <strong>{value}</strong>
      {hint ? <small className="muted-copy">{hint}</small> : null}
    </article>
  )
}

export function Table({ headers, rows }: Readonly<{ headers: string[]; rows: ReactNode[][] }>) {
  const columns = `repeat(${Math.max(headers.length, 1)}, minmax(0, 1fr))`
  return (
    <div className="table card-table" style={{ "--table-columns": columns } as CSSProperties}>
      <div className="table-row table-head">
        {headers.map((header) => (
          <span key={header}>{header}</span>
        ))}
      </div>
      {rows.map((row, index) => (
        <div key={index} className="table-row">
          {row.map((cell, cellIndex) => (
            <span key={cellIndex}>{cell}</span>
          ))}
        </div>
      ))}
    </div>
  )
}

export function Field({ label, children, hint }: Readonly<{ label: string; children: ReactNode; hint?: string }>) {
  return (
    <label className="field">
      <span>{label}</span>
      {children}
      {hint ? <small>{hint}</small> : null}
    </label>
  )
}

export function Chip({ children }: Readonly<{ children: ReactNode }>) {
  return <span className="chip">{children}</span>
}

export function SplitLayout({ left, right }: Readonly<{ left: ReactNode; right: ReactNode }>) {
  return (
    <div className="split-layout">
      <div className="split-pane">{left}</div>
      <div className="split-pane">{right}</div>
    </div>
  )
}
