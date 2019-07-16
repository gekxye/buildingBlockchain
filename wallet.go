package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/ripemd160"
)

const version  = byte(0x00)

type Wallet struct {

	Privatekey ecdsa.PrivateKey
	Publickey []byte

}

func Newwallet() *Wallet {

	private,public:=newkeyPair()
	wallet:=Wallet{private,public}
	return &wallet

}


//生成私钥和公钥函数
func newkeyPair()  (ecdsa.PrivateKey,[]byte){

	//生成椭圆曲线-secp256r1，比特币用的是secp256k1
	curve :=elliptic.P256()

	private,err :=ecdsa.GenerateKey(curve,rand.Reader)

	if err !=nil{
		fmt.Println("error")
	}
	pubkey :=append(private.PublicKey.X.Bytes(),private.PublicKey.Y.Bytes()...)
	return *private,pubkey

}


func (w Wallet) GetAddress() []byte{

	pubkeyHash:= HashPubkey(w.Publickey)
	versionPayload := append([]byte{version},pubkeyHash...)
	check:=checksum(versionPayload)
	fullPayload := append(versionPayload,check...)
	//返回地址
	address:=Base58Encode(fullPayload)
	return address
}


func HashPubkey(pubkey []byte) []byte{
	pubkeyHash256:=sha256.Sum256(pubkey)
	PIPEMD160Hasher := ripemd160.New()

	_,err:=	PIPEMD160Hasher.Write(pubkeyHash256[:])

	if err!=nil{
		fmt.Println("error")
	}

	publicRIPEMD160 := PIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160

}


func checksum(payload []byte) []byte{
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	//checksum 是前面的4个字节
	checksum:=secondSHA[:4]

	return checksum
}


func ValidateAddress(address []byte) bool{

	pubkeyHash := Base58Decode(address)

	actualCheckSum := pubkeyHash[len(pubkeyHash)-4:]

	publickeyHash  := pubkeyHash[1:len(pubkeyHash)-4]

	targetChecksum := checksum(append([]byte{0x00},publickeyHash...))


	return bytes.Compare(actualCheckSum,targetChecksum)==0
}
