Rol-o-Decks
===========
Index and search slides in Powerpoint presentations.

Setup
-----
Run the rol-o-decks binary - this will start the server on http://localhost:8000. The first time you run the application, you'll need to set up directories that you want to index for presentations - you can do this by:
- Clicking on the `...` button next to the search box
- Entering the path to where your presentations are (eg: /path/to/presentations)
- Clicking save

Rol-o-decks will index the presentations in the background and watch for changes over time. Once the presentations are indexed, you should be able to search for them, the click on the search results to open them.

Building
--------
Rol-o-decks is built using [please](https://please.build) - to build it simply run the `./pleasew` wrapper script. Building from source requires both [Go](https://golang.org/) and [elm](https://elm-lang.org/) to be installed.
