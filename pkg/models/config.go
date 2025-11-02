package models

// AppConfig - Uygulama konfigürasyonu (Viper ve Validator ile çalışır)
type AppConfig struct {
    // Google Cloud Ayarları
    GoogleCredentialsPath string `mapstructure:"google_credentials_path" validate:"required,file"`
    ProjectID             string `mapstructure:"project_id" validate:"required"`
    GCSBucket             string `mapstructure:"gcs_bucket" validate:"required"`
    
    // API Temel Ayarları
    LanguageCode string `mapstructure:"language_code" validate:"required,oneof=en-US en-GB tr-TR de-DE fr-FR es-ES"`
    Model        string `mapstructure:"model" validate:"required,oneof=default video telephony medical command_and_search"`
    UseEnhanced  bool   `mapstructure:"use_enhanced"`
    
    // Diarization (Konuşmacı Ayırma)
    EnableDiarization bool `mapstructure:"enable_diarization"`
    MinSpeakers       int  `mapstructure:"min_speakers" validate:"omitempty,min=1,max=6"`
    MaxSpeakers       int  `mapstructure:"max_speakers" validate:"omitempty,min=1,max=6,gtefield=MinSpeakers"`
    
    // Speech Contexts (Özel Kelimeler) - Dosya Yolları
    SpeechContextsFile string  `mapstructure:"speech_contexts_file" validate:"omitempty"`
    KeywordsFile       string  `mapstructure:"keywords_file" validate:"omitempty"`
    BoostValue         float64 `mapstructure:"boost_value" validate:"omitempty,min=0,max=20"`
    
    // Runtime'da TXT dosyalarından yüklenir (JSON'da yok)
    SpeechContexts []string `mapstructure:"-" json:"-"`
    Keywords       []string `mapstructure:"-" json:"-"`
    
    // KRİTİK: Yapısal Veri Elde Etme (Prompt.md'de zorunlu!)
    EnableAutomaticPunctuation bool `mapstructure:"enable_automatic_punctuation" validate:"required,eq=true"`
    EnableWordTimeOffsets      bool `mapstructure:"enable_word_time_offsets" validate:"required,eq=true"`
    EnableWordConfidence       bool `mapstructure:"enable_word_confidence"`
    
    // API Ekstra Ayarlar
    MinConfidence   float64 `mapstructure:"min_confidence" validate:"omitempty,min=0,max=1"`
    MaxAlternatives int     `mapstructure:"max_alternatives" validate:"omitempty,min=1,max=30"`
    ProfanityFilter bool    `mapstructure:"profanity_filter"`
    
    // Ses İşleme Ayarları
    TargetSampleRate int  `mapstructure:"target_sample_rate" validate:"required,min=8000,max=48000"`
    ConvertToMono    bool `mapstructure:"convert_to_mono"`
    ChunkSize        int  `mapstructure:"chunk_size" validate:"required,min=1024,max=65536"`
    
    // Çıktı Ayarları
    OutputDir    string `mapstructure:"output_dir" validate:"required"`
    GenerateJSON bool   `mapstructure:"generate_json"`
    GenerateSRT  bool   `mapstructure:"generate_srt"`
    GenerateTXT  bool   `mapstructure:"generate_txt"`
    
    // Logging
    EnableLogging bool   `mapstructure:"enable_logging"`
    LogLevel      string `mapstructure:"log_level" validate:"omitempty,oneof=debug info warn error"`
}
