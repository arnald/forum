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

type UserQueries struct {
	UserRegister      userQueries.UserRegisterRequestHandler
	UserLoginEmail    userQueries.UserLoginEmailRequestHandler
	UserLoginUsername userQueries.UserLoginUsernameRequestHandler
}

type PostQueries struct {
	CreatePost         postQueries.CreatePostRequestHandler
	GetAllPosts        postQueries.GetAllPostsRequestHandler
	GetPostByID        postQueries.GetPostByIDRequestHandler
	GetPostsByCategory postQueries.GetPostsByCategoryRequestHandler
	UpdatePost         postQueries.UpdatePostRequestHandler
}

type CommentQueries struct {
	CreateComment      commentQueries.CreateCommentRequestHandler
	GetCommentsByPost  commentQueries.GetCommentsByPostRequestHandler
	GetCommentTree     commentQueries.GetCommentTreeRequestHandler
	UpdateComment      commentQueries.UpdateCommentRequestHandler
}

type VoteQueries struct {
	CastVote        voteQueries.CastVoteRequestHandler
	GetVoteStatus   voteQueries.GetVoteStatusRequestHandler
	GetUserVotes    voteQueries.GetUserVotesRequestHandler
}

type CategoryQueries struct {
	CreateCategory         categoryQueries.CreateCategoryRequestHandler
	GetAllCategories       categoryQueries.GetAllCategoriesRequestHandler
	GetCategoryByID        categoryQueries.GetCategoryByIDRequestHandler
	UpdateCategory         categoryQueries.UpdateCategoryRequestHandler
	DeleteCategory         categoryQueries.DeleteCategoryRequestHandler
	GetCategoryWithPosts   categoryQueries.GetCategoryWithPostsRequestHandler
}

type UserServices struct {
	Queries UserQueries
}

type PostServices struct {
	Queries PostQueries
}

type CommentServices struct {
	Queries CommentQueries
}

type VoteServices struct {
	Queries VoteQueries
}

type CategoryServices struct {
	Queries CategoryQueries
}

type Services struct {
	UserServices     UserServices
	PostServices     PostServices
	CommentServices  CommentServices
	VoteServices     VoteServices
	CategoryServices CategoryServices
}

func NewServices(userRepo user.Repository, postRepo post.Repository, commentRepo comment.Repository, voteRepo vote.Repository, categoryRepo category.Repository) Services {
	uuidProvider := uuid.NewProvider()
	encryption := bcrypt.NewProvider()
	return Services{
		UserServices: UserServices{
			Queries: UserQueries{
				UserRegister:      userQueries.NewUserRegisterHandler(userRepo, uuidProvider, encryption),
				UserLoginEmail:    userQueries.NewUserLoginEmailHandler(userRepo, encryption),
				UserLoginUsername: userQueries.NewUserLoginUsernameHandler(userRepo, encryption),
			},
		},
		PostServices: PostServices{
			Queries: PostQueries{
				CreatePost:         postQueries.NewCreatePostHandler(postRepo, uuidProvider),
				GetAllPosts:        postQueries.NewGetAllPostsHandler(postRepo),
				GetPostByID:        postQueries.NewGetPostByIDHandler(postRepo),
				GetPostsByCategory: postQueries.NewGetPostsByCategoryHandler(postRepo),
				UpdatePost:         postQueries.NewUpdatePostHandler(postRepo),
			},
		},
		CommentServices: CommentServices{
			Queries: CommentQueries{
				CreateComment:     commentQueries.NewCreateCommentHandler(commentRepo, uuidProvider),
				GetCommentsByPost: commentQueries.NewGetCommentsByPostHandler(commentRepo),
				GetCommentTree:    commentQueries.NewGetCommentTreeHandler(commentRepo),
				UpdateComment:     commentQueries.NewUpdateCommentHandler(commentRepo),
			},
		},
		VoteServices: VoteServices{
			Queries: VoteQueries{
				CastVote:      voteQueries.NewCastVoteHandler(voteRepo, uuidProvider),
				GetVoteStatus: voteQueries.NewGetVoteStatusHandler(voteRepo),
				GetUserVotes:  voteQueries.NewGetUserVotesHandler(voteRepo),
			},
		},
		CategoryServices: CategoryServices{
			Queries: CategoryQueries{
				CreateCategory:       categoryQueries.NewCreateCategoryHandler(categoryRepo, uuidProvider),
				GetAllCategories:     categoryQueries.NewGetAllCategoriesHandler(categoryRepo),
				GetCategoryByID:      categoryQueries.NewGetCategoryByIDHandler(categoryRepo),
				UpdateCategory:       categoryQueries.NewUpdateCategoryHandler(categoryRepo),
				DeleteCategory:       categoryQueries.NewDeleteCategoryHandler(categoryRepo),
				GetCategoryWithPosts: categoryQueries.NewGetCategoryWithPostsHandler(categoryRepo),
			},
		},
	}
}
