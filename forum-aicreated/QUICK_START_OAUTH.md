# Quick Start: Enable GitHub OAuth

The 500 error you're seeing is because OAuth is not configured yet. Follow these steps to enable GitHub authentication:

## Step 1: Create a GitHub OAuth App

1. Go to: https://github.com/settings/developers
2. Click **"New OAuth App"** (or "Register a new application")
3. Fill in the form:
   ```
   Application name: My Forum App
   Homepage URL: http://localhost:8080
   Authorization callback URL: http://localhost:8080/auth/github/callback
   ```
   **Note**: You can only set ONE callback URL initially. After creating the app, you can add additional callback URLs.
4. Click **"Register application"**
5. You'll see your **Client ID** displayed
6. Click **"Generate a new client secret"** to get your **Client Secret**
7. **Add the network IP callback URL**:
   - In your OAuth App settings, add this additional callback URL:
   - `http://192.168.56.103:8080/auth/github/callback`
   - This allows OAuth to work from both localhost and remote machines

## Step 2: Update Your .env File

Open `/home/steven/Desktop/forum-ai/forum-aicreated/.env` and add your credentials:

```bash
# GitHub OAuth Configuration
GITHUB_CLIENT_ID=Iv1.abc123def456  # Replace with your Client ID
GITHUB_CLIENT_SECRET=ghp_yoursecrethere123456789  # Replace with your Client Secret
```

## Step 3: Restart the Application

```bash
cd /home/steven/Desktop/forum-ai/forum-aicreated

# Stop any running instances
pkill forum-app

# Start the application
./forum-app
```

Or if you haven't built it yet:
```bash
go run cmd/forum/main.go
```

## Step 4: Test GitHub Login

1. Open your browser and go to: http://192.168.56.103:8080/login
2. Click the **"Login with GitHub"** button
3. You should be redirected to GitHub's authorization page
4. After authorizing, you'll be logged into the forum!

## Current Status

- ✅ Application is running correctly
- ✅ OAuth endpoints are working
- ❌ OAuth credentials not configured (causing the 500 error)
- ✅ BASE_URL set to: `http://localhost:8080` (with alternative for `http://192.168.56.103:8080`)

## What's Happening

The 500 error is **intentional security behavior**. The code at `internal/handlers/oauth.go:88-92` checks:

```go
if OAuthConfig.GitHub.ClientID == "" {
    http.Error(w, "GitHub OAuth not configured", http.StatusInternalServerError)
    return
}
```

This prevents users from attempting OAuth login when it's not properly configured.

## Testing Other Providers

To enable Google OAuth:
- Get credentials from: https://console.cloud.google.com/apis/credentials
- Add both callback URLs:
  - `http://localhost:8080/auth/google/callback`
  - `http://192.168.56.103:8080/auth/google/callback`

To enable Facebook OAuth:
- Get credentials from: https://developers.facebook.com/apps/
- Add both callback URLs:
  - `http://localhost:8080/auth/facebook/callback`
  - `http://192.168.56.103:8080/auth/facebook/callback`

Add the credentials to your `.env` file and restart the app.

## Important Notes

1. **Callback URLs must match exactly** - If you change the IP or port, update:
   - The `BASE_URL` in `.env`
   - The callback URL in your OAuth provider settings

2. **Restart required** - Changes to `.env` only take effect when you restart the application

3. **Local network access** - The IP `192.168.56.103` suggests you're accessing from another machine on your network. Make sure:
   - The forum server can be reached from that IP
   - Firewall allows connections on port 8080

4. **Security** - For production use:
   - Use HTTPS
   - Set `Secure: true` on cookies
   - Keep your `.env` file private (never commit to git)
