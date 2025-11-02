package output

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"spt2/pkg/models"
)

//json çıktısı için wrapper struct
type JSONOutput struct {
	TranscriptionResult *models.TranscriptionResult `json:"transcription_result"`
	Metadata 			OutputMetadata				`json:"metadata"` 
}

//çıktı dosyası hakkında metadata bilgileri
type OutputMetadata struct {
	GeneratedAt	 string `json:"generated_at"`
	AudioFile 	 string `json:"audio_file"`
	OutputFormat string `json:"output_format"`
	Version		 string `json:"version"`
}

//transcriptionresult'u JSON formatında dosyaya kaydetme
func ExportJSON(result *models.TranscriptionResult, audioFilePath string, outputDir string) (string, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("output dizini oluşturulamadı: %w", err)
	}

	metadata := OutputMetadata{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		AudioFile:	 filepath.Base(audioFilePath),
		OutputFormat: "json",
		Version: 	  "1.0.0",
	}

	jsonOutput := JSONOutput{
		TranscriptionResult: result,
		Metadata:			 metadata,
	}

	jsonData, err := json.MarshalIndent(jsonOutput, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON oluşturulamadı: %w", err)
	}

	audioFileName := filepath.Base(audioFilePath)
	audioFileNameWithoutExt := audioFileName[:len(audioFileName)-len(filepath.Ext(audioFileName))]
	outputFileName := audioFileNameWithoutExt + ".json"
	outputPath := filepath.Join(outputDir, outputFileName)

	if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
		return "", fmt.Errorf("JSON dosyası yazılamadı: %w", err)
	}
	return outputPath, nil
}