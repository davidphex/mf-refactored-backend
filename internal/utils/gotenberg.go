package utils

import (
	"fmt"
	"os"
	"os/exec"
)

func GenerateGotenbergPDF(htmlFilePath, cssFilePath string) ([]byte, error) {
	// Run curl to generate PDF from index.html and style.css
	cmd := exec.Command("curl",
		"--request", "POST",
		"https://demo.gotenberg.dev/forms/chromium/convert/html",
		"--form", fmt.Sprintf("files=@%s", htmlFilePath),
		"--form", fmt.Sprintf("files=@%s", cssFilePath),
		//"--form", "paperWidth=512",
		//"--form", "paperHeight=512",
		"--form", "preferCssPageSize=true",
		"-o", "/tmp/generated-album.pdf",
	)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run curl: %v", err)
	}

	// Read the PDF file into a byte slice
	pdfBytes, err := os.ReadFile("/tmp/generated-album.pdf")
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF file: %v", err)
	}

	return pdfBytes, nil
}
