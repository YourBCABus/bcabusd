CREATE TABLE authorizedclients(
    client_id uuid NOT NULL REFERENCES authclients(client_id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    scopes text[] NOT NULL,
    CONSTRAINT authorizedclients_pkey PRIMARY KEY (client_id, user_id)
);