package service

import (
	"context"
	"enuma-elish/internal/teacher/repository"
	"enuma-elish/internal/teacher/service/data/request"
	"enuma-elish/internal/teacher/service/data/response"
	commonError "enuma-elish/pkg/error"
	commonHttp "enuma-elish/pkg/http"
	"enuma-elish/pkg/jwt"
	"fmt"
	"net/smtp"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) InviteTeacher(ctx context.Context, data request.InviteTeacherRequest) error {

	claim, err := jwt.ExtractContext(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to extract claims")
		return err
	}

	var teachers []repository.User
	var allTeacherSubjects []repository.TeacherSubject
	var allTeacherClasses []repository.TeacherClass
	now := time.Now().UnixMilli()

	// Prepare teachers and assignments
	for _, teacherRequest := range data.Teachers {
		teacherID := uuid.New()

		// Create teacher
		teachers = append(teachers, repository.User{
			Email:      teacherRequest.Email,
			ID:         teacherID,
			Name:       teacherRequest.Name,
			Password:   "",
			IsVerified: false,
			CreatedAt:  now,
			CreatedBy:  claim.User.ID,
			UpdatedAt:  0,
		})

		// Prepare subject assignments
		for _, subjectIDStr := range teacherRequest.SubjectIDs {
			subjectID, err := uuid.Parse(subjectIDStr)
			if err != nil {
				log.Err(err).Str("subject_id", subjectIDStr).Msg("invalid subject ID")
				continue
			}

			allTeacherSubjects = append(allTeacherSubjects, repository.TeacherSubject{
				ID:        uuid.New(),
				TeacherID: teacherID,
				SubjectID: subjectID,
				CreatedAt: now,
				CreatedBy: claim.User.ID,
				UpdatedAt: 0,
				IsDeleted: false,
			})
		}

		// Prepare class assignments
		for _, classIDStr := range teacherRequest.ClassIDs {
			classID, err := uuid.Parse(classIDStr)
			if err != nil {
				log.Err(err).Str("class_id", classIDStr).Msg("invalid class ID")
				continue
			}

			allTeacherClasses = append(allTeacherClasses, repository.TeacherClass{
				ID:        uuid.New(),
				TeacherID: teacherID,
				ClassID:   classID,
				CreatedAt: now,
				CreatedBy: claim.User.ID,
				UpdatedAt: 0,
				IsDeleted: false,
			})
		}
	}

	// Create teachers with assignments in a single transaction
	err = s.repository.CreateTeachersWithAssignments(ctx, teachers, data.SchoolID, allTeacherSubjects, allTeacherClasses)
	if err != nil {
		log.Error().Err(err).Msg("create teachers with assignments error")
		return err
	}

	// Send email invitations
	go func([]repository.User) {
		for _, v := range teachers {
			token, err := s.repository.CreateTeacherVerifyToken(context.Background(), v.Email)
			if err != nil {
				log.Err(err).Str("email", v.Email).Msg("create teacher verify token")
			}

			link := fmt.Sprintf("%s/email-verification?email=%s&token=%s&type=invite_teacher", s.config.Http.FrontendHost, v.Email, token)
			err = s.sendEmail(v.Email, link, "Email Verification")
			if err != nil {
				log.Err(err).Str("email", v.Email).Msg("send email")
			}
		}
	}(teachers)

	return nil
}

func (s *service) sendEmail(to string, msg, subject string) error {
	auth := smtp.PlainAuth("", s.config.SMTP.Username, s.config.SMTP.Password, s.config.SMTP.Host)
	message := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, msg))

	addr := fmt.Sprintf("%s:%d", s.config.SMTP.Host, s.config.SMTP.Port)
	err := smtp.SendMail(addr, auth, s.config.SMTP.Username, []string{to}, message)
	if err != nil {
		log.Err(err).Msg("failed to send email")
		return err
	}

	return nil
}

