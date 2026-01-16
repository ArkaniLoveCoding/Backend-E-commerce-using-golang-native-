CREATE TABLE public.product_clients (
    id          INT PRIMARY KEY,
    name        TEXT NOT NULL, 
    stock       INT NOT NULL,
    price       TEXT NOT NULL,
    expired     TEXT NOT NULL,
    category    TEXT NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);