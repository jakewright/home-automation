package output

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/logrusorgru/aurora"
)

const (
	checkMark    = "\xE2\x9C\x94"
	heavyBallotX = "\xE2\x9C\x98"
)

var (
	// Verbose can be set to show debug lines
	Verbose   bool
	currentOp *Operation
)

// Operation represents something currently happening that can
// be completed, failed or abandoned.
type Operation struct {
	completed bool
}

func newOp() *Operation {
	// If an operation has been abandoned, we need a new line
	if currentOp != nil {
		currentOp.Abandon()
	}
	currentOp = &Operation{}
	return currentOp
}

// Success prints the result of running fmt.Sprintf on args and
// a new line. If args is empty then a default check mark is used.
func (o *Operation) Success(args ...interface{}) {
	if o.completed || o != currentOp {
		return
	}

	end := fmt.Sprintf(" %s", aurora.Green(checkMark))

	switch len(args) {
	case 0:
	case 1:
		end = fmt.Sprintf("%s", args[0])
	default:
		end = fmt.Sprintf(fmt.Sprintf("%s", args[0]), args[1:]...)
	}

	fmt.Printf("%s\n", end)

	o.completed = true
	currentOp = nil
}

// Failed prints a heavy ballot X and a new line
func (o *Operation) Failed() {
	if o.completed || o != currentOp {
		return
	}

	fmt.Printf(" %s\n", aurora.Red(heavyBallotX))
	o.completed = true
	currentOp = nil
}

// Abandon just prints a new line without a symbol
func (o *Operation) Abandon() {
	if o.completed || o != currentOp {
		return
	}

	fmt.Printf("\n")
	o.completed = true
	currentOp = nil
}

// Info prints the string and returns an Operation that can
// be used to control whether a check mark or ballot x is placed
// at the end of the line.
func Info(format string, a ...interface{}) *Operation {
	op := newOp()
	fmt.Printf(format, a...)
	return op
}

// InfoLn prints the string on its own line without
// returning an Operation.
func InfoLn(format string, a ...interface{}) {
	Info(format, a...).Abandon()
}

// Debug will output only if the Verbose flag is set.
func Debug(format string, a ...interface{}) {
	if Verbose {
		if currentOp != nil {
			newOp().Abandon() // Print a blank line
		}

		f := fmt.Sprintf(format, a...)
		fmt.Printf("%s\n", aurora.Gray(12, f))
	}
}

// Fatal exists with status code 1 after printing the line in red
func Fatal(format string, a ...interface{}) {
	if currentOp != nil {
		newOp().Abandon() // Print a blank line
	}

	fmt.Printf(aurora.Sprintf(aurora.Red(format), a...))
	fmt.Printf("\n")
	os.Exit(1)
}

// Confirm will prompt the user to enter yes or no, returning the result.
func Confirm(def bool, format string, a ...interface{}) bool {
	newOp().Abandon()      // Print blank line
	defer fmt.Printf("\n") // Print another blank line at the end

	reader := bufio.NewReader(os.Stdin)
	msg := fmt.Sprintf(format, a...)

	yes, no := "y", "n"

	if def {
		yes = strings.ToUpper(yes)
	} else {
		no = strings.ToUpper(no)
	}

	opts := aurora.Sprintf(aurora.Gray(16, "(%s/%s)"), yes, no)

	for {
		fmt.Printf("%s %s ", msg, opts)

		input, _ := reader.ReadString('\n')
		input = strings.ToLower(strings.TrimSpace(input))

		switch input {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		case "":
			return def
		}
	}
}
