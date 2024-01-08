# go-test-api
Learning Go with a simple API. It exposes some test data and allows the user to upload valid data to the API.

It runs on port 8117 by default, but you can customize which port it runs on by setting the `listenPort` environment variable in your system.

## How to use
Clone the repo. If you have Go installed, do:

`go run .`

Or for Docker,

`docker build -t go-test-api .`

`docker run --rm -it -p 8117:8117 go-test-api`

You can now reach the API at `http://localhost:8117/albums` to see the data or use the `/upload` endpoint to send data in the same schema. It is not necessary to specify a row ID when uploading, just the album name, artist, and price.

### Set custom listening port

The server runs on port 8117 by default, but to change that you just set the `listenPort` environment variable in your system. On Linux/macOS, that's `export listenPort=1234` and for Windows that's `$env:listenPort='1234'`. For Docker, set an environment variable with the `-e` flag like so:

`docker run --rm -it -e listenPort=1234 -p 1234:1234 go-test-api`

Be sure to set it to a valid port number or the server will not start.