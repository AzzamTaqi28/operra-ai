import { Card, Field } from "@/components/ui"

export default function SetupPage() {
  return (
    <main className="page-shell">
      <section className="hero narrow">
        <p className="eyebrow">Organization setup</p>
        <h1>Create the first tenant.</h1>
        <p className="lede">
          This first-time setup screen matches the API contract for organization registration and owner user creation.
        </p>

        <Card title="Register organization">
          <div className="form-grid">
            <Field label="Organization name"><input className="input" placeholder="Demo Company" /></Field>
            <Field label="Organization slug"><input className="input" placeholder="demo-company" /></Field>
            <Field label="Owner name"><input className="input" placeholder="Taqi" /></Field>
            <Field label="Owner email"><input className="input" type="email" placeholder="taqi@example.com" /></Field>
            <Field label="Password"><input className="input" type="password" placeholder="Choose a password" /></Field>
          </div>
          <div className="action-row">
            <button className="button button-solid" type="button">Create organization</button>
          </div>
        </Card>
      </section>
    </main>
  )
}
