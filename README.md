# geocodio-cli

Command line application to create and manage spreadsheet uploads on Geocodio

## Usage

```
$ ./geocodio
NAME:
   Geocodio - Geocode lists using the Geocodio API

USAGE:
   geocodio [global options] command [command options] [arguments...]

VERSION:
   dev

COMMANDS:
   create      Geocode a new spreadsheet
   list        List existing geocoding jobs
   status      Query the status for a specific geocoding job
   remove, rm  Delete an existing geocoding job
   download    Download output data for a specific geocoding job
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --hostname value, -n value  Geocodio hostname to use, change this for Geocodio+HIPAA or on-premise environments (default: "api.geocod.io") [$GEOCODIO_HOSTNAME]
   --apikey value, -k value    Geocodio API Key to use. Generate a new one in the Geocodio Dashboard (default: "84a871a4ed113cc616ae328ae731a371c6dad14") [$GEOCODIO_API_KEY]
   --help, -h                  show help (default: false)
   --version, -v               print the version (default: false)
```
