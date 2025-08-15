package wallet

import (
	"github.com/divan/qrlogo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skip2/go-qrcode"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func CreateQR(pubkey string) (nameOfFile string, err error) {
	filePath := "./qr_codes/" + pubkey + ".png"
	err = qrcode.WriteFile(pubkey, qrcode.Medium, 256, filePath)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

func CreateQRWithLogo(str string) (filePath string, err error) {
	// Open the logo.png image from the current directory
	logoFile, err := os.Open("logo.png")
	if err != nil {
		return "", err
	}
	defer logoFile.Close()

	logo, err := png.Decode(logoFile)
	if err != nil {
		return "", err
	}

	// Assuming size is defined or passed somewhere else in your code
	size := 256

	encode, err := qrlogo.Encode(str, logo, size)
	if err != nil {
		return "", err
	}

	// Define the filePath
	filePath = "./qr_codes/" + str + ".png"

	// Log the filePath for debugging
	log.Println("Writing QR to:", filePath)

	// Check if the qr_codes directory exists; if not, create it
	if _, err := os.Stat("./qr_codes/"); os.IsNotExist(err) {
		os.Mkdir("./qr_codes/", 0755)
	}

	// Write buffer contents to the file
	err = ioutil.WriteFile(filePath, encode.Bytes(), 0644)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

type FileRequestData struct {
	FilePath string
}

func (f *FileRequestData) NeedsUpload() bool {
	return f.FilePath != ""
}

func (f *FileRequestData) UploadData() (string, io.Reader, error) {
	file, err := os.Open(f.FilePath)
	if err != nil {
		return "", nil, err
	}
	return filepath.Base(f.FilePath), file, nil
}

func (f *FileRequestData) SendData() string {
	// Return the attachment format for sending
	return "attach://" + filepath.Base(f.FilePath)
}

func NewInputMediaPhotoFromFile(qrFileName string) tgbotapi.InputMediaPhoto {
	attachmentName := filepath.Base(qrFileName) // Extract just the filename without directory paths
	fileData := &FileRequestData{
		FilePath: "qr_codes/" + attachmentName, // Storing the correct relative path
	}

	// Constructing the BaseInputMedia for InputMediaPhoto
	media := tgbotapi.BaseInputMedia{
		Type:  "photo",
		Media: fileData,
	}

	return tgbotapi.InputMediaPhoto{
		BaseInputMedia: media,
	}
}
