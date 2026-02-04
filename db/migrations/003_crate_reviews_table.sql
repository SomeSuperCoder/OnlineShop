-- +goose Up
CREATE TABLE reviews (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

  product UUID NOT NULL REFERENCES products(id),
  comment TEXT,
  stars INT NOT NULL CHECK ( stars > 0 AND stars <= 5 ),
  author UUID NOT NULL REFERENCES users(id),

  upvotes INT NOT NULL DEFAULT 0 CHECK ( upvotes >= 0 ),
  downvotes INT NOT NULL DEFAULT 0 CHECK (downvotes >= 0),
  
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE reviews;
