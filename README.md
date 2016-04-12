# serve

Simple Go file server for command line development use, a la Python's SimpleHTTPServer.

### Installation

    go get -u github.com/fogleman/serve

### Run the Server

By default, serves the current working directory on port 8000.

    $ serve

[http://localhost:8000/](http://localhost:8000/)

### Command-line Arguments

| Flag | Default | Description |
| --- | --- | --- |
| -port | `8000` | port to listen on |
| -dir | `"."` | directory to serve |
