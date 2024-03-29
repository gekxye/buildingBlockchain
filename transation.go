package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
)

const subsidy  = 100

type Transation struct {
	ID []byte
	Vin []TXInput
	Vout []TXOutput
}

type TXInput struct {
	TXid []byte
	Voutindex int
	Signature []byte
	Pubkey []byte
}

type TXOutput struct {
	Value int
	PubkeyHash []byte
}

type TXOutputs struct {
	Outputs []TXOutput
}


func (outs TXOutputs) Serialize() []byte{

	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(outs)

	if err !=nil{
		log.Panic(err)
	}

	return buff.Bytes()
}

func DeserializeOutputs(data []byte) TXOutputs{

	var outputs TXOutputs
	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)

	if err !=nil{
		log.Panic(err)
	}

	return outputs
}




func (out *TXOutput) Lock(address []byte){
	decodeAddress := Base58Decode(address)
	pubkeyhash := decodeAddress[1:len(decodeAddress)-4]
	out.PubkeyHash =  pubkeyhash
}


//将交易变成字符串，并做拼接，便于打印查看
func (tx Transation) String() string {
	var lines []string

	lines = append(lines,fmt.Sprintf("---Transation %x:",tx.ID))

	for i,input :=range tx.Vin{
		lines = append(lines,fmt.Sprintf(" Input %d:",i))
		lines = append(lines,fmt.Sprintf(" Txid %x:",input.TXid))
		lines = append(lines,fmt.Sprintf(" Voutindex %d:",input.Voutindex))
		lines = append(lines,fmt.Sprintf(" Signature %x:",input.Signature))
	}

	for i,output :=range tx.Vout{
		lines = append(lines,fmt.Sprintf(" Output %d:",i))
		lines = append(lines,fmt.Sprintf(" Value %d:",output.Value))
		lines = append(lines,fmt.Sprintf(" Script %x:",output.PubkeyHash))
	}

	return strings.Join(lines,"\n")

}


//将交易的结构体序列化
func (tx Transation) Serialize() []byte {
	var encoded bytes.Buffer
	enc :=gob.NewEncoder(&encoded)
	err:=enc.Encode(tx)
	if err!=nil{
		log.Panic(err)
	}
	return encoded.Bytes()
}

//将交易进行hash,作为交易ID
func (tx *Transation) Hash() []byte{
	txcopy :=*tx
	txcopy.ID =[]byte{}
	hash :=sha256.Sum256(txcopy.Serialize())
	return hash[:]
}

//定义交易输出函数（交易和地址）
func NewTxOutput(value int, address string) *TXOutput{
	txo :=&TXOutput{value,nil}
	//txo.PubkeyHash = []byte(address)
	txo.Lock([]byte(address))
	return txo
}

//创世区块函数（第一笔coinbase交易）
func NewCoinbaseTX(to,data string) *Transation{
	txin := TXInput{[]byte{},-1,nil,[]byte(data)}
	txout := NewTxOutput(subsidy,to)
	tx := Transation{nil,[]TXInput{txin},[]TXOutput{*txout}}
	tx.ID = tx.Hash()
	return &tx
}


func (out *TXOutput) CanBeUnlockedWith(pubkeyhash []byte) bool{

	return bytes.Compare(out.PubkeyHash,pubkeyhash) == 0
}

func (in * TXInput) canUnlockOutputWith(unlockdata []byte) bool{

	lockinghash :=HashPubkey(in.Pubkey)

	return bytes.Compare(lockinghash,unlockdata)==0

}

func (tx Transation) IsCoinBase() bool{
	return len(tx.Vin) == 1 && len(tx.Vin[0].TXid) ==0 &&  tx.Vin[0].Voutindex == -1
}


