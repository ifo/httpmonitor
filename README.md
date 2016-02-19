httpmonitor
===========

## Installation and Usage

### Prereqs

To run `httpmonitor`, you must install
[the Go programming language](https://golang.org/).

### Installation

- **Extract the zipfile** `unzip httpmonitor.zip`
- **Navigate to the project directory** `cd httpmonitor/`
- **Build the project** `go build`

**Note**: for Go 1.5, you must use experimental vendoring:
`GO15VENDOREXPERIMENT=1 go build`

This should create a binary named `httpmonitor` in the project directory.

### Running

- **Execute the binary** `./httpmonitor --file [path/to/file.log]`

**Note**: The file.log must be a w3c formatted HTTP access log in the
[Common Log](https://www.w3.org/Daemon/User/Config/Logging.html#common-logfile-format)
format.

The Common Log format looks like this:
```
remotehost rfc931 authuser [date] "request" status bytes
```

### Options

Options are available for httpmonitor.
To display them, type `./httpmonitor --help`
