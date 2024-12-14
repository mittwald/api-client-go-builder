package main

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/mittwald/api-client-go-builder/pkg/generator"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	cmd := &cli.App{
		Name:  "generate",
		Usage: "generate mittwald mStudio v2 API client",
		Action: func(ctx *cli.Context) error {
			log.SetLevel(log.DebugLevel)
			
			gen := generator.Generator{
				SpecLoader: generator.NewURLSpecLoader(nil),
				SchemaGenerator: generator.SchemaGenerator{
					SchemaNamingStrategy: generator.MittwaldV1Strategy,
				},
			}

			genOpts := generator.GeneratorOpts{
				SpecSource: ctx.Args().Get(0),
				Target:     ctx.Args().Get(1),
			}

			fmt.Println(genOpts)

			return gen.Build(genOpts)
		},
	}

	if err := cmd.Run(os.Args); err != nil {
		panic(err)
	}
}
