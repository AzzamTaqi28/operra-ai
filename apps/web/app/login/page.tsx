import Link from "next/link"

export default function LoginPage() {
  return (
    <main className="page-shell">
      <section className="hero">
        <p className="eyebrow">Operra access</p>
        <h1>Sign in to manage approvals.</h1>
        <p className="lede">
          This route is ready for the app login flow. The first implementation pass keeps the UI shell in place while the API integration is connected.
        </p>
        <div className="panel">
          <h2>Next step</h2>
          <p className="muted-copy">Connect this form to the `/api/v1/auth/login` endpoint and persist the JWT for protected routes.</p>
          <Link href="/dashboard" className="nav-link" style={{ display: "inline-flex", marginTop: 16 }}>
            Continue to dashboard
          </Link>
        </div>
      </section>
    </main>
  )
}
