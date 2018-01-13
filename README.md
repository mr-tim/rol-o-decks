Rol-o-Decks
===========

Index and search slides in Powerpoint presentations

Currently only works on macOS, and requires Python 3.6 and Node 9.3.

Installation
------------
Set up the python virtualenv, and install all python and node dependencies by running the install script `install.sh`.

You can specify any number of paths where presentations will be indexed using `config.json`, for example:

```json
{
    "paths": [
        "/path/to/presentations/directory1", "/path/to/presentations/directory2"
    ]
}
```

You can then start and stop the servers using `start.sh` and `stop.sh`.

