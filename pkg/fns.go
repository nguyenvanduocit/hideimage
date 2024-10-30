package pkg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

// decryptImage giải mã hình ảnh, có thể từ một file đã giấu
func DecryptImage(inputPath, outputPath string, key []byte) error {
	// Đọc file đầu vào
	input, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	// Tìm chunk "stEG"
	stegIndex := bytes.Index(input, []byte("stEG"))
	if stegIndex == -1 {
		return fmt.Errorf("không tìm thấy chunk stEG")
	}

	// Trích xuất dữ liệu đã mã hóa
	dataLength := binary.BigEndian.Uint32(input[stegIndex+4 : stegIndex+8])
	encryptedData := input[stegIndex+8 : stegIndex+8+int(dataLength)]
	
	// Giải mã dữ liệu đã trích xuất
	return DecryptImageFromBytes(encryptedData, outputPath, key)
}

// decryptImageFromBytes giải mã dữ liệu hình ảnh từ []byte
func DecryptImageFromBytes(data []byte, outputPath string, key []byte) error {
	// Tìm vị trí bắt đầu của chunk IDAT
	idatStart := bytes.Index(data, []byte("IDAT"))
	if idatStart == -1 {
		return fmt.Errorf("không tìm thấy chunk IDAT")
	}

	// Đọc độ dài của chunk IDAT
	idatLength := binary.BigEndian.Uint32(data[idatStart-4 : idatStart])

	// Tính vị trí kết thúc của chunk IDAT
	idatEnd := idatStart + 4 + int(idatLength) + 4

	// Tạo slice chứa dữ liệu đã giải mã
	decrypted := make([]byte, len(data))
	copy(decrypted, data)

	// Giải mã phần dữ liệu IDAT
	for i := idatStart + 4; i < idatEnd - 4; i++ {
		decrypted[i] = data[i] ^ key[(i-(idatStart+4))%len(key)]
	}

	// Ghi dữ liệu đã giải mã vào file đầu ra
	return os.WriteFile(outputPath, decrypted, 0644)
}

// encryptImage mã hóa hình ảnh và tùy chọn giấu nó vào một hình ảnh khác
func EncryptImage(inputPath, outputPath string, coverPath string, key []byte) error {
	// Mã hóa hình ảnh đầu vào
	encryptedData, err := EncryptImageToBytes(inputPath, key)
	if err != nil {
		return err
	}

	// Đọc hình ảnh cover
	coverImage, err := os.ReadFile(coverPath)
	if err != nil {
		return err
	}

	// Tìm vị trí cuối cùng của file PNG (IEND chunk)
	iendIndex := bytes.LastIndex(coverImage, []byte("IEND"))
	if iendIndex == -1 {
		return fmt.Errorf("không tìm thấy IEND chunk trong hình ảnh cover")
	}

	// Tạo chunk mới để chứa dữ liệu đã mã hóa
	newChunk := make([]byte, 8+len(encryptedData))
	copy(newChunk[:4], []byte("stEG")) // Tên chunk tùy chỉnh
	binary.BigEndian.PutUint32(newChunk[4:8], uint32(len(encryptedData)))
	copy(newChunk[8:], encryptedData)

	// Tạo slice mới để chứa hình ảnh kết quả
	result := make([]byte, len(coverImage)+len(newChunk))
	copy(result[:iendIndex], coverImage[:iendIndex])
	copy(result[iendIndex:], newChunk)
	copy(result[iendIndex+len(newChunk):], coverImage[iendIndex:])

	// Ghi kết quả vào file đầu ra
	return os.WriteFile(outputPath, result, 0644)
}

// encryptImageToBytes mã hóa hình ảnh và trả về dữ liệu dưới dạng []byte
func EncryptImageToBytes(inputPath string, key []byte) ([]byte, error) {
	input, err := os.ReadFile(inputPath)
	if err != nil {
		return nil, err
	}

	// Tìm vị trí bắt đầu của chunk IDAT
	idatStart := bytes.Index(input, []byte("IDAT"))
	if idatStart == -1 {
		return nil, fmt.Errorf("không tìm thấy chunk IDAT")
	}

	// Đọc độ dài của chunk IDAT
	idatLength := binary.BigEndian.Uint32(input[idatStart-4 : idatStart])

	// Tính vị trí kết thúc của chunk IDAT
	idatEnd := idatStart + 4 + int(idatLength) + 4

	// Tạo slice chứa dữ liệu đã xử lý
	processed := make([]byte, len(input))
	copy(processed, input)

	// Mã hóa phần dữ liệu IDAT
	for i := idatStart + 4; i < idatEnd - 4; i++ {
		processed[i] = input[i] ^ key[(i-(idatStart+4))%len(key)]
	}

	return processed, nil
}
