package repository

import (
	"context"
	"enuma-elish/internal/ppdb/service/data/request"
	"enuma-elish/internal/ppdb/service/data/response"
	commonHttp "enuma-elish/pkg/http"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Repository interface {
	CreatePPDB(ctx context.Context, ppdb *PPDB) error
	UpdatePPDB(ctx context.Context, ppdb *PPDB) error
	DeletePPDB(ctx context.Context, id uuid.UUID) error
	GetPPDBByID(ctx context.Context, id uuid.UUID) (*PPDB, error)
	GetListPPDB(ctx context.Context, query request.GetListPPDBQuery) ([]response.PPDBResponse, *commonHttp.Meta, error)

	RegisterPPDB(ctx context.Context, ppdbStudent *PPDBStudent) error
	GetPPDBRegistrants(ctx context.Context, query request.GetPPDBRegistrantsQuery) ([]response.PPDBStudentResponse, *commonHttp.Meta, error)
	UpdatePPDBStudentStatus(ctx context.Context, ppdbID uuid.UUID, studentIDs []uuid.UUID, status string) error
	GetPPDBStudentByPPDBIDAndEmail(ctx context.Context, ppdbID uuid.UUID, email string) (*PPDBStudent, error)
	GetPPDBStudentByPPDBIDAndUserID(ctx context.Context, ppdbID uuid.UUID, userID uuid.UUID) (*PPDBStudent, error) // New method
}

type PPDB struct {
	ID        uuid.UUID `db:"id"`
	SchoolID  uuid.UUID `db:"school_id"`
	StartAt   int64     `db:"start_at"`
	EndAt     int64     `db:"end_at"`
	CreatedAt int64     `db:"created_at"`
	UpdatedAt int64     `db:"updated_at"`
}

type PPDBStudent struct {
	ID        uuid.UUID `db:"id"`
	PPDBID    uuid.UUID `db:"ppdb_id"`
	StudentID uuid.UUID `db:"student_id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Status    string    `db:"status"` // registered, accepted, rejected
	CreatedAt int64     `db:"created_at"`
	UpdatedAt int64     `db:"updated_at"`
}

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreatePPDB(ctx context.Context, ppdb *PPDB) error {
	query := `INSERT INTO ppdb (id, school_id, start_at, end_at, created_at, updated_at) 
			  VALUES (:id, :school_id, :start_at, :end_at, :created_at, :updated_at)`

	_, err := r.db.NamedExecContext(ctx, query, ppdb)
	return err
}

func (r *repository) UpdatePPDB(ctx context.Context, ppdb *PPDB) error {
	query := `UPDATE ppdb SET start_at = :start_at, end_at = :end_at, updated_at = :updated_at WHERE id = :id`

	_, err := r.db.NamedExecContext(ctx, query, ppdb)
	return err
}

func (r *repository) DeletePPDB(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM ppdb WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *repository) GetPPDBByID(ctx context.Context, id uuid.UUID) (*PPDB, error) {
	var ppdb PPDB
	query := `SELECT id, school_id, start_at, end_at, created_at, updated_at FROM ppdb WHERE id = $1`

	err := r.db.GetContext(ctx, &ppdb, query, id)
	if err != nil {
		return nil, err
	}
	return &ppdb, nil
}

func (r *repository) GetListPPDB(ctx context.Context, query request.GetListPPDBQuery) ([]response.PPDBResponse, *commonHttp.Meta, error) {
	httpQuery, filters := query.Get()

	baseWhere := "WHERE school_id = :school_id"
	params := map[string]interface{}{
		"school_id": filters["school_id"],
		"limit":     httpQuery.PageSize,
		"offset":    httpQuery.GetOffset(),
	}

	now := time.Now().UnixMilli()
	if status, ok := filters["status"].(string); ok && status != "all" {
		if status == "active" {
			baseWhere += " AND start_at <= :now AND end_at >= :now"
			params["now"] = now
		} else if status == "inactive" {
			baseWhere += " AND (start_at > :now OR end_at < :now)"
			params["now"] = now
		}
	}

	countQuery := "SELECT COUNT(*) FROM ppdb " + baseWhere
	rows, err := r.db.NamedQueryContext(ctx, countQuery, params)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var total int
	if rows.Next() {
		err = rows.Scan(&total)
		if err != nil {
			return nil, nil, err
		}
	}

	selectQuery := `SELECT id, school_id, start_at, end_at, created_at, updated_at 
					FROM ppdb ` + baseWhere + ` 
					ORDER BY created_at DESC 
					LIMIT :limit OFFSET :offset`

	var ppdbs []PPDB
	err = r.db.SelectContext(ctx, &ppdbs, r.db.Rebind(selectQuery),
		params["school_id"], params["now"], params["limit"], params["offset"])
	if err != nil {
		return nil, nil, err
	}

	var result []response.PPDBResponse
	for _, ppdb := range ppdbs {
		status := "inactive"
		if ppdb.StartAt <= now && ppdb.EndAt >= now {
			status = "active"
		}

		result = append(result, response.PPDBResponse{
			ID:        ppdb.ID,
			SchoolID:  ppdb.SchoolID,
			StartAt:   ppdb.StartAt,
			EndAt:     ppdb.EndAt,
			Status:    status,
			CreatedAt: ppdb.CreatedAt,
			UpdatedAt: ppdb.UpdatedAt,
		})
	}

	meta := commonHttp.NewMetaFromQuery(httpQuery, total)
	return result, meta, nil
}

func (r *repository) RegisterPPDB(ctx context.Context, ppdbStudent *PPDBStudent) error {
	query := `INSERT INTO ppdb_student (id, ppdb_id, student_id, name, email, status, created_at, updated_at) 
			  VALUES (:id, :ppdb_id, :student_id, :name, :email, :status, :created_at, :updated_at)`

	_, err := r.db.NamedExecContext(ctx, query, ppdbStudent)
	return err
}

func (r *repository) GetPPDBRegistrants(ctx context.Context, query request.GetPPDBRegistrantsQuery) ([]response.PPDBStudentResponse, *commonHttp.Meta, error) {
	httpQuery, filters := query.Get()

	baseWhere := "WHERE ppdb_id = :ppdb_id"
	params := map[string]interface{}{
		"ppdb_id": filters["ppdb_id"],
		"limit":   httpQuery.PageSize,
		"offset":  httpQuery.GetOffset(),
	}

	if status, ok := filters["status"].(string); ok && status != "all" {
		baseWhere += " AND status = :status"
		params["status"] = status
	}

	countQuery := "SELECT COUNT(*) FROM ppdb_student " + baseWhere
	rows, err := r.db.NamedQueryContext(ctx, countQuery, params)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var total int
	if rows.Next() {
		err = rows.Scan(&total)
		if err != nil {
			return nil, nil, err
		}
	}

	selectQuery := `SELECT id, ppdb_id, student_id, name, email, status, created_at, updated_at 
					FROM ppdb_student ` + baseWhere + ` 
					ORDER BY created_at DESC 
					LIMIT :limit OFFSET :offset`

	var args []interface{}
	var reboundQuery string

	if _, hasStatus := params["status"]; hasStatus {
		reboundQuery = r.db.Rebind(selectQuery)
		args = []interface{}{params["ppdb_id"], params["status"], params["limit"], params["offset"]}
	} else {
		reboundQuery = r.db.Rebind(selectQuery)
		args = []interface{}{params["ppdb_id"], params["limit"], params["offset"]}
	}

	var registrants []PPDBStudent
	err = r.db.SelectContext(ctx, &registrants, reboundQuery, args...)
	if err != nil {
		return nil, nil, err
	}

	var result []response.PPDBStudentResponse
	for _, registrant := range registrants {
		result = append(result, response.PPDBStudentResponse{
			ID:        registrant.ID,
			PPDBID:    registrant.PPDBID,
			StudentID: registrant.StudentID,
			Name:      registrant.Name,
			Email:     registrant.Email,
			Status:    registrant.Status,
			CreatedAt: registrant.CreatedAt,
			UpdatedAt: registrant.UpdatedAt,
		})
	}

	meta := commonHttp.NewMetaFromQuery(httpQuery, total)
	return result, meta, nil
}

func (r *repository) UpdatePPDBStudentStatus(ctx context.Context, ppdbID uuid.UUID, studentIDs []uuid.UUID, status string) error {
	if len(studentIDs) == 0 {
		return nil // No student IDs to update
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE ppdb_student SET status = :status, updated_at = :updated_at 
			  WHERE ppdb_id = :ppdb_id AND student_id = ANY(:student_ids)`

	params := map[string]interface{}{
		"status":      status,
		"updated_at":  time.Now().UnixMilli(),
		"ppdb_id":     ppdbID,
		"student_ids": pq.Array(studentIDs),
	}

	_, err = tx.NamedExecContext(ctx, query, params)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *repository) GetPPDBStudentByPPDBIDAndEmail(ctx context.Context, ppdbID uuid.UUID, email string) (*PPDBStudent, error) {
	var ppdbStudent PPDBStudent
	query := `SELECT id, ppdb_id, student_id, name, email, status, created_at, updated_at 
			  FROM ppdb_student WHERE ppdb_id = $1 AND email = $2`

	err := r.db.GetContext(ctx, &ppdbStudent, query, ppdbID, email)
	if err != nil {
		return nil, err
	}
	return &ppdbStudent, nil
}

func (r *repository) GetPPDBStudentByPPDBIDAndUserID(ctx context.Context, ppdbID uuid.UUID, userID uuid.UUID) (*PPDBStudent, error) {
	var ppdbStudent PPDBStudent
	query := `SELECT id, ppdb_id, student_id, name, email, status, created_at, updated_at 
			  FROM ppdb_student WHERE ppdb_id = $1 AND student_id = $2`

	err := r.db.GetContext(ctx, &ppdbStudent, query, ppdbID, userID)
	if err != nil {
		return nil, err
	}
	return &ppdbStudent, nil
}
