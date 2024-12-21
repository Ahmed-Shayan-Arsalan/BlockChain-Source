package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	port := flag.String("port", "8001", "Port to run this node on")
	peersFlag := flag.String("peers", "", "Comma-separated list of other peer URLs")
	flag.Parse()

	peers := []string{}
	if *peersFlag != "" {
		for _, p := range splitAndTrim(*peersFlag, ",") {
			peers = append(peers, p)
		}
	}

	node := &Node{
		Blockchain: Blockchain{
			Blocks:     []Block{createGenesisBlock()},
			Difficulty: 2, // Adjust difficulty as needed
		},
		ModelCID:   "bafkreihnj6a5xfcjetmej3mcr4xg364je6vrzbolzz6fmtoxu46ldqlhgm",
		ScalerCID:  "bafkreibo2eo3lx2talokcz6hx44i2e4qxh2vn2tkge52m7vouv2bvqjt3m",
		DatasetCID: "bafkreifkzuwvvohltv4loy5kfgagfhlawrkpjbscbh2qsrz7v6vdy666ae",
		Peers:      peers,
		Port:       *port,
	}

	go node.startServer() // start HTTP server for block validation

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Press Enter to fetch a transaction (need total of 3 transactions).")
	for len(node.PendingTransactions) < 3 {
		_, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		// Fetch one transaction from `predict.py`
		txs, err := node.runPredictionScript(node.DatasetCID, node.ModelCID, node.ScalerCID)
		if err != nil {
			fmt.Println("Error running predictions:", err)
			continue
		}
		if len(txs) > 0 {
			// Take only the first transaction from the returned list
			node.PendingTransactions = append(node.PendingTransactions, txs[0])
			fmt.Printf("Transaction received: %.6f | %s\n", txs[0].Prediction, txs[0].Details)
		} else {
			fmt.Println("No transaction returned, press enter to try again.")
		}
		if len(node.PendingTransactions) < 3 {
			fmt.Printf("%d/3 transactions collected. Press Enter for next transaction.\n", len(node.PendingTransactions))
		}
	}

	fmt.Println("3 transactions collected. Sleeping for 20 seconds before mining...")
	time.Sleep(20 * time.Second)

	// Now proceed with mining
	node.startMiningProcess()

	fmt.Println("Final Blockchain:")
	node.Blockchain.printBlockchain()
}
