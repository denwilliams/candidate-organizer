# Heroku Deployment Setup

This repository uses a monorepo structure with the backend in the `backend/` directory.

## Required Heroku Configuration

### 1. Set the buildpacks (in order):

```bash
heroku buildpacks:clear
heroku buildpacks:add https://github.com/lstoll/heroku-buildpack-monorepo
heroku buildpacks:add heroku/go
```

### 2. Configure the app base directory:

```bash
heroku config:set APP_BASE=backend
```

### 3. Set required environment variables:

```bash
heroku config:set PORT=8080
heroku config:set DATABASE_URL='<your-postgres-url>'
heroku config:set GOOGLE_CLIENT_ID='<your-google-client-id>'
heroku config:set GOOGLE_CLIENT_SECRET='<your-google-client-secret>'
heroku config:set GOOGLE_REDIRECT_URL='https://your-app.herokuapp.com/api/v1/auth/callback'
heroku config:set WORKSPACE_DOMAIN='<your-workspace-domain>'
heroku config:set JWT_SECRET='<your-secret-key>'
heroku config:set FRONTEND_URL='<your-frontend-url>'
```

### 4. Add Heroku Postgres:

```bash
heroku addons:create heroku-postgresql:essential-0
```

The DATABASE_URL will be set automatically by the addon.

### 5. Deploy:

```bash
git push heroku main
```

## Frontend Deployment

The frontend (Next.js app in `frontend/`) should be deployed separately:
- **Recommended**: Deploy to Vercel (optimized for Next.js)
- Set `NEXT_PUBLIC_API_URL` to your Heroku backend URL
- Update backend's `FRONTEND_URL` to match your frontend deployment URL

## Notes

- The monorepo buildpack allows Heroku to treat the `backend/` directory as the app root
- The Go buildpack will find `go.mod` in the backend directory and build the application
- The Procfile runs relative to the APP_BASE directory
