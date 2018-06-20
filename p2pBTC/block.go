package main

import (
	"time"
	"encoding/gob"
	"bytes"
)

type Block struct {
	Version int64
	PrevBlockHash []byte //前区块的hash值

	Hash []byte //为了方便实现 做的简化 ，正常比特币区块 不包含自己的hash值
	TimeStamp int64
	TargetBits int64 //难度值
	Nonce int64     //随机值
	MerkelRoot []byte //


	Transactions []*Transaction  //区块体  （交易）
}

func (block *Block)Serialize()[]byte {
	var buffef bytes.Buffer
	encoder := gob.NewEncoder(&buffef)
	err := encoder.Encode(block)
	CheckErr(err)
	return buffef.Bytes()
}

func Deserialize(data []byte) * Block  {
	decoder := gob.NewDecoder(bytes.NewReader(data))

	var block Block
	err := decoder.Decode(&block)
	CheckErr(err)
	return &block
}

//func NewBlock(data string,prevBlickHash []byte) * Block {
func NewBlock(transactions []*Transaction,prevBlickHash []byte) * Block {
	block := &Block{

		Version:1,
		PrevBlockHash:  prevBlickHash,
		TimeStamp: 		time.Now().Unix(),
		TargetBits:		targetBits,
		Nonce:			0,
		MerkelRoot:		[]byte{},
		Transactions:	transactions,
	}

	pow := NewProofOfWork(block)
	nonce,hash := pow.Run()
	block.Nonce = nonce
	block.Hash = hash

	return block
}

func NewGenesisBlock(coinbase *Transaction) *Block  {

	return NewBlock([]*Transaction{coinbase},[]byte{})
	//return NewBlock("Gebesis Block",[]byte{})
}