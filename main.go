package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/0glabs/0g-storage-client/common/blockchain"
	"github.com/0glabs/0g-storage-client/indexer"
	"github.com/0glabs/0g-storage-client/transfer"
	"github.com/openweb3/web3go"
)

// Network configuration for 0G Testnet
const (
	// EVM RPC endpoint for interacting with the blockchain
	EvmRPC = "https://evmrpc-testnet.0g.ai"

	// Storage indexer endpoints
	// Standard is cheaper but slower, Turbo is faster but more expensive
	IndexerRPCStandard = "https://indexer-storage-testnet-standard.0g.ai"
	IndexerRPCTurbo    = "https://indexer-storage-testnet-turbo.0g.ai"

	// Number of storage replicas to maintain
	DefaultReplicas = 3
)

// StorageClient wraps the necessary clients and configuration
type StorageClient struct {
	web3Client    *web3go.Client
	indexerClient *indexer.Client
	ctx           context.Context
}

// NewStorageClient creates a new instance of StorageClient
func NewStorageClient(ctx context.Context, privateKey string, useTurbo bool) (*StorageClient, error) {
	// Initialize Web3 client for blockchain interactions
	web3Client := blockchain.MustNewWeb3(EvmRPC, privateKey)

	// Select appropriate indexer RPC based on speed preference
	indexerRPC := IndexerRPCStandard
	if useTurbo {
		indexerRPC = IndexerRPCTurbo
	}

	// Initialize indexer client for node selection
	indexerClient, err := indexer.NewClient(indexerRPC)
	if err != nil {
		web3Client.Close()
		return nil, fmt.Errorf("failed to create indexer client: %v", err)
	}

	return &StorageClient{
		web3Client:    web3Client,
		indexerClient: indexerClient,
		ctx:           ctx,
	}, nil
}

// Close cleans up client connections
func (c *StorageClient) Close() {
	if c.web3Client != nil {
		c.web3Client.Close()
	}
}

// UploadFile uploads a file to 0G Storage
func (c *StorageClient) UploadFile(filePath string) (string, error) {

	// Select storage nodes with desired replicas
	nodes, err := c.indexerClient.SelectNodes(c.ctx, 1, DefaultReplicas, nil)
	if err != nil {
		return "", fmt.Errorf("failed to select storage nodes: %v", err)
	}

	// Create uploader instance
	uploader, err := transfer.NewUploader(c.ctx, c.web3Client, nodes)
	if err != nil {
		return "", fmt.Errorf("failed to create uploader: %v", err)
	}

	// Upload the file with timeout
	ctx, cancel := context.WithTimeout(c.ctx, 5*time.Minute)
	defer cancel()

	txHash, rootHash, err := uploader.UploadFile(ctx, filePath)
	if err != nil {
		return "", fmt.Errorf("upload failed: %v", err)
	}
	fmt.Printf("Upload successful!\nTx hash: %s\nRoot hash: %s\n", txHash, rootHash)

	return rootHash.String(), nil
}

// DownloadFile downloads a file from 0G Storage
func (c *StorageClient) DownloadFile(rootHash, outputPath string) error {
	// Select storage nodes
	nodes, err := c.indexerClient.SelectNodes(c.ctx, 1, DefaultReplicas, nil)
	if err != nil {
		return fmt.Errorf("failed to select storage nodes: %v", err)
	}

	// Create downloader instance
	downloader, err := transfer.NewDownloader(nodes)
	if err != nil {
		return fmt.Errorf("failed to create downloader: %v", err)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Download with proof verification and timeout
	ctx, cancel := context.WithTimeout(c.ctx, 5*time.Minute)
	defer cancel()

	if err := downloader.Download(ctx, rootHash, outputPath, true); err != nil {
		return fmt.Errorf("download failed: %v", err)
	}

	return nil
}

func main() {
	// Parse command line flags
	privateKey := flag.String("key", "", "Private key for transactions (required)")
	uploadPath := flag.String("upload", "", "Path to file to upload")
	downloadHash := flag.String("download", "", "Root hash of file to download")
	outputPath := flag.String("output", "", "Path to save downloaded file")
	useTurbo := flag.Bool("turbo", false, "Use Turbo endpoint for faster but more expensive operations")
	flag.Parse()

	// Validate required flags
	if *privateKey == "" {
		log.Fatal("Private key is required. Use -key flag.")
	}

	if (*uploadPath == "") == (*downloadHash == "") {
		log.Fatal("Specify either -upload or -download (with -output), but not both or neither")
	}

	if *downloadHash != "" && *outputPath == "" {
		log.Fatal("Output path (-output) is required when downloading")
	}

	// Create storage client
	ctx := context.Background()
	client, err := NewStorageClient(ctx, *privateKey, *useTurbo)
	if err != nil {
		log.Fatalf("Failed to initialize storage client: %v", err)
	}
	defer client.Close()

	// Handle upload
	if *uploadPath != "" {
		rootHash, err := client.UploadFile(*uploadPath)
		if err != nil {
			log.Fatalf("Upload failed: %v", err)
		}
		fmt.Printf("Upload successful!\nRoot hash: %s\n", rootHash)
		return
	}

	// Handle download
	if *downloadHash != "" {
		if err := client.DownloadFile(*downloadHash, *outputPath); err != nil {
			log.Fatalf("Download failed: %v", err)
		}
		fmt.Printf("Download successful! File saved to: %s\n", *outputPath)
	}
}
