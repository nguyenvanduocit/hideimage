# Image Encryption and Steganography Tool

A Go-based tool for encrypting PNG images and hiding them within cover images using steganography techniques.

## Features

- Image encryption using XOR cipher
- Steganography support to hide encrypted images in cover images
- Automatic cover image selection from a pool
- Image dimension and size comparison
- Support for PNG format

## Prerequisites

- Go 1.23.2 or higher
- PNG images for input and cover images

## Installation

Clone the repository and install dependencies:

```bash
git clone https://github.com/nguyenvanduocit/hideimage
cd hideimage
go mod download
```

## Project Structure

```
.
├── main.go         # Main program entry point
├── fns.go         # Core encryption/decryption functions
├── assets/
│   └── covers/    # Directory containing cover images
├── input.png      # Your input image
└── README.md
```

## Usage

### Encrypting and Hiding an Image

```go
key := []byte("YourSecretKey")
err := encryptImage("./input.png", "./encrypted.png", "./cover.png", key)
```

### Decrypting a Hidden Image

```go
key := []byte("YourSecretKey")
err := decryptImage("./encrypted.png", "./decrypted.png", key)
```

### Running the Example

1. Place your input image as `input.png` in the project root
2. Add cover images to `./assets/covers/`
3. Run the program:

```bash
go run .
```

## How It Works

1. **Encryption**: The tool encrypts the IDAT chunk of PNG images using XOR cipher with the provided key.
2. **Steganography**: The encrypted image is embedded within a cover image using a custom PNG chunk named "stEG".
3. **Decryption**: The tool can extract and decrypt the hidden image using the same key.

## Technical Details

The encryption process targets the IDAT chunk of PNG images, which contains the actual image data. The program:

1. Locates the IDAT chunk in the PNG file
2. Applies XOR encryption with the provided key
3. Embeds the encrypted data in a custom chunk within the cover image
4. Preserves PNG file structure and integrity

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.