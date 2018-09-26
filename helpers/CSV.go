package v1

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var catOut, catError *os.File
var catWriter, errorWriter *bufio.Writer

func InitializeWriters(outFile string, errFile string, headers []string) {
	log.Println(outFile)
	absPath, _ := filepath.Abs(outFile)
	var err error
	catOut, err = os.Create(absPath)
	log.Println(err)
	catWriter = bufio.NewWriter(catOut)

	catError, _ = os.Create(errFile)
	errorWriter = bufio.NewWriter(catError)

	catWriter.WriteString(strings.Join(headers, ",") + "\n")

}

func SaveRow(row string) {

	fmt.Println(row)

	if catWriter != nil {

		catWriter.WriteString(row + "\n")
		catWriter.Flush()

	}

}

func CloseWriters() {

	catWriter.Flush()
	errorWriter.Flush()

	catError.Sync()
	catError.Close()

	catOut.Sync()
	catOut.Close()

}
