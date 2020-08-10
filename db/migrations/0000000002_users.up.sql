CREATE TABLE users(
    id uuid UNIQUE PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    user_token bytea,
    is_superadmin boolean NOT NULL DEFAULT false,
    meta jsonb
);