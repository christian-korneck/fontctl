package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	docs "github.com/urfave/cli-docs/v3"
	cli "github.com/urfave/cli/v3"
)

var (
	dbg debugLogger
)

func main() {
	var app *cli.Command
	app = &cli.Command{
		Name:        "fontctl",
		Usage:       "Install or uninstall a font on MS Windows",
		Description: "Copyright (C) 2025 Christian Korneck <christian@korneck.de>",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "enable verbose debug logging",
			},
		},
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			if c.Bool("debug") {
				dbg = CliDebugLogger{}
			}
			return ctx, nil
		},
		Commands: []*cli.Command{
			{
				Name:      "install",
				Usage:     "Install a font",
				UsageText: "fontctl install [--systemwide] <Font File>",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "systemwide",
						Aliases: []string{"s"},
						Usage:   "Install in the system font dir (default: install in the current user's userprofile). Requires Admin privileges.",
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					if c.NArg() != 1 {
						cli.ShowCommandHelpAndExit(ctx, app, c.Name, 1)
					}
					installSystemWide := c.Bool("systemwide")
					err := InstallFontFromFile(c.Args().First(), installSystemWide)
					if err != nil {
						return cli.Exit(fmt.Sprintf("Error - %s", err), 1)
					}
					return nil
				},
			},
			{
				Name:      "uninstall",
				Usage:     "Uninstall a font",
				UsageText: "fontctl uninstall [--systemwide] <Font File>",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "systemwide",
						Aliases: []string{"s"},
						Usage:   "Uninstall from the system font dir (default: uninstall from the user font dir). Requires Admin privileges.",
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					if c.NArg() != 1 {
						cli.ShowCommandHelpAndExit(ctx, app, c.Name, 1)
					}
					uninstallSystemWide := c.Bool("systemwide")
					err := UninstallFontFromFile(c.Args().First(), uninstallSystemWide)
					if err != nil {
						return cli.Exit(fmt.Sprintf("Error - %s", err), 1)
					}
					return nil
				},
			},
			{
				Name:      "getname",
				Usage:     "Get the font name from a file",
				UsageText: "fontctl getname <Font File>",
				Action: func(ctx context.Context, c *cli.Command) error {
					if c.NArg() != 1 {
						cli.ShowCommandHelpAndExit(ctx, app, c.Name, 1)
					}
					fontName, err := GetFontNameFromFile(c.Args().First())
					if err != nil {
						return cli.Exit(fmt.Sprintf("Error - %s", err), 1)
					}
					fmt.Println(fontName)
					return nil
				},
			},
			{
				Name:        "load",
				Usage:       "Load a font into memory",
				UsageText:   "fontctl load <Font File>",
				Description: "This makes a font temporarily available to applications, until the font gets unloaded or the next reboot",
				Action: func(ctx context.Context, c *cli.Command) error {
					if c.NArg() != 1 {
						cli.ShowCommandHelpAndExit(ctx, app, c.Name, 1)
					}
					err := LoadFontFromFile(c.Args().First())
					if err != nil {
						return cli.Exit(fmt.Sprintf("Error - %s", err), 1)
					}
					return nil
				},
			},
			{
				Name:      "unload",
				Usage:     "Unload a font from memory",
				UsageText: "fontctl unload <Font File>",
				Action: func(ctx context.Context, c *cli.Command) error {
					if c.NArg() != 1 {
						cli.ShowCommandHelpAndExit(ctx, app, c.Name, 1)
					}
					err := UnloadFontFromFile(c.Args().First())
					if err != nil {
						return cli.Exit(fmt.Sprintf("Error - %s", err), 1)
					}
					return nil
				},
			},
			{
				Name:        "refresh",
				Usage:       "Refresh known fonts for current user session",
				UsageText:   "fontctl refresh",
				Description: "Sends a WM_FONTCHANGE broadcast so currently running applications become aware of font changes.",
				Action: func(ctx context.Context, c *cli.Command) error {
					err := NotifyFontChange()
					if err != nil {
						return cli.Exit(fmt.Sprintf("Error - %s", err), 1)
					}
					return nil
				},
			},
			{
				Name:  "preview",
				Usage: "Preview a font",
				Commands: []*cli.Command{
					{
						Name:      "file",
						Usage:     "Preview a font file using Windows Font Viewer",
						ArgsUsage: "<Font File>",
						Action: func(ctx context.Context, c *cli.Command) error {
							if c.NArg() < 1 {
								cli.ShowSubcommandHelpAndExit(c, 1)
							}
							err := PreviewFontWithFontview(c.Args().First())
							if err != nil {
								return cli.Exit(err, 1)
							}
							return nil
						},
					},
					{
						Name:  "font",
						Usage: "Preview a loaded font using Windows GDI",
						Description: `Render font preview in a gdi window. Please note that when a font is missing or can't be used, Windows will fall back to displaying the text in a default font. So what you are seeing might not always be the font that you've requested.

Example:
fontctl preview font "Comic Sans MS" "regular"`,
						ArgsUsage: "<Font Name> <Font Style: regular|bold|italic|bold-italic>",
						Action: func(ctx context.Context, c *cli.Command) error {
							if c.NArg() < 1 {
								cli.ShowSubcommandHelpAndExit(c, 1)
							}
							if c.NArg() < 2 {
								return cli.Exit("Missing Font Name or Font Style", 1)
							}
							fontName := c.Args().Get(0)
							fontStyle := strings.ToLower(c.Args().Get(1))
							allowed := map[string]bool{
								"regular":     true,
								"bold":        true,
								"italic":      true,
								"bold-italic": true,
							}
							if !allowed[fontStyle] {
								return cli.Exit(fmt.Sprintf("invalid Font Style: %s (must be one of regular, bold, italic, bold-italic)", fontStyle), 1)
							}
							PreviewFontWithGDI(fontName, fontStyle)
							return nil
						},
					},
				},
			},
			{
				Name:   "mddocs",
				Hidden: true,
				Usage:  "Print CLI Markdown Docs",
				Action: func(ctx context.Context, c *cli.Command) error {
					md, err := docs.ToMarkdown(app)
					if err != nil {
						return err
					}

					fmt.Println("# " + app.FullName() + "\n\n" + md)

					return nil
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
