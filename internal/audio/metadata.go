//ses dosya formatını tespit etme
//süre bilgisi çıkarma gibi işlemler için
package audio

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	
	"spt2/pkg/models"
)

//dosya uzantısından ses formatını tespit edecek
func detectFormat(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))

	if len(ext) > 0 && ext[0] == '.' {
		ext = ext[1:]
	}

	return ext
}

//ses dosyasından metadata bilgilerini çıkaracak fonk
func ExtractMetadata(filePath string) (*models.AudioMetadata, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("dosya bulunamadı: %w", err)
	}

	//metadata structı
	metadata := &models.AudioMetadata{
		FilePath:		  filePath,
		OriginalFormat:   detectFormat(filePath),
		FileSize:		  fileInfo.Size(),
		CreatedAt:		  time.Now(),
		IsValid:		  true,
		ConversionStatus: "pending",

		Duration:		  0.0,
		SampleRate:		  16000,
		Channels:		  1,
		BitDepth:		  16,
		Codec:			  "unknown",
	}
	return metadata, nil
}


