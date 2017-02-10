package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"
	"time"
)

const usageString = `
Spiro is an template structure generator that uses Golangs text/template library. It accepts both single files as well 
as directory trees as input and will interpret any template calls found inside the files and the file/directory names.

The rule-set is probably a bit complex to display here, but the following links are useful:

- https://golang.org/pkg/text/template
- https://gohugo.io/templates/go-templates/

Some additional template functions are supplied:

- 'title': capitalise string
- 'upper': convert string to upper case 
- 'lower': convert string to lower case
- 'now': return current time object (time.Time)

The spec file should be in JSON form and will be passed to each template invocation.

$ spiro [options] {input template} {spec file} {output directory}
`

const logoImage = `
  _________      .__               
 /   _____/_____ |__|______  ____  
 \_____  \\____ \|  \_  __ \/  _ \ 
 /        \  |_> >  ||  | \(  <_> )
/_______  /   __/|__||__|   \____/ 
        \/|__|               
`

var Version = "<unofficial build>"
var GitSummary = "<changes unknown>"
var BuildDate = "<no date>"

func copyFileContents(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

func doTemplate(templateString string, spec *map[string]interface{}, funcMap *template.FuncMap) (string, error) {
	t := template.New("").Option("missingkey=error").Funcs(*funcMap)
	if t, err := t.Parse(templateString); err != nil {
		return "", err
	} else {
		var buf bytes.Buffer
		err := t.Execute(&buf, spec)
		return buf.String(), err
	}
}

func isTemplatedName(name string) bool {
	return strings.Contains(name, "{{") && strings.Contains(name, "}}")
}

func isTemplatedFile(name string) bool {
	return strings.HasSuffix(name, ".templated")
}

func processDir(templateString string, spec *map[string]interface{}, outputDir string, funcMap *template.FuncMap) error {
	base := path.Base(templateString)
	if isTemplatedName(base) {
		var err error
		base, err = doTemplate(base, spec, funcMap)
		if err != nil {
			return fmt.Errorf("Error while processing '%s': %s", templateString, err.Error())
		}
	}

	newOutputDir := path.Join(outputDir, base)
	fmt.Printf("Processing '%s/' -> '%s/'\n", templateString, newOutputDir)
	if err := os.Mkdir(newOutputDir, 0755); err != nil && !os.IsExist(err) {
		return fmt.Errorf("Error while processing '%s': %s", templateString, err.Error())
	}

	items, _ := ioutil.ReadDir(templateString)
	for _, item := range items {
		if err := process(path.Join(templateString, item.Name()), spec, newOutputDir, funcMap); err != nil {
			return err
		}
	}
	return nil
}

func processFile(templateString string, spec *map[string]interface{}, outputDir string, funcMap *template.FuncMap) error {
	fromBase := path.Base(templateString)
	toBase := fromBase
	if isTemplatedName(fromBase) {
		var err error
		toBase, err = doTemplate(fromBase, spec, funcMap)
		if err != nil {
			return fmt.Errorf("Error while processing '%s': %s", templateString, err.Error())
		}
	}

	fmt.Printf("Processing '%s' -> '%s'\n", templateString, path.Join(outputDir, toBase))
	if isTemplatedFile(toBase) {
		inputBytes, err := ioutil.ReadFile(templateString)
		if err != nil {
			return fmt.Errorf("Error while processing '%s': %s", templateString, err.Error())
		}
		outputBytes, err := doTemplate(string(inputBytes), spec, funcMap)
		if err != nil {
			return fmt.Errorf("Error while processing '%s': %s", templateString, err.Error())
		}
		if err := ioutil.WriteFile(path.Join(outputDir, toBase), []byte(outputBytes), 0644); err != nil {
			return fmt.Errorf("Error while processing '%s': %s", templateString, err.Error())
		}
	} else {
		return copyFileContents(templateString, path.Join(outputDir, toBase))
	}
	return nil
}

func process(templateString string, spec *map[string]interface{}, outputDir string, funcMap *template.FuncMap) error {
	stat, err := os.Stat(templateString)
	if err != nil {
		return fmt.Errorf("Error processing template %s: %s", templateString, err.Error())
	}
	if stat.IsDir() {
		return processDir(templateString, spec, outputDir, funcMap)
	} else {
		return processFile(templateString, spec, outputDir, funcMap)
	}
}

func tfNow() time.Time {
	return time.Now()
}

func getFuncMap() *template.FuncMap {
	tm := template.FuncMap{
		"title": strings.Title,
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
		"now":   tfNow,
	}
	return &tm
}

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
		fmt.Println("Project: github.com/AstromechZA/spiro")
		return nil
	}
	if flag.NArg() != 3 {
		flag.Usage()
		os.Exit(1)
	}

	inputTemplate := flag.Arg(0)
	specFile := flag.Arg(1)
	outputDirectory := flag.Arg(2)

	if _, err := os.Stat(inputTemplate); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("Input template '%s' does not exist!", inputTemplate)
		}
		return fmt.Errorf("Input template '%s' cannot be read! (%s)", inputTemplate, err.Error())
	}
	if stat, err := os.Stat(specFile); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("Spec file '%s' does not exist!", specFile)
		}
		return fmt.Errorf("Spec file '%s' cannot be read! (%s)", specFile, err.Error())
	} else if stat.IsDir() {
		return fmt.Errorf("Spec file '%s' cannot be a directory!", specFile)
	}
	if stat, err := os.Stat(outputDirectory); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("Output directory '%s' does not exist!", outputDirectory)
		}
		return fmt.Errorf("Output directory '%s' cannot be read! (%s)", outputDirectory, err.Error())
	} else if !stat.IsDir() {
		return fmt.Errorf("Output directory '%s' cannot be a file!", specFile)
	}

	specBytes, err := ioutil.ReadFile(specFile)
	if err != nil {
		return fmt.Errorf("Could not read json spec file: %s", err.Error())
	}
	var spec map[string]interface{}
	err = json.Unmarshal(specBytes, &spec)
	if err != nil {
		return fmt.Errorf("Could not parse json spec file: %s", err.Error())
	}

	return process(inputTemplate, &spec, outputDirectory, getFuncMap())
}

func main() {
	if err := mainInner(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}
