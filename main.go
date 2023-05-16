package main

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type Config struct {
	TotalFiles int     `json:"total_files"`
	SameRate   float64 `json:"same_rate"`
	TargetDir  string  `json:"target_dir"`
}

func main() {
	configPath := flag.String("config", "config.json", "Path to the configuration file")
	flag.Parse()

	configFile, err := os.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("Failed to unmarshal config file: %v", err)
	}

	err = os.MkdirAll(config.TargetDir, fs.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create target directory: %v", err)
	}

	sameCount := 0
	totalSameCount := 0
	var baseFilename string
	var lastContent []byte
	for i := 0; i < config.TotalFiles; i++ {
		var filename string
		var content []byte
		if i > 0 && randFloat() < config.SameRate {
			content = lastContent
			sameCount++
			totalSameCount++
		} else {
			content = generateRandomContent()
			lastContent = content
			sameCount = 0
			baseFilename = filepath.Join(config.TargetDir, generateRandomFilename())
		}
		filename = baseFilename
		if sameCount > 0 {
			filename += "-" + strconv.Itoa(sameCount)
		}

		err = os.WriteFile(filename, content, 0644)
		if err != nil {
			log.Fatalf("Failed to write file: %v", err)
		}
	}

	fmt.Printf("Number of identical files generated: %d\n", totalSameCount)
}

func generateRandomContent() []byte {
	b := make([]byte, 1024)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("Failed to generate random content: %v", err)
	}
	return b
}

func generateRandomFilename() string {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("Failed to generate random filename: %v", err)
	}
	return hex.EncodeToString(b)
}

func randFloat() float64 {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		log.Fatalf("Failed to generate random float: %v", err)
	}
	return float64(binary.BigEndian.Uint64(b[:])) / float64(1<<64)
}
