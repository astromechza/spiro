package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/AstromechZA/spiro/templatefactory"
)

const usageString = `
Spiro is an template structure generator that uses Golangs text/template library. It accepts both single files as well
as directory trees as input and will interpret any template calls found inside the files and the file/directory names.

The rule-set is probably a bit complex to display here, but the following links are useful:

- https://golang.org/pkg/text/template
- https://gohugo.io/templates/go-templates/

See the project homepage for more documentation: https://github.com/AstromechZA/spiro

The spec file should be in JSON or YAML form and will be passed to each template invocation. The specfile can be "-" to
indicate that YAML should be read from stdin.

$ spiro [options] {input template} {spec file} {output directory}
`

const logoImage = `
  _________      .__
 /   _____/_____ |__|______  ____
 \_____  \\____ \|  \_  __ \/  _ \
 /        \  |_) |  ||  | \(  (_) )
/_______  /   __/|__||__|   \____/
        \/|__|
`

// Version is a combination of version information (tag/commit/date/etc)
var Version = "<unofficial build>"

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

	items, err := ioutil.ReadDir(templateString)
	if err != nil {
		return fmt.Errorf("Error while reading '%s': %s", templateString, err.Error())
	}
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

	info, err := os.Stat(templateString)
	if err != nil {
		return fmt.Errorf("Error while checking file permissions for '%s': %s", templateString, err.Error())
	}
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
	var f *os.File
	var err error
	if specFile == "-" {
		f = os.Stdin
	} else {
		f, err = os.Open(specFile)
		if err != nil {
			return nil, fmt.Errorf("Could not read spec file: %s", err.Error())
		}
	}
	var spec map[string]interface{}
	if strings.HasSuffix(specFile, ".json") {
		dec := json.NewDecoder(f)
		err = dec.Decode(&spec)
		if err != nil {
			return nil, fmt.Errorf("Could not parse json spec file: %s", err.Error())
		}
		return &spec, nil
	} else if specFile == "-" || strings.HasSuffix(specFile, ".yaml") || strings.HasSuffix(specFile, ".yml") {
		dec := yaml.NewDecoder(f)
		err = dec.Decode(&spec)
		if err != nil {
			return nil, fmt.Errorf("Could not parse yaml spec file: %s", err.Error())
		}
		return &spec, nil
	} else {
		return nil, fmt.Errorf("I do not know how to parse the spec, expected .json, .yaml, or .yml")
	}
}

// Build an integer from a version string. The version string can contain 3 numbers and each number can be a maximum
// of 999. 1.2.3 -> 100200300.
func buildVersionInt(versionString string) (uint64, error) {
	parts := strings.Split(versionString, ".")
	var value uint64
	for index := 0; index < 3; index++ {
		if len(parts) > index {
			v, err := strconv.Atoi(parts[index])
			if v < 0 {
				v = 0
			} else if v > 999 {
				v = 999
			}
			if err != nil {
				return value, fmt.Errorf("Could not parse version part %s in %s", parts[index], versionString)
			}
			value += uint64(v)
		}
		value *= 1000
	}
	return value, nil
}

// This function compares the current version to the one in the spec file. If the running version is too low, return
// an error. The version numbers are compared as integers.
func checkVersionIfNecessary(spec *map[string]interface{}) error {
	if minVersion, ok := (*spec)["_spiro_min_version_"]; ok {
		if minVersionString, ok := minVersion.(string); ok {

			// extract 3 digit version from Version
			match := regexp.MustCompile(`v(\d+\.\d+\.\d+)`).FindStringSubmatch(Version)
			if match == nil {
				return fmt.Errorf("You are running an unofficial build of Spiro: we cannot handle version matches")
			}
			currentVersion := match[1]

			if minVersionValue, err := buildVersionInt(minVersionString); err != nil {
				return err
			} else if currentVersionValue, err := buildVersionInt(currentVersion); err != nil {
				return err
			} else if currentVersionValue < minVersionValue {
				return fmt.Errorf("Spiro template lists minimum version %s but you're using %s!", minVersionString, Version)
			}
		}
	}
	return nil
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
		fmt.Printf("Version: %s\n", Version)
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
	if specFile == "-" {
	} else if stat, err := os.Stat(specFile); err != nil {
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
	} else if err := checkVersionIfNecessary(spec); err != nil {
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
		tf.RegisterTemplateFunction("stringreplace", StringReplace)
		tf.RegisterTemplateFunction("regexreplace", RegexReplace)
		return process(inputTemplate, spec, outputDirectory, tf)
	}
}

func main() {
	if err := mainInner(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}
