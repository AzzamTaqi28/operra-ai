import { SetupForm } from "@/components/setup-form"

export default function SetupPage() {
  return (
    <main className="page-shell">
      <section className="hero narrow">
        <p className="eyebrow">Organization setup</p>
        <h1>Create the first tenant.</h1>
        <p className="lede">
          This first-time setup screen matches the API contract for organization registration and owner user creation.
        </p>

        <SetupForm />
      </section>
    </main>
  )
}
