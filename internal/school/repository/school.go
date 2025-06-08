package repository

import (
	"context"
	"database/sql"
	"enuma-elish/internal/school/service/data/request"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type School struct {
	ID          uuid.UUID      `db:"id"`
	Name        string         `db:"name"`
	Level       string         `db:"level"`
	Description string         `db:"description"`
	Address     string         `db:"address"`
	City        string         `db:"city"`
	Province    string         `db:"province"`
	PostalCode  string         `db:"postal_code"`
	Phone       string         `db:"phone"`
	Email       string         `db:"email"`
	Website     string         `db:"website"`
	Logo        string         `db:"logo"`
	Banner      string         `db:"banner"`
	CreatedAt   int64          `db:"created_at"`
	CreatedBy   uuid.UUID      `db:"created_by"`
	UpdatedAt   int64          `db:"updated_at"`
	UpdatedBy   sql.NullString `db:"updated_by"`
	DeletedAt   int64          `db:"deleted_at"`
	DeletedBy   sql.NullString `db:"deleted_by"`
}

type UserSchoolRole struct {
	ID        uuid.UUID      `db:"id"`
	UserID    uuid.UUID      `db:"user_id"`
	SchoolID  uuid.UUID      `db:"school_id"`
	RoleID    string         `db:"role_id"`
	IsDeleted bool           `db:"is_deleted"`
	CreatedAt int64          `db:"created_at"`
	CreatedBy uuid.UUID      `db:"created_by"`
	UpdatedAt int64          `db:"updated_at"`
	UpdatedBy sql.NullString `db:"updated_by"`
	DeletedAt int64          `db:"deleted_at"`
	DeletedBy sql.NullString `db:"deleted_by"`
}

func (r *repository) CreateSchool(ctx context.Context, userID uuid.UUID, school School) error {

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func(tx *sqlx.Tx) {
		err := tx.Rollback()
		if err != nil {
			log.Err(err).Msg("error rolling back transaction")
		}
	}(tx)

	_, err = tx.NamedExecContext(ctx, `INSERT INTO school (id, name, level, created_at, created_by) VALUES (:id, :name, :level, :created_at, :created_by)`, school)
	if err != nil {
		return err
	}

	adminSchoolRole := UserSchoolRole{
		UserID:    userID,
		SchoolID:  school.ID,
		RoleID:    "admin",
		CreatedBy: userID,
	}

	_, err = tx.NamedExecContext(ctx, `INSERT INTO user_school_role (user_id, school_id, role_id, created_by) VALUES (:user_id, :school_id, :role_id, :created_by)`, adminSchoolRole)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetSchoolByID(ctx context.Context, id uuid.UUID) (*School, error) {
	school := School{}
	err := r.db.GetContext(ctx, &school, `SELECT id, name, level, description, address, city, province, postal_code, phone, email, website, logo, banner, created_at, updated_at FROM school WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}

	return &school, nil
}

func (r *repository) GetListSchool(ctx context.Context, userID uuid.UUID, httpQuery request.GetListSchoolQuery) ([]School, int, error) {
	var schools []School
	selectQuery := `SELECT school.id, school.name, school.level, school.description, school.address, school.city, school.province, school.postal_code, school.phone, school.email, school.website, school.logo, school.banner, school.created_at, school.updated_at FROM user_school_role JOIN school on user_school_role.school_id = school.id WHERE user_id = ?`
	countQuery := `SELECT COUNT(*) FROM user_school_role JOIN school on user_school_role.school_id = school.id WHERE user_id = ?`

	filterParams := []interface{}{userID}
	filterQuery := ""

	if httpQuery.Search != "" && len(httpQuery.SearchBy) > 0 {
		filterQuery += " AND ("
		for i, v := range httpQuery.SearchBy {
			if i > 0 {
				filterQuery += " OR "
			}
			filterQuery += fmt.Sprintf("school.%s LIKE ?", v)
			filterParams = append(filterParams, "%"+httpQuery.Search+"%")
		}
		filterQuery += ")"
	}

	if httpQuery.Level != "" {
		filterQuery += " AND school.level = ?"
		filterParams = append(filterParams, httpQuery.Level)
	}

	if httpQuery.StartDate > 0 && httpQuery.EndDate > 0 {
		filterQuery += " AND school.created_at BETWEEN ? AND ?"
		filterParams = append(filterParams, httpQuery.StartDate, httpQuery.EndDate)
	} else if httpQuery.StartDate > 0 {
		filterQuery += " AND school.created_at >= ?"
		filterParams = append(filterParams, httpQuery.StartDate)
	} else if httpQuery.EndDate > 0 {
		filterQuery += " AND school.created_at <= ?"
		filterParams = append(filterParams, httpQuery.EndDate)
	}

	limitOrderQuery := fmt.Sprintf(" ORDER BY school.%s %s LIMIT ? OFFSET ?", httpQuery.OrderBy, httpQuery.Order)
	limitOrderParams := []interface{}{httpQuery.PageSize, httpQuery.GetOffset()}

	selectParams := append(filterParams, limitOrderParams...)
	err := r.db.SelectContext(ctx, &schools, r.db.Rebind(selectQuery+filterQuery+limitOrderQuery), selectParams...)
	if err != nil {
		return nil, 0, err
	}

	total := 0
	err = r.db.GetContext(ctx, &total, r.db.Rebind(countQuery+filterQuery), filterParams...)
	if err != nil {
		return nil, 0, err
	}

	return schools, total, nil
}

func (r *repository) GetSchoolRoleByUserIDAndSchoolID(ctx context.Context, userID uuid.UUID, schoolID uuid.UUID) (*UserSchoolRole, error) {
	schoolRole := UserSchoolRole{}
	err := r.db.GetContext(ctx, &schoolRole, "SELECT id, user_id, school_id, role_id, created_at, updated_at FROM user_school_role WHERE user_id = $1 AND school_id = $2", userID, schoolID)
	if err != nil {
		return nil, err
	}

	return &schoolRole, nil
}

func (r *repository) DeleteSchool(ctx context.Context, schoolID uuid.UUID) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func(tx *sqlx.Tx) {
		err := tx.Rollback()
		if err != nil {
			log.Err(err).Msg("error rolling back transaction")
		}
	}(tx)

	// Delete related records first (foreign key constraints)
	_, err = tx.ExecContext(ctx, "DELETE FROM user_school_role WHERE school_id = $1", schoolID)
	if err != nil {
		return err
	}

	// Delete the school
	_, err = tx.ExecContext(ctx, "DELETE FROM school WHERE id = $1", schoolID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateSchoolProfile(ctx context.Context, schoolID uuid.UUID, school School) (*School, error) {
	now := time.Now().UnixMilli()
	school.UpdatedAt = now

	query := `UPDATE school SET 
		name = $2, level = $3, description = $4, address = $5, city = $6, 
		province = $7, postal_code = $8, phone = $9, email = $10, 
		website = $11, logo = $12, updated_at = $13 
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		schoolID, school.Name, school.Level, school.Description, school.Address,
		school.City, school.Province, school.PostalCode, school.Phone,
		school.Email, school.Website, school.Logo, school.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return r.GetSchoolByID(ctx, schoolID)
}
