package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func (r *repository) UpdateStudentClass(ctx context.Context, studentID, oldClassID, newClassID uuid.UUID) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	committed := false
	defer func() {
		if !committed {
			if err := tx.Rollback(); err != nil {
				log.Error().Err(err).Msg("error rolling back transaction")
			}
		}
	}()

	// Update the class_student record
	updateQuery := `UPDATE class_student 
					SET class_id = $1, updated_at = $2 
					WHERE student_id = $3 AND class_id = $4`

	now := time.Now().UnixMilli()
	result, err := tx.ExecContext(ctx, updateQuery, newClassID, now, studentID, oldClassID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no class assignment found for student %s in class %s", studentID, oldClassID)
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	committed = true

	return nil
}
