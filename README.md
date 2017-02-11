# `spiro` - a project template runner

Spiro is an template structure generator that uses Golangs text/template library. It accepts both single files as well
as directory trees as input and will interpret any template calls found inside the files and the file/directory names.

The rule-set is probably a bit complex to display here, but the following links are useful:

- https://golang.org/pkg/text/template
- https://gohugo.io/templates/go-templates/

The only additional rule is the rule that controls whether a file or directory is processed or not. If a file name is
templated like `{{ if blah }}filename.txt{{ end }}` then that file will only be processed _if_ the name evaluates to a
non-empty string.

The contents of a file will only be treated as templated if the file name has a `.templated` suffix. If
it does, the contents will be evaluated and the `.templated` suffix will be removed.

Templating _inside_ the file is evaluated after any template in the file name. So if you want an optional file that has
templated content you'll need to use a name like `{{ if blah }}filename.txt.templated{{ end }}`. If the `.templated`
declaration is outside the condition the behaviour should be similar but is probably not the convention.

Some additional template functions are supplied:

- 'title': capitalise string
- 'upper': convert string to upper case
- 'lower': convert string to lower case
- 'now': return current time object (time.Time)

The spec file should be in JSON or Yaml form and will be passed to each template invocation.

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

By default the normal Golang template characters `{{` are used but sometimes the files you're working with containing
and you have to laboriously escape them. 

You can provide the special key `_spiro_delimiters_` in your spec file in order to override them:

```yaml 
_spiro_delimiters_: 
  - "<<<"
  - ">>>"
```

### What should you use this for:

- Does your team have a template project that gets copied and modified by hand? Why not use this!

## Demo

A demo exists in the `/demo` directory. Run it as follows:

```
$ rm -rfv demo/output/project

$ ./spiro demo/project demo/spec.json demo/output
Processing 'demo/project/' -> 'demo/output/project/'
Processing 'demo/project/demo-{{upper .animal}}.templated' -> 'demo/output/project/demo-BEAR'
Processing 'demo/project/{{ if .x }}dontskip.txt{{ end }}' -> 'demo/output/project/dontskip.txt'
Skipping 'demo/project/{{ if not .x }}skipthis.txt{{ end }}' since the name evaluated to ''
Processing 'demo/project/{{.subdir}}-thing/' -> 'demo/output/project/Elephant-thing/'
Processing 'demo/project/{{.subdir}}-thing/noop' -> 'demo/output/project/Elephant-thing/noop'
Processing 'demo/project/{{.subdir}}-thing/{{.subfile.name}}.{{.subfile.type}}' -> 'demo/output/project/Elephant-thing/snake.xml'

$ find demo/output
demo/output
demo/output/project
demo/output/project/demo-BEAR
demo/output/project/dontskip.txt
demo/output/project/Elephant-thing
demo/output/project/Elephant-thing/noop
demo/output/project/Elephant-thing/snake.xml
```

## Download

The best option is to download the latest binaries from the [releases page](https://github.com/AstromechZA/spiro/releases).
Extract the one for your platform and put it in any directory where you have access.

If a binary is not available for your platform, you'll need to build one yourself.

## Development

This project uses only two development time dependency:

- `govendor`: for managing the `vendor/` directory
- `govvv`: for embedding build versions and dates into the binaries

You'll want to add these into your GOPATH using `go get`.

Also, run `govendor sync` to synchronise your vendor folder once you've pulled the repository.

## Future features

- More useful template functions (need feedback from users)
- Syntax to split a single file into multiple
