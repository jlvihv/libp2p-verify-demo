package main

import (
	"crypto/rand"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
)

func main() {}

type Node struct {
	host.Host
}

func NewNode() (*Node, error) {
	privKey, _, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return nil, err
	}
	h, err := libp2p.New(libp2p.Identity(privKey))
	return &Node{
		Host: h,
	}, nil
}

// 使用私钥签名数据，返回签名后的数据
func (n *Node) signData(data []byte) ([]byte, error) {
	privKey := n.Peerstore().PrivKey(n.ID())
	return privKey.Sign(data)
}

// 验证传入的数据，验证成功返回true，否则返回false
// data 是要验证的原始数据，不签名的
// sig 是签名后的数据
// peerId 是消息发送者的peerId
// pubKey 是消息发送者的公钥
func verifyData(data []byte, sig []byte, peerId peer.ID, pubKey []byte) bool {
	key, err := crypto.UnmarshalPublicKey(pubKey)
	if err != nil {
		fmt.Println("UnmarshalPublicKey err:", err)
		return false
	}
	idFromKey, err := peer.IDFromPublicKey(key)
	if err != nil {
		fmt.Println("IDFromPublicKey err:", err)
		return false
	}
	if idFromKey != peerId {
		fmt.Printf("IDFromPublicKey != peerId,\nIDFromPublicKey: %s,\npeerId: %s\n", idFromKey, peerId)
		return false
	}
	verifyResult, err := key.Verify(data, sig)
	if err != nil {
		fmt.Println("Verify err:", err)
		return false
	}
	if !verifyResult {
		fmt.Println("Verify result:", verifyResult)
	}
	return verifyResult
}

// TODO: 完成加密解密功能

// 使用公钥加密数据
func encryptData(data []byte, pubKey []byte) {}

// 使用私钥解密数据
func decryptData() {}
