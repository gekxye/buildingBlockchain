package main

import "crypto/sha256"



type MerkleTree struct {
	RootNode *MerkleNode
}

type MerkleNode struct {
	Left *MerkleNode
	Right *MerkleNode
	Data []byte
}

func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode{
	mnode :=MerkleNode{}

	if left ==nil && right == nil{
		mnode.Data = data
	}else {
		prevhashes := append(left.Data,right.Data...)
		firsthash :=sha256.Sum256(prevhashes)
		hash := sha256.Sum256(firsthash[:])
		mnode.Data = hash[:]
	}

	mnode.Left = left
	mnode.Right = right

	return &mnode
}


//构建默克尔树
func NewMerkleTree(data [][]byte)  *MerkleTree{

	var nodes []MerkleNode

	for _,datum :=range data{
		node := NewMerkleNode(nil,nil,datum)
		nodes = append(nodes,*node)
	}

	j :=0
	for nSize :=len(data);nSize>1;nSize=(nSize+1)/2{
		for i:=0;i<nSize;i+=2{
			i2 :=min(i+1,nSize-1)

			node :=NewMerkleNode(&nodes[j+1],&nodes[j+i2],nil)
			nodes = append(nodes,*node)
		}
		j +=nSize
	}

	mTree :=MerkleTree{&(nodes[len(nodes)-1])}
	return &mTree
}
