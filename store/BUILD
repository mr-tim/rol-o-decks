go_library(
    name = 'store',
    srcs = [
        'models.go',
        'sqlite.go',
        'store.go',
    ],
    deps = [
        '//third_party:sqlite3',
    ],
    visibility = [
        '//cmd/main',
        '//indexer',
        '//server/handlers',
    ]
)