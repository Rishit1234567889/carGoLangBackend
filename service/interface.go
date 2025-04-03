package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rishit1234567889/carZone/models"
)

type CarServiceInterface interface {
	GetCarById(ctx context.Context, id string) (*models.Car, error)
	GetCarByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error)
	CreateCar(ctx context.Context, carReq *models.CarRequest) (*models.Car, error)
	UpdateCar(ctx context.Context, id uuid.UUID, carReq *models.CarRequest) (*models.Car, error)
	DeleteCar(ctx context.Context, id string) (*models.Car, error)
}

type EngineServiceInterface interface {
	GetEngineById(ctx context.Context, id string) (*models.Engine, error)
	CreateEngine(ctx context.Context, engineReq *models.EngineRequest) (*models.Engine, error)
	UpdateEngine(ctx context.Context, id uuid.UUID, engineReq *models.EngineRequest) (*models.Engine, error)
	DeleteEngine(ctx context.Context, id string) (*models.Engine, error)
}
