package cli

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/rombintu/GophKeeper/internal/client/models"
	kpb "github.com/rombintu/GophKeeper/internal/proto/keeper"
	"github.com/urfave/cli/v3"
)

const (
	appDirName = "gophkeeper"
)

type App struct {
	Cmd *cli.Command
}

func NewApp(man *Manager) *App {
	app := &App{
		Cmd: &cli.Command{
			Commands: []*cli.Command{
				{
					Name:  "config",
					Usage: "Manage config",
					Commands: []*cli.Command{
						{
							Name:  "set",
							Usage: "Set global configuration",
							Flags: []cli.Flag{
								&cli.StringFlag{
									Name:  "auth-address",
									Usage: "Address to Auth service",
								},
								&cli.StringFlag{
									Name:  "sync-address",
									Usage: "Address to Sync service",
								},
							},
							Action: func(ctx context.Context, cmd *cli.Command) error {
								// Собираем все значения флагов в map
								configValues := make(map[string]string)

								for _, k := range []string{"auth-address", "sync-address"} {
									if cmd.IsSet(k) {
										configValues[k] = cmd.String(k)
									}
								}
								// Передаем значения в общую функцию
								return man.ConfigSet(ctx, configValues)
							},
						},
						{
							Name:  "get",
							Usage: "Get global configuration",

							Action: func(ctx context.Context, cmd *cli.Command) error {
								// Передаем значения в общую функцию
								data, err := man.ConfigGetMap(ctx)
								if err != nil {
									return err
								}

								if len(data) == 0 {
									slog.Info("config is empty")
									return nil
								}
								for k, v := range data {
									fmt.Println(k, "-->", v)
								}
								return nil
							},
						},
					},
				},
				{
					Name:    "list",
					Aliases: []string{"l"},
					Usage:   "list all secrets",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return man.SecretList(ctx)
					},
				},
				{
					Name:  "login",
					Usage: "Login and get token",
					// Flags: []cli.Flag{

					// },
					Action: func(ctx context.Context, cmd *cli.Command) error {
						addr, err := man.ConfigGet(ctx, "auth-address")
						if err != nil {
							return err
						}
						return man.Login(ctx, addr)
					},
				},
				{
					Name:    "register",
					Aliases: []string{"reg"},
					Usage:   "Register and get token",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						addr, err := man.ConfigGet(ctx, "auth-address")
						if err != nil {
							return err
						}
						return man.Register(ctx, addr)
					},
				},
				{
					Name:  "sync",
					Usage: "Pull and Push cloud secrets",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						addr, err := man.ConfigGet(ctx, "sync-address")
						if err != nil {
							return err
						}
						return man.Sync(ctx, addr)
					},
				},
				{
					Name:    "create",
					Aliases: []string{"new", "add"},
					Usage:   "Create new secret",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "title",
							Usage:    "Title of secret",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "type",
							Value:    "",
							Usage:    "Type of secret",
							Required: true,
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
						&cli.StringFlag{
							Name:  "data",
							Usage: "Data of secret",
						},
						&cli.StringFlag{
							Name:  "url",
							Usage: "URL for new secret",
						},
						&cli.StringFlag{
							Name:  "login",
							Usage: "Your login",
						},
						&cli.StringFlag{
							Name:  "password",
							Usage: "Your password",
						},
						&cli.StringFlag{
							Name:  "owner",
							Usage: "Owner lastname of card",
						},
						&cli.StringFlag{
							Name:  "number",
							Usage: "Number of card",
						},
						&cli.StringFlag{
							Name:  "code",
							Usage: "code of card",
							Validator: func(s string) error {
								_, err := strconv.Atoi(s)
								if err != nil {
									return err
								}
								if len(s) < 3 || len(s) > 3 {
									return errors.New("code is a three digit number")
								}
								return nil
							},
						},
						&cli.StringFlag{
							Name:  "expire",
							Usage: "Expire date",
						},
					},
					Action: func(ctx context.Context, cmd *cli.Command) error {

						secret := models.Secret{
							Title: cmd.String("title"),
						}
						switch cmd.String("type") {
						case strings.ToLower(kpb.Secret_TEXT.String()):
							st := &models.SecretText{
								Secret: secret,
								Text:   cmd.String("data"),
							}
							return man.SecretCreate(ctx, st)
						case strings.ToLower(kpb.Secret_CRED.String()):
							st := &models.SecretCreds{
								Secret: secret,
								Creds: models.Creds{
									URL:      cmd.String("url"),
									Login:    cmd.String("login"),
									Password: cmd.String("password"),
								},
							}
							return man.SecretCreate(ctx, st)
						case strings.ToLower(kpb.Secret_DATA.String()):
							st := &models.SecretBinary{
								Secret:     secret,
								BinaryData: []byte(cmd.String("data")),
							}
							return man.SecretCreate(ctx, st)
						case strings.ToLower(kpb.Secret_CARD.String()):
							st := &models.SecretCard{
								Secret: secret,
								Card: models.Card{
									Owner:      cmd.String("owner"),
									ExpireDate: cmd.String("expire"),
									Number:     cmd.String("number"),
									Code:       cmd.String("code"),
								},
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
