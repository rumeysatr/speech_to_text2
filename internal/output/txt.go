package output

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"spt2/pkg/models"
)

func ExportTXT(result *models.TranscriptionResult, audioFilePath string, outputDir string) (string, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil{
		return "", fmt.Errorf("output dizini oluşturulamadı: %w", err)
	}

	var txtContent strings.Builder

	//metadata bilgileri
	txtContent.WriteString(fmt.Sprintf("Ses Dosyası: %s\n", filepath.Base(audioFilePath)))
	txtContent.WriteString(fmt.Sprintf("Tarih: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	txtContent.WriteString(fmt.Sprintf("Dil: %s\n", result.LanguageCode))
	txtContent.WriteString(fmt.Sprintf("Toplam Kelime: %d\n\n", len(result.Words)))

	txtContent.WriteString("--- TAM DEŞİFRE METNİ ---\n\n")
	txtContent.WriteString(result.Transcript)
	txtContent.WriteString("\n\n")

	if len(result.Words) > 0 {
		var totalConfidence float64
		for _, word := range result.Words {
			totalConfidence += word.Confidence
		}

		avgConfidence := (totalConfidence / float64(len(result.Words))) * 100

		txtContent.WriteString("--- İSTATİSTİKLER ---\n\n")
		txtContent.WriteString(fmt.Sprintf("Ortalama Güven: %.1f\n", avgConfidence))
		txtContent.WriteString(fmt.Sprintf("Toplam Kelime Sayısı: %d\n", len(result.Words)))
	}

	// Çıktı dosya adı oluştur (her zaman, Words olsun olmasın)
	audioFileName := filepath.Base(audioFilePath)
	audioFileNameWithoutExt := audioFileName[:len(audioFileName)-len(filepath.Ext(audioFileName))]
	outputFileName := audioFileNameWithoutExt + ".txt"
	outputPath := filepath.Join(outputDir, outputFileName)

	// Dosyaya yaz
	if err := os.WriteFile(outputPath, []byte(txtContent.String()), 0644); err != nil {
		return "", fmt.Errorf("TXT dosyası yazılamadı: %w", err)
	}

	return outputPath, nil
}
