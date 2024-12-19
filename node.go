package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type Node struct {
	Blockchain Blockchain
	ModelCID   string
	ScalerCID  string
	DatasetCID string
}

// runPredictionScript executes the Python script and returns up to 3 predictions
func (n *Node) runPredictionScript(datasetCID, modelCID, scalerCID string) ([]Transaction, error) {
	cmd := exec.Command("python", "predict.py", datasetCID, modelCID, scalerCID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Python script error:\n%s\n", string(output))
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	var transactions []Transaction
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		value, err := strconv.ParseFloat(line, 64)
		if err != nil {
			// Ignore non-numeric lines
			continue
		}
		transactions = append(transactions, Transaction{
			Prediction: value,
			Details:    fmt.Sprintf("Prediction for Row %d", i+1),
		})
		if len(transactions) == 3 {
			break
		}
	}
	return transactions, nil
}

// mineBlock gets predictions and mines a single block synchronously
func (node *Node) mineBlock() {
	fmt.Println("Running predictions...")
	transactions, err := node.runPredictionScript(node.DatasetCID, node.ModelCID, node.ScalerCID)
	if err != nil {
		fmt.Println("Error running predictions:", err)
		return
	}

	if len(transactions) == 0 {
		fmt.Println("No transactions were added. Check prediction output.")
		return
	}

	fmt.Println("Mining new block...")
	previousBlock := node.Blockchain.Blocks[len(node.Blockchain.Blocks)-1]
	newBlock := node.Blockchain.mineBlock(transactions, previousBlock)
	if !node.Blockchain.addBlock(newBlock) {
		fmt.Println("Failed to add new block.")
	}
}
