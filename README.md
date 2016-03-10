tapi: Google Translate API from the terminal
==============================

`tapi` is a command line tool to access google translate from terminal.

![tapi setup movie](https://s3.amazonaws.com/com-gariany-translate/tapi-movie.gif)

## Usage

### setup

Run setup dialog using:
```
tapi --setup
```
Provide the API Key and languages preferences and you're ready to go.

> **Notice:** Setup will save a configuration file at: `$HOME/.tapirc` with your API key, and preferences.


### Translate
```
tapi "some string to translate"
```

> **Notice:** `tapi` will perform a **reverse** request if Google Translate replied with similar text as user's input. In that matter, translating between 2 languages is as easy and does not require different input from the user.

## Install

### Obtaining a Google Translate API Key

For `tapi` to work you should have access to [Google Developer Console](https://console.developers.google.com) to be able to create a Google Translate API key. Follow Google's instruction on [Translate API Getting Started]( https://cloud.google.com/translate/v2/getting_started) page.

Also important, you **must** [enable billing](https://cloud.google.com/console/help/console/#billing) on your google developer account in order for the API to reply with a proper response. With that in mind, you should probably review [Translate API Pricing ](https://cloud.google.com/translate/v2/pricing) page. At this moment, current pricing is: **$20 per 1,000,000 characters of text.**

Once obtaining the key, set up your preferred **target** and **origin** languages and the reset is pretty straightforward.

### Building

#### From Github Releases
TBD.

#### From Source

`tapi` is written in [Go](https://golang.org/). To have a working go environment, follow Go's [Getting Started](https://golang.org/doc/install) guide. tapi uses [glide](https://github.com/Masterminds/glide) - Vendor Package Management for Golang, follow the short install guide their excellent [README](https://github.com/Masterminds/glide#install) file.

- configure: `export GOPATH=YOUR_GOLANG_PATH`
- `cd $GOPATH`
- `go get github.com/eladg/tapi`
- `cd src/github.com/eladg/tapi`
- `glide install`
- `go build -o $GOPATH/bin/tapi tapi.go`

To execute, make sure compiled file is on `$PATH`. It's recommended to set something like: `export PATH=$PATH:$GOPATH/bin` for your builds.

## Contributing
See [CONTRIBUTING](CONTRIBUTING.md)

## License
The MIT License (MIT)

See [LICENSE](LICENSE)
