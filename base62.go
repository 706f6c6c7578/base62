package main

import (
	"flag"
	"fmt"
	"io"
	"math/big"
	"os"
	"strings"
)

const DICTIONARY = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func Encode(data []byte) string {
	bi := new(big.Int).SetBytes(data)
	encoded := ""
	base := big.NewInt(62)
	zero := big.NewInt(0)
	mod := new(big.Int)

	for bi.Cmp(zero) > 0 {
		bi.DivMod(bi, base, mod)
		encoded = string(DICTIONARY[mod.Int64()]) + encoded
	}
	if encoded == "" {
		encoded = "0"
	}
	return encoded
}

func Decode(encoded string) ([]byte, error) {
	bi := new(big.Int)
	base := big.NewInt(62)
	for _, c := range encoded {
		index := strings.IndexRune(DICTIONARY, c)
		if index < 0 {
			return nil, fmt.Errorf("invalid character in input: %c", c)
		}
		bi.Mul(bi, base)
		bi.Add(bi, big.NewInt(int64(index)))
	}
	return bi.Bytes(), nil
}

func wrapText(text string, width int) string {
	if width <= 0 {
		return text
	}
	var result strings.Builder
	for i, char := range text {
		if i > 0 && i%width == 0 {
			result.WriteString("\n")
		}
		result.WriteRune(char)
	}
	return result.String()
}

func main() {
	decode := flag.Bool("d", false, "Decode mode")
	wrap := flag.Int("w", 0, "Wrap lines after specified number of characters (0 for no wrapping)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-d] [-w width] [input]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  If no input is provided, the program reads from stdin.\n")
		fmt.Fprintf(os.Stderr, "  Use -d flag for decode mode, otherwise encode mode is used.\n")
		fmt.Fprintf(os.Stderr, "  Use -w flag to wrap output (or unwrap input when decoding).\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	var input []byte
	var err error
	if flag.NArg() > 0 {
		input = []byte(flag.Arg(0))
	} else {
		input, err = io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}
	}

	if *decode {
		// Remove newlines and carriage returns when decoding
		cleanInput := strings.ReplaceAll(string(input), "\n", "")
		cleanInput = strings.ReplaceAll(cleanInput, "\r", "")
		decoded, err := Decode(strings.TrimSpace(cleanInput))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error decoding: %v\n", err)
			os.Exit(1)
		}
		os.Stdout.Write(decoded)
	} else {
		encoded := Encode(input)
		if *wrap > 0 {
			encoded = wrapText(encoded, *wrap)
		}
		fmt.Println(encoded)
	}
}