package user

type RegisterUserRequest struct {
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Occupation string `json:"occupation" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type CheckEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type RegisterUserResponse struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Occupation string `json:"occupation"`
	Token      string `json:"token"`
}

func NewMapperUserResponse(user User, token string) RegisterUserResponse {
	return RegisterUserResponse{
		ID:         user.Id,
		Name:       user.Name,
		Email:      user.Email,
		Occupation: user.Occupation,
		Token:      token,
	}
}
