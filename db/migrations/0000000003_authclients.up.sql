CREATE TYPE granttype AS ENUM ('refresh_token', 'authorization_code');

CREATE TABLE authclients(
    client_id uuid UNIQUE PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    client_secret bytea NOT NULL,
    grant_types granttype[] NOT NULL,
    redirect_uris text[],
    meta jsonb
);