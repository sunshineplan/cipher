package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sunshineplan/utils/ste"
	"golang.org/x/crypto/ssh/terminal"
)

func menu() (choice string) {
	fmt.Print(`   * * * * * * * * * * * * * * * * * * * * * *
   *                                         *
   *  Welcome to use Simple Text Encryption  *
   *                                         *
   *     1. Encrypt                          *
   *     2. Decrypt                          *
   *     Q. Quit                             *
   *                                         *
   * * * * * * * * * * * * * * * * * * * * * *

Please choose one function: `)
	fmt.Scanln(&choice)
	return
}

func clear() {
	c := make(map[string]func())
	c["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	c["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	c[OS]()
}

func multilineInput() string {
	scn := bufio.NewScanner(os.Stdin)
	var lines []string
	for scn.Scan() {
		line := scn.Text()
		if line == "EOC" {
			break
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func client() {
	var content string
	var key []byte
	for {
		clear()
		switch menu() {
		case "1":
			fmt.Println("\nMultiline Content(end with EOC):")
			content = multilineInput()
			fmt.Print("\nKey: ")
			key, _ = terminal.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println("\nEncrypted Content:")
			fmt.Println(ste.Encrypt(string(key), content))
			fmt.Println("\nPress enter to continue...")
			terminal.ReadPassword(int(os.Stdin.Fd()))
		case "2":
			fmt.Println("\nContent: ")
			fmt.Scanln(&content)
			fmt.Print("\nKey: ")
			key, _ = terminal.ReadPassword(int(os.Stdin.Fd()))
			ct, err := ste.Decrypt(string(key), content)
			if err != nil {
				fmt.Println("\nEmpty or Malformed content!")
			} else {
				fmt.Println("\nDecrypted Content:")
				fmt.Println(ct)
			}
			fmt.Println("\nPress enter to continue...")
			terminal.ReadPassword(int(os.Stdin.Fd()))
		case "Q", "q":
			clear()
			return
		default:
			fmt.Println("\nWrong choice! Press enter to continue...")
			terminal.ReadPassword(int(os.Stdin.Fd()))
		}
	}
}
