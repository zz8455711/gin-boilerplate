package repository

import (
	"gin-boilerplate/infra/database"
	"gin-boilerplate/infra/logger"
)

func Save(model interface{}) interface{} {
	err := database.DB.Create(model).Error
	if err != nil {
		logger.Errorf("error, not save data %v", err)
	}
	return err
}

func Get(model interface{}) interface{} {
	err := database.DB.Find(model).Error
	return err
}

func GetOne(model interface{}) interface{} {
	err := database.DB.Last(model).Error
	return err
}

func Update(model interface{}) interface{} {
	err := database.DB.Find(model).Error
	return err
}

// GetOneByAddress 函数用于从数据库中获取具有特定地址的模型数据
// 参数 model 是一个空接口，可以传递任何符合该接口的结构体实例
// 参数 Address 是要匹配的地址字符串
func GetAllByAddress(model interface{}, Address string) error {
	// 使用数据库查询，通过 payerAddress 字段匹配给定的地址
	// 如果有错误发生，将错误返回
	if err := database.DB.Where("payer_address = ? AND status = true", Address).Find(model).Error; err != nil {
		return err
	}

	// 没有错误发生，返回 nil 表示成功
	return nil
}
