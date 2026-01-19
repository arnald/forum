package validator

const (
	MinUsernameLength       = 3
	MaxUsernameLength       = 50
	MinPasswordLength       = 8
	MaxPasswordLength       = 72
	MinTopicTitleLength     = 5
	MaxTopicTitleLength     = 100
	MinTopicContentLength   = 10
	MaxTopicContentLength   = 1000
	MaxPageSize             = 100
	MinCategoryNameLength   = 3
	MaxCategoryNameLength   = 50
	MinCommentContentLength = 1
	MaxCommentContentLength = 1000
)

func ValidateUserRegistration(v *Validator, data any) {
	rules := []ValidationRule{
		{
			Field: "Username",
			Rules: []func(any) (bool, string){
				required,
				minLength(MinUsernameLength),
				maxLength(MaxUsernameLength),
			},
		},
		{
			Field: "Email",
			Rules: []func(any) (bool, string){
				required,
				validEmail,
			},
		},
		{
			Field: "Password",
			Rules: []func(any) (bool, string){
				required,
				minLength(MinPasswordLength),
				maxLength(MaxPasswordLength),
				hasLower,
				hasUpper,
				hasDigit,
				hasSpecial,
			},
		},
	}

	ValidateStruct(v, data, rules)
}

func ValidateUserLoginEmail(v *Validator, data any) {
	rules := []ValidationRule{
		{
			Field: "Email",
			Rules: []func(any) (bool, string){
				required,
				validEmail,
			},
		},
		{
			Field: "Password",
			Rules: []func(any) (bool, string){
				required,
				minLength(MinPasswordLength),
				maxLength(MaxPasswordLength),
			},
		},
	}

	ValidateStruct(v, data, rules)
}

func ValidateUserLoginUsername(v *Validator, data any) {
	rules := []ValidationRule{
		{
			Field: "Username",
			Rules: []func(any) (bool, string){
				required,
				minLength(MinUsernameLength),
				maxLength(MaxUsernameLength),
			},
		},
		{
			Field: "Password",
			Rules: []func(any) (bool, string){
				required,
				minLength(MinPasswordLength),
				maxLength(MaxPasswordLength),
			},
		},
	}

	ValidateStruct(v, data, rules)
}

func ValidateCreateTopic(v *Validator, data any) {
	rules := []ValidationRule{
		{
			Field: "Title",
			Rules: []func(any) (bool, string){
				required,
				minLength(MinTopicTitleLength),
				maxLength(MaxTopicTitleLength),
			},
		},
		{
			Field: "Content",
			Rules: []func(any) (bool, string){
				required,
				minLength(MinTopicContentLength),
				maxLength(MaxTopicContentLength),
			},
		},
		{
			Field: "ImagePath",
			Rules: []func(any) (bool, string){
				optional(validImagePath),
			},
		},
		// TODO: figure out validation with categoryID or Name
		// {
		// 	Field: "Category",
		// 	Rules: []func(any) (bool, string){
		// 		required,
		// 		validCategory,
		// 	},
		// },
	}

	ValidateStruct(v, data, rules)
}

func ValidateGetAllTopics(v *Validator, data any) {
	rules := []ValidationRule{
		{
			Field: "OrderBy",
			Rules: []func(any) (bool, string){
				optional(validOrderBy),
			},
		},
		{
			Field: "Page",
			Rules: []func(any) (bool, string){
				required,
				isPositiveInt,
			},
		},
		{
			Field: "PageSize",
			Rules: []func(any) (bool, string){
				required,
				isPositiveInt,
				maxInt(MaxPageSize),
			},
		},
	}

	ValidateStruct(v, data, rules)
}

func ValidateGetTopic(v *Validator, data any) {
	rules := []ValidationRule{
		{
			Field: "TopicID",
			Rules: []func(any) (bool, string){
				required,
				isPositiveInt,
			},
		},
	}

	ValidateStruct(v, data, rules)
}

func ValidateDeleteTopic(v *Validator, data any) {
	rules := []ValidationRule{
		{
			Field: "TopicID",
			Rules: []func(any) (bool, string){
				required,
				isPositiveInt,
			},
		},
	}

	ValidateStruct(v, data, rules)
}

func ValidateDeleteCategory(v *Validator, data any) {
	rules := []ValidationRule{
		{
			Field: "CategoryID",
			Rules: []func(any) (bool, string){
				required,
				isPositiveInt,
			},
		},
	}

	ValidateStruct(v, data, rules)
}

func ValidateGetCategoryByID(v *Validator, data any) {
	rules := []ValidationRule{
		{
			Field: "CategoryID",
			Rules: []func(any) (bool, string){
				required,
				isPositiveInt,
			},
		},
	}

	ValidateStruct(v, data, rules)
}

func ValidateCreateCategory(v *Validator, data any) {
	rules := []ValidationRule{
		{
			Field: "Name",
			Rules: []func(any) (bool, string){
				required,
				minLength(MinCategoryNameLength),
				maxLength(MaxCategoryNameLength),
			},
		},
	}

	ValidateStruct(v, data, rules)
}

func ValidateUpdateCategory(v *Validator, data any) {
	rules := []ValidationRule{
		{
			Field: "Name",
			Rules: []func(any) (bool, string){
				required,
				minLength(MinCategoryNameLength),
				maxLength(MaxCategoryNameLength),
			},
		},
	}

	ValidateStruct(v, data, rules)
}

func ValidateGetCommentsByTopic(v *Validator, data any) {
	rules := []ValidationRule{
		{
			Field: "TopicID",
			Rules: []func(any) (bool, string){
				required,
				isPositiveInt,
			},
		},
	}

	ValidateStruct(v, data, rules)
}

func ValidateGetComment(v *Validator, data any) {
	rules := []ValidationRule{
		{
			Field: "CommentID",
			Rules: []func(any) (bool, string){
				required,
				isPositiveInt,
			},
		},
	}

	ValidateStruct(v, data, rules)
}

func ValidateUpdateComment(v *Validator, data any) {
	rules := []ValidationRule{
		{
			Field: "CommentID",
			Rules: []func(any) (bool, string){
				required,
				isPositiveInt,
			},
		},
		{
			Field: "Content",
			Rules: []func(any) (bool, string){
				required,
				minLength(MinCommentContentLength),
				maxLength(MaxCommentContentLength),
			},
		},
	}

	ValidateStruct(v, data, rules)
}

func ValidateDeleteComment(v *Validator, data any) {
	rules := []ValidationRule{
		{
			Field: "CommentID",
			Rules: []func(any) (bool, string){
				required,
				isPositiveInt,
			},
		},
	}

	ValidateStruct(v, data, rules)
}

func ValidateCreateComment(v *Validator, data any) {
	rules := []ValidationRule{
		{
			Field: "TopicID",
			Rules: []func(any) (bool, string){
				required,
				isPositiveInt,
			},
		},
		{
			Field: "Content",
			Rules: []func(any) (bool, string){
				required,
				minLength(MinCommentContentLength),
				maxLength(MaxCommentContentLength),
			},
		},
	}

	ValidateStruct(v, data, rules)
}

func ValidateDeleteVote(v *Validator, data any) {
	rules := []ValidationRule{
		{
			Field: "TopicID",
			Rules: []func(any) (bool, string){
				required,
				isPositiveInt,
			},
		},
		{
			Field: "CommentID",
			Rules: []func(any) (bool, string){
				required,
				isPositiveInt,
			},
		},
	}

	ValidateStruct(v, data, rules)
}
