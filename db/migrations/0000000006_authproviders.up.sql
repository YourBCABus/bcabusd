CREATE TABLE authproviders(
    id uuid UNIQUE PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider varchar(255) NOT NULL,
    email text,
    subject text,
    meta jsonb
);