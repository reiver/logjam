package users

type CreateUserDTO struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	UserName string `json:"username"`
	Bio      string `json:"bio"`
}

type CreateOTPDTO struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required"` // assuming 6-digit numeric
}

type CompleteSignUpDTO struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,numeric"`
	Name  string `json:"name" validate:"required"`
}

type CompleteSignUpResponse struct {
	UserID string `json:"userId"`
	Token  string `json:"token"` // session/jwt/etc
}

type UpdateProfileDTO struct {
	UserID   string  `json:"-"`
	Name     *string `json:"name,omitempty"`
	Username *string `json:"username,omitempty"`
	Bio      *string `json:"bio,omitempty"`
	Avatar   []byte  `json:"avatar"`
}

type UserDTO struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Bio      string `json:"bio"`
	UserName string `json:"userName"`
}
