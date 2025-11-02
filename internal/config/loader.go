package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"speech_to_text2/pkg/models"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// LoadConfig - Viper ve Validator kullanarak modern, esnek config yönetimi
//
// ÖZELLİKLER:
// - JSON config dosyasından okuma
// - Environment variable desteği (APP_ prefix ile)
// - TXT dosyalarından speech contexts ve keywords yükleme
// - Otomatik validation (required, range, enum kontrolü)
// - Custom validator (dosya varlığı kontrolü)
// - Default değer desteği
//
// KULLANIM:
//   cfg, err := LoadConfig("./configs/default.json")
//   if err != nil { ... }
func LoadConfig(path string) (*models.AppConfig, error) {
	// --- BÖLÜM 1: DEFAULT DEĞERLER (Viper ile) ---
	
	// Zorunlu olmayan field'lar için varsayılan değerler
	viper.SetDefault("use_enhanced", false)
	viper.SetDefault("enable_diarization", false)
	viper.SetDefault("min_speakers", 1)
	viper.SetDefault("max_speakers", 6)
	viper.SetDefault("boost_value", 10.0)
	viper.SetDefault("min_confidence", 0.7)
	viper.SetDefault("max_alternatives", 1)
	viper.SetDefault("profanity_filter", false)
	viper.SetDefault("target_sample_rate", 16000)
	viper.SetDefault("convert_to_mono", true)
	viper.SetDefault("chunk_size", 4096)
	viper.SetDefault("output_dir", "./output")
	viper.SetDefault("generate_json", true)
	viper.SetDefault("generate_srt", true)
	viper.SetDefault("generate_txt", true)
	viper.SetDefault("enable_logging", true)
	viper.SetDefault("log_level", "info")

	// --- BÖLÜM 2: VIPER İLE CONFIG DOSYASI YÜKLEME ---

	// Config dosyasının tam yolunu belirt
	viper.SetConfigFile(path)

	// Environment variable desteği
	// Örnek: APP_LANGUAGE_CODE=tr-TR ortam değişkeni config'i geçersiz kılar
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Config dosyasını oku
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config dosyası okunamadı '%s': %w", path, err)
	}

	// Viper'dan AppConfig struct'ına unmarshal et
	var cfg models.AppConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config parse edilemedi: %w", err)
	}

	// --- BÖLÜM 3: TXT DOSYALARINDAN KELİMELERİ YÜKLE ---

	// Speech contexts dosyasını yükle (varsa)
	if cfg.SpeechContextsFile != "" {
		contexts, err := loadTextFile(cfg.SpeechContextsFile)
		if err != nil {
			return nil, fmt.Errorf("speech contexts dosyası yüklenemedi: %w", err)
		}
		cfg.SpeechContexts = contexts
	}

	// Keywords dosyasını yükle (varsa)
	if cfg.KeywordsFile != "" {
		keywords, err := loadTextFile(cfg.KeywordsFile)
		if err != nil {
			return nil, fmt.Errorf("keywords dosyası yüklenemedi: %w", err)
		}
		cfg.Keywords = keywords
	}

	// --- BÖLÜM 4: VALIDATOR İLE DOĞRULAMA ---

	// Custom validator oluştur (dosya varlığı kontrolü için)
	validate := validator.New()
	
	// Custom validation: dosya varlığı kontrolü
	validate.RegisterValidation("file", validateFileExists)

	// Struct validation
	if err := validate.Struct(&cfg); err != nil {
		// Validation hatalarını kullanıcı dostu formata çevir
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return nil, formatValidationErrors(validationErrors)
		}
		return nil, fmt.Errorf("validation hatası: %w", err)
	}

	// --- BÖLÜM 5: OUTPUT DİZİNİ OLUŞTUR ---
	
	// Output dizini yoksa oluştur
	if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
		return nil, fmt.Errorf("output dizini oluşturulamadı: %w", err)
	}

	return &cfg, nil
}

// loadTextFile - TXT dosyasından satır satır kelimeleri okur
//
// NEDEN: Speech contexts ve keywords'leri ayrı dosyalarda tutmak
// projeyi daha esnek yapar. Kullanıcı kod değiştirmeden kelimeleri güncelleyebilir.
//
// ÇALIŞMA MANTIĞI:
// 1. Dosyayı aç
// 2. Satır satır oku (bufio.Scanner ile)
// 3. Her satırı trim et (boşlukları temizle)
// 4. Boş satırları atla
// 5. String slice döndür
func loadTextFile(filePath string) ([]string, error) {
	// Dosyayı aç
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("dosya açılamadı: %w", err)
	}
	defer file.Close()

	// Scanner ile satır satır oku
	scanner := bufio.NewScanner(file)
	var lines []string

	// Her satırı işle
	for scanner.Scan() {
		// Satırı al ve boşlukları temizle
		line := strings.TrimSpace(scanner.Text())
		
		// Boş satırları ve yorum satırlarını atla
		if line != "" && !strings.HasPrefix(line, "#") {
			lines = append(lines, line)
		}
	}

	// Scanner hatası kontrolü
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("dosya okuma hatası: %w", err)
	}

	return lines, nil
}

// validateFileExists - Custom validator: dosya varlığı kontrolü
//
// KULLANIM: `validate:"file"` tag'i ile
// Örnek: GoogleCredentialsPath string `validate:"required,file"`
func validateFileExists(fl validator.FieldLevel) bool {
	filePath := fl.Field().String()
	
	// Boş path geçerli sayılır (required tag'i bunu kontrol eder)
	if filePath == "" {
		return true
	}

	// Dosya var mı kontrol et
	_, err := os.Stat(filePath)
	return err == nil
}

// formatValidationErrors - Validation hatalarını okunabilir formata çevir
//
// ÖRNEK ÇIKTI:
// "Config validation hatası:
//  - GoogleCredentialsPath: dosya bulunamadı
//  - LanguageCode: 'xyz' geçersiz, geçerli değerler: en-US, tr-TR
//  - TargetSampleRate: 0 geçersiz, minimum: 8000"
func formatValidationErrors(errs validator.ValidationErrors) error {
	var messages []string
	
	for _, err := range errs {
		field := err.Field()
		tag := err.Tag()
		
		var msg string
		switch tag {
		case "required":
			msg = fmt.Sprintf("%s: zorunlu field", field)
		case "file":
			msg = fmt.Sprintf("%s: dosya bulunamadı '%s'", field, err.Value())
		case "oneof":
			msg = fmt.Sprintf("%s: '%v' geçersiz, geçerli değerler: %s", field, err.Value(), err.Param())
		case "min":
			msg = fmt.Sprintf("%s: '%v' çok küçük, minimum: %s", field, err.Value(), err.Param())
		case "max":
			msg = fmt.Sprintf("%s: '%v' çok büyük, maksimum: %s", field, err.Value(), err.Param())
		case "eq":
			msg = fmt.Sprintf("%s: '%v' olmalı, '%s' değeri bekleniyor", field, err.Value(), err.Param())
		case "gtefield":
			msg = fmt.Sprintf("%s: %s field'ından büyük veya eşit olmalı", field, err.Param())
		default:
			msg = fmt.Sprintf("%s: validation hatası (%s)", field, tag)
		}
		
		messages = append(messages, " - "+msg)
	}
	
	return fmt.Errorf("config validation hatası:\n%s", strings.Join(messages, "\n"))
}