# Operra Deployment

## 1. Deployment goal

Operra v0.1 must run with Docker Compose.

Target developer/self-hosted command:

```bash
docker compose up -d
```

## 2. Required services

```text
web       - Next.js frontend
api       - Go backend API
postgres  - PostgreSQL database
minio     - S3-compatible object storage
```

Optional later:

```text
redis
worker
ollama
```

## 3. Example docker-compose.yml outline

```yaml
services:
  postgres:
    image: postgres:16
    environment:
      POSTGRES_DB: operra
      POSTGRES_USER: operra
      POSTGRES_PASSWORD: operra
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  minio:
    image: minio/minio:latest
    command: server /data --console-address ":9001"
    environment:
      MINIO_ROOT_USER: operra
      MINIO_ROOT_PASSWORD: operra-secret
    ports:
      - "9000:9000"
      - "9001:9001"
    volumes:
      - minio_data:/data

  api:
    build:
      context: ./apps/api
    env_file:
      - .env
    depends_on:
      - postgres
      - minio
    ports:
      - "8080:8080"

  web:
    build:
      context: ./apps/web
    env_file:
      - .env
    depends_on:
      - api
    ports:
      - "3000:3000"

volumes:
  postgres_data:
  minio_data:
```

## 4. Environment variables

`.env.example` should include:

```env
APP_ENV=development
APP_URL=http://localhost:3000
API_URL=http://localhost:8080

DATABASE_URL=postgres://operra:operra@postgres:5432/operra?sslmode=disable
JWT_SECRET=change-me-please

STORAGE_DRIVER=s3
S3_ENDPOINT=http://minio:9000
S3_BUCKET=operra
S3_ACCESS_KEY=operra
S3_SECRET_KEY=operra-secret
S3_REGION=us-east-1
S3_FORCE_PATH_STYLE=true

AI_PROVIDER=openai
AI_BASE_URL=https://api.openai.com/v1
AI_API_KEY=
AI_MODEL=

SMTP_HOST=
SMTP_PORT=
SMTP_USER=
SMTP_PASSWORD=
SMTP_FROM=
```

## 5. MinIO setup

On first startup, the app should either:

- create the configured bucket automatically, or
- provide a documented command to create it.

Recommended bucket name:

```text
operra
```

## 6. Database migrations

The API startup should not silently destroy data.

Recommended commands:

```bash
make migrate-up
make migrate-down
make seed
```

If using GORM AutoMigrate during early development, add a note:

> AutoMigrate is allowed during early dev, but real migrations are required before private beta.

## 7. Health checks

Required endpoints:

```text
GET /health
GET /ready
```

`/health` can return OK if API is running.

`/ready` should check:

- database connection
- storage connectivity if possible

## 8. Local development setup

Suggested setup:

```bash
cp .env.example .env
docker compose up -d postgres minio
cd apps/api && go run ./cmd/server
cd apps/web && npm run dev
```

Full Docker setup:

```bash
cp .env.example .env
docker compose up -d
```

## 9. Production self-hosting notes

Production users should configure:

- Strong JWT secret.
- Secure database password.
- External S3-compatible storage or hardened MinIO.
- HTTPS reverse proxy.
- Database backups.
- Object storage backups.
- SMTP for email notification.

## 10. One-click deployment future

Future versions may support one-click deployment on container deployment platforms.

Do not build a custom one-click deployment platform in v0.1.

v0.1 should focus on:

- clean Docker Compose
- clear environment variables
- reliable README
- minimal setup steps
