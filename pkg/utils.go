package otlh

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"net/mail"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

/*
Prettify takes a byte slice of JSON data and returns a formatted string
with indentation. It uses json.Indent to format the JSON data.

Parameters:
- b ([]byte): The JSON data to format.

Returns:
- string: The formatted JSON string.
*/
func Prettify(b []byte) string {
	buf := new(bytes.Buffer)
	json.Indent(buf, b, "", "  ")
	return buf.String()
}

func StructToString(i interface{}) (string, error) {
	b, err := json.Marshal(i)
	return Prettify(b), err
}

/*
CreateZipArchive creates a zip archive containing the provided files.

Parameters:
- zipPath: path to the zip file to create
- files: list of file paths to include in the zip

Returns:
- error: any error encountered while creating the zip file, or nil if successful

This function creates a new zip archive at the provided zipPath location containing
the files specified in the files parameter. Each file path in files will be added
to the zip archive using the base name of the file as the path inside the zip.
*/
func CreateZipArchive(zipPath string, files []string) error {

	log.Debug().Msgf("zip: %d files to be processed - %v", len(files), files)

	// Create a new zip file
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	log.Debug().Msgf("zip: %s created", zipPath)

	// Create a new zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Iterate over the files and add them to the zip archive
	for _, filePath := range files {
		// Open the file
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// Get the file information
		info, err := file.Stat()
		if err != nil {
			return err
		}

		// Create a new zip file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Set the name of the file inside the zip
		header.Name = filepath.Base(filePath)

		// Set compression method
		header.Method = zip.Deflate

		// Create a writer for the file header
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// Copy the file content to the zip writer
		_, err = io.Copy(writer, file)
		if err != nil {
			log.Error().Msgf("failed to write zip writer - %s", filePath)
			return err
		}
	}

	return nil
}

func IsValidEmailAddress(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func GetTimezoneLocation(tz string) string {
	switch tz {
	case "CST":
		return "America/Chicago"
	case "EST":
		return "America/New_York"
	case "MST":
		return "America/Denver"
	case "PST":
		return "America/Los_Angeles"
	// Add more cases for other timezones if needed
	default:
		return "UTC"
	}
}
