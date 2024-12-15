package generator

import (
	"fmt"
	"github.com/moznion/gowrtr/generator"
	"os"
	"path"
)

func EmitToFile(targetPath, pkgPath string, gen *generator.Root) error {
	outputPath := path.Join(targetPath, pkgPath)
	dirname := path.Dir(outputPath)
	if err := os.MkdirAll(dirname, 0o755); err != nil {
		return fmt.Errorf("error while mkdir'ing '%s': %w", dirname, err)
	}

	code, err := gen.Generate(0)
	if err != nil {
		return fmt.Errorf("error while generating code: %w", err)
	}

	if err := os.WriteFile(outputPath, []byte(code), 0o644); err != nil {
		return fmt.Errorf("error writing to file %s: %w", outputPath, err)
	}

	return nil
}
