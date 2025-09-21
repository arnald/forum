/*
Package app provides the application layer services for the forum application.

This package implements the application layer in Clean Architecture, containing:
- Service aggregation and dependency injection for business logic
- CQRS (Command Query Responsibility Segregation) pattern implementation
- Application service factories and configuration
- Cross-domain operation coordination

The application layer orchestrates domain operations and coordinates between
the infrastructure and domain layers. It contains no business rules but
manages the flow of data and operations.

Organization:
- UserServices: User registration, authentication, and management
- PostServices: Post creation, retrieval, and modification
- CommentServices: Comment operations and threading
- VoteServices: Like/dislike functionality
- CategoryServices: Category management and organization
*/
package app

import (
	categoryQueries "github.com/arnald/forum/internal/app/category/queries"
	commentQueries "github.com/arnald/forum/internal/app/comment/queries"
	postQueries "github.com/arnald/forum/internal/app/post/queries"
	userQueries "github.com/arnald/forum/internal/app/user/queries"
	voteQueries "github.com/arnald/forum/internal/app/vote/queries"
	"github.com/arnald/forum/internal/domain/category"
	"github.com/arnald/forum/internal/domain/comment"
	"github.com/arnald/forum/internal/domain/post"
	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/domain/vote"
	"github.com/arnald/forum/internal/pkg/bcrypt"
	"github.com/arnald/forum/internal/pkg/uuid"
)

// UserQueries aggregates all user-related query handlers
// Implements CQRS pattern for user operations including authentication and registration
type UserQueries struct {
	UserRegister      userQueries.UserRegisterRequestHandler      // Handle user registration requests
	UserLoginEmail    userQueries.UserLoginEmailRequestHandler    // Handle email-based login requests
	UserLoginUsername userQueries.UserLoginUsernameRequestHandler // Handle username-based login requests
}

// PostQueries aggregates all post-related query handlers
// Implements CQRS pattern for post operations including CRUD and filtering
type PostQueries struct {
	CreatePost         postQueries.CreatePostRequestHandler         // Handle post creation requests
	GetAllPosts        postQueries.GetAllPostsRequestHandler        // Handle requests for all posts
	GetPostByID        postQueries.GetPostByIDRequestHandler        // Handle requests for specific posts
	GetPostsByCategory postQueries.GetPostsByCategoryRequestHandler // Handle category-filtered post requests
	UpdatePost         postQueries.UpdatePostRequestHandler         // Handle post modification requests
}

// CommentQueries aggregates all comment-related query handlers
// Implements CQRS pattern for comment operations including threading and CRUD
type CommentQueries struct {
	CreateComment     commentQueries.CreateCommentRequestHandler     // Handle comment creation requests
	GetCommentsByPost commentQueries.GetCommentsByPostRequestHandler // Handle requests for post comments
	GetCommentTree    commentQueries.GetCommentTreeRequestHandler    // Handle hierarchical comment structure requests
	UpdateComment     commentQueries.UpdateCommentRequestHandler     // Handle comment modification requests
}

// VoteQueries aggregates all vote-related query handlers
// Implements CQRS pattern for voting operations including casting and status checking
type VoteQueries struct {
	CastVote      voteQueries.CastVoteRequestHandler      // Handle vote casting requests (like/dislike)
	GetVoteStatus voteQueries.GetVoteStatusRequestHandler // Handle vote status and count requests
	GetUserVotes  voteQueries.GetUserVotesRequestHandler  // Handle user voting history requests
}

// CategoryQueries aggregates all category-related query handlers
// Implements CQRS pattern for category operations including CRUD and post association
type CategoryQueries struct {
	CreateCategory       categoryQueries.CreateCategoryRequestHandler       // Handle category creation requests
	GetAllCategories     categoryQueries.GetAllCategoriesRequestHandler     // Handle requests for all categories
	GetCategoryByID      categoryQueries.GetCategoryByIDRequestHandler      // Handle requests for specific categories
	UpdateCategory       categoryQueries.UpdateCategoryRequestHandler       // Handle category modification requests
	DeleteCategory       categoryQueries.DeleteCategoryRequestHandler       // Handle category deletion requests
	GetCategoryWithPosts categoryQueries.GetCategoryWithPostsRequestHandler // Handle category with posts requests
}

// UserServices encapsulates all user-related application services
// Provides a clean interface for user operations in the application layer
type UserServices struct {
	Queries UserQueries // Query handlers for user operations
}

// PostServices encapsulates all post-related application services
// Provides a clean interface for post operations in the application layer
type PostServices struct {
	Queries PostQueries // Query handlers for post operations
}

// CommentServices encapsulates all comment-related application services
// Provides a clean interface for comment operations in the application layer
type CommentServices struct {
	Queries CommentQueries // Query handlers for comment operations
}

// VoteServices encapsulates all vote-related application services
// Provides a clean interface for voting operations in the application layer
type VoteServices struct {
	Queries VoteQueries // Query handlers for vote operations
}

