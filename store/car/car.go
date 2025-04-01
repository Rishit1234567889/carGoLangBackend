package car

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rishit1234567889/carZone/models"
)

type Store struct { // a types of constructor
	db *sql.DB
}

func New(db *sql.DB) Store {
	return Store{db: db}
}

func (s Store) GetCarById(ctx context.Context, id string) (models.Car, error) {
	var car models.Car // single car data returned
	//query := `SELECT c.id, c.name, c.year, c.brand, c.fuel_type, c.engine_id, c.price, c.created_at, c.updated_at,
	//e.id, e.displacement, e.no_of_cylinders, e.car_range FROM car c LEFT JOIN engine e ON c.engine_id = e.id WHERE c.id = $1`

	query := `SELECT c.id, c.name, c.year, c.brand, c.fuel_type, c.engine_id, c.price, c.created_at,
       c.updated_at, e.displacement, e.no_of_cylinders, e.car_range FROM car c LEFT JOIN engine e ON c.engine_id = e.id WHERE c.id = $1;`

	row := s.db.QueryRowContext(ctx, query, id) // single row fetched
	err := row.Scan(
		&car.ID,
		&car.Name,
		&car.Year,
		&car.Brand,
		&car.FuelType,
		&car.Engine.EngineID,
		&car.Price,
		&car.CreatedAt,
		&car.UpdatedAt,
		&car.Engine.Displacement,
		&car.Engine.NoOfCylinders,
		&car.Engine.CarRange,
	)
	if err != nil {
		if err == sql.ErrNoRows { //
			return car, nil
		}
		return car, err
	}
	return car, nil
}

func (s Store) GetCarByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error) {
	var cars []models.Car // multiple car can be in single Brand
	var query string

	if isEngine { // if engine info also required
		query = `SELECT c.id, c.name, c.year, c.brand, c.fuel_type, c.engine_id, c.price, c.created_at, c.updated_at, 
        e.id, e.displacement, e.no_of_cylinders, e.car_eange FROM car c LEFT JOIN engine e ON c.engine_id = e.id WHERE c.brand = $1`
	} else {
		query = `SELECT id, name, year, brand, fuel_type, engine_id, price, created_at, updated_at FROM car WHERE brand = $1`
	}

	rows, err := s.db.QueryContext(ctx, query, brand) // multiple rows may fetched
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() { // another way is checking for loop in isEngine cndtion, true case and false case
		var car models.Car
		if isEngine {
			var engine models.Engine
			err := rows.Scan(
				&car.ID,
				&car.Name,
				&car.Year,
				&car.Brand,
				&car.FuelType,
				&car.Engine.EngineID,
				&car.Price,
				&car.CreatedAt,
				&car.UpdatedAt,
				&car.Engine.Displacement,
				&car.Engine.NoOfCylinders,
				&car.Engine.CarRange,
			)
			if err != nil {
				return nil, err
			}
			car.Engine = engine //
		} else {
			err := rows.Scan(
				&car.ID,
				&car.Name,
				&car.Year,
				&car.Brand,
				&car.FuelType,
				&car.Engine.EngineID,
				&car.Price,
				&car.CreatedAt,
				&car.UpdatedAt,
			)
			if err != nil {
				return nil, err
			}
		}
		cars = append(cars, car)
	}
	return cars, nil
}

func (s Store) CreateCar(ctx context.Context, carReq *models.CarRequest) (models.Car, error) {
	var createdCar models.Car
	var engineID uuid.UUID
	err := s.db.QueryRowContext(ctx, "SELECT id FROM engine WHERE id = $1", carReq.Engine.EngineID).Scan(&engineID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return createdCar, errors.New("engine_id does not exists in the engine table")
		}
		return createdCar, err
	}

	carID := uuid.New()
	createdAt := time.Now()
	updatedAt := createdAt

	newCar := models.Car{
		ID:        carID,
		Name:      carReq.Name,
		Year:      carReq.Year,
		Brand:     carReq.Brand,
		FuelType:  carReq.FuelType,
		Engine:    carReq.Engine,
		Price:     carReq.Price,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	// begin the transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return createdCar, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return

		}
		err = tx.Commit()
	}()

	query := `
	INSERT INTO car (id, name, year, brand, fuel_type, engine_id, price, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id, name, year, brand, fuel_type, engine_id, price, created_at, updated_at
	`
	err = tx.QueryRowContext(ctx, query,
		newCar.ID,
		newCar.Name,
		newCar.Year,
		newCar.Brand,
		newCar.FuelType,
		newCar.Engine.EngineID,
		newCar.Price,
		newCar.CreatedAt,
		newCar.UpdatedAt,
	).Scan(
		&createdCar.ID,
		&createdCar.Name,
		&createdCar.Year,
		&createdCar.Brand,
		&createdCar.FuelType,
		&createdCar.Engine.EngineID,
		&createdCar.Price,
		&createdCar.CreatedAt,
		&createdCar.UpdatedAt,
	)

	if err != nil {
		return createdCar, err
	}
	return createdCar, nil
}

func (s Store) UpdateCar(ctx context.Context, id string, carReq *models.CarRequest) (models.Car, error) {
	var updatedCar models.Car
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return updatedCar, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return

		}
		err = tx.Commit()
	}()
	query := `
	UPDATE car
	SET name = $2, year = $3, brand = $4, fuel_type = $5, engine_id = $6, price = $7, updated_at = $8
	WHERE id = $1
	RETURNING id, name, year, brand, fuel_type, engine_id, price, created_at, updated_at
	`

	err = tx.QueryRowContext(ctx, query,
		id,
		carReq.Name,
		carReq.Year,
		carReq.Brand,
		carReq.FuelType,
		carReq.Engine.EngineID,
		carReq.Price,
		time.Now(),
	).Scan(
		&updatedCar.ID,
		&updatedCar.Name,
		&updatedCar.Year,
		&updatedCar.Brand,
		&updatedCar.FuelType,
		&updatedCar.Engine.EngineID,
		&updatedCar.Price,
		&updatedCar.CreatedAt,
		&updatedCar.UpdatedAt,
	)

	if err != nil {
		return updatedCar, err
	}

	return updatedCar, nil

}

func (s Store) DeleteCar(ctx context.Context, id string) (models.Car, error) {
	var deletedCar models.Car // singel car should be returned after deleting

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return deletedCar, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// find rows, delete rows and return deleted rows

	err = tx.QueryRowContext(ctx, "SELECT id, name, year, brand, fuel_type, engine_id, price, created_at, updated_at FROM car WHERE id = $1", id).Scan(
		&deletedCar.ID,
		&deletedCar.Name,
		&deletedCar.Year,
		&deletedCar.Brand,
		&deletedCar.FuelType,
		&deletedCar.Engine.EngineID,
		&deletedCar.Price,
		&deletedCar.CreatedAt,
		&deletedCar.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // if no rows with given id found case
			return models.Car{}, errors.New("car not found")
		}
		return models.Car{}, err
	}

	results, err := tx.ExecContext(ctx, "DELETE FROM car WHERE id = $1", id)
	if err != nil {
		return models.Car{}, err
	}
	rowsAffected, err := results.RowsAffected()

	if rowsAffected == 0 {
		return models.Car{}, errors.New("no rows were deleted")
	}
	return deletedCar, nil
}
