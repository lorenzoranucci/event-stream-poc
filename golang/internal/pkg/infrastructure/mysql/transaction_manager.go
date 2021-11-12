package mysql

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/application/commands"
)

type TransactionManager struct {
	db *sqlx.DB
}

func NewTransactionManager(db *sqlx.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

func (t *TransactionManager) Create() (commands.Transaction, error) {
	return &Transaction{db: t.db}, nil
}

type Transaction struct {
	db *sqlx.DB
	tx *sql.Tx
}

func (t *Transaction) Tx() *sql.Tx {
	return t.tx
}

func (t *Transaction) Begin() error {
	var err error
	t.tx, err = t.db.Begin()

	return err
}

func (t *Transaction) Commit() error {
	if t.tx == nil {
		return fmt.Errorf("transaction is nil. Begin transaction before commit")
	}

	return t.tx.Commit()
}

func (t *Transaction) Rollback() {
	if t.tx == nil {
		return
	}

	_ = t.tx.Rollback()
}

func getTransaction(
	transaction commands.Transaction,
	db *sqlx.DB,
) (*sql.Tx, bool, error) {
	if shouldCreateNewTransaction(transaction) {
		tx, err := db.Begin()
		return tx, true, err
	}

	t := transaction.(*Transaction)

	return t.Tx(), false, nil
}

func shouldCreateNewTransaction(transaction commands.Transaction) bool {
	if transaction == nil {
		return true
	}

	_, ok := transaction.(*Transaction)

	return !ok
}
