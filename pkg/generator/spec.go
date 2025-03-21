package generator

import (
	"fmt"
	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"io"
	"net/http"
	"os"
)

type SpecLoader interface {
	LoadSpec(source string) (*libopenapi.DocumentModel[v3.Document], error)
}

type urlSpecLoader struct {
	client *http.Client
}

func NewURLSpecLoader(client *http.Client) SpecLoader {
	l := urlSpecLoader{
		client: client,
	}

	if l.client == nil {
		l.client = http.DefaultClient
	}

	return &l
}

func (l *urlSpecLoader) LoadSpec(url string) (*libopenapi.DocumentModel[v3.Document], error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error building request for %s: %w", url, err)
	}

	res, err := l.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request to %s: %w", url, err)
	}

	specJson, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response from %s: %w", url, err)
	}

	return buildSpec(specJson)
}

func buildSpec(specByteArray []byte) (*libopenapi.DocumentModel[v3.Document], error) {
	config := datamodel.DocumentConfiguration{
		AllowFileReferences:   true,
		AllowRemoteReferences: false,
	}

	document, err := libopenapi.NewDocumentWithConfiguration(specByteArray, &config)
	if err != nil {
		return nil, fmt.Errorf("error while opening spec as OpenAPI document: %w", err)
	}

	model, errs := document.BuildV3Model()
	if len(errs) > 0 {
		return nil, fmt.Errorf("many errors while building model: %s", errs)
	}

	return model, nil
}

type fileSpecLoader struct{}

func NewFileSpecLoader() SpecLoader {
	return &fileSpecLoader{}
}

func (f *fileSpecLoader) LoadSpec(source string) (*libopenapi.DocumentModel[v3.Document], error) {
	file, err := os.ReadFile(source)
	if err != nil {
		return nil, err
	}

	return buildSpec(file)
}
