package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Usage string printed as part of the -help text. Keep the width below 80 characters to facilitate terminal
// printing.
const usageString = `
A detailed description of the project that will be inserted above the help text 
when the usage information is printed out.
`

// I like putting logo's in the -version information as a bit of a easter egg/signature. Remove if not required.
// Sources for logos:
// http://www.chris.com/ascii/
// http://www.ascii-code.com/ascii-art/
const logoImage = `
 _        _______  _______  _______ 
( \      (  ___  )(  ____ \(  ___  )
| (      | (   ) || (    \/| (   ) |
| |      | |   | || |      | |   | |
| |      | |   | || | ____ | |   | |
| |      | |   | || | \_  )| |   | |
| (____/\| (___) || (___) || (___) |
(_______/(_______)(_______)(_______)
`

// These variables are filled by the `govvv` tool at compile time.
// There are a few more granular variables available if necessary.
var Version = "<unofficial build>"
var GitSummary = "<changes unknown>"
var BuildDate = "<no date>"


func mainInner() error {

	// first set up config flag options
	versionFlag := flag.Bool("version", false, "Print the version string")

	// set a more verbose usage message.
	flag.Usage = func() {
		os.Stderr.WriteString(strings.TrimSpace(usageString) + "\n\n")
		flag.PrintDefaults()
	}
	// parse them
	flag.Parse()

	// do arg checking
	if *versionFlag {
		fmt.Printf("Version: %s (%s) on %s \n", Version, GitSummary, BuildDate)
		fmt.Println(logoImage)
		fmt.Println("Project: <project url here>")
		return nil
	}

	fmt.Println("Hello World")

	return nil
}

func main() {
	if err := mainInner(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}
