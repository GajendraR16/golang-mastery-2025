package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Sale struct {
	Category string
	Amount   int
}

func parseSale(line string) (Sale, error) {
	parts := strings.Split(line, ",")
	if len(parts) != 2 {
		return Sale{}, fmt.Errorf("bad line %q: want category, amount", line)
	}

	amount, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return Sale{}, fmt.Errorf("bad amount in %q: %w", amount, err)
	}
	return Sale{Category: strings.TrimSpace(parts[0]), Amount: amount}, nil
}

func main() {

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage : sales <filename>")
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	totals := map[string]int{}
	for scanner.Scan() {
		line := scanner.Text()
		sale, err := parseSale(line)
		if err != nil {
			fmt.Println("skipping:", err)
			continue
		}
		totals[strings.ToLower(sale.Category)] += sale.Amount
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "read error", err)
		os.Exit(1)
	}

	for category, total := range totals {
		fmt.Printf("%s: %d\n", category, total)
	}
}
