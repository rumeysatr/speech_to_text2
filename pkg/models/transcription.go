package models

import "time"

//tüm deşifre işleminin sonucu için
type TranscriptionResult struct {
	Transcript		string		   `json:"transcript"`
	Confidence		float64		   `json:"confidence"`
	LanguageCode 	string		   `json:"language_code"`
	Words			[]WordInfo	   `json:"words"`
	AudioDuration	float64		   `json:"file_length"`
	ProcessedAt		time.Time 	   `json:"processed_time"`
	Speakers		[]SpeakerInfo  `json:"speakers,omitempty"` //opsiyonel olduğundan omiempty
	KeywordMatches	[]KeywordMatch `json:"keyword_matches,omitempty"`

}

//her kelime için API'dan gelen zaman damgası ve güven bilgisi
type WordInfo struct {
	Word			string  	    `json:"word"`
	StartTime		float64			`json:"start_time"`
	EndTime			float64			`json:"end_time"`
	Confidence		float64			`json:"confidence"`
	SpeakerTag		int32			`json:"speaker_tag,omitempty"`

}

//A konuşmacısı kaç dakika konuştuğu gibi analizler 
type SpeakerInfo struct {
	SpeakerTag		int32			`json:"speaker_tag"`
	TotalDuration 	float64			`json:"total_duration"`
	WordCount		int 			`json:"word_count"`
	Transcript		string			`json:"text_of_speaker"`

}

//istenen keywordün nerede geçtiğini bulmak için
type KeywordMatch struct {
	Keyword			string 			`json:"keyword"`
	Timestamp 		float64			`json:"timestamp"`
	Context			string 			`json:"context"`
	SpeakerTag		int32			`json:"speaker_tag,omitempty"` //opsiyonel

}



