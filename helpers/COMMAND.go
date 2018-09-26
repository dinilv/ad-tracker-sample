package v1

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

var mongoExportPath string

//dev
//var mongoExportPath = "/Users/dinilv/Downloads/platform/mongodb-osx-x86_64-3.4.5/bin/mongoexport"

//prod
func init() {
	mongoExportPath = "mongoexport"
}
func RunMongoExport(arguments []string) {

	log.Println("Start writing to CSV:-", time.Now().UTC())

	cmd := exec.Command(mongoExportPath, arguments...)
	log.Println(cmd.Args)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err.Error())
	}
	log.Println("Done writing to CSV:-", time.Now().UTC())

}

func RunSystemRestart() {

	log.Println("Starting Sytem Restart:-", time.Now().UTC())

	cmd := exec.Command("sudo", "reboot")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err.Error())
	}
	log.Println("Done Restart", time.Now().UTC())

}
