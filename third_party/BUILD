go_get(
    name = 'gorilla_mux',
    get = 'github.com/gorilla/mux',
    visibility = ['PUBLIC'],
)

go_get(
    name = 'fsnotify',
    get = 'github.com/fsnotify/fsnotify',
    visibility = ['PUBLIC'],
    deps = [':x_sys'],
)

go_get(
    name = 'sqlite3',
    get = 'github.com/mattn/go-sqlite3',
    visibility = ['//store:store'],
)

go_get(
    name = 'x_sys',
    get = 'golang.org/x/sys/...',
)

go_get(
    name = 'xmlpath',
    get = 'gopkg.in/xmlpath.v2',
    visibility = ['PUBLIC'],
    deps = [':x_net_html'],
)

go_get(
    name = 'x_net_html',
    get = 'golang.org/x/net/html',
    install = [
        "",
        "atom",
    ],
)
