package engine

import (
	"context"

	"github.com/google/uuid"
	"github.com/rishit1234567889/carZone/models"
	"github.com/rishit1234567889/carZone/store"
	"go.opentelemetry.io/otel"
)

type EngineService struct {
	store store.EngineStoreInterface
}

func NewEngineService(store store.EngineStoreInterface) *EngineService {
	return &EngineService{
		store: store,
	}
}

func (s *EngineService) GetEngineById(ctx context.Context, id string) (*models.Engine, error) {
	tracer := otel.Tracer("EngineService")

	ctx, span := tracer.Start(ctx, "GetEngineById-Service")

	defer span.End()

	engine, err := s.store.GetEngineById(ctx, id)

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &engine, nil
}

func (s *EngineService) CreateEngine(ctx context.Context, engineReq *models.EngineRequest) (*models.Engine, error) {
	tracer := otel.Tracer("EngineService")

	ctx, span := tracer.Start(ctx, "CreateEngine-Service")

	defer span.End()

	if err := models.ValidateEngineRequest(*engineReq); err != nil {
		span.RecordError(err)
		return nil, err
	}

	createdEngine, err := s.store.CreateEngine(ctx, engineReq)

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &createdEngine, nil
}

func (s *EngineService) UpdateEngine(ctx context.Context, id uuid.UUID, engineReq *models.EngineRequest) (*models.Engine, error) {
	tracer := otel.Tracer("EngineService")

	ctx, span := tracer.Start(ctx, "UpdateEngine-Service")

	defer span.End()
	if err := models.ValidateEngineRequest(*engineReq); err != nil {
		span.RecordError(err)
		return nil, err
	}

	updatedEngine, err := s.store.EngineUpdate(ctx, id, engineReq)

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &updatedEngine, nil
}

func (s *EngineService) DeleteEngine(ctx context.Context, id string) (*models.Engine, error) {
	tracer := otel.Tracer("EngineService")

	ctx, span := tracer.Start(ctx, "DeleteEngine-Service")

	defer span.End()

	deletedEngine, err := s.store.EngineDelete(ctx, id)

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &deletedEngine, nil
}
