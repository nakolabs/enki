package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type School struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Level     string    `db:"level"`
	CreatedAt int64     `db:"created_at"`
	UpdatedAt int64     `db:"updated_at"`
}

type UserSchoolRole struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	SchoolID  uuid.UUID `db:"school_id"`
	RoleID    uuid.UUID `db:"role_id"`
	CreatedAt int64     `db:"created_at"`
	UpdatedAt int64     `db:"updated_at"`
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

	_, err = tx.NamedExecContext(ctx, `INSERT INTO school (id, name, level, created_at, updated_at) VALUES (:id, :name, :level, :created_at, :updated_at)`, school)
	if err != nil {
		return err
	}

	adminRoleID := uuid.UUID{}
	err = tx.GetContext(ctx, &adminRoleID, `SELECT id FROM role WHERE name = 'admin'`)
	if err != nil {
		return err
	}

	adminSchoolRole := UserSchoolRole{
		UserID:   userID,
		SchoolID: school.ID,
		RoleID:   adminRoleID,
	}

	_, err = tx.NamedExecContext(ctx, `INSERT INTO user_school_role (user_id, school_id, role_id) VALUES (:user_id, :school_id, :role_id)`, adminSchoolRole)
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
	err := r.db.GetContext(ctx, &school, `SELECT id, name, level, created_at, updated_at FROM school WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}

	return &school, nil
}

func (r *repository) GetListSchool(ctx context.Context, userID uuid.UUID) ([]School, error) {
	var schools []School
	err := r.db.SelectContext(ctx, &schools, `SELECT school.* FROM user_school_role JOIN school on user_school_role.school_id = school.id WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	return schools, nil
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
