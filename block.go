package main

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"
)

var (
	maxnonce int32 = math.MaxInt32
)

//定义区块结构体
type Block struct {
	Version int32
	PrevBlockHash []byte
	Merkleroot []byte
	Hash []byte
	Time int32
	Bits int32
	Nonce int32
	Transations []*Transation
	Height int32
}


//结构体序列化
func (block *Block) serialize() []byte{
	result :=bytes.Join(
		[][]byte{
			IntToHex(block.Version),
			block.PrevBlockHash,
			block.Merkleroot,
			IntToHex(block.Time),
			IntToHex(block.Bits),
			IntToHex(block.Nonce)},
		[]byte{},
	)
	return result
}

//区块结构体序列化的优化
func (b *Block) Serialize()  []byte{

	var encoded bytes.Buffer
	enc :=gob.NewEncoder(&encoded)
	err:=enc.Encode(b)
	if err!=nil{
		log.Panic(err)
	}
	return encoded.Bytes()
}

//区块结构体反序列化
func DeserializeBlock(d []byte)  *Block{
	var block Block
	decode :=gob.NewDecoder(bytes.NewReader(d))
	err:=decode.Decode(&block)
	if err!=nil{
		log.Panic(err)
	}
	return &block
}


//计算目标hash函数
func CalculateTargetFast(bits []byte)  []byte{

	var result []byte

	//计算指数
	exponent :=bits[:1]
	fmt.Printf("%x\n",exponent)

	//计算系数（后面3个）
	coeffient :=bits[1:]
	fmt.Printf("%x\n",coeffient)

	//将字节，16进制为18，转化成了字符串型
	str:=hex.EncodeToString(exponent)
	fmt.Printf("str=%s\n",str)

	//将字符串型转化成了10进制int型
	exp, _ := strconv.ParseInt(str,16,8)
	fmt.Printf("%d\n",exp)

	//拼接在一起，计算出目标hash
	result = append(bytes.Repeat([]byte{0x00},32-int(exp)),coeffient...)
	result = append(result,bytes.Repeat([]byte{0x00},32-len(result))...)

	return result

}


//计算默克尔树根节点函数
func (b*Block) createMerkleTreeRoot(transations []*Transation)  {
    var tranHash [][]byte

    for _,tx := range transations{
    	tranHash = append(tranHash,tx.Hash())
	}

    mTree := NewMerkleTree(tranHash)

    b.Merkleroot = mTree.RootNode.Data
}



func (b*Block) String() {
	fmt.Printf("version:%s\n",strconv.FormatInt(int64(b.Version),10))
	fmt.Printf("PrevBlockHash:%x\n",b.PrevBlockHash)
	fmt.Printf("Merkleroot:%x\n",b.Merkleroot)
	fmt.Printf("Hash:%x\n",b.Hash)
	fmt.Printf("Time:%s\n",strconv.FormatInt(int64(b.Time),10))
	fmt.Printf("Bits:%s\n",strconv.FormatInt(int64(b.Bits),10))
	fmt.Printf("Nonce:%s\n",strconv.FormatInt(int64(b.Nonce),10))
}


//创建新区块
func NewBlock(transations []*Transation,prevBlockHash []byte,height int32) *Block  {
	block :=&Block{
		2,
		prevBlockHash,
		[]byte{},
		[]byte{},
		int32(time.Now().Unix()),
		404454260,
		0,
		transations,
		height,
	}

	pow:=NewProofOfWork(block)
	nonce,hash:=pow.Run()
	block.Nonce=nonce
	block.Hash=hash

	return block
}



//创建创世区块
func NewGensisBlock(transations []*Transation) *Block  {
	block :=&Block{
		2,
		[]byte{},
		[]byte{},
		[]byte{},
		int32(time.Now().Unix()),
		404454260,
		0,
		transations,
		0,
	}

	pow:=NewProofOfWork(block)
	nonce,hash:=pow.Run()
	block.Nonce=nonce
	block.Hash=hash

	//block.String()

	return block
}


