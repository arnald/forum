# Forum Application Architecture Documentation

This document explains the complete architecture and design patterns used in this Go forum application. It's designed to help students understand how all the components work together.

## 📁 Project Structure

```
forum-aicreated/
├── cmd/forum/              # Application entry point
│   └── main.go            # HTTP server setup and routing
├── internal/              # Private application code
│   ├── auth/              # Authentication & authorization
│   ├── database/          # Database connection & schema
│   ├── handlers/          # HTTP request handlers
│   └── models/            # Data structures
├── static/                # Static files (CSS, JS, images)
├── templates/             # HTML templates
├── data/                  # Database and uploaded files
└── tests/                 # Test files
```

## 🏗️ Architecture Patterns

### 1. **Dependency Injection Pattern**
- The database connection is injected into handlers and auth services
- This makes the code testable and loosely coupled
- Example: `handlers.NewHandler(db)` injects the database dependency

### 2. **Repository Pattern**
- Database operations are abstracted into specific methods
- Each domain (users, posts, comments) has its own query methods
- Example: `GetUserByEmail()`, `CreatePost()`, `GetComments()`

### 3. **MVC-like Structure**
- **Models**: Data structures in `/internal/models/`
- **Views**: HTML templates in `/templates/`
- **Controllers**: HTTP handlers in `/internal/handlers/`

### 4. **Middleware Pattern**
- Rate limiting is implemented as HTTP middleware
- Wraps handler functions to add cross-cutting concerns
- Example: `h.RateLimitedHandler(h.Register)`

## 🔐 Authentication System

### Password Security
```go
// Passwords are hashed using bcrypt with automatic salt generation
func (a *Auth) HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}
```

### Session Management
1. **Session Creation**: UUID-based session IDs stored in HTTP cookies
2. **Session Storage**: Sessions saved in database with expiration times
3. **Session Validation**: Every request checks session validity and expiration
4. **Security Features**:
   - HttpOnly cookies (prevent XSS)
   - Automatic expiration (24 hours)
   - Session cleanup on logout

### Authorization System
```go
// Role-based permissions
const (
    RoleGuest     = "guest"     // View only
    RoleUser      = "user"      // Create content
    RoleModerator = "moderator" // Moderate content
    RoleAdmin     = "admin"     // Full access
)
```

## 🗄️ Database Design

### Core Tables
1. **users**: User accounts and authentication
2. **sessions**: Login session tracking
3. **posts**: Forum posts with moderation status
4. **comments**: User comments on posts
5. **categories**: Post categorization system

### Relationship Tables
1. **post_categories**: Many-to-many posts ↔ categories
2. **post_likes**: User votes on posts
3. **comment_likes**: User votes on comments

### Advanced Features
1. **notifications**: Real-time user notifications
2. **reports**: Content moderation system
3. **role_requests**: User role upgrade requests
4. **rate_limits**: DoS protection

### Key Design Decisions
- **Foreign Keys**: Ensure data integrity with CASCADE deletes
- **Computed Fields**: Like/dislike counts calculated via SQL aggregation
- **Status Fields**: Enable content moderation workflow
- **Timestamps**: Track creation and modification times

## 🌐 HTTP Request Flow

### 1. Request Routing
```go
// main.go sets up all routes
http.HandleFunc("/post/", h.ViewPost)           // View single post
http.HandleFunc("/create-post", h.CreatePost)   // Create new post
http.HandleFunc("/admin", h.AdminPanel)         // Admin dashboard
```

### 2. Handler Pattern
```go
func (h *Handler) ViewPost(w http.ResponseWriter, r *http.Request) {
    // 1. Extract URL parameters
    // 2. Authenticate user (optional)
    // 3. Fetch data from database
    // 4. Prepare template data
    // 5. Render HTML template
}
```

### 3. Template Rendering
```go
func (h *Handler) render(w http.ResponseWriter, tmpl string, data interface{}) {
    // 1. Parse base template + specific template
    // 2. Execute template with data
    // 3. Send HTML response to client
}
```

## 🔒 Security Features

### 1. **Rate Limiting**
- Prevents brute force attacks on login/registration
- Limits requests per IP address (100 requests/minute)
- Sliding window implementation

### 2. **Input Validation**
- All user inputs are sanitized and validated
- SQL injection prevention via prepared statements
- XSS prevention via template escaping

### 3. **File Upload Security**
- File type validation (MIME type checking)
- File size limits (20MB max)
- Secure file naming with UUIDs
- Extension whitelist for images

### 4. **Session Security**
- HttpOnly cookies prevent JavaScript access
- Secure session ID generation (UUID)
- Automatic session expiration
- Session cleanup on logout

## 📊 Data Flow Examples

### User Registration Flow
1. User submits registration form
2. Handler validates input (email/username uniqueness)
3. Password is hashed with bcrypt
4. User record created in database
5. Session created and cookie set
6. User redirected to homepage

### Post Creation Flow
1. User submits post form with optional image
2. Handler checks authentication and permissions
3. Image uploaded and validated if present
4. Post saved with appropriate status (pending for users, approved for admins)
5. Categories associated with post
6. User redirected to homepage

### Notification Flow
1. User performs action (like, comment, etc.)
2. System checks if notification should be created
3. Notification record inserted into database
4. Real-time notification shown on next page load
5. User can mark notifications as read

## 🎨 Frontend Integration

### Template System
- **Base Template**: Common layout with navigation
- **Specific Templates**: Page-specific content
- **Data Binding**: Go structs passed to templates
- **Helper Functions**: Custom template functions for formatting

### Static Assets
- CSS for styling
- JavaScript for interactivity
- Image uploads stored in `/static/uploads/`
- Served via Go's built-in file server

## 🧪 Testing Strategy

### Test Structure
```go
func TestCreateUser(t *testing.T) {
    // 1. Setup test database
    // 2. Create test data
    // 3. Execute function
    // 4. Assert results
    // 5. Cleanup
}
```

### What's Tested
- User authentication functions
- Database operations
- Password hashing/verification
- Session management
- Permission checking

## 🚀 Performance Considerations

### Database Optimization
- Proper indexing on frequently queried fields
- JOIN queries to minimize database round trips
- Pagination for large result sets
- Connection pooling via `database/sql`

### Caching Strategy
- Static assets served with appropriate headers
- Template parsing cached where possible
- Session data cached in database

### Security vs Performance
- Rate limiting adds small overhead but prevents abuse
- Bcrypt is intentionally slow to prevent brute force
- Session validation on every request ensures security

## 🔧 Configuration

### Environment Variables
- `DB_PATH`: Database file location
- Production vs development settings

### Default Settings
- Database: SQLite for simplicity
- Port: 8080
- Session timeout: 24 hours
- Rate limit: 100 requests/minute
- Upload limit: 20MB

## 📚 Learning Objectives

After studying this codebase, students should understand:

1. **Go Web Development**: HTTP handlers, routing, middleware
2. **Database Design**: Relational modeling, foreign keys, indexes
3. **Security**: Authentication, authorization, input validation
4. **Architecture**: Clean code organization, dependency injection
5. **Testing**: Unit tests, test database setup
6. **Deployment**: Static file serving, environment configuration

## 🎯 Extension Points

Areas where students can add features:
1. **OAuth Integration**: Google/GitHub login
2. **Real-time Features**: WebSocket notifications
3. **Advanced Moderation**: Content filtering, auto-moderation
4. **Search**: Full-text search across posts
5. **API**: REST endpoints for mobile app
6. **Caching**: Redis for session storage
7. **Email**: Notification emails, password reset

This architecture provides a solid foundation for learning modern web development patterns while maintaining simplicity and readability.