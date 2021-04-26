package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	heroku "github.com/heroku/heroku-go/v5"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var pairs arrayFlags

var (
	username = flag.String("username", "", "api username")
	password = flag.String("password", "", "api password")
	appname  = flag.String("appname", "", "app name")

	show   = flag.Bool("show", false, "show config vars")
	update = flag.Bool("update", false, "update config vars")
)

func main() {
	log.SetFlags(0)
	flag.Var(&pairs, "pairs", "value of config var")
	flag.Parse()

	args := os.Args
	if len(args) < 3 {
		flag.Usage()
		os.Exit(1)
	}

	heroku.DefaultTransport.Username = *username
	heroku.DefaultTransport.Password = *password
	h := heroku.NewService(heroku.DefaultClient)

	switch {
	case *show:
		vars, err := h.ConfigVarInfoForApp(context.TODO(), *appname)
		if err != nil {
			log.Fatal(err)
		}
		for k, v := range vars {
			fmt.Print(k, "=", *v, "\n")
		}
	case *update:
		if len(pairs) == 0 {
			log.Fatal("Missing key pairs")
		}
		var newVars = make(map[string]*string)
		for _, pair := range pairs {
			p := strings.Split(pair, "=")
			if len(p) == 0 {
				log.Fatal("Pair format not specified")
			}
			key := p[0]
			val := p[1]
			newVars[key] = &val
		}
		vars, err := h.ConfigVarUpdate(context.TODO(), *appname, newVars)
		log.Println(vars, err)
	default:
		log.Fatal("Supported target: update, show")
	}
}
