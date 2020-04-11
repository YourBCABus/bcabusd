CREATE TABLE accesstokens(
    token_hash bytea UNIQUE PRIMARY KEY NOT NULL,
    client_id uuid NOT NULL REFERENCES authclients(client_id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    scopes text[] NOT NULL,
    expires timestamp NOT NULL
);

CREATE TABLE refreshtokens(
    token_hash bytea UNIQUE PRIMARY KEY NOT NULL,
    client_id uuid NOT NULL REFERENCES authclients(client_id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    scopes text[] NOT NULL,
    expires timestamp NOT NULL
);