package detector

import (
	"flag"
	"os"

	"github.com/nano-interactive/go-fasttext"
)

var modelPath string

func New() (*Detector, error) {
	flag.StringVar(&modelPath, "modelpath", "", "Path to lang model")
	flag.Parse()

	if modelPath == "" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		modelPath = wd + "/lid.176.bin"
	}

	ff, err := fasttext.Open(modelPath)
	if err != nil {
		return nil, err
	}

	return &Detector{
		ff: ff,
	}, nil
}

type Detector struct {
	ff fasttext.Model
}

type Prediction struct {
	CleanedText string
	Label       string
	Probability float32
}

func (p *Detector) Detect(text string) ([]Prediction, error) {
	predictions, err := p.ff.Predict(text, 1, 0)
	if err != nil {
		return nil, err
	}

	pr := make([]Prediction, len(predictions))
	for i := range predictions {
		pr[i] = Prediction{
			Label:       predictions[i].Label,
			Probability: predictions[i].Probability,
		}
	}

	return pr, nil
}
