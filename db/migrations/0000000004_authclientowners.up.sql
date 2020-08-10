CREATE TABLE authclientowners(
    client_id uuid NOT NULL REFERENCES authclients(client_id) ON DELETE CASCADE,
    owner_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT authclientowners_pkey PRIMARY KEY (client_id, owner_id)
);