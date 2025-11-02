package models

import "time"
// ses dosyası bilgileri
type AudioMetadata struct {
    //dosya bilgileri
	FilePath      		string 		 `json:"file_path"`
    OriginalFormat 		string		 `json:"original_format"`  // "mp3", "wav", "m4a"
    FileSize      		int64  		 `json:"file_size"`        // Byte cinsinden
    ConvertedPath 		string 		 `json:"converted_path"`   // FLAC dosya yolu
  
	//ses özellikleri
	Duration      		float64 	`json:"duration"`         // Saniye cinsinden
    SampleRate    		int     	`json:"sample_rate"`      // Hz (örn: 16000)
    Channels      		int     	`json:"channels"`         // 1=mono, 2=stereo
    BitDepth      		int     	`json:"bit_depth"`        // 16, 24, 32 bit
	Codec				string		`json:"codec"`			  //"aac", "mp3", "pcm"

	//hata ve durum yönetimleri
	IsValid				bool		`json:"is_valid"`
	ConversionStatus 	string		`json:"conversion_status"` //"completed" or "failure"
	ValidationError		string  	`json:"validation_error,omitempty"`
	CreatedAt			time.Time 	`json:"created_at"`

}
