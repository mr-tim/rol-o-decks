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

python_binary(
    name='elm_helper',
    main='elm_helper.py',
)


def elm_build(name, srcs=[], elm_json):
    elm_stuff = genrule(
        name=f'{name}#elmstuff',
        tools=[':elm_helper'],
        srcs={'json':[elm_json]},
        cmd=[
            f'$TOOL --stuff',
            'cp -r $PKG_DIR/elm-stuff .'
        ],
        outs=['.elm'],
    )

    compiled_elm_output=f'{name}.js'

    compiled_elm_module = genrule(
        name=f'{name}',
        tools=[':elm_helper'],
        srcs={'elm': srcs, 'json': [elm_json]},
        deps=[elm_stuff],
        cmd=[
            'cp -r $PKG_DIR/.elm $HOME',
            f'$TOOL --build',
            f'cp $PKG_DIR/$NAME.js {compiled_elm_output}'
        ],
        outs=[compiled_elm_output],
    )


bundled_files(
    name='ui_bundle',
    srcs=['index.html', ':ui'],
    visibility=['//cmd/main:rol-o-decks']
)

elm_build(
    name='ui',
    srcs=['src/Main.elm'],
    elm_json='elm.json'
)