package tron

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

// SendTransaction 发起转账
func SendTransaction(toAddress string, amount int64) (string, error) {
	// 从配置文件读取私钥和发起交易的地址
	tornPrivate := viper.GetString("tornPrivate")
	fromAddress := viper.GetString("fromAddress")
	grpcLink := viper.GetString("grpcLink")
	tronGrid := viper.GetString("tronGrid")

	// 解码私钥字符串为字节切片
	privateKeyBytes, _ := hex.DecodeString(tornPrivate)

	// 创建 gRPC 客户端
	c := client.NewGrpcClient(grpcLink)

	// 设置 gRPC 连接
	creds := grpc.WithTransportCredentials(insecure.NewCredentials())
	c.SetAPIKey(tronGrid)
	if err := c.Start(creds); err != nil {
		return "", err
	}

	// 发起交易
	tx, err := c.Transfer(fromAddress, toAddress, amount)
	if err != nil {
		return "", err
	}

	// 计算交易的哈希
	rawData, err := proto.Marshal(tx.Transaction.GetRawData())
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(rawData)

	// 使用私钥对交易哈希进行签名
	sk, _ := btcec.PrivKeyFromBytes(privateKeyBytes)
	signature, err := crypto.Sign(hash[:], sk.ToECDSA())
	if err != nil {
		return "", err
	}
	tx.Transaction.Signature = append(tx.Transaction.Signature, signature)

	// 广播交易
	if _, err := c.Broadcast(tx.Transaction); err != nil {
		return "", err
	}

	// 获取交易ID，并返回
	txID := common.BytesToHexString(tx.GetTxid())
	return removeHexPrefix(txID), nil
}

func removeHexPrefix(hexString string) string {
	// 检查字符串是否以 "0x" 开头
	if len(hexString) >= 2 && hexString[:2] == "0x" {
		// 删除 "0x" 前缀
		return hexString[2:]
	}
	// 如果没有 "0x" 前缀，直接返回原字符串
	return hexString
}
