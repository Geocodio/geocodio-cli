# geocodio-cli

Command line application to create and manage spreadsheet uploads on Geocodio

## Download

Download our [latest release](https://github.com/Geocodio/geocodio-cli/releases), or use `go install github.com/geocodio/geocodio-cli@latest`

## Authentication

All commands require a valid API Key, which can be generated from the [Geocodio Dashboard](https://dash.geocod.io/apikey).

It can either be supplied using the global `--apikey` flag, or using an environment variable:

```bash
export GEOCODIO_API_KEY=__MY_API_KEY__
```

## Examples

### Create a new spreadsheet job

Upload a spreadsheet for geocoding. In this example, we're using a test file called [sample_list.csv](https://www.geocod.io/sample_list.csv).

The full address consists of column `A` (address), `B` (city), `C` (state) and `D` (zip):

```bash
$ ./geocodio create sample_list.csv "{{A}} {{B}} {{C}} {{D}} USA" --fields=cd
✅ Success: Spreadsheet job created
Id: 11472143
Filename: sample_list.csv
Headers: address | city | state | zip
Rows: 24
```

> Read more about the field format parameter in the [Geocodio API Docs](https://www.geocod.io/docs/#format-syntax).

### Check the job status

Check the current job progress using the spreadsheet id returned from the previous command:

```bash
$ ./geocodio status 11472143
Id: 11472143
Filename: sample_list.csv
Fields:
Rows: 24
State: COMPLETED
Progress: 100%
Message: Completed
Time left: -
Expires: Jan 21 15:18:31 2022
Download URL: https://api.geocod.io/v1.9/lists/11472143/download
```

You can also use the `--follow` flag to monitor the progress in real-time with a progress bar:

```bash
$ ./geocodio status 11472143 --follow
```

### List all spreadsheet jobs

View all your existing geocoding jobs:

```bash
$ ./geocodio list
```

### Download the geocoded file

Once the job has completed, it can be downloaded using the spreadsheet id.

The geocoded spreadsheet data will be outputted to `stdout`, so you can pipe it to a file to get a completed spreadsheet:

```bash
$ ./geocodio download 11472143 > sample_list_geocoded.csv
```

### Delete a spreadsheet job

Remove a spreadsheet job by id:

```bash
$ ./geocodio remove 11472143
```

## Global Options

```
--hostname value, -n value  Geocodio hostname to use, change this for Geocodio Enterprise or on-premise environments (default: "api.geocod.io") [$GEOCODIO_HOSTNAME]
--apikey value, -k value    Geocodio API Key to use. Generate a new one in the Geocodio Dashboard [$GEOCODIO_API_KEY]
--help, -h                  show help (default: false)
--version, -v               print the version (default: false)
```
