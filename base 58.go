package main

import (
	"bytes"
	"math/big"
)

//使用切片存储base58字母
var b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

//base58编码函数
func Base58Encode(input []byte) []byte {

	//定义一个字节切片
	var result []byte

	//把字节数组input转化为大整数
	x := big.NewInt(0).SetBytes(input)
	//长度58的大整数
	base := big.NewInt(int64(len(b58Alphabet)))
	//0的大整数
	zero :=big.NewInt(0)
	//大整数的指针
	mod :=&big.Int{}

	//循环，不停对x取余，大小为58，直到除数为0
	for x.Cmp(zero) !=0{
		x.DivMod(x,base,mod)
		//将余数添加到数组当中
		result = append(result,b58Alphabet[mod.Int64()])
	}


	//反转字节数组
	ReverseBytes(result)

	//如果字节数组的前面为字节0，会把它替换为1，比特币中的特殊做法
	for _,b:=range input{
		if b == 0x00{
			result = append([]byte{b58Alphabet[0]},result...)
		}else {
			break
		}
	}

	return result

}


//base58解码函数
func Base58Decode(input []byte) []byte{

	result := big.NewInt(0)

	//将前面的1变成0
	zeroBytes :=0
	for _,b :=range input{
		if b =='1'{
			zeroBytes++
		}else {
			break
		}
	}

	//除去前面的1
	payload := input[zeroBytes:]

	//循环，逆推出结果
	for _,b := range payload{
		charIndex := bytes.IndexByte(b58Alphabet,b)  //反推余数
		result.Mul(result,big.NewInt(58))         //之前的结果乘以58
		result.Add(result,big.NewInt(int64(charIndex)))  //加上余数
	}

	//将大整数转化成字节数组
	decoded := result.Bytes()

	//在前面填充0字节
	decoded = append(bytes.Repeat([]byte{0x00},zeroBytes),decoded...)
	return decoded

}
