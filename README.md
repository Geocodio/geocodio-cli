# geocodio-cli

Command line application to create and manage spreadsheet uploads on Geocodio

## Download

Download our [latest release](https://github.com/Geocodio/geocodio-cli/releases), or use `go install github.com/geocodio/geocodio-cli@latest`

## Usage

```
$ ./geocodio
NAME:
   Geocodio - Geocode lists using the Geocodio API

USAGE:
   geocodio [global options] command [command options] [arguments...]

VERSION:
   0.2.2

COMMANDS:
   create      Geocode a new spreadsheet
   list        List existing geocoding jobs
   status      Query the status for a specific geocoding job
   remove, rm  Delete an existing geocoding job
   download    Download output data for a specific geocoding job
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --hostname value, -n value  Geocodio hostname to use, change this for Geocodio Enterprise or on-premise environments (default: "api.geocod.io") [$GEOCODIO_HOSTNAME]
   --apikey value, -k value    Geocodio API Key to use. Generate a new one in the Geocodio Dashboard [$GEOCODIO_API_KEY]
   --help, -h                  show help (default: false)
   --version, -v               print the version (default: false)
```
