import { LoginForm } from "@/components/login-form"

export default function LoginPage() {
  return (
    <main className="page-shell">
      <section className="hero narrow">
        <p className="eyebrow">Operra access</p>
        <h1>Sign in to manage approvals.</h1>
        <p className="lede">
          Purchase request flow, approvals, audit logs, and exports all sit behind tenant-scoped authentication.
        </p>

        <LoginForm />
      </section>
    </main>
  )
}
