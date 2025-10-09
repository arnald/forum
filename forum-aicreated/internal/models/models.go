// Package models defines all data structures used throughout the forum application.
// These structs represent database entities and are used for JSON serialization
// and template rendering. They follow Go naming conventions and use struct tags
// for JSON marshaling and database field mapping.
package models

import (
	"time"
)

// ===== USER SYSTEM =====

// UserRole represents the different permission levels in the forum
type UserRole string

const (
	RoleGuest     UserRole = "guest"     // Not logged in users (view only)
	RoleUser      UserRole = "user"      // Regular registered users
	RoleModerator UserRole = "moderator" // Can moderate content and delete posts
	RoleAdmin     UserRole = "admin"     // Full system access
)

// User represents a forum user account
// The Password field is tagged with `json:"-"` to exclude it from JSON responses
// for security reasons
type User struct {
	ID            int       `json:"id"`                      // Unique user identifier
	Email         string    `json:"email"`                   // User's email address
	Username      string    `json:"username"`                // Display name
	Password      string    `json:"-"`                       // Bcrypt hashed password (never sent to frontend)
	Role          UserRole  `json:"role"`                    // User's permission level
	Provider      string    `json:"provider,omitempty"`      // OAuth provider ('google', 'github', empty for local)
	ProviderID    string    `json:"provider_id,omitempty"`   // OAuth provider user ID
	AvatarURL     string    `json:"avatar_url,omitempty"`    // Profile picture URL (from OAuth or uploaded)
	EmailVerified bool      `json:"email_verified"`          // Whether email has been verified
	CreatedAt     time.Time `json:"created_at"`              // Account creation timestamp
}

// Session represents a user authentication session
// Sessions are stored as UUIDs in HTTP cookies and have expiration times
type Session struct {
	ID        string    `json:"id"`         // UUID session identifier
	UserID    int       `json:"user_id"`    // Links to User.ID
	ExpiresAt time.Time `json:"expires_at"` // When this session becomes invalid
	CreatedAt time.Time `json:"created_at"` // Session creation timestamp
}

// ===== CONTENT ORGANIZATION =====

// Category represents a post category for organizing content
type Category struct {
	ID   int    `json:"id"`   // Unique category identifier
	Name string `json:"name"` // Category display name
}

// PostStatus represents the moderation state of a post
type PostStatus string

const (
	PostStatusPending  PostStatus = "pending"  // Awaiting moderation approval
	PostStatusApproved PostStatus = "approved" // Visible to all users
	PostStatusRejected PostStatus = "rejected" // Hidden from public view
)

// Post represents a forum post with all its metadata
// This struct includes computed fields like Likes/Dislikes counts and Username
// that are populated by JOIN queries rather than stored directly
type Post struct {
	ID         int        `json:"id"`                      // Unique post identifier
	UserID     int        `json:"user_id"`                 // Author's user ID
	Title      string     `json:"title"`                   // Post title
	Content    string     `json:"content"`                 // Post body content
	ImagePath  string     `json:"image_path,omitempty"`    // Optional uploaded image path
	Categories []string   `json:"categories"`              // List of category names (populated by JOIN)
	Likes      int        `json:"likes"`                   // Total like count (computed)
	Dislikes   int        `json:"dislikes"`                // Total dislike count (computed)
	Status     PostStatus `json:"status"`                  // Moderation status
	CreatedAt  time.Time  `json:"created_at"`              // Post creation timestamp
	UpdatedAt  time.Time  `json:"updated_at"`              // Last modification timestamp
	Username   string     `json:"username,omitempty"`      // Author's username (populated by JOIN)
	UserLiked  *bool      `json:"user_liked,omitempty"`    // Current user's like status (nil=no vote, true=like, false=dislike)
}

// Comment represents a user comment on a post
// Similar to Post, includes computed fields for likes and username
type Comment struct {
	ID        int       `json:"id"`                      // Unique comment identifier
	PostID    int       `json:"post_id"`                 // Which post this comment belongs to
	UserID    int       `json:"user_id"`                 // Comment author's user ID
	Content   string    `json:"content"`                 // Comment text content
	Likes     int       `json:"likes"`                   // Total like count (computed)
	Dislikes  int       `json:"dislikes"`                // Total dislike count (computed)
	CreatedAt time.Time `json:"created_at"`              // Comment creation timestamp
	UpdatedAt time.Time `json:"updated_at"`              // Last modification timestamp
	Username  string    `json:"username,omitempty"`      // Author's username (populated by JOIN)
	UserLiked *bool     `json:"user_liked,omitempty"`    // Current user's like status
}

// ===== RELATIONSHIP TABLES =====

// PostCategory represents the many-to-many relationship between posts and categories
// This allows posts to be associated with multiple categories
type PostCategory struct {
	PostID     int `json:"post_id"`     // Links to Post.ID
	CategoryID int `json:"category_id"` // Links to Category.ID
}

// ===== VOTING SYSTEM =====

