import argparse
import os
import subprocess
import sys

parser = argparse.ArgumentParser(description='Elm build helper for plz')
group = parser.add_mutually_exclusive_group()
group.add_argument('--stuff', action='store_true')
group.add_argument('--build', action='store_true')


def download_stuff():
    src_dir = os.path.join(os.environ['PKG_DIR'], 'src')
    os.makedirs(src_dir, exist_ok=True)

    with open(os.path.join(src_dir, 'Blank.elm'), 'w+') as f:
        f.write("""module Blank exposing (blank)

blank : ()
blank = ()
""")
    return subprocess.call(['elm', 'make', '--output', '/dev/null', 'src/Blank.elm'],
                          cwd=os.environ['PKG_DIR'])


def build():
    pkg_dir = os.environ['PKG_DIR']
    srcs = os.environ['SRCS_ELM']
    unprefixed_srcs = [s[len(pkg_dir)+1:] if s.startswith(pkg_dir) else s for s in srcs.split(' ') ]
    return subprocess.call(['elm', 'make', '--optimize', '--output', os.environ['NAME']+'.js',
                           ' '.join(unprefixed_srcs)], cwd=os.environ['PKG_DIR'])


if __name__ == '__main__':
    args = parser.parse_args()

    if args.stuff:
        sys.exit(download_stuff())
    elif args.build:
        sys.exit(build())
    else:
        parser.error('Please specify a valid action')
