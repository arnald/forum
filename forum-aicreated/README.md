# 🗣️ Go Forum Application

A comprehensive forum web application built with Go's standard library, featuring user authentication, post management, comments, likes/dislikes, notifications, and admin functionality.

## ✨ Features

### 🔐 User System
- **Registration & Login**: Secure user accounts with bcrypt password hashing
- **OAuth Authentication**: Fully implemented Google, GitHub, and Facebook login
- **Session Management**: UUID-based sessions with automatic expiration
- **Role-Based Access**: Guest, User, Moderator, and Admin roles

### 📝 Content Management
- **Posts**: Create, edit, delete posts with optional image uploads
- **Comments**: Threaded discussions on posts with edit/delete functionality
- **Categories**: Organize posts by topics (General, Technology, Sports, etc.)
- **Like/Dislike System**: Vote on posts and comments
- **Content Moderation**: Admin approval workflow for posts

### 🔔 Interactive Features
- **Real-time Notifications**: Get notified when someone interacts with your content
- **User Activity Tracking**: View your posts, comments, and liked content
- **Filtering**: Browse posts by category or user-specific criteria

### 👑 Admin Features
- **User Management**: Change user roles, view all users
- **Content Moderation**: Approve/reject pending posts
- **Category Management**: Create and delete post categories
- **Admin Dashboard**: Comprehensive overview of forum activity

### 🛡️ Security Features
- **Rate Limiting**: Prevents spam and brute force attacks (100 requests/minute)
- **Input Validation**: Protects against SQL injection and XSS
- **Secure File Uploads**: Type validation and size limits for images
- **Permission System**: Role-based access control throughout the application

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- SQLite (automatically handled by Go driver)

### Installation

1. **Clone the repository**
```bash
git clone <repository-url>
cd forum-aicreated
```

2. **Install dependencies**
```bash
go mod tidy
```

3. **Run the application**
```bash
go run cmd/forum/main.go
```

4. **Access the application**
Open your browser and visit: `http://localhost:8080`

### Default Admin Account
- **Email**: `admin@forum.com`
- **Username**: `admin`
- **Password**: `admin123`

## 🗂️ Project Structure

```
forum-aicreated/
├── cmd/forum/main.go          # Application entry point with routing
├── internal/                  # Private application code
│   ├── auth/auth.go          # Authentication & authorization
│   ├── database/database.go  # Database connection & schema
│   ├── handlers/             # HTTP request handlers
│   │   ├── handlers.go       # Core handlers (home, auth)
│   │   ├── queries.go        # Database query methods
│   │   ├── actions.go        # User action handlers (like, comment)
│   │   ├── notifications.go  # Notification system
│   │   ├── admin.go          # Admin panel functionality
│   │   ├── edit.go           # Edit/delete operations
│   │   ├── activity.go       # User activity tracking
│   │   ├── upload.go         # File upload handling
│   │   └── ratelimit.go      # Rate limiting middleware
│   └── models/models.go      # Data structures and types
├── static/                   # Static files (CSS, JS, images)
├── templates/                # HTML templates
├── data/                     # Database and uploaded files
├── tests/                    # Test files
├── ARCHITECTURE.md           # Detailed architecture documentation
└── README.md                 # This file
```

## 📊 Database Schema

The application automatically creates a SQLite database with comprehensive tables:

### Core Tables
- **users**: User accounts with roles and OAuth support
- **sessions**: Session management with expiration
- **posts**: Forum posts with moderation status
- **comments**: User comments with edit history
- **categories**: Post categorization system

### Feature Tables
- **post_likes/comment_likes**: Voting system
- **notifications**: Real-time user notifications
- **reports**: Content reporting for moderation
- **role_requests**: User role upgrade system
- **rate_limits**: DoS protection tracking

## 🎯 Usage Guide

