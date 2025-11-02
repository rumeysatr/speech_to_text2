package speechclient

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings" // Bu satırı ekleyin

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

//ses dosyasını streaming şekilde Google Speech Api'a gönderme
//ses dosyasını chunklar halinde okuyup api'a gönderecek
func (sc *SpeechClient) StreamingRecognize(ctx context.Context, audioPath string, recognitionConfig *speechpb.RecognitionConfig) (*models.TranscriptionResult, error) {
	stream, err := sc.client.StreamingRecognize(ctx)
	if err != nil {
		return nil, fmt.Errorf("Streaming client oluşturulamadı: %w", err)
	}

	if err := stream.Send(&speechpb.StreamingRecognizeRequest{
		StreamingRequest: &speechpb.StreamingRecognizeRequest_StreamingConfig{
			StreamingConfig: &speechpb.StreamingRecognitionConfig{
				Config: recognitionConfig,
			},
		},
	}); err != nil {
		return nil, fmt.Errorf("config gönderilemedi: %w", err)
	}

	audioFile, err := os.Open(audioPath)
	if err != nil {
		return nil, fmt.Errorf("ses dosyası açılamadı: %w", err)
	}
	defer audioFile.Close()

	//ses dosyalarını chunklar halinde gönder
	go func() {
		buffer := make([]byte, 1024*16) //16K chunks
		for {
			n, err := audioFile.Read(buffer)
			if n > 0 {
				stream.Send(&speechpb.StreamingRecognizeRequest{
					StreamingRequest: &speechpb.StreamingRecognizeRequest_AudioContent{
						AudioContent: buffer[:n],
					},
				})
			}
			if err == io.EOF {
				stream.CloseSend()
				break
			}
			if err != nil {
				return
			}
		}
	}()

	// --- DÜZELTME BAŞLANGICI ---

	var transcriptBuilder strings.Builder // 'transcripts' dizisi yerine 'strings.Builder' kullanmak daha verimli
	var allWords []models.WordInfo
	var lastEndTime float64 = 0 // timeOffset'i kaldırıp sadece son bitiş zamanını takip et

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("sonuç alınamadı: %w", err)
		}

		for _, result := range resp.Results {
			if result.IsFinal {
				alternative := result.Alternatives[0]
				
				// Transkripti ekle
				transcriptBuilder.WriteString(alternative.Transcript + " ")

				// Kelime listesini işle
				for _, wordInfo := range alternative.Words {
					startTime := float64(wordInfo.StartTime.Seconds) + float64(wordInfo.StartTime.Nanos)/1e9
					endTime := float64(wordInfo.EndTime.Seconds) + float64(wordInfo.EndTime.Nanos)/1e9

					// **ASIL DÜZELTME BURADA:**
					// Sadece, zaman damgası bir önceki kelimenin bitişinden
					// *sonra* olan kelimeleri ekle. Bu, örtüşmeyi (overlap) engeller.
					if startTime >= lastEndTime {
						allWords = append(allWords, models.WordInfo{
							Word:       wordInfo.Word,
							StartTime:  startTime, // Artık 'adjustedStartTime'a gerek yok
							EndTime:    endTime,   // Artık 'adjustedEndTime'a gerek yok
							Confidence: float64(wordInfo.Confidence),
							SpeakerTag: wordInfo.SpeakerTag,
						})
						
						// Son bitiş zamanını sadece *eklenen* kelimenin bitiş zamanıyla güncelle
						lastEndTime = endTime
					}
					// 'else' durumu (startTime < lastEndTime) basitçe atlanır,
					// çünkü bu bir örtüşmedir ve kelime zaten 'allWords' içinde vardır.
				}
			}
		}
	}
	
	fullTranscript := strings.TrimSpace(transcriptBuilder.String())

	// --- DÜZELTME SONU ---

	return &models.TranscriptionResult{
		Transcript:   fullTranscript,
		LanguageCode: recognitionConfig.LanguageCode,
		Words:        allWords,
	}, nil

}