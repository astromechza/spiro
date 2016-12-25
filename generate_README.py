#!/usr/bin/env python

import os
from textwrap import dedent
import subprocess

"""
This script can be launched to generate the README.md file. This script is a useful way of ensuring that your README
document is up to date and the relevant examples still work as intended.

```
$ ./generate_README.py -o
or
$ ./generate_README.py > README.md
```
"""

# the file that will be written
PROJECT_DIRECTORY = os.path.dirname(__file__)


class Generator(object):
    def __init__(self):
        self.lines = []

    def heading(self, text, level):
        self.lines.append('#' * level + " " + text)
        self.lines.append("")

    def h1(self, text):
        return self.heading(text, 1)

    def h2(self, text):
        return self.heading(text, 2)

    def h3(self, text):
        return self.heading(text, 3)

    def h4(self, text):
        return self.heading(text, 4)

    def paragraph(self, text):
        text = text.rstrip()
        self.lines.append(dedent(text))
        self.lines.append("")

    def command_example(self, command):
        self.lines.append("```")
        self.lines.append("$ {}".format(command))

        try:
            output = subprocess.check_output(command, stderr=subprocess.STDOUT, shell=True, cwd=PROJECT_DIRECTORY)
        except subprocess.CalledProcessError as e:
            output = e.output

        self.lines.append(output.strip())
        self.lines.append("```")
        self.lines.append("")

    def __str__(self):
        text = "\n".join(self.lines)
        if not text.endswith("\n"):
            text += "\n"
        return text


def main():
    g = Generator()
    g.h1("`go-cli-template` - an example Go CLI application")
    g.paragraph("""\
    This is an example repository that can be cloned and adapted for a new application. It contains useful scripts and
    automation that I have found useful while building simple CLI-based applications.
    """)

    g.h3("Example usage:")
    g.command_example("./go-cli-template -version")
    g.command_example("./go-cli-template -help")

    g.h3("Steps for setting up a new project from this template")
    g.paragraph("""\
    1. Clone the repository

    ```
    $ git clone https://github.com/AstromechZA/go-cli-template.git
    ```
    """)

    g.h3("Building the application")
    g.paragraph("""\
    The provided `make_official.sh` script will build official builds for both Linux and OSX with an official version
    number baked in.
    """)
    g.command_example("./make_official.sh")

    print str(g)


if __name__ == '__main__':
    main()
