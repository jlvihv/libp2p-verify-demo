package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/libp2p/go-libp2p-core/crypto"
	"io"
	"net/http"
)

type Body struct {
	Data VerifyData `json:"data"`
}

func client(port int) error {
	data, err := getData(port)
	if err != nil {
		return err
	}
	err = sign(&data)
	if err != nil {
		return err
	}
	return postData(port, &data)
}

func getData(port int) (VerifyData, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/", port))
	if err != nil {
		return VerifyData{}, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return VerifyData{}, err
	}
	resp.Body.Close()
	var b Body
	err = json.Unmarshal(body, &b)
	if err != nil {
		fmt.Println("Unmarshal err:", err)
		return VerifyData{}, err
	}
	return b.Data, nil
}

// 签名，就地修改data
func sign(data *VerifyData) error {
	node, err := NewNode()
	if err != nil {
		return err
	}
	dataBytes, err := base64.StdEncoding.DecodeString(data.Data)
	if err != nil {
		fmt.Println("base64 DecodeString err:", err)
		return err
	}
	signData, err := node.signData(dataBytes)
	if err != nil {
		fmt.Println("signData err:", err)
		return err
	}
	data.Signature = base64.StdEncoding.EncodeToString(signData)
	data.PeerID = node.ID().String()
	key, err := crypto.MarshalPublicKey(node.Peerstore().PubKey(node.ID()))
	if err != nil {
		fmt.Println("MarshalPublicKey err:", err)
		return err
	}
	data.PublicKey = base64.StdEncoding.EncodeToString(key)
	return nil
}

func postData(port int, data *VerifyData) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("json Marshal err:", err)
		return err
	}
	resp, err := http.Post(fmt.Sprintf("http://localhost:%d%s", port, data.Url), "application/json", bytes.NewReader(dataBytes))
	if err != nil {
		fmt.Println("http Post err:", err)
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("http Post err: %d", resp.StatusCode)
	}
	resp.Body.Close()
	return nil
}
