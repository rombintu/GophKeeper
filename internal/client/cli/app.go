package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/rombintu/GophKeeper/internal/client/models"
	kpb "github.com/rombintu/GophKeeper/internal/proto/keeper"
	"github.com/urfave/cli/v3"
)

type App struct {
	Cmd *cli.Command
}

func NewApp(man *Manager) *App {
	app := &App{
		Cmd: &cli.Command{
			Commands: []*cli.Command{
				{
					Name:    "list",
					Aliases: []string{"l"},
					Usage:   "list all secrets",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return man.SecretList(ctx)
					},
				},
				{
					Name:    "create",
					Aliases: []string{"new", "add"},
					Usage:   "Create new secret",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:  "type",
							Value: "text",
							Usage: "Type of secret",
							Validator: func(s string) error {
								var validValues []string
								for _, name := range kpb.Secret_SecretType_name {
									if strings.EqualFold(name, s) {
										return nil
									}
									validValues = append(validValues, name)
								}

								return fmt.Errorf("unknown type of secret. Valid types: %s",
									strings.ToLower(strings.Join(validValues, ",")))
							},
						},
					},
					Action: func(ctx context.Context, cmd *cli.Command) error {
						if cmd.NArg() < 1 {
							return errors.New("too low arguments")
						}
						switch cmd.String("type") {
						case strings.ToLower(kpb.Secret_TEXT.String()):
							st := &models.SecretText{
								Text: strings.Join(cmd.Args().Slice(), " "),
							}
							return man.SecretCreate(ctx, st)
						}
						return man.SecretList(ctx)
					},
				},
			},
		},
	}

	return app
}
