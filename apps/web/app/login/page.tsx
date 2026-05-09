import Link from "next/link"

import { Card, Field } from "@/components/ui"

export default function LoginPage() {
  return (
    <main className="page-shell">
      <section className="hero narrow">
        <p className="eyebrow">Operra access</p>
        <h1>Sign in to manage approvals.</h1>
        <p className="lede">
          Purchase request flow, approvals, audit logs, and exports all sit behind tenant-scoped authentication.
        </p>

        <Card title="Login">
          <div className="form-grid">
            <Field label="Email"><input className="input" type="email" placeholder="taqi@example.com" /></Field>
            <Field label="Password"><input className="input" type="password" placeholder="••••••••" /></Field>
          </div>
          <div className="action-row">
            <button className="button button-solid" type="button">Login</button>
            <Link className="button button-outline" href="/setup">First-time setup</Link>
          </div>
        </Card>
      </section>
    </main>
  )
}
