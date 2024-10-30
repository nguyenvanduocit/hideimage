package pkg

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestEncryptDecryptImage(t *testing.T) {
	// Test setup
	tempDir := t.TempDir()
	testKey := []byte("test-key-12345")
	
	// Test paths
	inputPath := filepath.Join(tempDir, "input.png")
	encryptedPath := filepath.Join(tempDir, "encrypted.png")
	decryptedPath := filepath.Join(tempDir, "decrypted.png")
	coverPath := filepath.Join(tempDir, "cover.png")

	// Create a simple valid PNG file for testing
	samplePNG := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG signature
		// IHDR chunk
		0x00, 0x00, 0x00, 0x0D, // length
		0x49, 0x48, 0x44, 0x52, // "IHDR"
		0x00, 0x00, 0x00, 0x01, // width=1
		0x00, 0x00, 0x00, 0x01, // height=1
		0x08, 0x06, 0x00, 0x00, 0x00, // bit depth=8, color type=6, etc
		0x1F, 0x15, 0xC4, 0x89, // CRC
		// IDAT chunk
		0x00, 0x00, 0x00, 0x0A, // length=10
		0x49, 0x44, 0x41, 0x54, // "IDAT"
		0x08, 0xD7, 0x63, 0x60, 0x60, 0x60, // some sample data
		0x00, 0x00, 0x00, 0x00, // CRC
		// IEND chunk
		0x00, 0x00, 0x00, 0x00, // length
		0x49, 0x45, 0x4E, 0x44, // "IEND"
		0xAE, 0x42, 0x60, 0x82, // CRC
	}

	// Write test files
	if err := os.WriteFile(inputPath, samplePNG, 0644); err != nil {
		t.Fatalf("Failed to create test input file: %v", err)
	}
	if err := os.WriteFile(coverPath, samplePNG, 0644); err != nil {
		t.Fatalf("Failed to create test cover file: %v", err)
	}

	tests := []struct {
		name    string
		testFn  func() error
		wantErr bool
	}{
		{
			name: "Basic encryption and decryption",
			testFn: func() error {
				// First encrypt
				if err := EncryptImage(inputPath, encryptedPath, coverPath, testKey); err != nil {
					return err
				}
				// Then decrypt
				return DecryptImage(encryptedPath, decryptedPath, testKey)
			},
			wantErr: false,
		},
		{
			name: "Decrypt non-existent file",
			testFn: func() error {
				return DecryptImage("nonexistent.png", decryptedPath, testKey)
			},
			wantErr: true,
		},
		{
			name: "Decrypt file without stEG chunk",
			testFn: func() error {
				return DecryptImage(inputPath, decryptedPath, testKey)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.testFn()
			if (err != nil) != tt.wantErr {
				t.Errorf("Test %s: got error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}

	// Verify the decrypted content matches the original for successful cases
	if _, err := os.Stat(decryptedPath); err == nil {
		original, err := os.ReadFile(inputPath)
		if err != nil {
			t.Fatalf("Failed to read original file: %v", err)
		}
		decrypted, err := os.ReadFile(decryptedPath)
		if err != nil {
			t.Fatalf("Failed to read decrypted file: %v", err)
		}
		if !bytes.Equal(original, decrypted) {
			t.Error("Decrypted content does not match original")
		}
	}
}

func TestEncryptImageToBytes(t *testing.T) {
	tempDir := t.TempDir()
	testKey := []byte("test-key-12345")
	inputPath := filepath.Join(tempDir, "input.png")

	// Create test PNG file
	samplePNG := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D,
		0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00,
		0x1F, 0x15, 0xC4, 0x89,
		0x00, 0x00, 0x00, 0x0A,
		0x49, 0x44, 0x41, 0x54,
		0x08, 0xD7, 0x63, 0x60, 0x60, 0x60,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x49, 0x45, 0x4E, 0x44,
		0xAE, 0x42, 0x60, 0x82,
	}

	if err := os.WriteFile(inputPath, samplePNG, 0644); err != nil {
		t.Fatalf("Failed to create test input file: %v", err)
	}

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "Valid PNG file",
			path:    inputPath,
			wantErr: false,
		},
		{
			name:    "Non-existent file",
			path:    "nonexistent.png",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := EncryptImageToBytes(tt.path, testKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptImageToBytes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
} 