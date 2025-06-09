package response

import (
	"github.com/google/uuid"
)

type ListSchool []ListSchoolItem
type ListSchoolItem struct {
	School
	StudentCount int `json:"student_count"`
	TeacherCount int `json:"teacher_count"`
}

type School struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Level       string    `json:"level"`
	Description string    `json:"description"`
	Address     string    `json:"address"`
	City        string    `json:"city"`
	Province    string    `json:"province"`
	PostalCode  string    `json:"postal_code"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	Website     string    `json:"website"`
	Logo        string    `json:"logo"`
	Banner      string    `json:"banner"`
	Status      string    `json:"status"`
	CreatedAt   int64     `json:"created_at"`
	CreatedBy   string    `json:"created_by,omitempty"`
	UpdatedAt   int64     `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by,omitempty"`
}

type DetailSchool struct {
	School
	Statistics SchoolStatistics `json:"statistics,omitempty"`
}

type SchoolStatistics struct {
	StudentCount    int     `json:"student_count"`
	TeacherCount    int     `json:"teacher_count"`
	ClassCount      int     `json:"class_count"`
	SubjectCount    int     `json:"subject_count"`
	ExamCount       int     `json:"exam_count"`
	PendingStudents int     `json:"pending_students"`
	TeacherRatio    float64 `json:"teacher_ratio"`
	AvgClassSize    float64 `json:"avg_class_size"`
}

type ListSchoolStatistics struct {
	TotalSchools  int `json:"total_schools"`
	TotalStudents int `json:"total_students"`
	TotalTeachers int `json:"total_teachers"`
	ActiveSchools int `json:"active_schools"`
}
