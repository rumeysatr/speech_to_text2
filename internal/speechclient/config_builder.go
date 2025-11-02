package speechclient

import (
	"cloud.google.com/go/speech/apiv1/speechpb"
	"speech_to_text2/pkg/models"
)

func BuildRecognitionConfig(cfg *models.AppConfig) *speechpb.RecognitionConfig {
	recognitionConfig := &speechpb.RecognitionConfig{
		LanguageCode: cfg.LanguageCode,
		Encoding:	  speechpb.RecognitionConfig_LINEAR16,
		SampleRateHertz: int32(cfg.TargetSampleRate),
		EnableAutomaticPunctuation: cfg.EnableAutomaticPunctuation,
		EnableWordTimeOffsets:		cfg.EnableWordTimeOffsets,
		Model: 			cfg.Model,
		UseEnhanced:	cfg.UseEnhanced,

	}

	if len(cfg.SpeechContexts) > 0 {
		recognitionConfig.SpeechContexts = []*speechpb.SpeechContext{
			{
				Phrases: cfg.SpeechContexts,
				Boost:	 float32(cfg.BoostValue),
			},
		}
	}

	if cfg.EnableDiarization {
		recognitionConfig.DiarizationConfig = &speechpb.SpeakerDiarizationConfig{
			EnableSpeakerDiarization: true,
			MinSpeakerCount:		  int32(cfg.MinSpeakers),
			MaxSpeakerCount:		  int32(cfg.MaxSpeakers),
		}
	}
	return recognitionConfig
}

