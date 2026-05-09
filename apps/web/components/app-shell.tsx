import type { ReactNode } from "react"
import Link from "next/link"

const navItems = [
  { href: "/dashboard", label: "Dashboard" },
  { href: "/requests", label: "Requests" },
  { href: "/workflows", label: "Workflows" },
  { href: "/ai-workflow", label: "AI Builder" },
  { href: "/users", label: "Users" },
  { href: "/departments", label: "Departments" },
  { href: "/audit-logs", label: "Audit Logs" },
  { href: "/exports", label: "Exports" },
]

export function AppShell({ title, description, children }: Readonly<{ title: string; description?: string; children: ReactNode }>) {
  return (
    <main className="app-shell">
      <aside className="sidebar">
        <div className="brand-block">
          <span className="brand-mark">O</span>
          <div>
            <p className="brand-name">Operra</p>
            <p className="brand-subtitle">Approval workflows</p>
          </div>
        </div>

        <nav className="nav">
          {navItems.map((item) => (
            <Link key={item.href} href={item.href} className="nav-link">
              {item.label}
            </Link>
          ))}
        </nav>
      </aside>

      <section className="content">
        <header className="page-header">
          <p className="eyebrow">Operra v0.1</p>
          <h1>{title}</h1>
          {description ? <p className="lede">{description}</p> : null}
        </header>
        {children}
      </section>
    </main>
  )
}
