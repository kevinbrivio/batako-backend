CREATE TABLE IF NOT EXISTS cement_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS cement_stocks(
    id VARCHAR(36) PRIMARY KEY,
    cement_type_id INT NOT NULL REFERENCES cement_types(id),
    quantity DOUBLE PRECISION NOT NULL CHECK (quantity >= 0),
    price_per_bag DOUBLE PRECISION NOT NULL,
    purchase_date DATE NOT NULL DEFAULT CURRENT_DATE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS production_cement_usage (
    id SERIAL PRIMARY KEY,
    production_id VARCHAR(36) NOT NULL REFERENCES productions(id) ON DELETE CASCADE,
    cement_type_id INT NOT NULL REFERENCES cement_types(id),
    cement_used DOUBLE PRECISION NOT NULL CHECK (cement_used >= 0)
);