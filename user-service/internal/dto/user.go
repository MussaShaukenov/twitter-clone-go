package dto

// DTO for user registration
type RegisterUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Age       int    `json:"age,omitempty"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

// DTO for login credentials
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// DTO for registration response
type RegisterUserResponse struct {
	Message string `json:"message"`
	UserID  int    `json:"user_id"`
}

// DTO for authorization response
type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type VerifyOTPRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

type FollowRequest struct {
	FollowerID int `json:"follower_id"`
	FollowedID int `json:"followed_id"`
}
