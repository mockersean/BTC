package main

import (
	"github.com/boltdb/bolt"
	"os"
	"fmt"
)

const dbfile  = "blockChainDB.db"
const blockBucket  = "block"
const lastHash = "lastHash"
const gensisBlockInfo  = "黄义庆"


type  BlockChain struct {
	//blocks []*Block
	db *bolt.DB
	lastHash []byte
}

func NewBlockChain(address string) *BlockChain  {

	if IsBlockChainExist() {
		fmt.Println("Block chain is olready extst!")
		os.Exit(1)
	}
	db,err := bolt.Open(dbfile,0600,nil)
	CheckErr(err)
	var lasthash []byte

	err = db.Update(func(tx *bolt.Tx) error {

		coinbase := NewCoinbaseTX(address,gensisBlockInfo)
		genesis := NewGenesisBlock(coinbase)

		bucket,err :=tx.CreateBucket([]byte(blockBucket))
		CheckErr(err)

		err = bucket.Put(genesis.Hash,genesis.Serialize())
		CheckErr(err)

		err = bucket.Put([]byte(lastHash), genesis.Hash)
		CheckErr(err)
		lasthash = genesis.Hash

		return nil
	})

	CheckErr(err)
	return &BlockChain{db,lasthash}
}

func GetBlockChainHandler() *BlockChain	 {
	if !IsBlockChainExist() {
		fmt.Println("block chain not exist,pleasse creat first")
		os.Exit(1)
	}
	db,err := bolt.Open(dbfile,0600,nil)
	CheckErr(err)
	var lasthash []byte
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket != nil{
			lasthash = bucket.Get([]byte(lastHash))
		}
		return nil
	})
	CheckErr(err)

	return &BlockChain{db,lasthash}
}

//添加区块
//func (bc * BlockChain) AddBlock(data string)  {
func (bc * BlockChain) AddBlock(transactions []*Transaction)  {
	var prevBlockHash []byte
	err := bc.db.View(func(tx *bolt.Tx) error { //只有取
		bucket := tx.Bucket([]byte(blockBucket))
		lasthash := bucket.Get([]byte(lastHash))
		prevBlockHash = lasthash
		return nil
	})
	CheckErr(err)
 	block := NewBlock(transactions,prevBlockHash)

 	err = bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))

		err := bucket.Put(block.Hash,block.Serialize())
		CheckErr(err)

		err = bucket.Put([]byte(lastHash), block.Hash)
		CheckErr(err)
		bc.lastHash = block.Hash
 		return nil
	})
 	CheckErr(err)
}

type BlockChainItertor struct {
	db *bolt.DB
	currentHsah []byte

}

func (bc *BlockChain)Itertor() *BlockChainItertor {
	return &BlockChainItertor{bc.db,bc.lastHash}
}

func (it *BlockChainItertor)Next() *Block{
	var block *Block
	err := it.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			os.Exit(1)
		}
		blockTmp := bucket.Get(it.currentHsah)
		block = Deserialize(blockTmp)
		it.currentHsah = block.PrevBlockHash
		return nil
	})
	CheckErr(err)
	return block
}

func IsBlockChainExist() bool  {
	_ ,err := os.Stat(dbfile)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (bc *BlockChain)FindUnspendTransacions(address string) []Transaction {
	bci := bc.Itertor()
	var transactons []Transaction
	var spentUTXOs = make(map[string/*交易txid*/][]int64)

	for  {
		block := bci.Next()

		for _,currenTx := range block.Transactions{
			txid := string(currenTx.TXID)

			//遍历当前交易的inputs 找到和当前地址消耗 utxo
			for _,input := range currenTx.TXInputs{
				if currenTx.IsCoinbase() == false {
					if input.CanUnlockUTXOByAddress(address) {
						spentUTXOs[string(input.Txid)] = append(spentUTXOs[string(input.Txid)],input.ReferOutputIndex)
					}
				}
			}

			LABEL1:
			//遍历当前交易的outputs,通过output的解锁条件，确定满足条件的交易
			for outputIndex, output := range currenTx.TXOutPut{
				 if spentUTXOs[txid] != nil {
					 for _,usedIndex := range spentUTXOs[txid]{

						 if int64(outputIndex) == usedIndex{
						 	  fmt.Println("user,no need to add again!")
							  continue LABEL1
						 }
					 }
				 }

				if output.CanBeUnlockedByAddress(address) {
					transactons = append(transactons,*currenTx)
				}
			}

		}


		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return transactons
}

func (bc *BlockChain)FindUTXOs(address string) []OutPut  {
	var	outputs []OutPut
	txs := bc.FindUnspendTransacions(address)
	for _,tx := range txs{
		for _,output := range tx.TXOutPut {
			if output.CanBeUnlockedByAddress(address) {
				outputs = append(outputs,output)
			}
		}
	}

	return outputs
}

func (bc *BlockChain)FindSuitableUTXOs(address string, amount float64) (float64,map[string][]int64)  {
	txs := bc.FindUnspendTransacions(address)
	var	countTotal float64
	var container = make(map[string][]int64)

	LABEL2:
	for _,tx := range txs{
		for index,output := range tx.TXOutPut {
			if countTotal < amount {
				if output.CanBeUnlockedByAddress(address) {
					countTotal += output.Value
					container[string(tx.TXID)] = append(container[string(tx.TXID)],int64(index))
				}
			}else {
				break LABEL2
			}

		}
	}

	return countTotal,container
}

