package server

import (
	iso6391 "github.com/emvi/iso-639-1"
	"github.com/goccy/go-json"
	"log"
	"net/http"
	"slices"
	"strings"
)

type languagesLang struct {
	Iso6391    string `json:"iso_639_1"`
	Name       string `json:"name"`
	NativeName string `json:"native_name"`
}

func (c *handlers) languages(w http.ResponseWriter, _ *http.Request) {
	allLangs := iso6391.Languages
	resp := make([]languagesLang, 0, len(allLangs))

	for _, lang := range allLangs {
		resp = append(
			resp,
			languagesLang{
				Name:       lang.Name,
				Iso6391:    lang.Code,
				NativeName: lang.NativeName,
			},
		)
	}

	slices.SortFunc(resp, func(a, b languagesLang) int {
		return strings.Compare(a.Name, b.Name)
	})

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println(err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
