package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, u *User) error
	CreateVerifyEmailToken(ctx context.Context, u *UserVerifyEmailToken) error
	VerifyEmailToken(ctx context.Context, email string) (*UserVerifyEmailToken, error)
	Redis() *redis.Client
	GetFirstUserSchoolRolByUserID(ctx context.Context, userID uuid.UUID) (*UserSchoolRole, error)
	CreateForgotPasswordToken(ctx context.Context, u *UserForgotPasswordToken) error
	VerifyForgotPasswordToken(ctx context.Context, email string) (*UserForgotPasswordToken, error)
	UpdatePassword(ctx context.Context, email, password string) error
	UpdateProfile(ctx context.Context, userID uuid.UUID, profile *Profile) (*Profile, error)
	GetProfileByUserID(ctx context.Context, userID uuid.UUID) (*Profile, error)
}

type repository struct {
	db  *sqlx.DB
	rdb *redis.Client
}

func New(db *sqlx.DB, rdb *redis.Client) Repository {
	return &repository{db: db, rdb: rdb}
}

type Profile struct {
	ID          uuid.UUID `db:"id"`
	UserID      uuid.UUID `db:"user_id"`
	FirstName   *string   `db:"first_name"`
	LastName    *string   `db:"last_name"`
	Phone       *string   `db:"phone"`
	DateOfBirth *string   `db:"date_of_birth"` // Using string for DATE type
	Gender      *string   `db:"gender"`
	Address     *string   `db:"address"`
	City        *string   `db:"city"`
	Country     *string   `db:"country"`
	Avatar      *string   `db:"avatar"`
	Bio         *string   `db:"bio"`
	ParentName  *string   `db:"parent_name"`
	ParentPhone *string   `db:"parent_phone"`
	ParentEmail *string   `db:"parent_email"`
	CreatedAt   int64     `db:"created_at"`
	UpdatedAt   int64     `db:"updated_at"`
}

func (r *repository) UpdateProfile(ctx context.Context, userID uuid.UUID, profile *Profile) (*Profile, error) {
	now := getCurrentTimestamp()

	// Use upsert to handle both create and update cases
	query := `INSERT INTO profiles (id, user_id, first_name, last_name, phone, date_of_birth, gender, address, city, country, avatar, bio, parent_name, parent_phone, parent_email, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
			  ON CONFLICT (user_id) 
			  DO UPDATE SET 
				first_name = EXCLUDED.first_name,
				last_name = EXCLUDED.last_name,
				phone = EXCLUDED.phone,
				date_of_birth = EXCLUDED.date_of_birth,
				gender = EXCLUDED.gender,
				address = EXCLUDED.address,
				city = EXCLUDED.city,
				country = EXCLUDED.country,
				avatar = EXCLUDED.avatar,
				bio = EXCLUDED.bio,
				parent_name = EXCLUDED.parent_name,
				parent_phone = EXCLUDED.parent_phone,
				parent_email = EXCLUDED.parent_email,
				updated_at = EXCLUDED.updated_at
			  RETURNING id, user_id, first_name, last_name, phone, date_of_birth, gender, address, city, country, avatar, bio, parent_name, parent_phone, parent_email, created_at, updated_at`

	var updatedProfile Profile
	err := r.db.GetContext(ctx, &updatedProfile, query,
		uuid.New(), // new ID for insert case
		userID,
		profile.FirstName,
		profile.LastName,
		profile.Phone,
		profile.DateOfBirth,
		profile.Gender,
		profile.Address,
		profile.City,
		profile.Country,
		profile.Avatar,
		profile.Bio,
		profile.ParentName,
		profile.ParentPhone,
		profile.ParentEmail,
		now, // created_at
		now, // updated_at
	)
	if err != nil {
		return nil, err
	}
	return &updatedProfile, nil
}

func (r *repository) GetProfileByUserID(ctx context.Context, userID uuid.UUID) (*Profile, error) {
	var profile Profile
	query := `SELECT id, user_id, first_name, last_name, phone, date_of_birth, gender, address, city, country, avatar, bio, parent_name, parent_phone, parent_email, created_at, updated_at FROM profiles WHERE user_id = $1`
	err := r.db.GetContext(ctx, &profile, query, userID)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

// Helper function following project convention
func getCurrentTimestamp() int64 {
	return time.Now().UnixMilli()
}
