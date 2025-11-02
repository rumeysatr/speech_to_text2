package main

import (
	"fmt"
	"log"
	"speech_to_text2/internal/config"
	"speech_to_text2/pkg/models"
)

func main() {
	fmt.Println("=== Config Loader Test Programı ===\n")

	// Test 1: İngilizce Config
	fmt.Println("📄 Test 1: İngilizce config (default.json) yükleniyor...")
	cfgEN, err := config.LoadConfig("./configs/default.json")
	if err != nil {
		log.Fatalf("❌ İngilizce config yüklenemedi: %v", err)
	}
	fmt.Println("✅ İngilizce config başarıyla yüklendi!")
	printConfigSummary(cfgEN, "EN-US")

	// Test 2: Türkçe Config
	fmt.Println("\n📄 Test 2: Türkçe config (config-tr.json) yükleniyor...")
	cfgTR, err := config.LoadConfig("./configs/config-tr.json")
	if err != nil {
		log.Fatalf("❌ Türkçe config yüklenemedi: %v", err)
	}
	fmt.Println("✅ Türkçe config başarıyla yüklendi!")
	printConfigSummary(cfgTR, "TR-TR")

	fmt.Println("\n🎉 Tüm testler başarıyla tamamlandı!")
}

func printConfigSummary(cfg *models.AppConfig, label string) {
	fmt.Printf("\n--- %s Config Özeti ---\n", label)
	fmt.Printf("Dil: %s\n", cfg.LanguageCode)
	fmt.Printf("Model: %s\n", cfg.Model)
	fmt.Printf("Sample Rate: %d Hz\n", cfg.TargetSampleRate)
	fmt.Printf("Chunk Size: %d bytes\n", cfg.ChunkSize)
	fmt.Printf("Output Dir: %s\n", cfg.OutputDir)
	fmt.Printf("Noktalama: %v\n", cfg.EnableAutomaticPunctuation)
	fmt.Printf("Kelime Zaman Damgası: %v\n", cfg.EnableWordTimeOffsets)
	
	fmt.Printf("\nSpeech Contexts (%d adet):\n", len(cfg.SpeechContexts))
	for i, ctx := range cfg.SpeechContexts {
		if i < 3 { // İlk 3'ünü göster
			fmt.Printf("  - %s\n", ctx)
		}
	}
	if len(cfg.SpeechContexts) > 3 {
		fmt.Printf("  ... ve %d adet daha\n", len(cfg.SpeechContexts)-3)
	}
	
	fmt.Printf("\nKeywords (%d adet):\n", len(cfg.Keywords))
	for i, kw := range cfg.Keywords {
		if i < 3 { // İlk 3'ünü göster
			fmt.Printf("  - %s\n", kw)
		}
	}
	if len(cfg.Keywords) > 3 {
		fmt.Printf("  ... ve %d adet daha\n", len(cfg.Keywords)-3)
	}
}
