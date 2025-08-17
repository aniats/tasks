package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"unicode"
)

var (
	ErrInvalidString = errors.New("некорректная строка")
)

const reverseSolidus = '\\'

// %s	the uninterpreted bytes of the string or slice
// %q	a double-quoted string safely escaped with Go syntax
// %x	base 16, lower-case, two characters per byte
// %X	base 16, upper-case, two characters per byte

func debugString(input string) {
	fmt.Printf("Входная строка: %s\n", input)
	fmt.Printf("Длина строки: %d\n", len(input))

	runes := []rune(input)
	fmt.Printf("Количество рун: %d\n", len(runes))

	for i, v := range runes {
		fmt.Printf("Позиция %d: %s (код: %d)\n", i, string(v), int(v))
	}
}

func unpackString(input string, escapeEnabled bool) (string, error) {
	if len(input) == 0 {
		return "", nil
	}

	runes := []rune(input)
	var result []rune

	for i := 0; i < len(runes); i++ {
		v := runes[i]

		if escapeEnabled && v == reverseSolidus {
			if i+1 >= len(runes) {
				return "", ErrInvalidString
			}

			nextChar := runes[i+1]
			if !unicode.IsDigit(nextChar) && nextChar != reverseSolidus {
				return "", ErrInvalidString
			}

			result = append(result, nextChar)
			i++
			continue
		}

		switch {
		case unicode.IsDigit(v) && len(result) == 0:
			return "", ErrInvalidString
		case unicode.IsDigit(v):
			if i+1 < len(runes) && unicode.IsDigit(runes[i+1]) {
				return "", ErrInvalidString
			}

			if v == '0' {
				if len(result) > 0 {
					result = result[:len(result)-1]
				}
			} else {
				number, err := strconv.Atoi(string(v))
				if err != nil {
					return "", ErrInvalidString
				}

				if len(result) > 0 {
					prevChar := result[len(result)-1]
					for j := 1; j < number; j++ {
						result = append(result, prevChar)
					}
				}
			}
		default:
			result = append(result, v)
		}
	}

	if escapeEnabled && len(runes) > 0 && runes[len(runes)-1] == reverseSolidus {
		return "", ErrInvalidString
	}

	return string(result), nil
}

func packString(input string) (string, error) {
	if len(input) == 0 {
		return "", nil
	}

	runes := []rune(input)
	var result []rune

	for i := 0; i < len(runes); {
		currentChar := runes[i]
		count := 1

		for i+count < len(runes) && runes[i+count] == currentChar {
			count++
		}

		result = append(result, currentChar)

		if count > 1 {
			result = append(result, rune('0'+count))
		}

		i += count
	}

	return string(result), nil
}

func processInput(input string, escapeEnabled bool, pack bool) {
	var result string
	var err error

	if pack {
		result, err = packString(input)
	} else {
		result, err = unpackString(input, escapeEnabled)
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
		debugString(*inputFlag)
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
