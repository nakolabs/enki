package service

import (
	"context"
	"enuma-elish/internal/school/repository"
	"enuma-elish/internal/school/service/data/request"
	"enuma-elish/internal/school/service/data/response"
	commonHttp "enuma-elish/pkg/http"
	"enuma-elish/pkg/jwt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func (s *service) CreatSchool(ctx context.Context, data request.CreateSchoolRequest) error {

	claim, err := jwt.ExtractContext(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to extract claims")
		return err
	}

	now := time.Now().UnixMilli()
	school := repository.School{
		ID:          uuid.New(),
		Name:        data.Name,
		Level:       data.Level,
		Description: data.Description,
		Address:     data.Address,
		City:        data.City,
		Province:    data.Province,
		PostalCode:  data.PostalCode,
		Phone:       data.Phone,
		Email:       data.Email,
		Website:     data.Website,
		Logo:        data.Logo,
		Banner:      data.Banner,
		CreatedAt:   now,
		CreatedBy:   claim.User.ID,
		UpdatedAt:   0,
	}

	err = s.repository.CreateSchool(ctx, claim.User.ID, school)
	if err != nil {
		log.Err(err).Msg("error creating school")
		return err
	}

	return nil
}

func (s *service) GetDetailSchool(ctx context.Context, schoolID uuid.UUID) (response.DetailSchool, error) {
	var detailSchool response.DetailSchool
	school, err := s.repository.GetSchoolByID(ctx, schoolID)
	if err != nil {
		return detailSchool, err
	}

	user, err := s.repository.GetUserByID(ctx, school.CreatedBy)
	if err != nil {
		log.Err(err).Msg("error getting user by ID")
		return detailSchool, err
	}

	updatedBy := ""
	if school.UpdatedBy.Valid {
		updatedByUser, err := s.repository.GetUserByID(ctx, school.UpdatedBy.UUID)
		if err != nil {
			log.Err(err).Msg("error getting updated by user")
			return detailSchool, err
		}
		updatedBy = updatedByUser.Name
	}

	// Get school statistics
	stats, err := s.repository.GetSchoolStatistics(ctx, schoolID)
	if err != nil {
		log.Err(err).Msg("error getting school statistics")
		return detailSchool, err
	}

	teacherRation := 0.0
	avgClassSize := 0.0
	if stats.TeacherCount > 0 {
		teacherRation = math.Round(float64(stats.StudentCount) / float64(stats.TeacherCount))
	}
	if stats.ClassCount > 0 {
		avgClassSize = math.Round(float64(stats.StudentCount) / float64(stats.ClassCount))
	}

	detailSchool = response.DetailSchool{
		School: response.School{
			ID:          school.ID,
			Name:        school.Name,
			Level:       school.Level,
			Description: school.Description,
			Address:     school.Address,
			City:        school.City,
			Province:    school.Province,
			PostalCode:  school.PostalCode,
			Phone:       school.Phone,
			Email:       school.Email,
			Website:     school.Website,
			Logo:        school.Logo,
			Banner:      school.Banner,
			CreatedAt:   school.CreatedAt,
			CreatedBy:   user.Name,
			UpdatedAt:   school.UpdatedAt,
			UpdatedBy:   updatedBy,
			Status:      s.setStatus(school.DeletedAt),
		},
		Statistics: response.SchoolStatistics{
			StudentCount:    stats.StudentCount,
			TeacherCount:    stats.TeacherCount,
			ClassCount:      stats.ClassCount,
			SubjectCount:    stats.SubjectCount,
			ExamCount:       stats.ExamCount,
			PendingStudents: stats.PendingStudents,
			TeacherRatio:    teacherRation,
			AvgClassSize:    avgClassSize,
		},
	}

	return detailSchool, nil
}

func (s *service) GetListSchool(ctx context.Context, httpQuery request.GetListSchoolQuery) (response.ListSchool, *commonHttp.Meta, error) {
	claim, err := jwt.ExtractContext(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to extract claims")
		return response.ListSchool{}, nil, err
	}

	listSchool, total, err := s.repository.GetListSchool(ctx, claim.User.ID, httpQuery)
	if err != nil {
		log.Err(err).Msg("error getting list school")
		return response.ListSchool{}, nil, err
	}

	// Get student and teacher counts for each school
	schoolCounts := make(map[uuid.UUID]repository.SchoolCounts)
	if len(listSchool) > 0 {
		var schoolIDs []uuid.UUID
		for _, school := range listSchool {
			schoolIDs = append(schoolIDs, school.ID)
		}

		counts, err := s.repository.GetSchoolCounts(ctx, schoolIDs)
		if err != nil {
			log.Err(err).Msg("failed to get school counts")
			return response.ListSchool{}, nil, err
		}
		schoolCounts = counts
	}

	res := response.ListSchool{}
	for _, school := range listSchool {
		res = append(res, response.ListSchoolItem{
			School: response.School{
				ID:          school.ID,
				Name:        school.Name,
				Level:       school.Level,
				Description: school.Description,
				Address:     school.Address,
				City:        school.City,
				Province:    school.Province,
				PostalCode:  school.PostalCode,
				Phone:       school.Phone,
				Email:       school.Email,
				Website:     school.Website,
				Logo:        school.Logo,
				Banner:      school.Banner,
				CreatedAt:   school.CreatedAt,
				UpdatedAt:   school.UpdatedAt,
				Status:      s.setStatus(school.DeletedAt),
			},
			StudentCount: schoolCounts[school.ID].StudentCount,
			TeacherCount: schoolCounts[school.ID].TeacherCount,
		},
		)
	}

	meta := commonHttp.NewMetaFromQuery(httpQuery, total)

	return res, meta, nil
}

