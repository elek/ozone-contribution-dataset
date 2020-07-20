package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/markbates/pkger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	app := cli.NewApp()
	app.Name = "ghprfetch"
	app.Usage = "Github Data retriever"
	app.Commands = []cli.Command{
		{
			Name:  "update",
			Usage: "Update all pull requests)",
			Action: func(c *cli.Context) error {
				return run(c.Args().Get(0))
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}

}

func run(dir string) error {
	query, err := pkger.Open("/query.graphql")
	if err != nil {
		return err
	}
	defer query.Close()
	originalGraphql, err := ioutil.ReadAll(query)
	if err != nil {
		return err
	}

	firstGraphql := strings.Replace(string(originalGraphql), "hadoop-ozone", "hadoop-ozone", -1)
	graphql := firstGraphql
	hasNextPage := true
	for hasNextPage {
		response, err := asJson(readGithubApiV4Query([]byte(graphql)))
		if err != nil {
			return err
		}
		if m(response, "errors") != nil {
			firstError := l(m(response, "errors"))[0]
			return errors.New(ms(firstError, "message"))
		}
		for _, predge := range l(m(response, "data", "repository", "pullRequests", "edges")) {
			pr := m(predge, "node")
			number := mn(pr, "number")
			println(number)
			err = persist(dir, number, pr)
			if err != nil {
				return err
			}
		}
		println("All the messages are persisted from this batch")
		//hasNextPage = false
		hasNextPage = m(response, "data", "repository", "pullRequests", "pageInfo", "hasNextPage").(bool)
		if hasNextPage {
			cursor := ms(response, "data", "repository", "pullRequests", "pageInfo", "endCursor")
			graphql = strings.Replace(firstGraphql, "pullRequests(", "pullRequests(after:\""+cursor+"\",", 1)
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

func persist(dir string, prnum int, pr interface{}) error {
	destFile := path.Join(dir, strconv.Itoa(prnum)+".json")
	json, err := json.MarshalIndent(pr, "", "   ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(destFile, json, 0644)
	if err != nil {
		return err
	}
	return nil
}