// PostLike represents a user's like or dislike on a post
// The IsLike field determines if it's a like (true) or dislike (false)
type PostLike struct {
	ID     int  `json:"id"`      // Unique like record identifier
	PostID int  `json:"post_id"` // Which post was liked/disliked
	UserID int  `json:"user_id"` // Who made the vote
	IsLike bool `json:"is_like"` // true = like, false = dislike
}

// CommentLike represents a user's like or dislike on a comment
// Same structure as PostLike but for comments
type CommentLike struct {
	ID        int  `json:"id"`         // Unique like record identifier
	CommentID int  `json:"comment_id"` // Which comment was liked/disliked
	UserID    int  `json:"user_id"`    // Who made the vote
	IsLike    bool `json:"is_like"`    // true = like, false = dislike
}

// ===== NOTIFICATION SYSTEM =====

// NotificationType represents the different types of notifications users can receive
// Each type corresponds to a specific forum interaction or moderation action
type NotificationType string

const (
	// User interaction notifications - triggered by other users' actions
	NotificationTypeLike    NotificationType = "like"    // Someone liked user's content
	NotificationTypeDislike NotificationType = "dislike" // Someone disliked user's content
	NotificationTypeComment NotificationType = "comment" // Someone commented on user's post

	// Moderation notifications - triggered by admin/moderator actions
	// Added to provide transparency in the content moderation process
	NotificationTypePostApproved NotificationType = "post_approved" // User's post was approved by moderator/admin
	NotificationTypePostRejected NotificationType = "post_rejected" // User's post was rejected by moderator/admin
)

// Notification represents a user notification about forum activity
// Uses pointer fields (*int) for optional foreign keys to posts/comments
type Notification struct {
	ID        int              `json:"id"`                      // Unique notification identifier
	UserID    int              `json:"user_id"`                 // Who receives this notification
	ActorID   int              `json:"actor_id"`                // Who triggered this notification
	Type      NotificationType `json:"type"`                    // What type of action happened
	PostID    *int             `json:"post_id,omitempty"`       // Optional: related post ID
	CommentID *int             `json:"comment_id,omitempty"`    // Optional: related comment ID
	Message   string           `json:"message"`                 // Human-readable notification text
	Read      bool             `json:"read"`                    // Whether user has seen this notification
	CreatedAt time.Time        `json:"created_at"`              // When the notification was created
	ActorName string           `json:"actor_name,omitempty"`    // Actor's username (populated by JOIN)
}

// ===== MODERATION SYSTEM =====

// ReportStatus represents the state of a content report
type ReportStatus string

const (
	ReportStatusPending  ReportStatus = "pending"  // Awaiting admin review
	ReportStatusReviewed ReportStatus = "reviewed" // Admin has looked at it
	ReportStatusResolved ReportStatus = "resolved" // Issue has been resolved
)

// Report represents a user report about inappropriate content
// Can report either posts or comments (one of PostID/CommentID will be set)
type Report struct {
	ID           int          `json:"id"`                       // Unique report identifier
	ReporterID   int          `json:"reporter_id"`              // Who filed the report
	PostID       *int         `json:"post_id,omitempty"`        // Optional: reported post
	CommentID    *int         `json:"comment_id,omitempty"`     // Optional: reported comment
	Reason       string       `json:"reason"`                   // Category of the issue
	Description  string       `json:"description"`              // Detailed description of the problem
	Status       ReportStatus `json:"status"`                   // Current state of the report
	AdminID      *int         `json:"admin_id,omitempty"`       // Which admin handled this report
	AdminNotes   string       `json:"admin_notes,omitempty"`    // Admin's notes on the resolution
	CreatedAt    time.Time    `json:"created_at"`               // When the report was filed
	UpdatedAt    time.Time    `json:"updated_at"`               // Last modification time
	ReporterName string       `json:"reporter_name,omitempty"`  // Reporter's username (populated by JOIN)
}

// ===== ROLE MANAGEMENT =====

// RoleRequest represents a user's request for role promotion
// Allows users to request moderator status which admins can approve
type RoleRequest struct {
	ID            int       `json:"id"`                      // Unique request identifier
	UserID        int       `json:"user_id"`                 // Who is requesting the role change
	RequestedRole UserRole  `json:"requested_role"`          // Which role they want (usually 'moderator')
	Reason        string    `json:"reason"`                  // Why they want the role
	Status        string    `json:"status"`                  // 'pending', 'approved', 'rejected'
	AdminID       *int      `json:"admin_id,omitempty"`      // Which admin processed this request
	AdminNotes    string    `json:"admin_notes,omitempty"`   // Admin's decision reasoning
	CreatedAt     time.Time `json:"created_at"`              // When the request was made
	UpdatedAt     time.Time `json:"updated_at"`              // Last modification time
	Username      string    `json:"username,omitempty"`      // Requester's username (populated by JOIN)
}

// ===== SECURITY SYSTEM =====

// RateLimitEntry tracks request rates per IP address for DoS protection
// Implements a simple rate limiting system to prevent abuse
type RateLimitEntry struct {
	IP        string    `json:"ip"`         // Client IP address
	Requests  int       `json:"requests"`   // Number of requests made in current window
	ResetTime time.Time `json:"reset_time"` // When the counter resets (sliding window)
}
