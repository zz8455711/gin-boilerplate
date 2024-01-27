package models

import (
	"time"
)

// User 结构体定义了模型的字段以及它们在 JSON 中的表示
type Transaction struct {
	ID              int64      `json:"id" primaryKey:"autoIncrement"` // 标识符，自增的主键
	PayerAddress    string     `json:"payerAddress"`                  // 付款方地址
	PayeeAddress    string     `json:"payeeAddress"`                  // 收款方地址
	UsdtAmount      float64    `json:"UsdtAmount"`                    // 金额
	TrxOverdraft    int64      `json:"TrxOverdraft"`                  // 透支额
	TrxAmount       float64    `json:"TrxAmount"`                     // 交易金额
	TrxPrice        float64    `json:"TrxPrice"`                      // 交易价格
	TransactionTime *time.Time `json:"transactionTime"`               // 交易时间
	Status          bool       `json:"status"`                        // 状态
}

// TableName 方法返回模型对应的数据库表的名称
func (e *Transaction) TableName() string {
	return "tansaction"
}
