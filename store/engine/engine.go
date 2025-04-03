package engine

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rishit1234567889/carZone/models"
	"go.opentelemetry.io/otel"
)

type EngineStore struct {
	db *sql.DB
}

func New(db *sql.DB) *EngineStore {
	return &EngineStore{db: db}
}

func (s EngineStore) GetEngineById(ctx context.Context, id string) (models.Engine, error) {
	tracer := otel.Tracer("EngineStore")

	ctx, span := tracer.Start(ctx, "GetEngineById-Store")

	defer span.End()

	var engine models.Engine

	// Begin Transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		span.RecordError(err)
		return models.Engine{}, err
	}

	defer func() {
		if err != nil {
			span.RecordError(err)
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("Transaction rollback error: %v\n", rbErr)
			}
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				span.RecordError(err)
				fmt.Printf("Transaction commit error: %v\v", cmErr)
			}
		}
	}()

	// Prepare the SQL query to select the engine by ID
	query := `
		SELECT 
			engine_id, displacement, no_of_cylinders, car_range 
		FROM 
			engines 
		WHERE 
			engine_id = $1`

	// Execute the query
	err = tx.QueryRowContext(ctx, query, id).Scan(
		&engine.EngineID,
		&engine.Displacement,
		&engine.NoOfCylinders,
		&engine.CarRange,
	)

	if err != nil {
		span.RecordError(err)
		if err == sql.ErrNoRows {
			return models.Engine{}, errors.New("engine not found")
		}
		return models.Engine{}, err
	}

	// Commit the transaction (optional for read operations)
	if err := tx.Commit(); err != nil {
		span.RecordError(err)
		return models.Engine{}, err
	}

	return engine, nil
}

func (s EngineStore) CreateEngine(ctx context.Context, engineReq *models.EngineRequest) (models.Engine, error) {
	tracer := otel.Tracer("EngineStore")

	ctx, span := tracer.Start(ctx, "CreateEngine-Store")

	defer span.End()
	var engine models.Engine

	// Begin Transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		span.RecordError(err)
		return engine, err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				span.RecordError(err)
				fmt.Printf("Transaction rollback error: %v\n", rbErr)
			}
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				span.RecordError(err)
				fmt.Printf("Transaction commit error: %v\v", cmErr)
			}
		}
	}()

	engineId := uuid.New()

	// Prepare the SQL query to insert a new engine
	query := `
		INSERT INTO engines (engine_id, displacement, no_of_cylinders, car_range) 
		VALUES ($1, $2, $3, $4)`

	// Execute the insert query
	_, err = tx.ExecContext(ctx, query, engineId, engineReq.Displacement, engineReq.NoOfCylinders, engineReq.CarRange)
	if err != nil {
		span.RecordError(err)
		return engine, err // Return error if the insertion fails
	}

	// Set the engine fields
	engine.Displacement = engineReq.Displacement
	engine.NoOfCylinders = engineReq.NoOfCylinders
	engine.CarRange = engineReq.CarRange

	return engine, nil // Return the created engine
}

func (s EngineStore) EngineUpdate(ctx context.Context, id uuid.UUID, engineReq *models.EngineRequest) (models.Engine, error) {
	tracer := otel.Tracer("EngineStore")

	ctx, span := tracer.Start(ctx, "EngineUpdate-Store")

	defer span.End()

	var updatedEngine models.Engine

	// Begin Transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		span.RecordError(err)
		return updatedEngine, err
	}

	defer func() {
		if err != nil {
			tx.Rollback() // Rollback on error
			span.RecordError(err)
			return
		}
		err = tx.Commit() // Commit if no error
	}()

	// Prepare the SQL query to update the engine
	query := `
		UPDATE engines 
		SET displacement = $1, no_of_cylinders = $2, car_range = $3 
		WHERE engine_id = $4`

	// Execute the update query
	_, err = tx.ExecContext(ctx, query, engineReq.Displacement, engineReq.NoOfCylinders, engineReq.CarRange, id)
	if err != nil {
		span.RecordError(err)
		return updatedEngine, err // Return error if the update fails
	}

	// Set the updated engine fields
	updatedEngine.EngineID = id // Assuming id is the engine's ID
	updatedEngine.Displacement = engineReq.Displacement
	updatedEngine.NoOfCylinders = engineReq.NoOfCylinders
	updatedEngine.CarRange = engineReq.CarRange

	return updatedEngine, nil
}

func (s EngineStore) EngineDelete(ctx context.Context, id string) (models.Engine, error) {
	tracer := otel.Tracer("EngineStore")

	ctx, span := tracer.Start(ctx, "EngineDelete-Store")

	defer span.End()

	var deletedEngine models.Engine

	// Begin Transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		span.RecordError(err)
		return deletedEngine, err
	}

	defer func() {
		if err != nil {
			tx.Rollback() // Rollback on error
			span.RecordError(err)
			return
		}
		err = tx.Commit() // Commit if no error
	}()

	// Prepare the SQL query to select the engine before deletion
	selectQuery := `SELECT engine_id, displacement, no_of_cylinders, car_range FROM engines WHERE engine_id = $1`
	err = tx.QueryRowContext(ctx, selectQuery, id).Scan(
		&deletedEngine.EngineID,
		&deletedEngine.Displacement,
		&deletedEngine.NoOfCylinders,
		&deletedEngine.CarRange,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return deletedEngine, errors.New("engine not found")
		}
		return deletedEngine, err
	}

	// Prepare the SQL query to delete the engine
	deleteQuery := `DELETE FROM engines WHERE engine_id = $1`
	_, err = tx.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		span.RecordError(err)
		return deletedEngine, err // Return error if the deletion fails
	}

	return deletedEngine, nil
}
