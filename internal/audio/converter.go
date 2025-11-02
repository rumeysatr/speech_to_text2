package audio

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"spt2/pkg/models"
)

//flac dosyası için çıktı yolu oluşturma
func generateOutputPath(inputPath string, outputDir string) string {
	fileName := filepath.Base(inputPath)

	fileNameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	flacFileName := fileNameWithoutExt + ".flac"
	outputPath := filepath.Join(outputDir, flacFileName)

	return outputPath
}

//ses dosyasını flac formatına dönüştürme
func ConvertToFLAC(metadata *models.AudioMetadata, outputDir string) error{
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("output dizini oluşturulamadı: %w", err)
	}

	outputPath := generateOutputPath(metadata.FilePath, outputDir)

	cmd := exec.Command("ffmpeg", "-i", metadata.FilePath, "-ar", "16000", "-ac", "1", "-y", outputPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("FFmpeg hatası: %w\nÇıktı: %s", err, string(output))
	}

	metadata.ConvertedPath = outputPath
	metadata.ConversionStatus = "completed"

	return nil
}
