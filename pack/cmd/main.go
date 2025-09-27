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
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	fmt.Println(result)
}

func runDaemon(escapeEnabled bool, pack bool) {
	mode := "распаковки"
	if pack {
		mode = "упаковки"
	}

	fmt.Printf("Запуск в режиме демона (%s). Ctrl+C для завершения\n", mode)
	if escapeEnabled {
		fmt.Println("Режим экранирования включен")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\nПолучен сигнал завершения. Выход...")
		os.Exit(0)
	}()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Введите строку: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Ошибка чтения ввода: %v\n", err)
			continue
		}

		input = strings.TrimSuffix(input, "\n")
		processInput(input, escapeEnabled, pack)
	}
}

func showUsage() {
	progName := filepath.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, "Использование:\n")
	fmt.Fprintf(os.Stderr, "  %s --pack [--input <строка>] [--daemon] [--escape]\n", progName)
	fmt.Fprintf(os.Stderr, "  %s --unpack [--input <строка>] [--daemon] [--escape]\n", progName)
	fmt.Fprintf(os.Stderr, "\nФлаги:\n")
	fmt.Fprintf(os.Stderr, "  --pack      Режим упаковки строк\n")
	fmt.Fprintf(os.Stderr, "  --unpack    Режим распаковки строк\n")
	fmt.Fprintf(os.Stderr, "  --input     Строка для обработки\n")
	fmt.Fprintf(os.Stderr, "  --daemon    Запуск в интерактивном режиме\n")
	fmt.Fprintf(os.Stderr, "  --escape    Включить поддержку экранирования\n")
	fmt.Fprintf(os.Stderr, "\nПримеры упаковки:\n")
	fmt.Fprintf(os.Stderr, "  %s --pack --input 'aaaabccddddde'\n", progName)
	fmt.Fprintf(os.Stderr, "  %s --pack --daemon\n", progName)
	fmt.Fprintf(os.Stderr, "\nПримеры распаковки:\n")
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
		inputFlag  = flag.String("input", "", "Строка для обработки")
		daemonFlag = flag.Bool("daemon", false, "Запуск в интерактивном режиме")
		escapeFlag = flag.Bool("escape", false, "Включить поддержку экранирования")
		packFlag   = flag.Bool("pack", false, "Режим упаковки строк")
		unpackFlag = flag.Bool("unpack", false, "Режим распаковки строк")
		helpFlag   = flag.Bool("help", false, "Показать справку")
	)

	flag.Usage = showUsage
	flag.Parse()

	if *helpFlag {
		showUsage()
		return
	}

	if !*packFlag && !*unpackFlag {
		fmt.Fprintf(os.Stderr, "Ошибка: необходимо указать режим --pack или --unpack\n\n")
		showUsage()
		os.Exit(1)
	}

	if *packFlag && *unpackFlag {
		fmt.Fprintf(os.Stderr, "Ошибка: нельзя одновременно использовать --pack и --unpack\n\n")
		showUsage()
		os.Exit(1)
	}

	pack := *packFlag

	if *inputFlag != "" && *daemonFlag {
		fmt.Fprintf(os.Stderr, "Ошибка: нельзя одновременно использовать --input и --daemon\n\n")
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
		fmt.Fprintf(os.Stderr, "Ошибка: не указан режим работы\n\n")
		showUsage()
		os.Exit(1)
	}
}
