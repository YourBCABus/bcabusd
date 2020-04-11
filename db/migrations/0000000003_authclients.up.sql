CREATE TABLE authclients(
    client_id uuid UNIQUE PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    client_secret bytea NOT NULL,
    meta jsonb
);