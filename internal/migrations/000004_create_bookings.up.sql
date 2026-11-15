CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE IF NOT EXISTS bookings (
    id           BIGSERIAL   PRIMARY KEY,
    table_id     BIGINT      NOT NULL REFERENCES tables (id),
    guest_id     BIGINT      NOT NULL REFERENCES guests (id),
    starts_at    TIMESTAMPTZ NOT NULL,
    ends_at      TIMESTAMPTZ NOT NULL,
    guests_count INT         NOT NULL CHECK (guests_count > 0),
    status       TEXT        NOT NULL DEFAULT 'active'
                 CHECK (status IN ('active', 'cancelled', 'completed')),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT bookings_duration_2h
        CHECK (ends_at - starts_at = interval '2 hours'),
    CONSTRAINT bookings_no_overlap EXCLUDE USING gist (
        table_id WITH =,
        tstzrange(starts_at, ends_at, '[)') WITH &&
    ) WHERE (status = 'active')
);
