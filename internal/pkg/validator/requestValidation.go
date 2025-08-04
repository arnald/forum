package validator

const (
	MinUsernameLength = 3
	MaxUsernameLength = 50
	MinPasswordLength = 8
	MaxPasswordLength = 72
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
			},
		},
	}

	ValidateStruct(v, data, rules)
}