func (tx *Transation) Sign(privkey ecdsa.PrivateKey, prevTXs map[string]Transation) {
	if tx.IsCoinBase(){
		return
	}
	//检查过程
	for _,vin :=range tx.Vin{
		if prevTXs[hex.EncodeToString(vin.TXid)].ID == nil{
			log.Panic("ERROR:")
		}
	}

	txcopy:=tx.TrimmedCopy()

	for inID,vin := range txcopy.Vin{
		prevTx := prevTXs[hex.EncodeToString(vin.TXid)] //前一笔交易的结果体

		txcopy.Vin[inID].Signature  = nil
		txcopy.Vin[inID].Pubkey  =  prevTx.Vout[vin.Voutindex].PubkeyHash // 这笔交易的这笔输入的引用的前一笔交易的输出的公钥哈希
		txcopy.ID = txcopy.Hash()

		r,s,err := ecdsa.Sign(rand.Reader,&privkey,txcopy.ID)

		if err !=nil{
			log.Panic(err)
		}

		signature := append(r.Bytes(),s.Bytes()...)

		tx.Vin[inID].Signature = signature

		//txcopy.Vin[inID].Pubkey  = nil
	}

}


func (tx *Transation) TrimmedCopy() Transation {

	var inputs []TXInput
	var outputs []TXOutput


	for _,vin := range tx.Vin {
		inputs = append(inputs,TXInput{vin.TXid,vin.Voutindex,nil,nil})
	}

	for _,vout := range tx.Vout{

		outputs = append(outputs,TXOutput{vout.Value,vout.PubkeyHash})
	}

	txCopy := Transation{tx.ID,inputs,outputs}

	return txCopy
}

func (tx Transation) Verify(prevTxs map[string]Transation) bool {

	if tx.IsCoinBase(){
		return true
	}

	for _,vin := range tx.Vin{

		if prevTxs[hex.EncodeToString(vin.TXid)].ID==nil{
			log.Panic("ERRor")
		}
	}

	txcopy := tx.TrimmedCopy()

	//椭圆曲线
	curve := elliptic.P256()

	for inID,vin := range tx.Vin{
		prevTx:= prevTxs[hex.EncodeToString(vin.TXid)]
		txcopy.Vin[inID].Signature = nil
		txcopy.Vin[inID].Pubkey = prevTx.Vout[vin.Voutindex].PubkeyHash
		txcopy.ID = txcopy.Hash()

		r:=big.Int{}
		s:=big.Int{}

		siglen:=len(vin.Signature)
		r.SetBytes(vin.Signature[:(siglen/2)])
		s.SetBytes(vin.Signature[(siglen/2):])

		x:=big.Int{}
		y := big.Int{}

		keylen :=len(vin.Pubkey)

		x.SetBytes(vin.Pubkey[:(keylen/2)])
		y.SetBytes(vin.Pubkey[(keylen/2):])

		rawPubkey := ecdsa.PublicKey{curve,&x,&y}

		if ecdsa.Verify(&rawPubkey,txcopy.ID,&r,&s) == false{
			return false
		}

		txcopy.Vin[inID].Pubkey =nil
	}

	return true
	
}





func NewUTXOTransation(from,to string,amount int, bc * Blockchain) *Transation{
	var inputs []TXInput
	var outputs []TXOutput

	wallets,err :=NewWallets()
	if err !=nil{
		log.Panic(err)
	}

	wallet :=wallets.GetWallet(from)


	acc,validoutputs := bc.FindSpendableOutputs(HashPubkey(wallet.Publickey),amount)

	if acc < amount{
		log.Panic("Error:Not enough funds")
	}

	for txid,outs := range validoutputs{
		txID ,err := hex.DecodeString(txid)
		if err !=nil{
			log.Panic(err)
		}

		for  _,out := range outs{

			input := TXInput{txID,out,nil,wallet.Publickey}
			inputs  = append(inputs,input)
		}

	}
	outputs  = append(outputs,*NewTxOutput(amount,to))


	if acc > amount{
		outputs = append(outputs,*NewTxOutput(acc-amount,from))
	}


	tx:= Transation{nil,inputs,outputs}
	tx.ID = tx.Hash()

	bc.SignTransation(&tx,wallet.Privatekey)
	return &tx
}


