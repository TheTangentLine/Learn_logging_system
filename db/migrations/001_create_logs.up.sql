CREATE TABLE IF NOT EXISTS logs (
    id           UUID         PRIMARY KEY,
    level        VARCHAR(5)   NOT NULL CHECK (level IN ('DEBUG','INFO','WARN','ERROR','FATAL')),
    message      TEXT         NOT NULL,
    service_name VARCHAR(255) NOT NULL,
    timestamp    TIMESTAMPTZ  NOT NULL,
    metadata     JSONB,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_logs_service_name ON logs (service_name);
CREATE INDEX IF NOT EXISTS idx_logs_level        ON logs (level);
CREATE INDEX IF NOT EXISTS idx_logs_timestamp    ON logs (timestamp DESC);
