package controllers

import (
	"fmt"
	"net/http"

	"gin-boilerplate/helpers"
	"gin-boilerplate/models"
	"gin-boilerplate/repository"
	"gin-boilerplate/tron"

	"github.com/gin-gonic/gin"
)

// GetData 获取所有用户数据
func GetData(ctx *gin.Context) {
	// 查询所有用户数据
	users := []*models.Transaction{}
	err := repository.Get(&users)
	if err != nil {
		// 处理获取用户数据时的错误
		ctx.JSON(http.StatusInternalServerError, helpers.Response{
			Code:    http.StatusInternalServerError,
			Message: "获取用户数据失败",
			Data:    nil,
		})
		return
	}

	// 返回成功的响应
	ctx.JSON(http.StatusOK, helpers.Response{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    users,
	})
}

// Create 创建新用户
func Create(ctx *gin.Context) {
	// 从请求体中绑定 JSON 数据到用户对象
	user := &models.Transaction{}
	if err := ctx.ShouldBindJSON(user); err != nil {
		fmt.Println("Error binding JSON data:", err)
		fmt.Println("Request JSON data:", ctx.Request.Body) // 输出请求体中的JSON数据
		// 如果绑定 JSON 数据出现错误，返回相应的错误信息给客户端
		ctx.JSON(http.StatusBadRequest, helpers.Response{
			Code:    http.StatusBadRequest,
			Message: "请求数据绑定失败",
			Data:    nil,
		})
		return
	}

	users := []*models.Transaction{}
	if err := repository.GetAllByAddress(&users, user.PayerAddress); err != nil {
		// 处理获取用户数据时的错误
		fmt.Println("Error getting user data:", err)
	}
	trxOverdraftSum := int64(0)
	for _, Overdraft := range users {
		trxOverdraftSum += Overdraft.TrxOverdraft
	}
	// fmt.Println("trxOverdraftSum:", trxOverdraftSum)
	// fmt.Println("user.TrxAmount:", user.TrxAmount)
	// fmt.Println("user,usdtAmount:", user.UsdtAmount)
	// fmt.Println("user.TrxOverdraft:", user.TrxOverdraft)
	var trxTransaction string
	var err error
	// 发起 TRX 转账并获取交易ID
	if trxOverdraftSum == 0 {
		trxTransaction, err = tron.SendTransaction(user.PayerAddress, int64(user.TrxAmount*1_000_000))
	} else {
		trxTransaction, err = tron.SendTransaction(user.PayerAddress, int64(user.TrxAmount*1_000_000)-trxOverdraftSum*1_000_000)
		user.TrxOverdraft = -int64(trxOverdraftSum)
	}

	if err != nil {
		// 处理交易失败时的错误
		ctx.JSON(http.StatusInternalServerError, helpers.Response{
			Code:    http.StatusInternalServerError,
			Message: "交易失败",
			Data:    fmt.Sprintf("%v", err),
		})

		// 保存用户到数据库（即使交易失败也保存用户数据）
		if saveErr := repository.Save(user); saveErr != nil {
			// 处理保存用户数据时的错误
			ctx.JSON(http.StatusInternalServerError, helpers.Response{
				Code:    http.StatusInternalServerError,
				Message: "保存用户数据失败",
				Data:    nil,
			})
		}

		return
	}

	// 保存用户到数据库
	user.Status = true
	if err := repository.Save(user); err != nil {
		// 处理保存用户数据时的错误
		ctx.JSON(http.StatusInternalServerError, helpers.Response{
			Code:    http.StatusInternalServerError,
			Message: "保存用户数据失败",
			Data:    nil,
		})
		return
	}

	// 返回成功的响应，包括交易ID
	ctx.JSON(http.StatusOK, helpers.Response{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    gin.H{"trxTransaction": trxTransaction, "trxOverdraftSum": trxOverdraftSum},
	})
}

// GetOneData 根据用户 ID 获取用户数据
func GetOneData(ctx *gin.Context) {
	// 从路径参数中获取用户的 ID
	userAddress := ctx.Param("address")

	// 使用 GetOneByID 函数获取特定用户的数据
	users := []*models.Transaction{}
	if err := repository.GetAllByAddress(&users, userAddress); err != nil {
		// 处理获取用户数据时的错误
		ctx.JSON(http.StatusInternalServerError, helpers.Response{
			Code:    http.StatusInternalServerError,
			Message: "获取用户数据失败",
			Data:    nil,
		})
		return
	}

	// 如果 users 是空切片，返回空数据的响应
	if len(users) == 0 {
		ctx.JSON(http.StatusNotFound, helpers.Response{
			Code:    http.StatusNotFound,
			Message: "Address does not exist",
			Data:    nil,
		})
		return
	}

	// 计算合计值
	var usdtSum float64
	var trxOverdraftSum int64
	for _, user := range users {
		usdtSum += user.UsdtAmount
		trxOverdraftSum += user.TrxOverdraft
	}

	// 返回成功的响应
	ctx.JSON(http.StatusOK, helpers.Response{
		Code:    http.StatusOK,
		Message: "Success",
		Data: map[string]interface{}{
			"address":         userAddress,
			"usdtAmountSum":   usdtSum,
			"trxOverdraftSum": trxOverdraftSum,
		},
	})
}