// CategoryServices encapsulates all category-related application services
// Provides a clean interface for category operations in the application layer
type CategoryServices struct {
	Queries CategoryQueries // Query handlers for category operations
}

// Services aggregates all application services for dependency injection
// This is the main application layer structure used throughout the application
// It provides access to all business logic operations through a single interface
type Services struct {
	UserServices     UserServices     // User-related operations (auth, registration)
	PostServices     PostServices     // Post-related operations (CRUD, filtering)
	CommentServices  CommentServices  // Comment-related operations (CRUD, threading)
	VoteServices     VoteServices     // Vote-related operations (like/dislike)
	CategoryServices CategoryServices // Category-related operations (organization)
}

// NewServices creates and configures all application services with their dependencies
// This factory function implements dependency injection for the application layer
//
// Parameters:
//   - userRepo: Repository for user data operations
//   - postRepo: Repository for post data operations
//   - commentRepo: Repository for comment data operations
//   - voteRepo: Repository for vote data operations
//   - categoryRepo: Repository for category data operations
//
// Returns:
//   - Services: Fully configured application services ready for use
//
// The function:
// 1. Creates shared utilities (UUID generation, encryption)
// 2. Initializes all service handlers with appropriate repositories
// 3. Configures dependency injection for consistent operation
func NewServices(userRepo user.Repository, postRepo post.Repository, commentRepo comment.Repository, voteRepo vote.Repository, categoryRepo category.Repository) Services {
	// Create shared utility providers used across services
	uuidProvider := uuid.NewProvider() // UUID generation for entity IDs
	encryption := bcrypt.NewProvider()  // Password encryption for user security

	return Services{
		// User services for authentication and user management
		UserServices: UserServices{
			Queries: UserQueries{
				UserRegister:      userQueries.NewUserRegisterHandler(userRepo, uuidProvider, encryption),      // Registration with UUID and encryption
				UserLoginEmail:    userQueries.NewUserLoginEmailHandler(userRepo, encryption),                  // Email login with password verification
				UserLoginUsername: userQueries.NewUserLoginUsernameHandler(userRepo, encryption),               // Username login with password verification
			},
		},
		// Post services for content creation and management
		PostServices: PostServices{
			Queries: PostQueries{
				CreatePost:         postQueries.NewCreatePostHandler(postRepo, uuidProvider),         // Post creation with UUID generation
				GetAllPosts:        postQueries.NewGetAllPostsHandler(postRepo),                      // Retrieve all posts for listing
				GetPostByID:        postQueries.NewGetPostByIDHandler(postRepo),                      // Retrieve specific posts by ID
				GetPostsByCategory: postQueries.NewGetPostsByCategoryHandler(postRepo),               // Retrieve posts filtered by category
				UpdatePost:         postQueries.NewUpdatePostHandler(postRepo),                       // Post modification operations
			},
		},
		// Comment services for discussion functionality
		CommentServices: CommentServices{
			Queries: CommentQueries{
				CreateComment:     commentQueries.NewCreateCommentHandler(commentRepo, uuidProvider), // Comment creation with UUID generation
				GetCommentsByPost: commentQueries.NewGetCommentsByPostHandler(commentRepo),           // Retrieve comments for specific posts
				GetCommentTree:    commentQueries.NewGetCommentTreeHandler(commentRepo),              // Retrieve hierarchical comment structure
				UpdateComment:     commentQueries.NewUpdateCommentHandler(commentRepo),               // Comment modification operations
			},
		},
		// Vote services for like/dislike functionality
		VoteServices: VoteServices{
			Queries: VoteQueries{
				CastVote:      voteQueries.NewCastVoteHandler(voteRepo, uuidProvider),      // Vote casting with UUID generation
				GetVoteStatus: voteQueries.NewGetVoteStatusHandler(voteRepo),               // Vote status and count retrieval
				GetUserVotes:  voteQueries.NewGetUserVotesHandler(voteRepo),                // User voting history retrieval
			},
		},
		// Category services for content organization
		CategoryServices: CategoryServices{
			Queries: CategoryQueries{
				CreateCategory:       categoryQueries.NewCreateCategoryHandler(categoryRepo, uuidProvider),       // Category creation with UUID generation
				GetAllCategories:     categoryQueries.NewGetAllCategoriesHandler(categoryRepo),                   // Retrieve all categories for browsing
				GetCategoryByID:      categoryQueries.NewGetCategoryByIDHandler(categoryRepo),                    // Retrieve specific categories by ID
				UpdateCategory:       categoryQueries.NewUpdateCategoryHandler(categoryRepo),                     // Category modification operations
				DeleteCategory:       categoryQueries.NewDeleteCategoryHandler(categoryRepo),                     // Category deletion operations
				GetCategoryWithPosts: categoryQueries.NewGetCategoryWithPostsHandler(categoryRepo),               // Retrieve categories with associated posts
			},
		},
	}
}
