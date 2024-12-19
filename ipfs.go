package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// DownloadFileFromIPFS fetches a file from IPFS using the gateway
func DownloadFileFromIPFS(cid string, outputPath string) error {
	url := fmt.Sprintf("https://gateway.pinata.cloud/ipfs/%s", cid)
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}
	defer response.Body.Close()

	// Create the output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	fmt.Printf("File downloaded successfully: %s\n", outputPath)
	return nil
}
