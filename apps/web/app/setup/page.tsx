import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"

export default function SetupPage() {
  return (
    <main className="page-shell">
      <section className="hero narrow">
        <p className="eyebrow">Organization setup</p>
        <h1>Create the first tenant.</h1>
        <p className="lede">
          This first-time setup screen matches the API contract for organization registration and owner user creation.
        </p>

        <Card>
          <CardHeader>
            <CardTitle>Register organization</CardTitle>
            <CardDescription>Create the first tenant and initial owner account.</CardDescription>
          </CardHeader>
          <CardContent className="stack">
            <div className="form-grid">
              <label className="field">
                <span>Organization name</span>
                <Input placeholder="Demo Company" />
              </label>
              <label className="field">
                <span>Organization slug</span>
                <Input placeholder="demo-company" />
              </label>
              <label className="field">
                <span>Owner name</span>
                <Input placeholder="Taqi" />
              </label>
              <label className="field">
                <span>Owner email</span>
                <Input type="email" placeholder="taqi@example.com" />
              </label>
              <label className="field">
                <span>Password</span>
                <Input type="password" placeholder="Choose a password" />
              </label>
            </div>
            <div className="action-row">
              <Button type="button">Create organization</Button>
            </div>
          </CardContent>
        </Card>
      </section>
    </main>
  )
}
