package parser

import (
	"errors"
	"strconv"
	"strings"
)

type Sign int

const (
	Negative Sign = iota
	NonNegative
)

type Parity int

const (
	Even Parity = iota
	Odd
)

type Result struct {
	N      int
	Parity Parity
	Sign   Sign
}

// ParseAndClassify парсит s в int (trim пробелы) и возвращает результат.
// Если парсинг не удался, возвращается ошибка (оригинальная ошибка strconv).
func ParseAndClassify(s string) (Result, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Result{}, errors.New("empty input")
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return Result{}, err
	}

	var p Parity
	if n%2 == 0 {
		p = Even
	} else {
		p = Odd
	}

	var sg Sign
	if n < 0 {
		sg = Negative
	} else {
		sg = NonNegative
	}

	return Result{
		N:      n,
		Parity: p,
		Sign:   sg,
	}, nil
}
