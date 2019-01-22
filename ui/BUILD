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
            'cp -r $PKG_DIR/elm-stuff .',
            'find .elm -type f'
        ],
        outs=['elm-stuff/0.19.0/summary.dat'],
        post_build=lambda rule_name, output: [add_out(rule_name, out) for out in output if out.startswith('.elm')],
    )

    compiled_elm_output=f'{name}_built.js'

    compiled_elm_module = genrule(
        name=f'{name}_built',
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

    minified = genrule(
        name=f'{name}',
        tools=['uglifyjs'],
        cmd=[
            """uglifyjs $PKG_DIR/""" + compiled_elm_output + """ --compress 'pure_funcs="F2,F3,F4,F5,F6,F7,F8,F9,A2,A3,A4,A5,A6,A7,A8,A9",pure_getters,keep_fargs=false,unsafe_comps,unsafe' | uglifyjs --mangle --output=$NAME.min.js"""
        ],
        deps=[compiled_elm_module],
        outs=[f'{name}.min.js']
    )


bundled_files(
    name='ui_bundle',
    srcs=['index.html', ':ui'],
    visibility=['//cmd/main']
)

elm_build(
    name='ui',
    srcs=['src/Main.elm'],
    elm_json='elm.json'
)