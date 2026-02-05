-- +goose Up
CREATE TYPE vote_type AS ENUM ('upvote', 'downvote');

CREATE TABLE votes (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

  review UUID NOT NULL REFERENCES reviews(id),
  voter UUID NOT NULL REFERENCES users(id),
  UNIQUE (review, voter),
  
  type vote_type NOT NULL,
  
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE votes;
DROP TYPE vote_type;
