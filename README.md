# otlh: OpenText LegalHold

`otlh` is a command-line interface (CLI) tool for managing legal holds service. It allows users to perform various operations related to legal holds, such as creating, updating, and listing legal holds , as well as managing custodians.

Key Features:
- Single executable file with no dependencies on Runtime env (e.g., java) or external libraries (e.g., 3rd party dlls)
- Cross platform support (Windows, Linux, and MacOS)

## Installation

To install `otlh`, you can download the latest release from the [GitHub releases page](https://github.com/xifanyan/otlh/releases) or build it from source using Go.

### Building from Source

1. Make sure you have Go installed on your system. You can download it from the official [Go website](https://golang.org/dl/).
2. Clone the repository: 
```
git clone https://github.com/xifanyan/otlh
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
### Top-level Commands
```
bin/otlh.exe
NAME:
   otlh - Command Line Interface to access Opentext LegalHold service

USAGE:
   otlh [global options] command [command options] 

VERSION:
   0.1.0-beta

COMMANDS:
   create   
   get      
   import   
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --domain value, -x value     domain name for Opentext legahold service (default: "api.otlegalhold.com") [%LHN_DOMAIN%]
   --port value, -p value       port (default: 443)
   --tenant value, -t value     tenant name [%LHN_TENANT%]
   --authToken value, -a value  token to access legalhold web service [%LHN_AUTHTOKEN%]
   --config value, -c value     LHN json config file (default: ".otlh.json") [%LHN_CONFIG%]
   --debug, -d                  Debug Mode: Log to stderr (default: false)
   --help, -h                   show help
   --version, -v                print the version
```

#### Notes
- domain name defaults to api.otlegalhold.com and port defaults to 443 (normally you do not have to specify them).
- tenant is mandatory and can be specified in the config file or via environment variable LHN_TENANT.
- authToken is mandatory and can be specified in the config file or via environment variable LHN_AUTHTOKEN.
- config file is optional and can be specified via environment variable LHN_CONFIG, it uses the following format:
```
{
    "domain": "api.otlegalhold.com",
    "port": 443,
    "tenant": "test",
    "authToken": "*************************"
}
```

### Legalhold Import
```
./otlh.exe import legalholds -h
NAME:
   otlh import legalholds

USAGE:
   otlh import legalholds [command options] [arguments...]

CATEGORY:
   import

OPTIONS:
   --attachmentDirectory value, --ad value  attachment directory (default: ".")
   --excel value, -e value                  excel file used for legalhold import
   --holdName value, --hn value             hold name
   --matterName value, --mn value           matter name
   --checkInputOnly, --ci                   check input only (default: false)
   --help, -h                               show help
```
#### Examples
- Validate input data in excel file only
```
./otlh.exe import legalholds --excel=../testdata/sample.xlsx --attachmentDirectory=../testdata/attachments --checkInputOnly
```

- Import all legalholds from excel file
```
./otlh.exe --tenant test --authToken [*****] import legalholds --excel=../testdata/sample.xlsx --attachmentDirectory=../testdata/attachments
```

- Debug mode
```
./otlh.exe --debug --tenant test --authToken [*****] import legalholds --excel=../testdata/sample.xlsx --attachmentDirectory=../testdata/attachments
```

- Partially import legalholds from excel file based on matter or hold names
```
./otlh.exe --tenant test --authToken [*****] import legalholds --excel=../testdata/sample.xlsx --attachmentDirectory=../testdata/attachments --matterName="Fargo vs Acme" --holdName="Fargo vs Acme Legal Hold"
```

#### Notes
- all date fields in excel need to follow pattern "1/2/06 3:04 PM", by default, they will be considered to be UTC time.