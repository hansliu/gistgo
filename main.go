package main

import (
	"os"
	"log"

	gist "github.com/hansliu/gistgo/gist"

	cli "github.com/urfave/cli/v2" // imports as package "cli"
)

func check(e error) {
    if e != nil {
		log.Fatalln(e)
    }
}

func main() {
	app := cli.NewApp()
	app.Name = "gistgo"
	app.Usage = "A cli tool to upload file to Github Gist"
	app.Version = "1.0.1"

	app.Commands = []*cli.Command{
		{
			Name: "get",
			Aliases: []string{"g"},
			Usage: "Get gist by gistID",
			Action: func(c *cli.Context) error {
				gist.GetGist(c.Args().First())
				return nil
			},
		},
		{
			Name: "upload",
			Aliases: []string{"u"},
			Usage: "Upload file to gist",
			Action: func(c *cli.Context) error {
				// log.Println(c.String("name"))
				gist.UploadGist(c.String("name"), c.Args().First(), c.Bool("public"))
				return nil
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name: "name",
					Aliases: []string{"n"},
					Value: "",
					Usage: "gistfile name",
				},
				&cli.BoolFlag{
					Name: "public",
					Value: false,
				},
			},
		},
	}

	err := app.Run(os.Args)
	check(err)
}
