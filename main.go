package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/AstromechZA/spiro/templatefactory"
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

See the project homepage for more documentation: https://github.com/AstromechZA/spiro

The spec file should be in JSON or Yaml form and will be passed to each template invocation.

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
var GitCommit = "<commit unknown>"
var GitState = "<changes unknown>"
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

func processDir(templateString string, spec *map[string]interface{}, outputDir string, tf *templatefactory.TemplateFactory) error {
	fromBase := path.Base(templateString)
	toBase := fromBase
	if tf.StringContainsTemplating(fromBase) {
		var err error
		toBase, err = tf.Render(fromBase)
		if err != nil {
			return fmt.Errorf("Error while processing '%s': %s", templateString, err.Error())
		}
	}
	toBase = strings.TrimSpace(toBase)
	if len(toBase) == 0 {
		fmt.Printf("Skipping '%s' since the name evaluated to ''\n", templateString)
		return nil
	}

	newOutputDir := path.Join(outputDir, toBase)
	fmt.Printf("Processing '%s/' -> '%s/'\n", templateString, newOutputDir)
	if err := os.Mkdir(newOutputDir, 0755); err != nil && !os.IsExist(err) {
		return fmt.Errorf("Error while processing '%s': %s", templateString, err.Error())
	}

	items, _ := ioutil.ReadDir(templateString)
	for _, item := range items {
		if err := process(path.Join(templateString, item.Name()), spec, newOutputDir, tf); err != nil {
			return err
		}
	}
	return nil
}

func processFile(templateString string, spec *map[string]interface{}, outputDir string, tf *templatefactory.TemplateFactory) error {
	fromBase := path.Base(templateString)
	toBase := fromBase
	if tf.StringContainsTemplating(fromBase) {
		var err error
		toBase, err = tf.Render(fromBase)
		if err != nil {
			return fmt.Errorf("Error while processing '%s': %s", templateString, err.Error())
		}
	}
	toBase = strings.TrimSpace(toBase)
	if len(toBase) == 0 {
		fmt.Printf("Skipping '%s' since the name evaluated to ''\n", templateString)
		return nil
	}

	if strings.HasSuffix(toBase, ".templated") {
		toBase = toBase[:len(toBase)-10]
		if len(toBase) == 0 {
			fmt.Printf("Skipping '%s' since the name evaluated to ''\n", templateString)
			return nil
		}
		fmt.Printf("Processing '%s' -> '%s'\n", templateString, path.Join(outputDir, toBase))
		inputBytes, err := ioutil.ReadFile(templateString)
		if err != nil {
			return fmt.Errorf("Error while reading '%s': %s", templateString, err.Error())
		}
		outputBytes, err := tf.Render(string(inputBytes))
		if err != nil {
			return fmt.Errorf("Error while rendering template for '%s': %s", templateString, err.Error())
		}
		if err := ioutil.WriteFile(path.Join(outputDir, toBase), []byte(outputBytes), 0644); err != nil {
			return fmt.Errorf("Error while writing file bytes for '%s': %s", templateString, err.Error())
		}
	} else {
		fmt.Printf("Processing '%s' -> '%s'\n", templateString, path.Join(outputDir, toBase))
		if err := copyFileContents(templateString, path.Join(outputDir, toBase)); err != nil {
			return fmt.Errorf("Error while copying file bytes for '%s': %s", templateString, err.Error())
		}
	}

	info, _ := os.Stat(templateString)
	if err := os.Chmod(path.Join(outputDir, toBase), info.Mode()); err != nil {
		return fmt.Errorf("Error while writing file permissions for '%s': %s", templateString, err.Error())
	}

	return nil
}

func process(templateString string, spec *map[string]interface{}, outputDir string, tf *templatefactory.TemplateFactory) error {
	stat, err := os.Stat(templateString)
	if err != nil {
		return fmt.Errorf("Error processing template %s: %s", templateString, err.Error())
	}
	if stat.IsDir() {
		return processDir(templateString, spec, outputDir, tf)
	}
	return processFile(templateString, spec, outputDir, tf)
}

func readSpec(specFile string) (*map[string]interface{}, error) {
	specBytes, err := ioutil.ReadFile(specFile)
	if err != nil {
		return nil, fmt.Errorf("Could not read json spec file: %s", err.Error())
	}
	var spec map[string]interface{}
	if strings.HasSuffix(specFile, ".json") {
		err = json.Unmarshal(specBytes, &spec)
		if err != nil {
			return nil, fmt.Errorf("Could not parse json spec file: %s", err.Error())
		}
		return &spec, nil
	} else if strings.HasSuffix(specFile, ".yaml") || strings.HasSuffix(specFile, ".yml") {
		err = yaml.Unmarshal(specBytes, &spec)
		if err != nil {
			return nil, fmt.Errorf("Could not parse yaml spec file: %s", err.Error())
		}
		return &spec, nil
	} else {
		return nil, fmt.Errorf("I do not know how to parse the spec, expected .json, .yaml, or .yml")
	}
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
		fmt.Printf("Version: %s (%s-%s) on %s \n", Version, GitCommit, GitState, BuildDate)
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

	if spec, err := readSpec(specFile); err != nil {
		return err
	} else {
		tf := templatefactory.NewTemplateFactory()
		if err := tf.SetSpec(spec); err != nil {
			return err
		}
		tf.RegisterTemplateFunction("title", strings.Title)
		tf.RegisterTemplateFunction("lower", strings.ToLower)
		tf.RegisterTemplateFunction("upper", strings.ToUpper)
		tf.RegisterTemplateFunction("now", time.Now)
		tf.RegisterTemplateFunction("json", Jsonify)
		tf.RegisterTemplateFunction("jsonindent", JsonifyIndent)
		tf.RegisterTemplateFunction("unescape", Unescape)
		return process(inputTemplate, spec, outputDirectory, tf)
	}
}

func main() {
	if err := mainInner(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}
