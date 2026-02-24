CREATE TABLE IF NOT EXISTS urls (
    id BIGSERIAL PRIMARY KEY,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS click_events (
    id BIGSERIAL PRIMARY KEY,
    short_code VARCHAR(10) NOT NULL,
    clicked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    referer TEXT,
    user_agent TEXT
);

CREATE INDEX idx_click_events_short_code ON click_events (short_code);
CREATE INDEX idx_click_events_clicked_at ON click_events (clicked_at);