func (s *service) VerifyTeacherEmail(ctx context.Context, data request.VerifyTeacherEmailRequest) error {
	token, err := s.repository.VerifyEmailToken(ctx, data.Email)
	if err != nil {
		log.Err(err).Msg("Failed to verify email token")
		return commonError.ErrInvalidToken
	}

	if token != data.Token {
		log.Err(commonError.ErrInvalidToken).Msg("Failed to verify email token")
		return commonError.ErrInvalidToken
	}

	return nil
}

func (s *service) UpdateTeacherAfterVerifyEmail(ctx context.Context, data request.UpdateTeacherAfterVerifyEmailRequest) error {

	token, err := s.repository.VerifyEmailToken(ctx, data.Email)
	if err != nil {
		log.Err(err).Msg("Failed to verify email token")
		return commonError.ErrInvalidToken
	}

	if token != data.Token {
		log.Err(commonError.ErrInvalidToken).Msg("Failed to verify email token")
		return commonError.ErrInvalidToken
	}

	teacher, err := s.repository.GetTeacherByEmail(ctx, data.Email)
	if err != nil {
		log.Err(err).Msg("Failed to get teacher")
		return commonError.ErrUserNotFound
	}

	hashPass, err := hashPassword(data.Password)
	if err != nil {
		log.Err(err).Msg("Failed to hash password")
		return err
	}

	teacher = &repository.User{
		ID:         teacher.ID,
		Name:       data.Name,
		Email:      data.Email,
		Password:   hashPass,
		IsVerified: true,
		UpdatedAt:  time.Now().UnixMilli(),
	}

	err = s.repository.UpdateTeacher(ctx, *teacher)
	if err != nil {
		log.Err(err).Msg("Failed to update teacher")
		return commonError.ErrInternal
	}

	s.repository.Redis().Del(ctx, repository.TeacherVerifyEmailTokenKey+":"+data.Email)
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *service) ListTeachers(ctx context.Context, httpQuery request.GetListTeacherQuery) (response.GetListTeacherResponse, *commonHttp.Meta, error) {

	listTeacher, total, err := s.repository.GetListTeachers(ctx, httpQuery)
	if err != nil {
		log.Err(err).Msg("list teachers")
		return response.GetListTeacherResponse{}, nil, err
	}

	meta := commonHttp.NewMetaFromQuery(httpQuery, total)

	// Get teacher IDs for assignments lookup
	var teacherIDs []uuid.UUID
	for _, teacher := range listTeacher {
		teacherIDs = append(teacherIDs, teacher.ID)
	}

	// Get assignments if we have teachers
	var subjectMap map[uuid.UUID][]repository.Subject
	var classMap map[uuid.UUID][]repository.Class

	if len(teacherIDs) > 0 {
		schoolIDUUID, err := uuid.Parse(httpQuery.SchoolID)
		if err != nil {
			log.Err(err).Msg("parse school ID")
			return response.GetListTeacherResponse{}, nil, err
		}

		subjectMap, classMap, err = s.repository.GetTeacherAssignments(ctx, teacherIDs, schoolIDUUID)
		if err != nil {
			log.Err(err).Msg("get teacher assignments")
			// Don't fail the request, just log the error and continue without assignments
			subjectMap = make(map[uuid.UUID][]repository.Subject)
			classMap = make(map[uuid.UUID][]repository.Class)
		}
	}

	teachers := make([]response.GetTeacherResponse, len(listTeacher))
	for i, teacher := range listTeacher {
		// Convert repository subjects to response subjects
		var subjects []response.Subject
		for _, subject := range subjectMap[teacher.ID] {
			subjects = append(subjects, response.Subject{
				ID:        subject.ID,
				Name:      subject.Name,
				SchoolID:  subject.SchoolID,
				CreatedAt: subject.CreatedAt,
				UpdatedAt: subject.UpdatedAt,
			})
		}

		// Convert repository classes to response classes
		var classes []response.Class
		for _, class := range classMap[teacher.ID] {
			classes = append(classes, response.Class{
				ID:        class.ID,
				Name:      class.Name,
				SchoolID:  class.SchoolID,
				CreatedAt: class.CreatedAt,
				UpdatedAt: class.UpdatedAt,
			})
		}

		teachers[i] = response.GetTeacherResponse{
			ID:         teacher.ID,
			Name:       teacher.Name,
			Email:      teacher.Email,
			IsVerified: teacher.IsVerified,
			CreateAt:   teacher.CreatedAt,
			UpdateAt:   teacher.UpdatedAt,
			Subjects:   subjects,
			Classes:    classes,
		}
	}

	return teachers, meta, nil
}

