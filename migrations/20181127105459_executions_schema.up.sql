CREATE TABLE IF NOT EXISTS executions (
  id                      BIGSERIAL PRIMARY KEY,
  timestamp               BIGINT NOT NULL,
  status                  TEXT NOT NULL,
  image                   TEXT NOT NULL,
  tag                     TEXT NOT NULL,
  jobId                   UUID NOT NULL
);

CREATE INDEX executions_index ON executions (jobId ASC, timestamp DESC);
