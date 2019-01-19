import sys


def filename_to_camel_case(filename):
    return ''.join([s.capitalize() for s in filename.replace('/', '_').split('.')])


def bundle_file(filename):
    method_name = filename_to_camel_case(filename)
    with open(filename, 'r') as f:
        content = f.read()

    return """\
func {method_name}() string {{
    return `{content}`
}}

""".format(method_name=method_name, content=escape_backticks(content))


def bundle_files(package_name, filenames):
    return f'package {package_name}\n\n' + ''.join([bundle_file(f) for f in filenames])


def escape_backticks(content):
    return content.replace('`', '`+"`"+`')


if __name__ == '__main__':
    package_name = sys.argv[1]
    print(bundle_files(package_name, sys.argv[2:]))
