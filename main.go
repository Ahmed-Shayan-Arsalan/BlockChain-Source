package main

import "fmt"

func main() {
	// Initialize the blockchain with a certain difficulty
	node := Node{
		Blockchain: Blockchain{
			Blocks:     []Block{createGenesisBlock()},
			Difficulty: 4, // Low difficulty for demonstration
		},
		ModelCID:   "bafkreihnj6a5xfcjetmej3mcr4xg364je6vrzbolzz6fmtoxu46ldqlhgm",
		ScalerCID:  "bafkreibo2eo3lx2talokcz6hx44i2e4qxh2vn2tkge52m7vouv2bvqjt3m",
		DatasetCID: "bafkreifkzuwvvohltv4loy5kfgagfhlawrkpjbscbh2qsrz7v6vdy666ae",
	}

	// Simulate mining a block after predictions
	node.mineBlock()

	// Print the final blockchain
	fmt.Println("Final Blockchain:")
	node.Blockchain.printBlockchain()
}
