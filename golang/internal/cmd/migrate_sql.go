package cmd

import "github.com/jmoiron/sqlx"

type Migrate struct {
	db *sqlx.DB
}

func (m *Migrate) Migrate() {
	m.db.MustExec(`CREATE TABLE IF NOT EXISTS reviews_write (
    id BIGINT NOT NULL AUTO_INCREMENT,
    uuid VARCHAR(256) NOT NULL,
    comment TEXT DEFAULT '',
    rating INT DEFAULT 0,
    PRIMARY KEY (id),
    UNIQUE (uuid)
);`,
	)

	m.db.MustExec(`CREATE TABLE IF NOT EXISTS reviews_read (
    uuid VARCHAR(256) NOT NULL,
    comment TEXT DEFAULT '',
    rating INT DEFAULT 0,
    PRIMARY KEY (uuid)
);`,
	)

	m.db.MustExec(`CREATE TABLE IF NOT EXISTS review_events_outbox (
    id BIGINT NOT NULL AUTO_INCREMENT,
    uuid VARCHAR(256) NOT NULL,
    aggregate_id VARCHAR(256) NOT NULL,
    name VARCHAR(256) NOT NULL,
    payload TEXT NOT NULL,
    version VARCHAR(256) NOT NULL,
    status INT NOT NULL,
    message_counter_by_aggregate INT NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (uuid),
    UNIQUE (aggregate_id, message_counter_by_aggregate)
);`,
	)
}
