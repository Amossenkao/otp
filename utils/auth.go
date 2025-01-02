package utils

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/term"
)

func Authenticate(password string) bool {

	fmt.Print("Enter your password: ")

	// Put the terminal into raw mode
	oldState, err := term.MakeRaw(int(syscall.Stdin))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error enabling raw mode: %v\n", err)
		return false
	}

	defer term.Restore(int(syscall.Stdin), oldState)

	var password_provided []byte
	for {
		// Read a single byte from stdin
		var buf [1]byte
		_, err := os.Stdin.Read(buf[:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			return false
		}

		if buf[0] == '\n' || buf[0] == '\r' {
			break
		}

		if buf[0] == 127 || buf[0] == 8 {
			if len(password_provided) > 0 {
				password_provided = password_provided[:len(password_provided)-1]
				fmt.Print("\b \b")
			}
			continue
		}
		password_provided = append(password_provided, buf[0])

		// Print the masking character
		fmt.Print("*")
	}

	fmt.Println()
	return string(password_provided) == password
}

