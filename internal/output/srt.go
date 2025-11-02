package output

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"spt2/pkg/models"
)

func formatSRTTime(seconds float64) string {
	hours := int(seconds / 3600)
	minutes := int((seconds - float64(hours*3600)) / 60)
	secs := int(seconds - float64(hours*3600) - float64(minutes*60))
	millis := int((seconds - float64(int(seconds))) * 1000)

	return fmt.Sprintf("%02d:%02d:%02d,%03d", hours, minutes, secs, millis)
}

func groupWordsIntoSubtitles(words []models.WordInfo) [][]models.WordInfo {
	if len(words) == 0 {
		return nil
	}

	var groups [][]models.WordInfo
	var currentGroup []models.WordInfo
	maxWordsPerSubtitle := 10
	maxDuration := 3.0

	for i, word := range words {
		currentGroup = append(currentGroup, word)

		groupDuration := word.EndTime - currentGroup[0].StartTime

		if len(currentGroup) >= maxWordsPerSubtitle || groupDuration >= maxDuration || i == len(words)-1 {
			groups = append(groups, currentGroup)
			currentGroup = nil
		}
	}
	return groups
}

func ExportSRT(result *models.TranscriptionResult, audioFilePath string, outputDir string) (string, error) {
	if len(result.Words) == 0 {
		return "", fmt.Errorf("SRT oluşturmak için kelime zaman damgaları gerekli (Words boş)")
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("output dizini oluşturulamadı: %w", err)
	}

	subtitleGroups := groupWordsIntoSubtitles(result.Words)

	var srtContent strings.Builder

	for i, group := range subtitleGroups {
		if len(group) == 0 {
			continue
		}

		srtContent.WriteString(fmt.Sprintf("%d\n", i+1))

		startTime := formatSRTTime(group[0].StartTime)
		endTime   := formatSRTTime(group[len(group)-1].EndTime)
		srtContent.WriteString(fmt.Sprintf("%s --> %s\n", startTime, endTime))

		var words []string
		for _, word := range group {
			words = append(words, word.Word)
		}
		srtContent.WriteString(strings.Join(words, " "))
		srtContent.WriteString("\n\n")
	}

	audioFileName := filepath.Base(audioFilePath)
	audioFileNameWithoutExt := audioFileName[:len(audioFileName) - len(filepath.Ext(audioFileName))]
	outputFileName := audioFileNameWithoutExt + ".srt"
	outputPath 	   := filepath.Join(outputDir, outputFileName)

	//dosyaya yazdırma işlemi
	if err := os.WriteFile(outputPath, []byte(srtContent.String()), 0644); err != nil {
		return "", fmt.Errorf("SRT dosyası yazılamadı: %w", err)
	}
	return outputPath, nil
}