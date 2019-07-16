package main

import (
	"bytes"
	"crypto/sha256"
	"math/big"
)

type ProofOfWork struct {
	block * Block
	target * big.Int
}

const targetBits = 16

//初始化
func NewProofOfWork(b *Block) *ProofOfWork {

	target :=big.NewInt(1)
	target.Lsh(target,uint(256-targetBits))
	pow := &ProofOfWork{b,target}

	return pow
}

//序列化
func (pow *ProofOfWork) prepareData(nonce int32) []byte {

	data :=bytes.Join(
		[][]byte{
			IntToHex(pow.block.Version),
			pow.block.PrevBlockHash,
			pow.block.Merkleroot,
			IntToHex(pow.block.Time),
			IntToHex(pow.block.Bits),
			IntToHex(nonce)},
		[]byte{},
	)
	return data

}

//挖矿
func (pow *ProofOfWork) Run() (int32,[]byte) {

	var nonce int32
	var secondHash [32]byte
	nonce = 0
    var currenthash big.Int

		for nonce<maxnonce {

			//序列化
			data := pow.prepareData(nonce)

			//double哈希
			firstHash := sha256.Sum256(data)
			secondHash = sha256.Sum256(firstHash[:])
			//fmt.Printf("%x\n",secondHash)

			currenthash.SetBytes(secondHash[:])
			//比较目标hash和当前hash，当前hash小于目标hash时，挖矿成功
			if currenthash.Cmp(pow.target) == -1 {
				break
			} else {
				nonce++
			}
		}

        return nonce,secondHash[:]

}


//验证挖矿是否有效
func (pow * ProofOfWork) Validdata() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)

	firstHash := sha256.Sum256(data)
	secondHash := sha256.Sum256(firstHash[:])
	hashInt.SetBytes(secondHash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}