# Heroku Deployment Setup

This repository uses a monorepo structure with the backend in the `backend/` directory and frontend in the `frontend/` directory. The deployment is configured to run both the Next.js server and Go backend API on a single Heroku dyno.

## How It Works

During deployment:
1. The Node.js buildpack builds the Next.js frontend using standalone output mode
2. The Go buildpack compiles the backend server
3. Both services start via a process manager script (`start-services.sh`)
4. Next.js runs as the web server and proxies API requests to the Go backend
5. If either service crashes, the entire dyno restarts to maintain system integrity

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

The frontend runs as a full Next.js server with all features enabled:

- Next.js is configured with `output: 'standalone'` for optimized production builds
- Supports all Next.js features including dynamic routes, SSR, ISR, and API routes
- Next.js proxies API requests from `/api/v1/*` to the Go backend on port 8080
- The frontend and API share the same domain, eliminating CORS issues
- Process monitoring ensures both services stay healthy

## Notes

- The monorepo buildpack allows Heroku to treat the `backend/` directory as the app root
- The Node.js buildpack runs first and executes the `heroku-prebuild` script to build the frontend
- The Go buildpack then compiles the backend server
- The `start-services.sh` script manages both processes and monitors their health
- If either service crashes, the script exits and Heroku automatically restarts the dyno
- No separate frontend deployment is needed - everything runs from a single Heroku dyno

## Process Management

The `start-services.sh` script provides robust process management:

- Starts the Go backend on port 8080
- Starts the Next.js server on the port specified by Heroku ($PORT)
- Monitors both processes every 5 seconds
- If either process dies, kills the other and exits (triggering Heroku restart)
- This ensures the application stays in a consistent state
