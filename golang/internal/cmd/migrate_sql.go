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
}
