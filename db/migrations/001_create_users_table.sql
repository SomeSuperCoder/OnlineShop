-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  email VARCHAR UNIQUE NOT NULL,
  username VARCHAR UNIQUE NOT NULL,

  name VARCHAR NOT NULL,
  balance INT NOT NULL DEFAULT 0 CHECK ( balance >= 0 ),
  password_hash TEXT NOT NULL,
  
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE users;
DROP EXTENSION "uuid-ossp";
DROP EXTENSION pgcrypto;
