//Some convenient Console interface functions:
//Color change, error logging

package mytrace

import (
	"bufio"
	"bytes"
	"fmt"
	"time"

	//    "log"
	"os"
	"path"
	"runtime"
	"strconv"
	//    "time"
)

var logflag bool

func getGID() uint64 { // using this will "send us straight to hell"
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

//we want our log and error messages to stand out, so we color code them. These constants help with that.
const (
	GREEN  = "Green"
	RESET  = "Reset"
	RED    = "Red"
	YELLOW = "Yellow"
	BLUE   = "Blue"

	ANSIRESET  string = "\x1b[0m"
	ANSIRED    string = "\x1b[31m"
	ANSIGREEN  string = "\x1b[92m"
	ANSIYELLOW string = "\x1b[33m"
	ANSIBLUE   string = "\u001b[34m"

	TRACEFILE_PATH   = "tracefile.txt"
	CONSOLEFILE_PATH = "consolefile.txt"
)

var (
	Indent_Level  []int    = make([]int, 100)
	traceFile     *os.File = os.Stderr
	traceWriter   *bufio.Writer
	consoleFile   *os.File = os.Stdout
	consoleWriter *bufio.Writer
)

func InitLogging(ec chan error) {
	var err error
	traceFile, err = os.Create(TRACEFILE_PATH)
	defer traceFile.Close()
	if err != nil {
		fmt.Printf("could not open trace file\n")
		ec <- err
		return
	}
	traceWriter = bufio.NewWriter(traceFile)
	consoleFile, err = os.Create(CONSOLEFILE_PATH)
	defer consoleFile.Close()
	if err != nil {
		fmt.Printf("could not open trace file\n")
		ec <- err
		return
	}
	consoleWriter = bufio.NewWriter(traceFile)

	ec <- nil

	for ; true; time.Sleep(time.Second) {
		//a broken flush is the only way for logging to fail, then, right?
		err = traceWriter.Flush()
		if err != nil {
			Info_Log(err)
			return
		}
		err = consoleWriter.Flush()
		if err != nil {
			Info_Log(err)
			return
		}
	}
	return
	//this function is designed to run asynchronosly, for the life of the program, for logging purposes
}

func indent(depth int) string {
	accum := ""

	for i := 0; i <= depth; i++ {
		accum += "    "
	}
	return accum
}

func LogEnter() {
	gid := getGID()
	if gid > uint64(len(Indent_Level)) {
		Indent_Level = append(Indent_Level, make([]int, 4*(gid-uint64(len(Indent_Level))))...)
	}
	ilp := &Indent_Level[gid]
	thingtoprint := "Entered"
	pc, filepath, line, _ := runtime.Caller(1)
	file := path.Base(filepath)
	details := runtime.FuncForPC(pc)
	fmt.Fprintf(traceFile, "\t%-12.12s %5d\t%-72.72s %4d(gr) %4d(d): %s %v\n", file, line, fmt.Sprintf("%s()", details.Name()), gid, *ilp, indent(*ilp), thingtoprint)
	(*ilp) += 1
}

func LogExit() {
	gid := getGID()
	ilp := &Indent_Level[gid]
	(*ilp) -= 1
	thingtoprint := "Exited"
	pc, filepath, line, _ := runtime.Caller(1)
	file := path.Base(filepath)
	details := runtime.FuncForPC(pc)
	fmt.Fprintf(traceFile, "\t%-12.12s %5d\t%-72.72s %4d(gr) %4d(d): %s %v\n", file, line, fmt.Sprintf("%s()", details.Name()), gid, *ilp, indent(*ilp), thingtoprint)
}

//In SetConsoleColor, we change the console color using this map as a lookup table.
var colormap map[string]string = map[string]string{RESET: ANSIRESET, RED: ANSIRED, GREEN: ANSIGREEN, YELLOW: ANSIYELLOW, BLUE: ANSIBLUE}

//We sometimes exit after errors. When we do, we call this function. Error messages are red here. These will be the last output from the program.
func Errhandle_Exit(err error, reason string) {
	fmt.Printf("%s:", reason)
	_, file, line, _ := runtime.Caller(1)
	if err != nil {
		SetConsoleColor(RED)
		fmt.Fprintf(traceFile, "\t%s %d\t failed: %v\n", file, line, err)
		SetConsoleColor(RESET)
		os.Exit(1)
	} else {
		SetConsoleColor(GREEN)
		fmt.Fprintf(traceFile, "\t %s %d \t successful.\n", file, line)
		SetConsoleColor(RESET)
	}
}

//We call this function for kind of trivial errors. It doesn't kill the program, error messages are yellow here.
func Errhandle_Log(err error, reason string) {
	if !logflag {
		return
	}
	pc, filepath, line, _ := runtime.Caller(1)
	file := path.Base(filepath)
	details := runtime.FuncForPC(pc)
	if err != nil {
		SetConsoleColor(RED)
		//fmt.Printf("\t%s %d\t failed: %v\n", file, line, err)
		outstring := fmt.Sprintf("\t%-12.12s %5d\t%-72.72s : Errhandle_Log(%s, %v)\n", file, line, fmt.Sprintf("%s()", details.Name()), reason, err)
		traceWriter.Write([]byte(outstring))
		SetConsoleColor(RESET)
	} else {
		SetConsoleColor(GREEN)
		//fmt.Printf("\t %s %d \t successful.\n", file, line)
		outstring := fmt.Sprintf("\t%-12.12s %5d\t%-72.72s : Errhandle_Log(%s, [not error])\n", file, line, fmt.Sprintf("%s()", details.Name()), reason)
		traceWriter.Write([]byte(outstring))
		SetConsoleColor(RESET)
	}
}

func Info_Log(thingtoprint interface{}, i ...interface{}) {

	if !logflag {
		return
	}
	switch thingtoprint.(type) {
	case string:
		spf_result := fmt.Sprintf(thingtoprint.(string), i...)

		pc, filepath, line, _ := runtime.Caller(1)
		file := path.Base(filepath)
		details := runtime.FuncForPC(pc)
		SetConsoleColor(YELLOW)
		outstring := fmt.Sprintf("\t%-12.12s %5d\t%-72.72s : %s\n", file, line, fmt.Sprintf("%s()", details.Name()), spf_result)
		traceWriter.Write([]byte(outstring))
		SetConsoleColor(RESET)

	default:
		pc, filepath, line, _ := runtime.Caller(1)
		file := path.Base(filepath)
		details := runtime.FuncForPC(pc)
		SetConsoleColor(YELLOW)
		outstring := fmt.Sprintf("\t%-12.12s %5d\t%-72.72s : %v\n", file, line, fmt.Sprintf("%s()", details.Name()), thingtoprint)
		traceWriter.Write([]byte(outstring))
		SetConsoleColor(RESET)
	}
}

//Sets consolo color according to the string:string map above.
func SetConsoleColor(color string) {
	for key, value := range colormap {
		if key == color {
			outstring := fmt.Sprintf("%s", value)
			traceWriter.Write([]byte(outstring))
		}
	}
}
