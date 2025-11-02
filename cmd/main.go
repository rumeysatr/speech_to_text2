package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"spt2/internal/audio"
	"spt2/internal/config"
	"spt2/internal/output"
	"spt2/internal/speechclient"
	"spt2/internal/storage"
)

func main() {
	configPath := flag.String("config", "configs/default.json", "Path to the configuration file.")
	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatal("Kullanım: go run cmd/main.go [options] <audio_file_path>")
	}
	audioFilePath := flag.Arg(0)

	fmt.Println("=== Google Cloud Speech-to-Text Deşifre Sistemi ===\n")
	fmt.Printf("Ses Dosyası: %s\n\n", audioFilePath)

	ctx := context.Background()

	fmt.Println("📄 Config yükleniyor...")
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Config yüklenemedi: %v", err)
	}
	fmt.Printf("✅ Config yüklendi (Dil: %s, Model: %s)\n\n", cfg.LanguageCode, cfg.Model)

	fmt.Println("🎵 Ses dosyası metadata'sı çıkarılıyor...")
	metadata, err := audio.ExtractMetadata(audioFilePath)
	if err != nil {
		log.Fatalf("Metadata çıkarılamadı: %v", err)
	}
	fmt.Printf("✅ Metadata çıkarıldı (Format: %s, Boyut: %d bytes)\n\n", metadata.OriginalFormat, metadata.FileSize)

	//validate etme
	fmt.Println("✔️  Ses dosyası validate ediliyor...")
	if err := audio.ValidateMetadata(metadata); err != nil {
		log.Fatalf("Validasyon hatası: %v", err)
	}
	fmt.Println("✅ Validasyon başarılı\n")

	//flac
	fmt.Println("🔄 Ses dosyası FLAC formatına dönüştürülüyor...")
	if err := audio.ConvertToFLAC(metadata, cfg.OutputDir); err != nil {
		log.Fatalf("FLAC dönüştürme hatası: %v", err)
	}
	fmt.Printf("✅ FLAC'e dönüştürüldü: %s\n\n", metadata.ConvertedPath)

	// GCS'ye yükle
	fmt.Println("☁️  FLAC dosyası Google Cloud Storage'a yükleniyor...")
	gcsURI, err := storage.UploadToGCS(ctx, metadata.ConvertedPath, cfg.GCSBucket, cfg.GoogleCredentialsPath)
	if err != nil {
		log.Fatalf("GCS'ye yükleme hatası: %v", err)
	}
	fmt.Printf("✅ Dosya GCS'ye yüklendi: %s\n\n", gcsURI)

	//recognitionConfig
	fmt.Println("⚙️  Google API konfigürasyonu oluşturuluyor...")
	recognitionConfig := speechclient.BuildRecognitionConfig(cfg)
	fmt.Printf("✅ RecognitionConfig hazır (Dil: %s, Sample Rate: %d Hz)\n\n", recognitionConfig.LanguageCode, recognitionConfig.SampleRateHertz)

	//speech client başlatma
	fmt.Println("🔌 Google Speech API'a bağlanılıyor...")
	client, err := speechclient.NewSpeechClient(ctx, cfg)
	if err != nil {
		log.Fatalf("Speech client başlatılamadı: %v", err)
	}
	defer client.Close()
	fmt.Println("✅ Google Speech API bağlantısı kuruldu\n")

	//long running recognize
	fmt.Println("🎤 Ses dosyası deşifre ediliyor (bu birkaç dakika sürebilir)...")
	result, err := client.LongRunningRecognize(ctx, gcsURI, recognitionConfig)
	if err != nil {
		log.Fatalf("Deşifre hatası: %v", err)
	}
	fmt.Printf("✅ Deşifre tamamlandı (%d karakter)\n\n", len(result.Transcript))

	//json output
	fmt.Println("💾 Sonuçlar JSON dosyasına kaydediliyor...")
	jsonPath, err := output.ExportJSON(result, audioFilePath, cfg.OutputDir)
	if err != nil {
		log.Fatalf("JSON kaydetme hatası: %v", err)
	}
	fmt.Printf("✅ JSON dosyası oluşturuldu: %s\n", jsonPath)

	// SRT altyazı export
	fmt.Println("\n📝 SRT altyazı dosyası oluşturuluyor...")
	srtPath, err := output.ExportSRT(result, audioFilePath, cfg.OutputDir)
	if err != nil {
		log.Fatalf("SRT kaydetme hatası: %v", err)
	}
	fmt.Printf("✅ SRT dosyası oluşturuldu: %s\n", srtPath)

	// TXT rapor export
	fmt.Println("\n📄 TXT rapor dosyası oluşturuluyor...")
	txtPath, err := output.ExportTXT(result, audioFilePath, cfg.OutputDir)
	if err != nil {
		log.Fatalf("TXT kaydetme hatası: %v", err)
	}
	fmt.Printf("✅ TXT dosyası oluşturuldu: %s\n", txtPath)

	fmt.Println("\n✅ İşlem tamamlandı!")
}
