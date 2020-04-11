CREATE TABLE botowners(
    bot_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    owner_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT botowners_pkey PRIMARY KEY (bot_id, owner_id)
);