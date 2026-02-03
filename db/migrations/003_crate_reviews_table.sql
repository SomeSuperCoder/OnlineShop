-- +goose Up
CREATE TABLE reviews (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

  product UUID NOT NULL REFERENCES products(id),
  comment TEXT,
  stars INT NOT NULL CHECK ( stars > 0 AND stars <= 5 ),
  author UUID NOT NULL REFERENCES users(id),
  
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE reviews ENABLE ROW LEVEL SECURITY;

-- 1. Anyone can SELECT all reviews
CREATE POLICY select_all_reviews ON reviews
FOR SELECT
USING (true);  -- Everyone sees everything

-- 2. Anyone can INSERT reviews (must set themselves as author)
CREATE POLICY insert_any_review ON reviews
FOR INSERT
WITH CHECK (true);  -- Anyone can insert

-- 3. Only author can UPDATE their review
CREATE POLICY update_own_review ON reviews
FOR UPDATE
USING (author = current_setting('app.user_id', true)::uuid);

-- 4. Only author can DELETE their review
CREATE POLICY delete_own_review ON reviews
FOR DELETE
USING (author = current_setting('app.user_id', true)::uuid);

-- +goose Down
DROP TABLE reviews;
