package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	packtool "pack/pack"
	"path/filepath"
	"strings"
	"syscall"
)

func processInput(input string, escapeEnabled bool, pack bool) {
	var result string
	var err error

	if pack {
		result, err = packtool.PackString(input)
	} else {
		result, err = packtool.UnpackString(input, escapeEnabled)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println(result)
}

func runDaemon(escapeEnabled bool, pack bool) {
	mode := "unpacking"
	if pack {
		mode = "packing"
	}

	fmt.Printf("Running in daemon mode (%s). Ctrl+C to exit\n", mode)
	if escapeEnabled {
		fmt.Println("Escape mode enabled")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\nReceived termination signal. Exiting...")
		os.Exit(0)
	}()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter string: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}

		input = strings.TrimSuffix(input, "\n")
		processInput(input, escapeEnabled, pack)
	}
}

func showUsage() {
	progName := filepath.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  %s --pack [--input <string>] [--daemon] [--escape]\n", progName)
	fmt.Fprintf(os.Stderr, "  %s --unpack [--input <string>] [--daemon] [--escape]\n", progName)
	fmt.Fprintf(os.Stderr, "\nFlags:\n")
	fmt.Fprintf(os.Stderr, "  --pack      String packing mode\n")
	fmt.Fprintf(os.Stderr, "  --unpack    String unpacking mode\n")
	fmt.Fprintf(os.Stderr, "  --input     String to process\n")
	fmt.Fprintf(os.Stderr, "  --daemon    Run in interactive mode\n")
	fmt.Fprintf(os.Stderr, "  --escape    Enable escape support\n")
	fmt.Fprintf(os.Stderr, "\nPacking examples:\n")
	fmt.Fprintf(os.Stderr, "  %s --pack --input 'aaaabccddddde'\n", progName)
	fmt.Fprintf(os.Stderr, "  %s --pack --daemon\n", progName)
	fmt.Fprintf(os.Stderr, "\nUnpacking examples:\n")
	fmt.Fprintf(os.Stderr, "  %s --unpack --input 'a4bc2d5e'\n", progName)
	fmt.Fprintf(os.Stderr, "  %s --unpack --input 'qwe\\4\\5' --escape\n", progName)
	fmt.Fprintf(os.Stderr, "  %s --unpack --daemon --escape\n", progName)
}

// go run main.go --help

// go run main.go --unpack --input 'a4bc2d5e'
// go run main.go --unpack --input $'d\n5abc' --escape
// go run main.go --unpack --input 'qwe\4\5' --escape
// go run main.go --unpack --daemon
// go run main.go --unpack --daemon --escape

// go run main.go --pack --input 'aaaabccddddde'
func main() {
	var (
		inputFlag  = flag.String("input", "", "String to process")
		daemonFlag = flag.Bool("daemon", false, "Run in interactive mode")
		escapeFlag = flag.Bool("escape", false, "Enable escape support")
		packFlag   = flag.Bool("pack", false, "String packing mode")
		unpackFlag = flag.Bool("unpack", false, "String unpacking mode")
		helpFlag   = flag.Bool("help", false, "Show help")
	)

	flag.Usage = showUsage
	flag.Parse()

	if *helpFlag {
		showUsage()
		return
	}

	if !*packFlag && !*unpackFlag {
		fmt.Fprintf(os.Stderr, "Error: must specify either --pack or --unpack mode\n\n")
		showUsage()
		os.Exit(1)
	}

	if *packFlag && *unpackFlag {
		fmt.Fprintf(os.Stderr, "Error: cannot use --pack and --unpack simultaneously\n\n")
		showUsage()
		os.Exit(1)
	}

	pack := *packFlag

	if *inputFlag != "" && *daemonFlag {
		fmt.Fprintf(os.Stderr, "Error: cannot use --input and --daemon simultaneously\n\n")
		showUsage()
		os.Exit(1)
	}

	if *inputFlag != "" {
		processInput(*inputFlag, *escapeFlag, pack)
		return
	}

	if *daemonFlag {
		runDaemon(*escapeFlag, pack)
		return
	}

	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "Error: no operating mode specified\n\n")
		showUsage()
		os.Exit(1)
	}
}
