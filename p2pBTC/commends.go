package main

import "fmt"

func (cli *CLI) CreactChain(addrss string)  {
	bc := NewBlockChain(addrss)
	defer bc.db.Close()
	fmt.Println("creact block chain successfully")

}


func (cli * CLI)PrintChain()  {
	bc := GetBlockChainHandler()
	it := bc.Itertor()
	for {
		block := it.Next()

		fmt.Printf("Transactions: %s\n", block.Transactions)
		fmt.Println("Version", block.Version)
		fmt.Printf("PrevBlockHash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("TimeStamp: %d\n", block.TimeStamp)
		fmt.Printf("MerkelRoot: %x\n", block.MerkelRoot)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		pow := NewProofOfWork(block)
		fmt.Printf("isvalid :%v\n",pow.IsValid())
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

}

func (cli *CLI)GetBalance(address string)  {
	bc := GetBlockChainHandler()
	defer bc.db.Close()
	var total float64
	utxos := bc.FindUTXOs(address)
	for _,utxo := range utxos{
		total += utxo.Value
	}
	fmt.Printf("the balance of %s  is %f\n",address,total)
}

func (cli *CLI) Send(from , to string, amount float64)  {
	bc := GetBlockChainHandler()
	tx := NewTransction(from,to,amount,bc)
	bc.AddBlock([]*Transaction{tx})
	fmt.Println("send successfilly! 成功")
}