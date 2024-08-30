# otlh: OpenText LegalHold

`otlh` is a command-line interface (CLI) tool for managing opentext legalhold service. It allows users to perform various operations related to legal holds, such as creating, updating, and listing legal holds , as well as managing matter, folders and custodians.

Key Features:
- Single executable file with no dependencies on Runtime env (e.g., java) or external libraries (e.g., 3rd party dlls)
- Cross platform support (Windows, Linux, and MacOS)

## Installation

- Download the latest binary release from the [GitHub releases page](https://github.com/xifanyan/otlh/releases)

### Building from Source (Optional)
- Make sure you have Go & git installed on your system. You can download it from the official [Go website](https://golang.org/dl/).
- Clone the repository: 
```
git clone https://github.com/xifanyan/otlh
```
- Navigate to the project directory:
```
cd otlh/cmd/cli
```
- Build the binary:
```
./build.sh (macos or linux) or ./build.bat (windows)
```
- The `otlh` binary will be created in the bin/ directory.

## Usage
### Top-level Commands
```
NAME:
   otlh - Command Line Interface to access Opentext LegalHold service

USAGE:
   otlh [global options] command [command options]

VERSION:
   0.4.0-beta

COMMANDS:
   create
   get
   import
   verify
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --domain value, -x value     domain name for Opentext legahold service (default: "api.otlegalhold.com") [%LHN_DOMAIN%]
   --port value, -p value       port (default: 443)
   --tenant value, -t value     tenant name [%LHN_TENANT%]
   --authToken value, -a value  token to access legalhold web service [%LHN_AUTHTOKEN%]
   --config value, -c value     LHN json config file (default: ".otlh.json") [%LHN_CONFIG%]
   --debug, -d                  Debug Mode (default: false)
   --trace, -z                  Trace Mode (default: false)
   --help, -h                   show help
   --version, -v                print the version
```

#### Notes
- domain name defaults to api.otlegalhold.com 
- port defaults to 443
- tenant is mandatory and can be specified in the config file or via environment variable LHN_TENANT.
- authToken is mandatory and can be specified in the config file or via environment variable LHN_AUTHTOKEN.
- config file is optional and can be specified via environment variable LHN_CONFIG, it defaults to .otlh.json with following format:
```
{
    "domain": "api.otlegalhold.com",
    "port": 443,
    "tenant": "test",
    "authToken": "*************************"
}
```

### Get - list otlh entities
```
NAME:
   otlh get

USAGE:
   otlh get command [command options]

COMMANDS:
   help, h  Shows a list of commands or help for one command
   get:
     custodians
     folders
     groups
     matters
     legalholds
     silentholds

OPTIONS:
   --help, -h  show help
```
#### Examples
optins --filterName, --filterTerm, --pageSize, --pageNumber, --sort, --id apply to all get commands

- Get custodian with specific id
```
./otlh.exe get custodians --id 100000383
```

- Get custodians with name filter
```
./otlh.exe get custodians --filterName paul
```

- Get custodians with pagination
```
./otlh.exe get custodians --pageSize 2 --pageNumber 2
```

### Import Legalholds/Silentholds
All of the options (except --attachmentDirectory) apply to Silenthold import as well
```
NAME:
   otlh import legalholds

USAGE:
   otlh import legalholds [command options] [arguments...]

CATEGORY:
   import

OPTIONS:
   --attachmentDirectory value, --ad value  attachment directory (default: ".")
   --excel value, -e value                  excel file used for legalhold import
   --zipfile value, -z value                zip package for importing holds e.g., legal_hold_details.zip
   --timezone value, --tz value             timezone for dates used in input file, supproted timezones: PST|EST|MST|CST (default: "UTC")
   --holdName value, --hn value             hold name
   --matterName value, --mn value           matter name
   --checkInputOnly, --ci                   check input only (default: false)
   --help, -h                               show help
```
#### Examples
- Import zip package directly via command line
```
./otlh.exe --debug import legalholds --zipfile=../testdata/legal_hold_details.zip
```

- Validate input data in excel file only
```
./otlh.exe --debug import legalholds --excel=../testdata/sample.xlsx --attachmentDirectory=../testdata/attachments --checkInputOnly
```

- Import all legalholds from excel file, also convert datetime fields in PST to UTC.
```
./otlh.exe --tenant test --authToken [*****] import legalholds --timezone=PST --excel=../testdata/sample.xlsx --attachmentDirectory=../testdata/attachments
```

- Trace mode
```
./otlh.exe --trace --tenant test --authToken [*****] import legalholds --excel=../testdata/sample.xlsx --attachmentDirectory=../testdata/attachments
```

- Partially import legalholds from excel file based on matter or hold names
```
./otlh.exe --tenant test --authToken [*****] import legalholds --excel=../testdata/sample.xlsx --attachmentDirectory=../testdata/attachments --matterName="Fargo vs Acme" --holdName="Fargo vs Acme Legal Hold"
```

#### Notes
- all datetime fields in excel need to follow pattern "1/2/06 3:04 PM", UTC is the default timezone, if you want to change it, please use --timezone option.
- Attachment files should be put under attachment directory, it currently does not support subfolders, so please make attachment file names unique.