package main

import (
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerifySuccess(t *testing.T) {
	// 先new一个节点，用来签名
	a, err := NewNode()
	assert.NoError(t, err)

	// 现在用a签名一个数据
	data := []byte("hello")
	sig, err := a.signData(data)
	assert.NoError(t, err)

	// 把a的公钥转换成字节数组
	publicKey, err := crypto.MarshalPublicKey(a.Peerstore().PubKey(a.ID()))
	assert.NoError(t, err)

	// 把sig和data和a的公钥传给b，b验证sig和data是否合法
	verifyResult := verifyData(data, sig, a.ID(), publicKey)
	assert.True(t, verifyResult)
}

// 用错误的数据验证，结果应该是false
func TestVerifyFailedWithDataError(t *testing.T) {
	a, err := NewNode()
	assert.NoError(t, err)

	// 现在用a签名一个数据
	data := []byte("hello")
	sig, err := a.signData(data)
	assert.NoError(t, err)

	// 把a的公钥转换成字节数组
	publicKey, err := crypto.MarshalPublicKey(a.Peerstore().PubKey(a.ID()))
	assert.NoError(t, err)

	// 把sig和data和a的公钥传给b，b验证sig和data是否合法
	// 在这里把data改一下，验证结果应该是false
	verifyResult := verifyData([]byte("hi"), sig, a.ID(), publicKey)
	assert.False(t, verifyResult)
}

// 用错误的PeerID验证，结果应该是false
func TestVerifyFailedWithErrorPeerId(t *testing.T) {
	a, err := NewNode()
	assert.NoError(t, err)

	// 现在用a签名一个数据
	data := []byte("hello")
	sig, err := a.signData(data)
	assert.NoError(t, err)

	// 把a的公钥转换成字节数组
	publicKey, err := crypto.MarshalPublicKey(a.Peerstore().PubKey(a.ID()))
	assert.NoError(t, err)

	// 把sig和data和a的公钥传给b，b验证sig和data是否合法
	// 在这里把PeerId改一下，验证结果应该是false
	verifyResult := verifyData(data, sig, a.ID()+"1", publicKey)
	assert.False(t, verifyResult)
}
