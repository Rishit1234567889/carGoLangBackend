package car

import (
	"context"

	"github.com/google/uuid"
	"github.com/rishit1234567889/carZone/models"
	"github.com/rishit1234567889/carZone/store"
	"go.opentelemetry.io/otel"
)

type CarService struct {
	store store.CarStoreInterface
}

func NewCarService(store store.CarStoreInterface) *CarService {
	return &CarService{
		store: store,
	}
}

func (s *CarService) GetCarById(ctx context.Context, id string) (*models.Car, error) {
	tracer := otel.Tracer("CarService")

	ctx, span := tracer.Start(ctx, "GetCarById-Service")

	defer span.End()

	car, err := s.store.GetCarById(ctx, id)

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &car, nil
}

func (s *CarService) GetCarByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error) {
	tracer := otel.Tracer("CarService")

	ctx, span := tracer.Start(ctx, "GetCarByBrand-Service")

	defer span.End()

	cars, err := s.store.GetCarByBrand(ctx, brand, isEngine)

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return cars, nil
}

func (s *CarService) CreateCar(ctx context.Context, car *models.CarRequest) (*models.Car, error) {
	tracer := otel.Tracer("CarService")

	ctx, span := tracer.Start(ctx, "CreateCar-Service")

	defer span.End()

	if err := models.ValidateRequest(*car); err != nil {
		span.RecordError(err)
		return nil, err
	}
	createdCar, err := s.store.CreateCar(ctx, car)

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &createdCar, nil
}

func (s *CarService) UpdateCar(ctx context.Context, id uuid.UUID, carReq *models.CarRequest) (*models.Car, error) {
	tracer := otel.Tracer("CarService")

	ctx, span := tracer.Start(ctx, "UpdateCar-Service")

	defer span.End()

	if err := models.ValidateRequest(*carReq); err != nil {
		span.RecordError(err)
		return nil, err
	}

	updatedcar, err := s.store.UpdateCar(ctx, id, carReq)

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &updatedcar, nil
}

func (s *CarService) DeleteCar(ctx context.Context, id string) (*models.Car, error) {
	tracer := otel.Tracer("CarService")

	ctx, span := tracer.Start(ctx, "DeleteCar-Service")

	defer span.End()

	deletedCar, err := s.store.DeleteCar(ctx, id)

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &deletedCar, nil
}