func (s *service) GetTeacherStatistics(ctx context.Context) (response.TeacherStatistics, error) {
	// Get teacher statistics
	jwtClaim, err := jwt.ExtractContext(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to extract claims")
		return response.TeacherStatistics{}, commonError.ErrUnauthorized
	}
	totalTeachers, verifiedTeachers, pendingTeachers, activeTeachers, err := s.repository.GetTeacherStatistics(ctx, jwtClaim.User.SchoolID.String())
	if err != nil {
		log.Err(err).Msg("get teacher statistics")
		return response.TeacherStatistics{}, err
	}

	return response.TeacherStatistics{
		TotalTeachers:    totalTeachers,
		VerifiedTeachers: verifiedTeachers,
		PendingTeachers:  pendingTeachers,
		ActiveTeachers:   activeTeachers,
	}, nil
}

func (s *service) DeleteTeacher(ctx context.Context, teacherID uuid.UUID, schoolID uuid.UUID) error {
	err := s.repository.DeleteTeacher(ctx, teacherID, schoolID)
	if err != nil {
		log.Err(err).Msg("Failed to delete teacher")
		return err
	}
	return nil
}

func (s *service) GetDetailTeacher(ctx context.Context, teacherID uuid.UUID) (response.GetDetailTeacherResponse, error) {
	teacher, err := s.repository.GetTeacherByID(ctx, teacherID)
	if err != nil {
		log.Err(err).Msg("Failed to get teacher")
		return response.GetDetailTeacherResponse{}, err
	}

	res := response.GetDetailTeacherResponse{
		ID:         teacher.ID,
		Name:       teacher.Name,
		Email:      teacher.Email,
		IsVerified: teacher.IsVerified,
		CreateAt:   teacher.CreatedAt,
		UpdateAt:   teacher.UpdatedAt,
	}

	return res, nil
}

func (s *service) UpdateTeacherClass(ctx context.Context, data request.UpdateTeacherClassRequest) error {
	err := s.repository.UpdateTeacherClass(ctx, data.TeacherID, data.OldClassID, data.NewClassID)
	if err != nil {
		log.Err(err).Msg("Failed to update teacher class assignment")
		return err
	}
	return nil
}

func (s *service) GetTeacherSubjects(ctx context.Context, teacherID uuid.UUID) ([]response.Subject, error) {
	subjects, err := s.repository.GetTeacherSubjects(ctx, teacherID)
	if err != nil {
		return nil, err
	}
	res := make([]response.Subject, len(subjects))
	for i, subject := range subjects {
		res[i] = response.Subject{
			ID:        subject.ID,
			Name:      subject.Name,
			SchoolID:  subject.SchoolID,
			CreatedAt: subject.CreatedAt,
			UpdatedAt: subject.UpdatedAt,
		}
	}
	return res, nil
}

func (s *service) GetTeacherClasses(ctx context.Context, teacherID uuid.UUID) ([]response.Class, error) {
	classes, err := s.repository.GetTeacherClasses(ctx, teacherID)
	if err != nil {
		return nil, err
	}
	res := make([]response.Class, len(classes))
	for i, class := range classes {
		res[i] = response.Class{
			ID:        class.ID,
			Name:      class.Name,
			SchoolID:  class.SchoolID,
			CreatedAt: class.CreatedAt,
			UpdatedAt: class.UpdatedAt,
		}
	}
	return res, nil
}
