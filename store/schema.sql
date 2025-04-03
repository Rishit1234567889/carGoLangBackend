
CREATE TABLE IF NOT EXISTS engines (
    engine_id UUID PRIMARY KEY,      
    displacement INT NOT NULL,       
    no_of_cylinders INT NOT NULL,    
    car_range INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cars (
    id UUID PRIMARY KEY,             
    name VARCHAR(100) NOT NULL,   
    year VARCHAR(4) NOT NULL,      
    brand VARCHAR(50) NOT NULL,     
    fuel_type VARCHAR(20) NOT NULL, 
    engine_id UUID NOT NULL,         
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP   
);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

INSERT INTO engines (engine_id, displacement, no_of_cylinders, car_range) VALUES
    (uuid_generate_v4(), 2000, 4, 500),
    (uuid_generate_v4(), 1500, 4, 450),
    (uuid_generate_v4(), 3000, 6, 600);


INSERT INTO cars (id, name, year, brand, fuel_type, engine_id, price) VALUES
    (uuid_generate_v4(), 'Toyota Camry', '2020', 'Toyota', 'Petrol', (SELECT engine_id FROM engines LIMIT 1), 24000.00),
    (uuid_generate_v4(), 'Honda Accord', '2019', 'Honda', 'Petrol', (SELECT engine_id FROM engines LIMIT 1 OFFSET 1), 22000.00),
    (uuid_generate_v4(), 'Ford Mustang', '2021', 'Ford', 'Petrol', (SELECT engine_id FROM engines LIMIT 1 OFFSET 2), 30000.00);