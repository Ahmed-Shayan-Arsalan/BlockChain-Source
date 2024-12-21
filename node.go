package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Node struct {
	Blockchain          Blockchain
	ModelCID            string
	ScalerCID           string
	DatasetCID          string
	Peers               []string
	Port                string
	PendingTransactions []Transaction
}

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
			continue
		}
		transactions = append(transactions, Transaction{
			Prediction: value,
			Details:    fmt.Sprintf("Prediction for Row %d", i+1),
		})
		// For our new requirement, we only need one transaction per run
		break
	}
	return transactions, nil
}

func (n *Node) startMiningProcess() {
	fmt.Println("Starting mining process...")

	// We already have 3 transactions in PendingTransactions
	transactions := n.PendingTransactions[:3]
	n.PendingTransactions = n.PendingTransactions[3:]

	previousBlock := n.Blockchain.Blocks[len(n.Blockchain.Blocks)-1]

	foundBlockChan := make(chan Block, 1)
	var wg sync.WaitGroup

	// Make at least 3 miners
	minerCount := 3
	for i := 0; i < minerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			newBlock := n.Blockchain.mineBlock(transactions, previousBlock)
			// Attempt to send the block to the channel
			select {
			case foundBlockChan <- newBlock:
				// This miner found the block first
			default:
				// Another miner already found a block
			}
		}()
	}

	// Wait for the first found block
	newBlock := <-foundBlockChan
	// At this point, a block is found; other miners will finish soon
	wg.Wait()

	fmt.Println("Block mined, requesting validation from peers...")
	if n.validateWithPeers(newBlock) {
		// Majority valid
		if n.Blockchain.addBlock(newBlock) {
			fmt.Println("Block accepted by majority and added to blockchain.")
		} else {
			fmt.Println("Block was valid by majority but failed local validation.")
		}
	} else {
		fmt.Println("Block rejected by majority.")
	}
}

func (n *Node) validateWithPeers(block Block) bool {
	if len(n.Peers) == 0 {
		// No peers, accept automatically
		return true
	}

	blockData, _ := json.Marshal(block)
	validCount := 0
	var wg sync.WaitGroup
	mu := &sync.Mutex{}

	for _, peer := range n.Peers {
		wg.Add(1)
		go func(peer string) {
			defer wg.Done()
			client := &http.Client{Timeout: 5 * time.Second}
			resp, err := client.Post(peer+"/validate", "application/json", bytes.NewBuffer(blockData))
			if err != nil {
				return
			}
			defer resp.Body.Close()
			var result map[string]string
			json.NewDecoder(resp.Body).Decode(&result)
			if result["status"] == "valid" {
				mu.Lock()
				validCount++
				mu.Unlock()
			}
		}(peer)
	}

	wg.Wait()
	return validCount > len(n.Peers)/2
}

func (n *Node) startServer() {
	http.HandleFunc("/validate", n.handleValidateBlock)
	server := &http.Server{Addr: ":" + n.Port}
	fmt.Println("Listening on port:", n.Port)
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Server error:", err)
	}
}

func (n *Node) handleValidateBlock(w http.ResponseWriter, r *http.Request) {
	var block Block
	err := json.NewDecoder(r.Body).Decode(&block)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"status": "invalid"})
		return
	}

	// Validate block against our current chain
	lastBlock := n.Blockchain.Blocks[len(n.Blockchain.Blocks)-1]
	if n.Blockchain.isBlockValid(block, lastBlock) {
		json.NewEncoder(w).Encode(map[string]string{"status": "valid"})
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "invalid"})
	}
}

// Utility function to split peers
func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