func (s *service) GetListSchoolStatistics(ctx context.Context) (*response.ListSchoolStatistics, error) {
	claim, err := jwt.ExtractContext(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to extract claims")
		return nil, err
	}

	// Get overall statistics
	stats, err := s.repository.GetListSchoolStatistics(ctx, claim.User.ID)
	if err != nil {
		log.Err(err).Msg("error getting list school statistics")
		// Don't fail the request, just use empty stats
		stats = &repository.ListSchoolStatistics{}
	}

	return &response.ListSchoolStatistics{
		TotalSchools:  stats.TotalSchools,
		TotalStudents: stats.TotalStudents,
		TotalTeachers: stats.TotalTeachers,
		ActiveSchools: stats.ActiveSchools,
	}, nil

}

func (s *service) SwitchSchool(ctx context.Context, schoolID uuid.UUID) (string, string, error) {
	claim, err := jwt.ExtractContext(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to extract claims")
		return "", "", err
	}

	userSchoolRole, err := s.repository.GetSchoolRoleByUserIDAndSchoolID(ctx, claim.User.ID, schoolID)
	if err != nil {
		log.Err(err).Msg("error getting school role")
		return "", "", err
	}

	now := time.Now()
	exp := now.Add(time.Hour * 2).Unix()
	nbf := now.Unix()
	iat := now.Unix()
	payload := jwt.Payload{
		Exp: exp,
		Iat: iat,
		Nbf: nbf,
		Iss: "genesis",
		Sub: claim.User.ID.String(),
		Aud: "genesis",
		User: jwt.User{
			ID:         claim.User.ID,
			Email:      claim.User.Email,
			SchoolID:   schoolID,
			SchoolRole: userSchoolRole.RoleID,
			UserRole:   claim.User.UserRole,
		},
	}

	token, err := jwt.GenerateToken(payload, s.config.JWT.Secret)
	if err != nil {
		log.Err(err).Msg("error generating token")
		return "", "", err
	}

	refreshPayload := payload
	refreshPayload.Exp = now.Add(time.Hour * 240).Unix()
	refreshToken, err := jwt.GenerateToken(refreshPayload, s.config.JWT.Secret)
	if err != nil {
		log.Err(err).Msg("Error generating refresh token")
		return "", "", err
	}

	return token, refreshToken, nil
}

func (s *service) DeleteSchool(ctx context.Context, schoolID uuid.UUID) error {
	err := s.repository.DeleteSchool(ctx, schoolID)
	if err != nil {
		log.Err(err).Msg("Failed to delete school")
		return err
	}
	return nil
}

func (s *service) UpdateSchoolProfile(ctx context.Context, schoolID uuid.UUID, data request.UpdateSchoolProfileRequest) (response.DetailSchool, error) {
	school := repository.School{
		Name:        data.Name,
		Level:       data.Level,
		Description: data.Description,
		Address:     data.Address,
		City:        data.City,
		Province:    data.Province,
		PostalCode:  data.PostalCode,
		Phone:       data.Phone,
		Email:       data.Email,
		Website:     data.Website,
		Logo:        data.Logo,
	}

	updatedSchool, err := s.repository.UpdateSchoolProfile(ctx, schoolID, school)
	if err != nil {
		log.Err(err).Msg("error updating school profile")
		return response.DetailSchool{}, err
	}

	return response.DetailSchool{
		School: response.School{
			ID:          updatedSchool.ID,
			Name:        updatedSchool.Name,
			Level:       updatedSchool.Level,
			Description: updatedSchool.Description,
			Address:     updatedSchool.Address,
			City:        updatedSchool.City,
			Province:    updatedSchool.Province,
			PostalCode:  updatedSchool.PostalCode,
			Phone:       updatedSchool.Phone,
			Email:       updatedSchool.Email,
			Website:     updatedSchool.Website,
			Logo:        updatedSchool.Logo,
			Banner:      updatedSchool.Banner,
			CreatedAt:   updatedSchool.CreatedAt,
			CreatedBy:   updatedSchool.CreatedBy.String(),
			UpdatedAt:   updatedSchool.UpdatedAt,
			UpdatedBy:   updatedSchool.UpdatedBy.UUID.String(),
			Status:      s.setStatus(updatedSchool.DeletedAt),
		},
	}, nil
}

func (s *service) setStatus(deteledAt int64) string {
	if deteledAt == 0 {
		return "active"
	}
	return "inactive"
}
