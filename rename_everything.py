#!/usr/bin/env python

import argparse
import os.path

PROJECT_DIRECTORY = os.path.dirname(__file__)
FILES_TO_CHECK = (
    'generate_README.py',
    'make_official.sh',
    'main.go',
    'Godeps/Godeps.json',
)


def main():
    p = argparse.ArgumentParser()
    p.add_argument("importpath", help="The project import path like 'github.com/<user>/repo'")
    args = p.parse_args()

    args.importpath = args.importpath.strip()

    replacers = [
        ('github.com/AstromechZA/go-cli-template', args.importpath),
        ('go-cli-template', os.path.basename(args.importpath))
    ]

    for f in FILES_TO_CHECK:
        path = os.path.join(PROJECT_DIRECTORY, f)
        if os.path.exists(path):
            content = str(open(path, 'r').read())
            for find_s, replace_s in replacers:
                content = content.replace(find_s, replace_s)
            print 'Rewriting %s content' % path
            with open(path, 'w') as f:
                f.write(content)

    os.remove(__file__)

if __name__ == '__main__':
    main()
