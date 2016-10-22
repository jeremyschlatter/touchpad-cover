package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"path"
	"time"
)

var (
	events  = flag.String("dev-input-keyboard", "", "Path to keyboard event file in /dev/input. (example: \"/dev/input/by-path/platform-i8042-serio-0-event-kbd\")")
	pad     = flag.String("touchpad-name", "", "Name of touchpad as listed in \"xinput list\". (example: \"Atmel maXTouch Touchpad\")")
	timeout = flag.Duration("delay", 500*time.Millisecond, "How long to wait after last keypress to re-enable touchpad.")
	help    = flag.Bool("help", false, "Print this help.")
	verbose = flag.Bool("verbose", false, "Print information about what is going on. Use it to check that --dev-input-keyboard and --touchpad-name are correct.")
)

const (
	// Could get this from <linux/input.h>, but that adds one more thing to go wrong, and this size doesn't even need to be correct anyway.
	sizeofCInputEvent = 24
	eventBufSize      = sizeofCInputEvent * 64
)

func main() {
	flag.Usage = func() {
		fmt.Println("This program disables your touchpad while typing.\n")
		fmt.Printf("Usage:\n\t%s <flags>\n", path.Base(os.Args[0]))
		fmt.Println()
		fmt.Println("Flags:")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(1)
	}

	if *events == "" || *pad == "" {
		fmt.Println("Please set --dev-input-keyboard and --touchpad-name. (Pass --help for help).\n")
		os.Exit(1)
	}

	u, err := user.Current()
	if err != nil || u.Uid != "0" {
		fmt.Println("You are not root! This probably won't work.")
	}

	f, err := os.Open(*events)
	mustBeNil(err)
	defer f.Close()

	// Listen for keypresses.
	keypress := make(chan struct{})
	go func() {
		buf := make([]byte, eventBufSize)
		for {
			// Ignore number of bytes read. We don't care what keys were pressed, just that something was pressed.
			_, err = f.Read(buf)
			mustBeNil(err)

			// Send a keypress event.
			select {
			case keypress <- struct{}{}:
				if *verbose {
					fmt.Println("pressed key")
				}
			default:
			}
		}
	}()

	// If user ctrl-c's the program, we'll sure touchpad is enabled.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill)

	for {
		// Disable touchpad on keypress.
		select {
		case <-keypress:
		case <-interrupt:
			os.Exit(1)
		}
		touchpad(false)

		// Wait for typing to finish.
	Typing:
		for {
			select {
			case <-keypress:
			case <-interrupt:
				touchpad(true)
				os.Exit(1)
			case <-time.After(*timeout):
				break Typing
			}
		}

		// Re-enable touchpad.
		touchpad(true)
	}
}

func touchpad(on bool) {
	cmd, verb := "enable", "Enabling"
	if !on {
		cmd, verb = "disable", "Disabling"
	}
	if *verbose {
		fmt.Printf("-- %s touchpad: \"xinput %s '%s'\" --\n", verb, cmd, *pad)
	}
	if b, err := exec.Command("xinput", cmd, *pad).CombinedOutput(); err != nil {
		log.Fatalf("%s\n%v\n", b, err)
	}
}

func mustBeNil(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
