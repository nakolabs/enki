package service

import (
	"context"
	"enuma-elish/internal/ppdb/repository"
	"enuma-elish/internal/ppdb/service/data/request"
	"github.com/rs/zerolog/log"
)

func (s *service) CreatSchool(ctx context.Context, data request.CreateSchoolRequest) error {

	school := repository.School{
		Name:  data.Name,
		Level: data.Level,
	}

	err := s.repository.CreateSchool(ctx, school)
	if err != nil {
		log.Err(err).Msg("error creating school")
		return err
	}

	return nil
}
