CREATE TABLE post (
  id uuid PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT,
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW()
);
