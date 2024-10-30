package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"server/pkg"
)

type imageInfo struct {
	width, height int
	size         int64
}

func getImageInfo(filepath string) (*imageInfo, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Get file size
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// Get image dimensions
	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return nil, err
	}

	return &imageInfo{
		width:  img.Width,
		height: img.Height,
		size:   stat.Size(),
	}, nil
}


func getRandomCoverPath() string {
	// get randome file from ./assets/covers
	covers, err := os.ReadDir("./assets/covers")
	if err != nil {
		return ""
	}
	return "./assets/covers/" + covers[rand.Intn(len(covers))].Name()
}

func main() {
	// Example encryption key
	key := []byte("MySecretKey123")

	// Get original image info
	originalInfo, err := getImageInfo("./input.png")
	if err != nil {
		log.Fatalf("Failed to get original image info: %v", err)
	}
	fmt.Printf("Original image size: %dx%d pixels (%d bytes)\n",
		originalInfo.width, originalInfo.height, originalInfo.size)

	// Example 1: Encrypt an image and hide it in a cover image
	err = pkg.EncryptImage(
		"./input.png",
		"./encrypted.png",
		getRandomCoverPath(),
		key,
	)
	if err != nil {
		log.Fatalf("Encryption failed: %v", err)
	}

	// Get encrypted image info and show comparison
	encryptedInfo, err := getImageInfo("./encrypted.png")
	if err == nil {
		fmt.Printf("Encrypted image size: %dx%d pixels (%d bytes)\n",
			encryptedInfo.width, encryptedInfo.height, encryptedInfo.size)
		fmt.Printf("Size ratio (encrypted/original): %.2fx\n",
			float64(encryptedInfo.size)/float64(originalInfo.size))
	}
	fmt.Println("Successfully encrypted and hidden the image!")

	// Example 2: Decrypt a hidden image
	err = pkg.DecryptImage(
			"./encrypted.png",
			"./decrypted.png",
			key,
	)
	if err != nil {
		log.Fatalf("Decryption failed: %v", err)
	}

	// Get decrypted image info and show comparison
	if decryptedInfo, err := getImageInfo("./decrypted.png"); err == nil {
		fmt.Printf("Decrypted image size: %dx%d pixels (%d bytes)\n",
			decryptedInfo.width, decryptedInfo.height, decryptedInfo.size)
		fmt.Printf("Size ratio (decrypted/original): %.2fx\n",
			float64(decryptedInfo.size)/float64(originalInfo.size))
	} else {
		fmt.Println("Failed to get decrypted image info:", err)
	}
	fmt.Println("Successfully decrypted the hidden image!")
}
