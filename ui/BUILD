python_binary(
    name='bundle_files',
    main='bundle_files.py',
)

def bundled_files(name, srcs=[], visibility=[]):
    tools=[':bundle_files']
    generated_sources = genrule(
        name=f'{name}#src',
        tools=tools,
        srcs=srcs,
        cmd=f'$TOOL {name} $SRCS > {name}.go',
        outs=[f'{name}.go']
    )

    go_library(
        name=name,
        srcs=[generated_sources],
        visibility=visibility,
        deps=[generated_sources]
    )

bundled_files(
    name='ui',
    srcs=['index.html', 'ui.min.js'],
    visibility=['//cmd/main']
)