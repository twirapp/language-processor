package server

import (
	"github.com/goccy/go-json"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
)

type detectResult struct {
	Text              string         `json:"text"`
	CleanedText       string         `json:"cleaned_text"`
	DetectedLanguages []detectedLang `json:"detected_languages"`
}

type detectedLang struct {
	Language    string  `json:"language"`
	Probability float32 `json:"probability"`
}

func (c *handlers) detect(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")
	if text == "" {
		http.Error(w, "no text provided", http.StatusBadRequest)
		return
	}

	cleanedText := cleanPredictText(text)

	languages, err := c.detector.Detect(cleanedText)
	if err != nil {
		slog.Error("cannot detect language", slog.Any("err", err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	resp := detectResult{
		CleanedText:       cleanedText,
		Text:              text,
		DetectedLanguages: make([]detectedLang, len(languages)),
	}

	for i, lang := range languages {
		resp.DetectedLanguages[i] = detectedLang{
			Language:    lang.Label,
			Probability: lang.Probability,
		}
	}

	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		slog.Error("cannot encode response", slog.Any("err", err))
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

var predictSymbolsRegexp = regexp.MustCompile(`(?mi)[\d\p{P}\p{S}]+`)

func cleanPredictText(text string) string {
	cleaned := predictSymbolsRegexp.ReplaceAllString(text, " ")
	cleaned = strings.Join(strings.Fields(cleaned), " ")
	cleaned = strings.ToLower(cleaned)
	cleaned = strings.TrimSpace(cleaned)

	return cleaned
}
