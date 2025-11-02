package audio

import (
	"fmt" //hata mesajları için
	"spt2/pkg/models"
)

//ses dosya formatları listesi
var supportedFormats = map[string]bool{
	"mp3" : true,
	"wav" : true,
	"flac": true,
	"m4a" : true,
	"ogg" : true,
	"opus": true,
	"aac" : true,
}

func isFormatSupported(format string) bool{
	return supportedFormats[format]
}

//ses dosya metadata geçerliliği için
func ValidateMetadata(metadata *models.AudioMetadata) error{
	if !isFormatSupported(metadata.OriginalFormat) {
		return fmt.Errorf("desteklenmeyen ses formatı: %s \n (desteklenenler: mp3, wav, flac, m4a, ogg, opus, aac)", metadata.OriginalFormat)
	}

	const maxFileSize = 500 * 1024 * 1024
	if metadata.FileSize <= 0 {
		return fmt.Errorf("geçersiz dosya boyutu: %d bytes", metadata.FileSize)
	}
	if metadata.FileSize > maxFileSize {
		return fmt.Errorf("dosya çok büyük: %d bytes (maksimum: %d bytes / 500MB)", metadata.FileSize, maxFileSize)
	}

	if metadata.FilePath == "" {
		return fmt.Errorf("dosya yolu boş")
	}

	return nil
}
