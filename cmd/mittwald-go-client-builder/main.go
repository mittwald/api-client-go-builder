package main

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/mittwald/api-client-go-builder/pkg/generator"
	"github.com/urfave/cli/v2"
	"os"
	"strconv"
	"strings"
)

func main() {
	cmd := &cli.App{
		Name:  "mittwald-go-client-builder",
		Usage: "Helper tool for generating the mittwald mStudio API client",
		Commands: []*cli.Command{
			{
				Name:  "generate",
				Usage: "Generate client code from OpenAPI spec",
				Action: func(ctx *cli.Context) error {
					log.SetLevel(log.DebugLevel)

					gen := generator.Generator{
						SpecLoader: generator.NewURLSpecLoader(nil),
						SchemaGenerator: generator.SchemaGenerator{
							SchemaNamingStrategy: generator.MittwaldV1Strategy,
						},
					}

					genOpts := generator.GeneratorOpts{
						SpecSource:      ctx.Args().Get(0),
						Target:          ctx.Args().Get(1),
						BasePackageName: ctx.Args().Get(2),
					}

					return gen.Build(genOpts)
				},
			},
			{
				Name:  "next-version",
				Usage: "Automatically determine next version for client release",
				Action: func(ctx *cli.Context) error {
					version := ctx.Args().Get(0)
					versionParts := strings.Split(version, ".")

					patchLevelPart, err := strconv.ParseInt(versionParts[len(versionParts)-1], 10, 64)
					if err != nil {
						return err
					}

					versionParts[len(versionParts)-1] = fmt.Sprintf("%d", patchLevelPart+1)
					fmt.Println(strings.Join(versionParts, "."))
					return nil
				},
			},
		},
	}

	if err := cmd.Run(os.Args); err != nil {
		panic(err)
	}
}
