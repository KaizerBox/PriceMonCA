-- migrations/0001_init.sql
CREATE TABLE users (
  id uuid PRIMARY KEY,
  email text UNIQUE NOT NULL,
  name text,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE subscriptions (
  id uuid PRIMARY KEY,
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  site text NOT NULL,
  sku text NOT NULL,
  hints jsonb NOT NULL DEFAULT '{}',
  rule jsonb NOT NULL,
  targets jsonb NOT NULL,
  cron_spec text,
  active boolean NOT NULL DEFAULT true,
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX ON subscriptions (active) WHERE active;

CREATE TABLE price_history (
  id bigserial PRIMARY KEY,
  site text NOT NULL,
  sku text NOT NULL,
  regular_price numeric NOT NULL,
  sale_price numeric NOT NULL,
  currency text NOT NULL,
  online_available boolean NOT NULL,
  online_qty int,
  fetched_at timestamptz NOT NULL,
  payload jsonb NOT NULL
);
CREATE INDEX ON price_history (site, sku, fetched_at DESC);
