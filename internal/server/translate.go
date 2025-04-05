package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/OwO-Network/DeepLX/translate"
	"github.com/goccy/go-json"
)

type translateRequest struct {
	Text          string   `json:"text"`
	Dest          string   `json:"dest"`
	Src           string   `json:"src"`
	ExcludedWords []string `json:"excluded_words"`
}

type translateResponse struct {
	SourceLanguage      string   `json:"source_language"`
	SourceText          string   `json:"source_text"`
	TranslatedText      []string `json:"translated_text"`
	DestinationLanguage string   `json:"destination_language"`
}

func (c *handlers) translate(w http.ResponseWriter, r *http.Request) {
	requestBody := translateRequest{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if requestBody.Src == requestBody.Dest {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if requestBody.Text == "" {
		http.Error(w, "no text provided", http.StatusBadRequest)
		return
	}

	if requestBody.Dest == "" {
		http.Error(w, "no destination language provided", http.StatusBadRequest)
		return
	}

	source := "auto"
	if requestBody.Src != "" {
		source = requestBody.Src
	}

	text := requestBody.Text

	// Replace excluded words with placeholders before translation
	if len(requestBody.ExcludedWords) > 0 {
		for idx, word := range requestBody.ExcludedWords {
			text = strings.Replace(text, word, fmt.Sprintf("__%d__", idx), -1)
		}
	}

	translated, err := translate.TranslateByDeepLX(source, requestBody.Dest, text, "", "", "")
	if err != nil {
		slog.Error("cannot make request to deepl", slog.Any("err", err))
		http.Error(w, "cannot make request to deepl", http.StatusInternalServerError)
	}
	if translated.Code != 200 {
		slog.Error("deepl returned non-200 code", slog.Any("code", translated.Code))
		http.Error(w, "deepl returned non-200 code", http.StatusInternalServerError)
	}

	// After translation, replace placeholders back with original words
	translatedText := translated.Data
	if len(requestBody.ExcludedWords) > 0 {
		for idx, word := range requestBody.ExcludedWords {
			translatedText = strings.Replace(translatedText, fmt.Sprintf("__%d__", idx), word, -1)
		}
	}

	result := translateResponse{
		SourceLanguage:      translated.SourceLang,
		DestinationLanguage: translated.TargetLang,
		SourceText:          translated.Message,
		TranslatedText:      []string{translatedText},
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&result); err != nil {
		slog.Error("cannot encode response", slog.Any("err", err))
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
