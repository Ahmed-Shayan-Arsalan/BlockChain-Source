package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Transaction struct {
	Prediction float64 `json:"prediction"`
	Details    string  `json:"details"`
}

type Block struct {
	Index        int           `json:"index"`
	Timestamp    string        `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
	PreviousHash string        `json:"previousHash"`
	Hash         string        `json:"hash"`
	Nonce        int           `json:"nonce"`
}

type Blockchain struct {
	Blocks     []Block
	Difficulty int // Number of leading zeros required
}

// calculateHash computes a SHA256 hash of the block's content
func calculateHash(block Block) string {
	record, _ := json.Marshal(block.Transactions)
	data := fmt.Sprintf("%d%s%s%s%d", block.Index, block.Timestamp, string(record), block.PreviousHash, block.Nonce)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// createGenesisBlock initializes the chain with a genesis block
func createGenesisBlock() Block {
	genesisTransaction := Transaction{
		Prediction: 0.0,
		Details:    "Genesis Block",
	}
	genesisBlock := Block{
		Index:        0,
		Timestamp:    time.Now().String(),
		Transactions: []Transaction{genesisTransaction},
		PreviousHash: "0",
		Nonce:        0,
	}
	genesisBlock.Hash = calculateHash(genesisBlock)
	return genesisBlock
}

// isBlockValid checks if a new block is properly linked and correct
func (bc *Blockchain) isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}
	if oldBlock.Hash != newBlock.PreviousHash {
		return false
	}
	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	// Check difficulty by verifying the hash's prefix
	prefix := strings.Repeat("0", bc.Difficulty)
	if !strings.HasPrefix(newBlock.Hash, prefix) {
		return false
	}
	return true
}

// mineBlock performs the proof-of-work until a valid hash is found
func (bc *Blockchain) mineBlock(transactions []Transaction, previousBlock Block) Block {
	newBlock := Block{
		Index:        previousBlock.Index + 1,
		Timestamp:    time.Now().String(),
		Transactions: transactions,
		PreviousHash: previousBlock.Hash,
		Nonce:        0,
	}

	prefix := strings.Repeat("0", bc.Difficulty)
	for {
		newBlock.Hash = calculateHash(newBlock)
		if strings.HasPrefix(newBlock.Hash, prefix) {
			break
		}
		newBlock.Nonce++
	}
	return newBlock
}

// addBlock adds a mined block to the chain if it's valid
func (bc *Blockchain) addBlock(newBlock Block) bool {
	previousBlock := bc.Blocks[len(bc.Blocks)-1]
	if bc.isBlockValid(newBlock, previousBlock) {
		bc.Blocks = append(bc.Blocks, newBlock)
		fmt.Printf("New Block Added: %+v\n", newBlock)
		return true
	}
	return false
}

// printBlockchain shows the entire blockchain with all transactions
func (bc *Blockchain) printBlockchain() {
	for _, block := range bc.Blocks {
		fmt.Printf("Index: %d\nTimestamp: %s\nPreviousHash: %s\nHash: %s\nTransactions:\n",
			block.Index, block.Timestamp, block.PreviousHash, block.Hash)
		for _, tx := range block.Transactions {
			fmt.Printf("  Prediction: %.6f, Details: %s\n", tx.Prediction, tx.Details)
		}
		fmt.Println()
	}
}
