package Internal

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
)

var (
	TONAME   string
	FROMNAME string
	EOL      string
)

func init() {
	if runtime.GOOS == "windows" {
		fmt.Println("pipe-test.go, running on windows")
		TONAME = `\\.\pipe\ToSrvPipe`
		FROMNAME = `\\.\pipe\FromSrvPipe`
		EOL = "\r\n"
	} else {
		fmt.Println("pipe-test.go, running on linux or mac")
		TONAME = fmt.Sprintf("/tmp/audacity_script_pipe.to.%d", os.Getuid())
		FROMNAME = fmt.Sprintf("/tmp/audacity_script_pipe.from.%d", os.Getuid())
		EOL = "\n"
	}
}

type Connection struct {
	send    *os.File
	recieve *os.File
}

type Audacity struct {
	connection Connection
}

// Establish a connection to audacity
func (a *Audacity) Connect() {
	fmt.Println("Write to  \"" + TONAME + "\"")
	if _, err := os.Stat(TONAME); err != nil {
		fmt.Println(" ..does not exist.  Ensure Audacity is running with mod-script-pipe.")
		os.Exit(1)
	}
	fmt.Println("Read from \"" + FROMNAME + "\"")
	if _, err := os.Stat(FROMNAME); err != nil {
		fmt.Println(" ..does not exist.  Ensure Audacity is running with mod-script-pipe.")
		os.Exit(1)
	}
	fmt.Println("-- Both pipes exist.  Good.")

	toFile, err := os.OpenFile(TONAME, os.O_RDWR, os.ModeNamedPipe)
	if err != nil {
		fmt.Println("-- Failed to open file to write to:", err)
		os.Exit(1)
	}
	fmt.Println("-- File to write to has been opened")
	a.connection.send = toFile

	fromFile, err := os.OpenFile(FROMNAME, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		fmt.Println("-- Failed to open file to read from:", err)
		os.Exit(1)
	}
	fmt.Println("-- File to read from has now been opened too\r")
	a.connection.recieve = fromFile
}

func (a Audacity) Close() {
	a.connection.recieve.Close()
	a.connection.send.Close()
}

// send custom command to audacity. reffer to https://manual.audacityteam.org/man/scripting_reference.html for formatting.
func (a Audacity) Do_command(command string) (res string) {

	//send command
	fmt.Println("Send: >>> \n" + command)
	a.connection.send.Write([]byte(command + EOL))

	//get response
	scanner := bufio.NewScanner(a.connection.recieve)
	for scanner.Scan() {
		text := scanner.Text()
		res += text
		if text == "" && len(res) != 0 {
			break
		}
	}
	return res
}