### For Regular Users
1. **Registration**: Create account with email and username
2. **Create Posts**: Write posts and assign categories
3. **Upload Images**: Add images to posts (20MB max)
4. **Interact**: Like/dislike posts and comments
5. **Track Activity**: View your activity in the Activity section
6. **Notifications**: Get notified of interactions

### For Moderators
- All user features plus:
- Approve/reject pending posts
- Delete inappropriate content
- Access moderation tools

### For Administrators
- All moderator features plus:
- Promote users to moderators
- Create/delete categories
- Access admin dashboard
- Manage all users and content

## 🧪 Testing

```bash
# Run all tests
go test ./tests/...

# Run tests with coverage
go test -cover ./tests/...

# Run specific test
go test ./tests/ -run TestCreateUser
```

## 🔧 Configuration

### Environment Variables
Create a `.env` file in the project root (see `.env.example` for reference):

- `BASE_URL`: Application base URL (default: `http://localhost:8080`)
- `DB_PATH`: Custom database file location (default: `./data/forum.db`)

**OAuth Credentials** (optional - leave empty to disable):
- `GOOGLE_CLIENT_ID` & `GOOGLE_CLIENT_SECRET`: Google OAuth credentials
- `GITHUB_CLIENT_ID` & `GITHUB_CLIENT_SECRET`: GitHub OAuth credentials
- `FACEBOOK_CLIENT_ID` & `FACEBOOK_CLIENT_SECRET`: Facebook OAuth credentials

📖 **See [OAUTH_SETUP.md](OAUTH_SETUP.md) for detailed OAuth configuration guide**

### Application Settings
- **Server Port**: 8080
- **Session Timeout**: 24 hours
- **Rate Limit**: 100 requests per minute per IP
- **Max Upload Size**: 20MB
- **Allowed Image Types**: JPG, JPEG, PNG, GIF

## 🔒 Security Features

1. **Authentication Security**
   - Bcrypt password hashing with cost 10
   - OAuth 2.0 integration with Google, GitHub, and Facebook
   - CSRF protection for OAuth flows using state tokens
   - UUID-based session IDs
   - HttpOnly cookies to prevent XSS
   - Automatic session expiration

2. **Input Validation**
   - SQL injection prevention via prepared statements
   - XSS protection through template escaping
   - File upload validation (type, size, content)

3. **Rate Limiting**
   - Per-IP request limiting
   - Sliding window algorithm
   - Protection against brute force attacks

4. **Authorization**
   - Role-based permission system
   - Ownership checks for content modification
   - Admin-only functionality protection

## 🚧 Future Enhancements

Planned features for educational extension:

- **HTTPS Support**: SSL/TLS certificates for production deployment
- **Email System**: Email verification and notification emails
- **Search**: Full-text search across posts and comments
- **REST API**: JSON API for mobile app support
- **Real-time**: WebSocket notifications for instant updates
- **Advanced Moderation**: Auto-moderation tools and content filtering
- **Two-Factor Authentication**: Additional security layer for accounts

## 🛠️ Troubleshooting

### Common Issues

1. **Port 8080 in use**: Change port in main.go or kill existing process
2. **Database errors**: Ensure ./data/ directory is writable
3. **Upload failures**: Check static/uploads/ directory permissions
4. **Rate limiting**: Wait 1 minute for reset or reduce request frequency

## 📚 Learning Resources

- **ARCHITECTURE.md**: Detailed system architecture explanation
- **Code Comments**: Extensive inline documentation
- **Go Documentation**: https://golang.org/doc/
- **Security Guidelines**: OWASP best practices

## 💡 Educational Value

This project teaches:

- **Clean Architecture**: Dependency injection, separation of concerns
- **Security**: Authentication, authorization, input validation
- **Database Design**: Relational modeling, foreign keys, indexing
- **Web Development**: HTTP handling, templating, middleware
- **Testing**: Unit tests, test databases, coverage analysis
- **Performance**: Efficient queries, rate limiting, caching strategies

Perfect for students learning modern web development with Go!