CREATE TABLE users(
    id uuid UNIQUE PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    is_bot boolean NOT NULL,
    bot_token_hash bytea,
    user_token bytea,
    is_superadmin boolean NOT NULL DEFAULT false,
    meta jsonb
);