package main

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/libp2p/go-libp2p-core/peer"
	"math/rand"
	"net/http"
	"time"
)

func start() error {
	r := gin.Default()
	r.GET("/", root)
	r.POST("/verify", verify)
	return r.Run()
}

// VerifyData 数据验证结构体
type VerifyData struct {
	Url       string // 回传地址
	PeerID    string // 对应的PeerID，应为string类型
	Data      string // 原始数据，应为[]byte类型，用base64编码为string
	Signature string // 签名，应为[]byte类型，用base64编码为string
	PublicKey string // 对应的公钥，应为[]byte类型，用base64编码为string
}

func root(c *gin.Context) {
	var verifyData VerifyData
	verifyData.Data = base64.StdEncoding.EncodeToString([]byte(RandStringRunes(10)))
	verifyData.Url = "/verify"
	err := c.Bind(&verifyData)
	if err != nil {
		c.String(http.StatusInternalServerError, "internal server error")
	}
	c.JSON(http.StatusOK, gin.H{"data": verifyData})
}

func verify(c *gin.Context) {
	var verifyData VerifyData
	err := c.Bind(&verifyData)
	if err != nil {
		c.String(http.StatusInternalServerError, "internal server error")
	}
	verifyResult := verifyDataWithStruct(&verifyData)
	if verifyResult {
		c.String(http.StatusOK, "ok")
	} else {
		c.String(http.StatusBadRequest, "bad request")
	}
}

func verifyDataWithStruct(data *VerifyData) bool {
	dataBytes, err := base64.StdEncoding.DecodeString(data.Data)
	if err != nil {
		fmt.Println("base64 DecodeString err:", err)
		return false
	}
	signData, err := base64.StdEncoding.DecodeString(data.Signature)
	if err != nil {
		fmt.Println("base64 DecodeString err:", err)
		return false
	}
	publicKey, err := base64.StdEncoding.DecodeString(data.PublicKey)
	if err != nil {
		fmt.Println("base64 DecodeString err:", err)
		return false
	}
	peerID, err := peer.Decode(data.PeerID)
	if err != nil {
		fmt.Println("peer.IDFromString err:", err)
		return false
	}
	return verifyData(dataBytes, signData, peerID, publicKey)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
