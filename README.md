# Fresh-dlv

Fresh-dlv is a command line tool that builds and (re)starts your web application everytime you save a Go or template file.

**This fork integrates delve debugger with [fresh](https://github.com/gravityblast/fresh).**

If the web framework you are using supports the Fresh-dlv runner, it will show build errors on your browser.

It currently works with [Traffic](https://github.com/pilu/traffic), [Martini](https://github.com/codegangsta/martini) and [gocraft/web](https://github.com/gocraft/web).

## Installation

    go get github.com/vishal-wadhwa/fresh-dlv

## Usage

    cd /path/to/myapp

Start fresh-dlv:

    fresh-dlv

Start debugging:
1. Start `fresh-dlv` in the root of your project and let the build complete.
2. To start a debugging session. Create an empty file `.debug` in project root.

    ```sh
    touch .debug
    ```
3. Connect to the delve debugger with the following command:

    ```sh
    dlv connect :40000
    ```
4. To disable debugging. Delete `.debug`

    ```sh
    rm .debug
    ```

Fresh-dlv will watch for file events, and every time you create/modify/delete a file it will build and restart the application.
If `go build` returns an error, it will log it in the tmp folder.

[Traffic](https://github.com/pilu/traffic) already has a middleware that shows the content of that file if it is present. This middleware is automatically added if you run a Traffic web app in dev mode with Fresh-dlv.
Check the `_examples` folder if you want to use it with Martini or Gocraft Web.

`fresh-dlv` uses `./runner.conf` for configuration by default, but you may specify an alternative config filepath using `-c`:

    fresh-dlv -c other_runner.conf

Here is a sample config file with the default settings:

    root:               .
    tmp_path:           ./tmp
    build_name:         runner-build
    build_log:          runner-build-errors.log
    valid_ext:          .go, .tpl, .tmpl, .html
    no_rebuild_ext:     .tpl, .tmpl, .html
    ignored:            assets, tmp
    build_delay:        600
    colors:             1
    log_color_main:     cyan
    log_color_build:    yellow
    log_color_runner:   green
    log_color_watcher:  magenta
    log_color_debugger: red
    log_color_app:


## Author

* [Andrea Franz](http://gravityblast.com)

## Delve integration
* [Vishal Wadhwa](https://github.com/vishal-wadhwa)


