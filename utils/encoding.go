package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
)

func EncodeAndCompress(data []byte) string {
	// Compress the JSON string
	compressedData := CompressData(string((data)))

	// Encode the compressed data to Base64
	encodedData := encodeBase64(compressedData)

	return encodedData
}

func CompressData(data string) []byte {
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)

	_, err := writer.Write([]byte(data))
	if err != nil {
		fmt.Println("Error compressing data:", err)
		return nil
	}

	writer.Close()
	return buffer.Bytes()
}

func DecompressData(data []byte) (string, error) {
	reader, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	defer reader.Close()

	var buffer bytes.Buffer
	_, err = buffer.ReadFrom(reader)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func encodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func decodeBase64(encodedData string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encodedData)
}
