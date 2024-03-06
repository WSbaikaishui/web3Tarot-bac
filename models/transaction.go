package models

const transactionTable = "transaction"

type Transaction struct {
	Model
	Address string `gorm:"column:address;type:VARCHAR(255);NOT NULL;DEFAULT:'';"`
	TxHash  string `gorm:"column:tx_hash;type:VARCHAR(255);NOT NULL;DEFAULT:''; unique"`
	Amount  int    `gorm:"column:number;type:int;NOT NULL;DEFAULT:0"`
}

func CreateTransaction(data interface{}, user_id int, count int) error {
	tx := Db.Begin()

	ddb := tx.Table(transactionTable).Create(data)
	if ddb.Error != nil {
		tx.Rollback()
	}
	if err := tx.Table(userTable).Where("user_id=?", user_id).Updates(map[string]interface{}{
		"token": count,
	}).Error; err != nil {
		// 回滚事务
		tx.Rollback()

	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func GetTransactionByAddress(address string) (transactions []*Transaction) {
	Db.Table(transactionTable).Where("address=?", address).Find(&transactions).Limit(10).Order("desc")
	return transactions
}

func GetTransactionByTxHash(txHash string) (transaction *Transaction) {
	Db.Table(transactionTable).Where("tx_hash=?", txHash).First(&transaction)
	return transaction
}

func GetTransactions() (transactions []*Transaction) {
	Db.Table(transactionTable).Find(&transactions)
	return transactions
}

func IsTransactionExist(txHash string) bool {
	var transaction Transaction
	Db.Table(transactionTable).Select("id").Where("tx_hash = ?", txHash).First(&transaction)
	if transaction.ID > 0 {
		return true
	}
	return false
}
