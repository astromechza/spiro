# `spiro` - a project template runner

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

**Important**: the contents of a file will only be treated as templated if the file name has a `.templated` suffix. If 
it does, the contents will be evaluated and the `.templated` suffix will be removed.

#### What should you use this for:

- Does your team have a template project that gets copied and modified by hand? Why not use this!

## Demo

A demo exists in the `/demo` directory. Run it as follows:

```
$ rm -rfv demo/output/project 

$ ./spiro demo/project demo/spec.json demo/output
Processing 'demo/project/' -> 'demo/output/project/'
Processing 'demo/project/demo-{{upper .animal}}.templated' -> 'demo/output/project/demo-BEAR'
Processing 'demo/project/{{.subdir}}-thing/' -> 'demo/output/project/Elephant-thing/'
Processing 'demo/project/{{.subdir}}-thing/noop' -> 'demo/output/project/Elephant-thing/noop'
Processing 'demo/project/{{.subdir}}-thing/{{.subfile.name}}.{{.subfile.type}}' -> 'demo/output/project/Elephant-thing/snake.xml'

$ find demo/project 
demo/project
demo/project/demo-{{upper .animal}}.templated
demo/project/{{.subdir}}-thing
demo/project/{{.subdir}}-thing/noop
demo/project/{{.subdir}}-thing/{{.subfile.name}}.{{.subfile.type}}
```

## Download

The best option is to download the latest binaries from the [releases page](https://github.com/AstromechZA/spiro/releases).
Extract the one for your platform and put it in any directory where you have access.

If a binary is not available for your platform, you'll need to build one yourself.

## Development

This project uses only one development time dependency:

- `govvv`: for embedding build versions and dates into the binaries

You'll want to add it into your GOPATH using `go get`.

## Future features

- More useful template functions (need feedback from users)
- Allow yaml spec file
- Syntax to split a single file into multiple
