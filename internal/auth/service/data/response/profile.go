package response

type ProfileResponse struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	Phone       *string `json:"phone"`
	DateOfBirth *string `json:"date_of_birth"`
	Gender      *string `json:"gender"`
	Address     *string `json:"address"`
	City        *string `json:"city"`
	Country     *string `json:"country"`
	Avatar      *string `json:"avatar"`
	Bio         *string `json:"bio"`
	ParentName  *string `json:"parent_name"`
	ParentPhone *string `json:"parent_phone"`
	ParentEmail *string `json:"parent_email"`
	CreatedAt   int64   `json:"created_at"`
	UpdatedAt   int64   `json:"updated_at"`
}
