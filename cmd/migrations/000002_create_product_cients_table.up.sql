CREATE TABLE public.product_clients (
    id          UUID PRIMARY KEY DEFAULT
                gen_random_uuid(),
    name        TEXT NOT NULL, 
    stock       INT NOT NULL,
    image       TEXT NOT NULL,
    price       TEXT NOT NULL,
    expired     TEXT NOT NULL,
    category    TEXT NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);