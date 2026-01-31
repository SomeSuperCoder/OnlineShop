-- +goose Up
CREATE TABLE reviews (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

  product UUID NOT NULL REFERENCES products(id),
  comment TEXT,
  stars INT NOT NULL CHECK ( stars > 0 AND stars <= 5 ),
  
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE reviews;
