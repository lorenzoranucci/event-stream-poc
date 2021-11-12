package commands

type TransactionManager interface {
	Create() (Transaction, error)
}

type Transaction interface {
	Begin() error
	Commit() error
	// Rollback must be idempotent
	Rollback()
}

func openTransaction(transactionManager TransactionManager) (Transaction, error) {
	transaction, err := transactionManager.Create()
	if err != nil {
		return nil, err
	}

	err = transaction.Begin()
	if err != nil {
		return nil, err
	}
	return transaction, nil
}
