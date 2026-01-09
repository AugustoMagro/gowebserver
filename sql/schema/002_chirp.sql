-- +goose Up
CREATE TABLE chirpy (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  body TEXT NOT NULL,
  FOREIGN KEY(user_id) 
  REFERENCES users(id)
);

-- +goose Down
DROP TABLE users;
