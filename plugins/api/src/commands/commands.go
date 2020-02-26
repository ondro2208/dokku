package main

import (
	"flag"
	"fmt"
	"github.com/dokku/dokku/plugins/api"
	"os"
	"strconv"
)

func main() {

	flag.Parse()

	cmd := flag.Arg(0)
	switch cmd {
	case "api":
		api.ApiRoute()
	default:
		dokkuNotImplementExitCode, err := strconv.Atoi(os.Getenv("DOKKU_NOT_IMPLEMENTED_EXIT"))
		if err != nil {
			fmt.Println("failed to retrieve DOKKU_NOT_IMPLEMENTED_EXIT environment variable")
			dokkuNotImplementExitCode = 10
		}
		os.Exit(dokkuNotImplementExitCode)
	}
}
