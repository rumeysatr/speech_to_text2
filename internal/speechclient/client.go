package speechclient

import (
	"context"
	"fmt"
	"strings"

	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
	"google.golang.org/api/option"

	"spt2/pkg/models"
)

//google speech to text api client wrapper
type SpeechClient struct {
	client *speech.Client
	config *models.AppConfig
}

//yeni speech client oluşturma
func NewSpeechClient(ctx context.Context, cfg *models.AppConfig) (*SpeechClient, error) {
	client, err := speech.NewClient(ctx, option.WithCredentialsFile(cfg.GoogleCredentialsPath))
	if err != nil {
		return nil, fmt.Errorf("Speech client oluşturulamadı: %w", err)
	}

	return &SpeechClient{
		client: client,
		config: cfg,
	}, nil
}

// client bağlantısını kapat
func (sc *SpeechClient) Close() error {
	return sc.client.Close()
}

// LongRunningRecognize sends a long audio file to Google Speech API for transcription.
func (sc *SpeechClient) LongRunningRecognize(ctx context.Context, gcsURI string, recognitionConfig *speechpb.RecognitionConfig) (*models.TranscriptionResult, error) {
	req := &speechpb.LongRunningRecognizeRequest{
		Config: recognitionConfig,
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Uri{Uri: gcsURI},
		},
	}

	op, err := sc.client.LongRunningRecognize(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("uzun süreli tanıma başlatılamadı: %w", err)
	}

	resp, err := op.Wait(ctx)
	if err != nil {
		return nil, fmt.Errorf("sonuç beklenirken hata oluştu: %w", err)
	}

	var transcriptBuilder strings.Builder
	var allWords []models.WordInfo

	for _, result := range resp.Results {
		alternative := result.Alternatives[0]
		transcriptBuilder.WriteString(alternative.Transcript + " ")

		for _, wordInfo := range alternative.Words {
			startTime := float64(wordInfo.StartTime.Seconds) + float64(wordInfo.StartTime.Nanos)/1e9
			endTime := float64(wordInfo.EndTime.Seconds) + float64(wordInfo.EndTime.Nanos)/1e9

			allWords = append(allWords, models.WordInfo{
				Word:       wordInfo.Word,
				StartTime:  startTime,
				EndTime:    endTime,
				Confidence: float64(wordInfo.Confidence),
				SpeakerTag: wordInfo.SpeakerTag,
			})
		}
	}

	fullTranscript := strings.TrimSpace(transcriptBuilder.String())

	return &models.TranscriptionResult{
		Transcript:   fullTranscript,
		LanguageCode: recognitionConfig.LanguageCode,
		Words:        allWords,
	}, nil
}