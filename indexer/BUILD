go_library(
    name = 'indexer',
    srcs = [
        'indexer.go',
    ],
    deps = [
        '//server/handlers:handlers',
        '//store:store',
        '//third_party:fsnotify',
        '//third_party:xmlpath',
    ],
    visibility = [
        '//cmd/main'
    ]
)