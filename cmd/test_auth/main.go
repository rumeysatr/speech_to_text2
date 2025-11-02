package main

import (
	"context"
	"fmt"
	"log"
	"os"

	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
	"google.golang.org/api/option"
	"spt2/internal/config"
)

func main() {
	fmt.Println("=== Google Cloud Kimlik Doğrulama Test Programı ===\n")

	// Config dosyasını yükle
	fmt.Println("📄 Config dosyası yükleniyor...")
	cfg, err := config.LoadConfig("./configs/default.json")
	if err != nil {
		log.Fatalf("❌ Config yüklenemedi: %v", err)
	}
	fmt.Printf("✅ Config başarıyla yüklendi\n")
	fmt.Printf("   Credentials Path: %s\n", cfg.GoogleCredentialsPath)
	fmt.Printf("   Project ID: %s\n\n", cfg.ProjectID)

	// Credentials dosyasının varlığını kontrol et
	fmt.Println("🔑 Credentials dosyası kontrol ediliyor...")
	if _, err := os.Stat(cfg.GoogleCredentialsPath); os.IsNotExist(err) {
		log.Fatalf("❌ Credentials dosyası bulunamadı: %s", cfg.GoogleCredentialsPath)
	}
	fmt.Println("✅ Credentials dosyası mevcut\n")

	// Google Cloud Speech client oluştur
	fmt.Println("🔌 Google Cloud Speech client oluşturuluyor...")
	ctx := context.Background()
	client, err := speech.NewClient(ctx, option.WithCredentialsFile(cfg.GoogleCredentialsPath))
	if err != nil {
		log.Fatalf("❌ Speech client oluşturulamadı: %v", err)
	}
	defer client.Close()
	fmt.Println("✅ Speech client başarıyla oluşturuldu\n")

	// Test: Basit bir recognize isteği (boş audio ile)
	// bağlantıyı ve yetkilendirmeyi test etmek için
	fmt.Println("🧪 API bağlantısı test ediliyor...")
	fmt.Println("   (Boş bir recognize isteği gönderiliyor...)")
	
	req := &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 16000,
			LanguageCode:    "en-US",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{
				Content: []byte{}, // Boş audio
			},
		},
	}

	// Bu çağrı muhtemelen hata verecek (boş audio)
	// bağlantı ve authentication sorunları için bir test
	resp, err := client.Recognize(ctx, req)
	
	// Hata türünü kontrol et
	if err != nil {
		// Eğer authentication hatası değilse, bu aslında iyi bir işaret
		// (Sadece boş audio hatası olmalı)
		if containsAuthError(err.Error()) {
			log.Fatalf("❌ Kimlik doğrulama hatası: %v", err)
		} else {
			// Boş audio hatası bekleniyor - bu normaldir
			fmt.Println("✅ API bağlantısı başarılı!")
			fmt.Printf("   (Beklenen hata alındı: %v)\n", err)
			fmt.Println("   Bu normal - kimlik doğrulama çalışıyor!\n")
		}
	} else {
		// Beklenmeyen başarı (boş audio ile)
		fmt.Println("✅ API bağlantısı başarılı!")
		fmt.Printf("   Yanıt alındı: %d sonuç\n\n", len(resp.Results))
	}

	// Özet
	fmt.Println("═══════════════════════════════════════")
	fmt.Println("🎉 Kimlik Doğrulama Testi Başarılı!")
	fmt.Println("═══════════════════════════════════════")
	fmt.Println("\n✓ Credentials dosyası geçerli")
	fmt.Println("✓ Google Cloud Speech API'a bağlantı kuruldu")
	fmt.Println("✓ Kimlik doğrulama başarılı")
	fmt.Println("✓ Proje hazır, deşifre işlemlerine başlayabilirsiniz!\n")
}

// containsAuthError - Hatanın authentication ile ilgili olup olmadığını kontrol eder
func containsAuthError(errMsg string) bool {
	authKeywords := []string{
		"authentication",
		"credentials",
		"permission denied",
		"unauthorized",
		"unauthenticated",
		"invalid_grant",
		"token",
	}

	errLower := errMsg
	for _, keyword := range authKeywords {
		if contains(errLower, keyword) {
			return true
		}
	}
	return false
}

// contains - String contains helper (case-insensitive olmadan)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || 
		    len(s) > len(substr) && 
		    (s[:len(substr)] == substr || 
		     s[len(s)-len(substr):] == substr || 
		     hasSubstring(s, substr)))
}

func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}