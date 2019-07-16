package main

import "fmt"

//默克尔值测试函数
func TestCreateMerkleTreeRoot() {

	//初始化区块
	block :=&Block{
		Version:2,
		PrevBlockHash:[]byte{},
		Merkleroot:[]byte{},
		Hash:[]byte{},
		Time:1418755780,
		Bits:404454260,
		Nonce:0,
		Transations:[]*Transation{},
		Height:0,
	}

	txin := TXInput{[]byte{},-1,nil,nil}
	txout := NewTxOutput(subsidy,"first")
	tx := Transation{nil,[]TXInput{txin},[]TXOutput{*txout}}

	txin2 := TXInput{[]byte{},-1,nil,nil}
	txout2 := NewTxOutput(subsidy,"second")
	tx2 := Transation{nil,[]TXInput{txin2},[]TXOutput{*txout2}}

	var Transations []*Transation
	Transations = append(Transations,&tx,&tx2)
	block.createMerkleTreeRoot(Transations)

	fmt.Printf("%x\n",block.Merkleroot)

}

//序列化新方式测试
func TestNewSerialize()  {
	//初始化区块
	block :=&Block{
		Version:2,
		PrevBlockHash:[]byte{},
		Merkleroot:[]byte{},
		Hash:[]byte{},
		Time:1418755780,
		Bits:404454260,
		Nonce:0,
		Transations:[]*Transation{},
		Height:0,
	}

	deBlock :=DeserializeBlock(block.Serialize())

	deBlock.String()

}


//pow测试
func TestPow() {
	//初始化区块
	block :=&Block{
		Version:2,
		PrevBlockHash:[]byte{},
		Merkleroot:[]byte{},
		Hash:[]byte{},
		Time:1418755780,
		Bits:404454260,
		Nonce:0,
		Transations:[]*Transation{},
		Height:0,
	}

	pow:=NewProofOfWork(block)
	nonce,_:=pow.Run()
	block.Nonce = nonce
	fmt.Println("Pow:",pow.Validdata())

}

func TestBoltDB()  {
	blockchain := NewBlockchain("1EPQE8qPDXt3tjAZa5zT6nCudc79KMxWC2")
	blockchain.MineBlock([]*Transation{})
	blockchain.MineBlock([]*Transation{})
	blockchain.printBlockchain()
	
}