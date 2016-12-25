# `go-cli-template` - an example Go CLI application

This is an example repository that can be cloned and adapted for a new application. It contains useful scripts and
automation that I have found useful while building simple CLI-based applications.

### Example usage:

```
$ ./go-cli-template -version
Version: 0.1 (commit cb39180 @ 2016-12-25)

 _        _______  _______  _______ 
( \      (  ___  )(  ____ \(  ___  )
| (      | (   ) || (    \/| (   ) |
| |      | |   | || |      | |   | |
| |      | |   | || | ____ | |   | |
| |      | |   | || | \_  )| |   | |
| (____/\| (___) || (___) || (___) |
(_______/(_______)(_______)(_______)

Project: <project url here>
```

```
$ ./go-cli-template -help
A detailed description of the project that will be inserted above the help text 
when the usage information is printed out.

  -version
    	Print the version string
```

### Steps for setting up a new project from this template

1. Clone the repository

```
$ git clone https://github.com/AstromechZA/go-cli-template.git
```

### Building the application

The provided `make_official.sh` script will build official builds for both Linux and OSX with an official version
number baked in.

```
$ ./make_official.sh
Building official darwin amd64 binary for version '0.1 (commit cb39180 @ 2016-12-25)'
Output Folder build/darwin_amd64
github.com/AstromechZA/go-cli-template
Done
-rwxr-xr-x  1 benmeier  staff   2.5M Dec 25 14:05 build/darwin_amd64/go-cli-template
build/darwin_amd64/go-cli-template: Mach-O 64-bit executable x86_64

Building official linux amd64 binary for version '0.1 (commit cb39180 @ 2016-12-25)'
Output Folder build/linux_amd64
github.com/AstromechZA/go-cli-template
Done
-rwxr-xr-x  1 benmeier  staff   2.5M Dec 25 14:05 build/linux_amd64/go-cli-template
build/linux_amd64/go-cli-template: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), statically linked, not stripped

-rw-r--r--  1 benmeier  staff   1.5M Dec 25 14:05 go-cli-template-0.1.tgz
go-cli-template-0.1.tgz: gzip compressed data, from Unix, last modified: Sun Dec 25 14:05:24 2016
```

