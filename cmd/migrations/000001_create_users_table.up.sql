CREATE TABLE public.users (
    id              UUID PRIMARY KEY DEFAULT
                    gen_random_uuid(),
    firstname       TEXT NOT NULL,
    lastname        TEXT NOT NULL,
    password        TEXT NOT NULL,
    email           TEXT NOT NULL,
    country         TEXT NOT NULL, 
    address         TEXT NOT NULL,
    role            TEXT NOT NULL,
    token           TEXT,
    refresh_token   TEXT,
    created_at      TIMESTAMP NOT NULL,
    updated_at      TIMESTAMP NOT NULL
);