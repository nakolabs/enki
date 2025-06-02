package response

type UserProfileResponse struct {
	ID         string           `json:"id"`
	Name       string           `json:"name"`
	Email      string           `json:"email"`
	IsVerified bool             `json:"is_verified"`
	CreatedAt  int64            `json:"created_at"`
	UpdatedAt  int64            `json:"updated_at"`
	Profile    *ProfileResponse `json:"profile"`
}
