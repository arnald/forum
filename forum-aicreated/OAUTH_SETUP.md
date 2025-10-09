# OAuth Setup Guide

This guide will help you configure third-party authentication (Google, GitHub, and Facebook) for your forum application.

## Prerequisites

- A `.env` file in the root directory (copy from `.env.example` if needed)
- Access to the OAuth provider's developer console

## Quick Start

1. **Create or edit your `.env` file** in the project root
2. **Add your OAuth credentials** for each provider you want to enable
3. **Restart the application** to load the new credentials

## Setting Up OAuth Providers

### Google OAuth

1. **Go to Google Cloud Console**: https://console.cloud.google.com/
2. **Create a new project** or select an existing one
3. **Enable Google+ API**:
   - Go to "APIs & Services" > "Library"
   - Search for "Google+ API"
   - Click "Enable"
4. **Create OAuth credentials**:
   - Go to "APIs & Services" > "Credentials"
   - Click "Create Credentials" > "OAuth client ID"
   - Choose "Web application"
   - Add authorized redirect URI: `http://localhost:8080/auth/google/callback`
   - For production, also add: `https://yourdomain.com/auth/google/callback`
5. **Copy credentials to `.env`**:
   ```
   GOOGLE_CLIENT_ID=your-client-id-here.apps.googleusercontent.com
   GOOGLE_CLIENT_SECRET=your-client-secret-here
   ```

### GitHub OAuth

1. **Go to GitHub Settings**: https://github.com/settings/developers
2. **Click "New OAuth App"**
3. **Fill in the application details**:
   - Application name: `Your Forum Name`
   - Homepage URL: `http://localhost:8080`
   - Authorization callback URL: `http://localhost:8080/auth/github/callback`
   - For production: `https://yourdomain.com/auth/github/callback`
4. **Register the application**
5. **Copy credentials to `.env`**:
   ```
   GITHUB_CLIENT_ID=your-github-client-id
   GITHUB_CLIENT_SECRET=your-github-client-secret
   ```

### Facebook OAuth

1. **Go to Facebook Developers**: https://developers.facebook.com/apps/
2. **Click "Create App"**
3. **Choose "Consumer" as the app type**
4. **Fill in app details** and create the app
5. **Add Facebook Login product**:
   - In the dashboard, click "Add Product"
   - Choose "Facebook Login" and click "Set Up"
6. **Configure OAuth settings**:
   - Go to "Facebook Login" > "Settings"
   - Add Valid OAuth Redirect URIs:
     - `http://localhost:8080/auth/facebook/callback`
     - For production: `https://yourdomain.com/auth/facebook/callback`
7. **Get your credentials**:
   - Go to "Settings" > "Basic"
   - Copy App ID and App Secret
8. **Copy credentials to `.env`**:
   ```
   FACEBOOK_CLIENT_ID=your-facebook-app-id
   FACEBOOK_CLIENT_SECRET=your-facebook-app-secret
   ```

## Configuration File

Your `.env` file should look like this:

```bash
# Application Settings
BASE_URL=http://localhost:8080

# Database Configuration
DB_PATH=./data/forum.db

# Google OAuth Configuration
GOOGLE_CLIENT_ID=1234567890-abc123def456.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=GOCSPX-your_secret_here

# GitHub OAuth Configuration
GITHUB_CLIENT_ID=Iv1.a1b2c3d4e5f6g7h8
GITHUB_CLIENT_SECRET=your_github_secret_here

# Facebook OAuth Configuration
FACEBOOK_CLIENT_ID=1234567890123456
FACEBOOK_CLIENT_SECRET=your_facebook_secret_here
```

## Running the Application

### From the project root directory:

```bash
# Option 1: Run directly with go run
go run cmd/forum/main.go

# Option 2: Build and run the binary
go build -o forum-app cmd/forum/main.go
./forum-app
```

### Important Notes:

1. **Do NOT run from parent directory** - The error you were seeing:
   ```
   package forum/internal/database is not in std
   ```
   happens when you run `go run forum-aicreated/cmd/forum/main.go` from outside the project directory.

   **Always run from within the `forum-aicreated` directory**:
   ```bash
   cd forum-aicreated
   go run cmd/forum/main.go
   ```

2. **Environment Variables**: The application will automatically load the `.env` file when it starts. You'll see this message if no `.env` file is found:
   ```
   No .env file found, using system environment variables
   ```

3. **OAuth Not Configured**: If you don't configure an OAuth provider, users will see an error message when trying to use that provider:
   ```
   Google OAuth not configured
   ```

## Testing OAuth

1. **Start the application**:
   ```bash
   cd forum-aicreated
   ./forum-app
   ```

2. **Open your browser** and go to: http://localhost:8080/login

3. **You should see three OAuth buttons**:
   - Login with Google
   - Login with GitHub
   - Login with Facebook

4. **Click on a provider button** - You should be redirected to that provider's login page

5. **After authentication**, you'll be redirected back to your forum and automatically logged in

## Troubleshooting

### "OAuth not configured" error
- Check that your `.env` file exists and has the correct credentials
- Make sure there are no extra spaces or quotes around the values
- Restart the application after updating the `.env` file

### "Invalid redirect URI" error
- Verify the callback URL in your OAuth provider settings matches exactly:
  - Google: `http://localhost:8080/auth/google/callback`
  - GitHub: `http://localhost:8080/auth/github/callback`
  - Facebook: `http://localhost:8080/auth/facebook/callback`
- Make sure the protocol (http/https) matches

### "Invalid state token" error
- This is a security feature to prevent CSRF attacks
- Clear your browser cookies and try again
- Make sure cookies are enabled in your browser

### Import path errors
- Always run the application from within the `forum-aicreated` directory
- If you still see import errors, run: `go mod tidy`

## Production Deployment

When deploying to production:

1. **Update BASE_URL** in `.env`:
   ```
   BASE_URL=https://yourdomain.com
   ```

2. **Add production callback URLs** to each OAuth provider:
   - Google: `https://yourdomain.com/auth/google/callback`
   - GitHub: `https://yourdomain.com/auth/github/callback`
   - Facebook: `https://yourdomain.com/auth/facebook/callback`

3. **Enable HTTPS cookies** by setting `Secure: true` in:
   - `internal/handlers/oauth.go` (lines 44, 103, 162, 241)

4. **Never commit `.env`** to version control - keep it in `.gitignore`

## Security Notes

- The `.env` file contains sensitive credentials and should **never** be committed to version control
- Always use HTTPS in production (set `Secure: true` on cookies)
- The application uses state tokens to prevent CSRF attacks
- Sessions expire after 24 hours for security
- OAuth users don't have local passwords - they authenticate through the provider

## Additional Resources

- [Google OAuth Documentation](https://developers.google.com/identity/protocols/oauth2)
- [GitHub OAuth Documentation](https://docs.github.com/en/developers/apps/building-oauth-apps)
- [Facebook Login Documentation](https://developers.facebook.com/docs/facebook-login)
