# Heroku Deployment Setup

This repository uses a monorepo structure with the backend in the `backend/` directory and frontend in the `frontend/` directory. The deployment is configured to serve both the API and the frontend static files from a single Heroku dyno.

## How It Works

During deployment:
1. The Node.js buildpack builds the Next.js frontend into static files
2. The static files are copied to `backend/static/`
3. The Go buildpack compiles the backend server
4. The Go server serves both the API endpoints and the frontend static files

## Required Heroku Configuration

### 1. Set the buildpacks (in order):

```bash
heroku buildpacks:clear
heroku buildpacks:add https://github.com/lstoll/heroku-buildpack-monorepo
heroku buildpacks:add heroku/nodejs
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
heroku config:set FRONTEND_URL='https://your-app.herokuapp.com'
```

**Note**: The `FRONTEND_URL` should be set to your Heroku app URL since the frontend is served from the same domain.

### 4. Add Heroku Postgres:

```bash
heroku addons:create heroku-postgresql:essential-0
```

The DATABASE_URL will be set automatically by the addon.

### 5. Deploy:

```bash
git push heroku main
```

## Frontend Configuration

The frontend is built as static files during deployment and served by the Go backend server:

- Next.js is configured with `output: 'export'` to generate static files
- Static files are automatically built during deployment via the `heroku-prebuild` script in `backend/package.json`
- The Go server serves the static files from the `backend/static/` directory
- All API requests go to `/api/v1/*` endpoints
- The frontend and API share the same domain, eliminating CORS issues

## Notes

- The monorepo buildpack allows Heroku to treat the `backend/` directory as the app root
- The Node.js buildpack runs first and executes the `heroku-prebuild` script to build the frontend
- The Go buildpack then compiles the backend server which serves both API and static files
- The Procfile runs relative to the APP_BASE directory
- No separate frontend deployment is needed - everything runs from a single Heroku dyno
