package main

import (
	"encoding/gob"
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
)

const rewaord float64 = 12.5

type Transaction struct {
	TXID	 []byte
	TXInputs []Input
	TXOutPut []OutPut
}

type Input struct {
	Txid []byte
	ReferOutputIndex int64
	UnlockScript string   //scriptSig

}

type OutPut struct {
	Value float64
	LockScript string   //scriptPubKey

}

func (tx *Transaction)SetTXID()  {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	encoder.Encode(tx)

	hash := sha256.Sum256(buffer.Bytes())
	tx.TXID = hash[:]

}

func (input *Input) CanUnlockUTXOByAddress(unlockdata string) bool {
	return input.UnlockScript == unlockdata
}

func (output *OutPut)CanBeUnlockedByAddress(unlockdata string) bool  {
	return output.LockScript == unlockdata
}



func NewCoinbaseTX(address ,msg string) *Transaction  {
	if msg == "" {
		fmt.Sprintf(msg,"Current Reword is : %f\n",rewaord)
	}

	input := Input{nil,-1,msg}
	output := OutPut{rewaord,address}
	
	tx := Transaction{nil, []Input{input}, []OutPut{output}}
	tx.SetTXID()
	return &tx
}

func NewTransction(from,to string,amount float64,bc *BlockChain) *Transaction{
	counted,container := bc.FindSuitableUTXOs(from,amount)

	if counted < amount{
		 fmt.Println("余额额不足")
		 os.Exit(1)
	}

	var inputs []Input
	var outputs []OutPut
	for txid,outputIndexs := range container {
		for _, index := range outputIndexs{
			input := Input{[]byte(txid),index,from}
			inputs = append(inputs,input)
		}
	}

	output := OutPut{amount,to}
	outputs = append(outputs,output)

	if counted > amount {
		 outputs = append(outputs,OutPut{counted - amount,from})
	}
	
	tx :=Transaction{nil,inputs,outputs}
	tx.SetTXID()
	return &tx
}

func (tx *Transaction)IsCoinbase() bool {
	if len(tx.TXInputs) == 1 {
		if tx.TXInputs[0].Txid == nil && tx.TXInputs[0].ReferOutputIndex == -1 {
			return true
		}
	}
	return false
}