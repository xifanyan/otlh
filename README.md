# otlh: OpenText LegalHold

`otlh` is a command-line interface (CLI) tool for managing legal holds service. It allows users to perform various operations related to legal holds, such as creating, updating, and listing legal holds , as well as managing custodians.

## Installation

To install `otlh`, you can download the latest release from the [GitHub releases page](https://gitlab.otxlab.net/rec-ps-dev/legalhold/otlh/-/releases) or build it from source using Go.

### Building from Source

1. Make sure you have Go installed on your system. You can download it from the official [Go website](https://golang.org/dl/).
2. Clone the repository: 
```
git clone https://gitlab.otxlab.net/rec-ps-dev/legalhold/otlh
```
3. Navigate to the project directory:
```
cd otlh/cmd/cli
```
4. Build the binary:
```
./build.sh (macos or linux) or ./build.bat (windows)
```
5. The `otlh` binary will be created in the bin/ directory.

## Usage
```
bin/otlh.exe
NAME:
   otlh - Command Line Interface to access Opentext LegalHold service

USAGE:
   otlh [global options] command [command options] 

VERSION:
   0.1-alpha

COMMANDS:
   create   
   get      
   import   
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --domain value, -x value     domain name for Opentext legahold service (default: "api.otlegalhold.com") [$LHN_DOMAIN]
   --port value, -p value       port (default: 443)
   --tenant value, -c value     tenant name [$LHN_TENANT]
   --authToken value, -t value  token to access legalhold web service (default: "") [$LHN_AUTHTOKEN]
   --debug, -d                  Debug Mode: Log to stderr (default: false)
   --help, -h                   show help
   --version, -v                print the version
```