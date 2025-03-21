package main

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/mittwald/api-client-go-builder/pkg/generator"
	"github.com/mittwald/api-client-go-builder/pkg/reference"
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
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "url",
						Usage: "The URL to the openapi.json",
					},
					&cli.StringFlag{
						Name:  "path",
						Usage: "The path to your local openapi.json",
					},
					&cli.StringFlag{
						Name:     "target",
						Usage:    "The target directory, into which the client shall be generated",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "pkg",
						Usage:    "The package name to be generated",
						Required: true,
					},
				},
				Before: func(c *cli.Context) error {
					url := c.String("url")
					path := c.String("path")

					if url == "" && path == "" {
						return fmt.Errorf("either --url or --path must be provided")
					}
					if url != "" && path != "" {
						return fmt.Errorf("only one of --url or --path can be provided, not both")
					}

					return nil
				},
				Action: func(ctx *cli.Context) error {
					source := ctx.String("url")
					specLoader := generator.NewURLSpecLoader(nil)

					if source == "" {
						source = ctx.String("path")
						specLoader = generator.NewFileSpecLoader()
					}

					log.SetLevel(log.DebugLevel)

					apiVersion := "v2"

					gen := generator.Generator{
						SpecLoader: specLoader,
						SchemaGenerator: generator.SchemaGenerator{
							SchemaNamingStrategy: generator.MittwaldAPIVersionSchemaStrategy(apiVersion),
						},
						ReferenceLinkBuilder: reference.NewMittwaldReferenceLinkBuilder(apiVersion),
					}

					genOpts := generator.GeneratorOpts{
						SpecSource:      source,
						Target:          ctx.String("target"),
						BasePackageName: ctx.String("pkg"),
						APIVersion:      apiVersion,
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
