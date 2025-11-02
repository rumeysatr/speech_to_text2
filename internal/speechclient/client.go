package speechclient

import (
	"context"
	"fmt"
	"io"
	"os"

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
	
	var transcripts []string
	var allWords []models.WordInfo
	
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
				transcripts = append(transcripts, alternative.Transcript)
				
				for _, wordInfo := range alternative.Words {
					allWords = append(allWords, models.WordInfo{
						Word: wordInfo.Word,
						StartTime: float64(wordInfo.StartTime.Seconds) + float64(wordInfo.StartTime.Nanos)/1e9,
						EndTime:   float64(wordInfo.EndTime.Seconds) + float64(wordInfo.EndTime.Nanos)/1e9,
						Confidence: float64(wordInfo.Confidence),
						SpeakerTag: wordInfo.SpeakerTag,
					})
				}
			}
		}
	}

	fullTranscript := ""
	for _, t := range transcripts {
		fullTranscript += t + " "
	}

	return &models.TranscriptionResult{
		Transcript:   fullTranscript,
		LanguageCode: recognitionConfig.LanguageCode,
		Words:        allWords,
	}, nil

}
