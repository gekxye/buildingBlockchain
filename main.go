package main

func main() {
	//测试默克尔树
	//TestCreateMerkleTreeRoot()

	//初始化目标hash测试
	//target :=big.NewInt(1)
	//target.Lsh(target,uint(256-targetBits))
	//fmt.Printf("%x\n",target.Bytes())

	//Pow测试
	//TestPow()

	//序列化新方式测试
	//TestNewSerialize()

	//NewGensisBlock()

	//TestBoltDB()

	bc := NewBlockchain("1EPQE8qPDXt3tjAZa5zT6nCudc79KMxWC2")
	cli :=CLI{bc}
	cli.Run()


	//wallet :=Newwallet()
	//fmt.Printf("私钥：%x\n",wallet.Privatekey.D.Bytes())
	//fmt.Printf("公钥：%x\n",wallet.Publickey)
	//fmt.Printf("地址：%x\n",wallet.GetAddress())
	//
	//address,_:=hex.DecodeString("3134695365467169313275417172347a6e724554396641784c7639527775534e6457")
	//fmt.Printf("%d\n",ValidateAddress(address))





}