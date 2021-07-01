package main

import (
	"flag"
	"fmt"
)

func main() {
	host := flag.String("host", "localhost", "Host name")
	user := flag.String("user", "root", "User name")
	password := flag.String("password", "123456789", "Password")
	database := flag.String("database", "jmu", "Database")
	shipInfoID := flag.String("shipInfoID", "1", "Ship information ID")
	startTime := flag.String("startTime", "", "Start time")
	endTime := flag.String("endTime", "", "End time")
	outDir := flag.String("outDir", "", "Output Directory")

	var svar string
	flag.StringVar(&svar, "svar", "bar", "a string var")

	flag.Parse()

	fmt.Println("host:", *host)
	fmt.Println("user:", *user)
	fmt.Println("password:", *password)
	fmt.Println("database:", *database)
	fmt.Println("shipInfoID:", *shipInfoID)
	fmt.Println("startTime:", *startTime)
	fmt.Println("endTime:", *endTime)
	fmt.Println("outDir:", *outDir)
	fmt.Println("svar:", svar)
	fmt.Println("tail:", flag.Args())

	// fooCmd := flag.NewFlagSet("foo", flag.ExitOnError)
	// fooEnable := fooCmd.Bool("enable", false, "enable")
	// fooName := fooCmd.String("name", "", "name")

	// barCmd := flag.NewFlagSet("bar", flag.ExitOnError)
	// barLevel := barCmd.Int("level", 0, "level")

	// if len(os.Args) < 2 {
	// 	fmt.Println("expected 'foo' or 'bar' subcommands")
	// 	os.Exit(1)
	// }

	// switch os.Args[1] {

	// case "foo":
	// 	fooCmd.Parse(os.Args[2:])
	// 	fmt.Println("subcommand 'foo'")
	// 	fmt.Println("  enable:", *fooEnable)
	// 	fmt.Println("  name:", *fooName)
	// 	fmt.Println("  tail:", fooCmd.Args())
	// case "bar":
	// 	barCmd.Parse(os.Args[2:])
	// 	fmt.Println("subcommand 'bar'")
	// 	fmt.Println("  level:", *barLevel)
	// 	fmt.Println("  tail:", barCmd.Args())
	// default:
	// 	fmt.Println("expected 'foo' or 'bar' subcommands")
	// 	os.Exit(1)
	// }
}
