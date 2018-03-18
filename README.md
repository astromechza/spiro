# `spiro` - a project template runner

![travis ci badge](https://travis-ci.org/AstromechZA/spiro.svg?branch=master)

Spiro is an template structure generator that uses Golangs text/template library. It accepts both single files as well as directory trees as input and will interpret any template calls found inside the files and the file/directory names.

The rule-set is probably a bit complex to display here, but the following links are useful:

- https://golang.org/pkg/text/template
- https://gohugo.io/templates/go-templates/

The only additional rule is the rule that controls whether a file or directory is processed or not. If a file name is templated like `{{ if blah }}filename.txt{{ end }}` then that file will only be processed _if_ the name evaluates to a non-empty string.

The contents of a file will only be treated as templated if the file name has a `.templated` suffix. If it does, the contents will be evaluated and the `.templated` suffix will be removed.

Templating _inside_ the file is evaluated after any template in the file name. So if you want an optional file that has templated content you'll need to use a name like `{{ if blah }}filename.txt.templated{{ end }}`. If the `.templated` declaration is outside the condition the behaviour should be similar but is probably not the convention.

Some additional template functions are supplied:

- `title`: capitalise string `(string) -> (string)`
- `upper`: convert string to upper case `(string) -> (string)`
- `lower`: convert string to lower case `(string) -> (string)`
- `now`: return current time object `() -> (time.Time)`
- `json`: output a structure as json `(object) -> (string)`
- `jsonindent`: output a structure as indented json `(object) -> (string)`
- `unescape`: unescape escaped html characters `(string) -> (string)`
- `stringreplace`: basic string replace `(subject, old, new) -> (string)`
- `regexreplace`: regular expression based string replace `(subject, pattern, repl) -> (string)`

The spec file should be in JSON or Yaml form and will be passed to each template invocation. The specfile can be "-" to indicate that YAML should be read from stdin.

Permission bits for any files, including `.templated` ones, **will** be copied to the destination files.

### Basic example of features:

You have a file on disk called `{{ lower .projectname }}.md.templated` with the following content:

```
# Heading for {{ .projectname }}

This project was started on {{ now.Format "2006-01-02" }} by {{ .author }}.
```

And if you feed it the following spec JSON:

```json
{
    "projectname": "HelloWorld",
    "author": "Joe Soap"
}
```

You'll end up with a file called `helloworld.md` containing:

```
# Heading for HelloWorld

This project was started on 2017-02-11 by Joe Soap.
```

### Overriding the template characters

By default the normal Golang template characters `{{` are used but sometimes the files you're working with containing and you have to laboriously escape them.

You can provide the special key `_spiro_delimiters_` in your spec file in order to override them:

```yaml
_spiro_delimiters_:
  - "<<<"
  - ">>>"
```

### Enforcing a `spiro` version

Sometimes new features are added to Spiro which are not supported by earlier versions. Sometimes templates rely on these features. By specifying a `_spiro_min_version_` in your spec file, an error will be thrown if an earlier version of `spiro` is used to build the template.

```yaml
_sprio_min_version_: 1.5
```

The version rules work as follows:

- 3 numbers are allowed (`major.minor.patch`)
- `X.Y` == `X.Y.0` and `X` == `X.0.0`
- `X.Z` >= `X.0`

### What should you use this project for:

- Does your team have a template project that gets copied and modified by hand? Use `spiro`!

## Demo

Some demos exist in the `/demo` directory. Run them as follows:

```
$ rm -rfv demo/output/project

$ ./spiro demo/example1 demo/example1.yml demo/output
Processing 'demo/example1/' -> 'demo/output/example1/'
Processing 'demo/example1/demo-{{upper .animal}}.templated' -> 'demo/output/example1/demo-BEAR'
Processing 'demo/example1/{{ if .x }}dontskip.txt{{ end }}' -> 'demo/output/example1/dontskip.txt'
Skipping 'demo/example1/{{ if not .x }}skipthis.txt{{ end }}' since the name evaluated to ''
Processing 'demo/example1/{{.subdir}}-thing/' -> 'demo/output/example1/Elephant-thing/'
Processing 'demo/example1/{{.subdir}}-thing/noop' -> 'demo/output/example1/Elephant-thing/noop'
Processing 'demo/example1/{{.subdir}}-thing/{{.subfile.name}}.{{.subfile.type}}' -> 'demo/output/example1/Elephant-thing/snake.xml'

$ find demo/output
demo/output
demo/output/example1
demo/output/example1/demo-BEAR
demo/output/example1/dontskip.txt
demo/output/example1/Elephant-thing
demo/output/example1/Elephant-thing/noop
demo/output/example1/Elephant-thing/snake.xml
```

## Download & Installation

The best option is to download the latest binaries from the [releases page](https://github.com/AstromechZA/spiro/releases). Extract the one for your platform and put it in any directory where you have access.

Alternatively, use the install script which will do this for you:

```
$ curl https://raw.githubusercontent.com/AstromechZA/spiro/master/install.sh | sh
```

If a binary is not available for your platform, you'll need to build one yourself.

## Changelog

**v1.8**

- Added `-edit` option to the CLI

**v1.7**

- Allow spec to be read from stdin when `-` is given ([#5](https://github.com/AstromechZA/spiro/issues/5))

**v1.6.0**

- ...

## Future features

- More useful template functions (need feedback from users)
- Syntax to split a single file into multiple
