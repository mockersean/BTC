package main

import (
	"math/big"
	"bytes"
	"math"
	"crypto/sha256"
	"fmt"
)

const targetBits  =	24
type ProofOfoWork struct {
	block *Block
	targetBit *big.Int
}

func NewProofOfWork(block *Block) *ProofOfoWork {
	var IntTarget = big.NewInt(1)
	//向左移动 256 - 24位
	IntTarget.Lsh(IntTarget,uint(256 - targetBits))
	return &ProofOfoWork{block,IntTarget}
}

func (pow *ProofOfoWork)PrepareRawData(nonce int64) []byte  {
	block := pow.block

	tmp := [][]byte{
		IntToByte(block.Version) ,
		block.PrevBlockHash,
		IntToByte(block.TimeStamp),
		block.MerkelRoot,
		IntToByte(nonce),
		IntToByte(targetBits),
		//block.Data,
		//block.Transactions  //TODO
	}
	data := bytes.Join(tmp,[]byte{})
	return data
}

func (pow *ProofOfoWork) Run() (int64,[]byte) {
	var nonce int64
	var hash [32]byte
	var hashInt big.Int
	fmt.Printf("正在查找\n")
	for nonce < math.MaxInt64 {
		data := pow.PrepareRawData(nonce)
		hash = sha256.Sum256(data)
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.targetBit) == -1 {
			fmt.Printf("found hash %x\n",hash)
			break
		}else {
			nonce++
		}
	}
	return nonce,hash[:]
}

func (pow *ProofOfoWork) IsValid() bool {
	data := pow.PrepareRawData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	var IntHash big.Int
	IntHash.SetBytes(hash[:])
	return IntHash.Cmp(pow.targetBit) == -1
}