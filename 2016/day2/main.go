package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var part = flag.Int("p", 1, "-p [12]")

func main() {
	flag.Parse()
	buf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	strs := strings.Split(string(buf), "\n")
	strs = append(strs, "O")
	source := gen(strs)
	data := &state{}
	var fn stateFn
	switch *part {
	case 1:
		fn = part1
	case 2:
		data.horizontal = -2
		fn = part2
	}
	for fn != nil {
		fn = fn(data, source)
	}
	fmt.Println()
}

func gen(src []string) <-chan string {
	res := make(chan string)
	go func() {
		for _, s := range src {
			res <- s
		}
		close(res)
	}()
	return res
}

type state struct {
	vertical   int
	horizontal int
	key        string
	num        int
	// table map[string]int
}

type stateFn func(*state, <-chan string) stateFn

func part1(s *state, c <-chan string) stateFn {
	command := <-c
	if command == "O" {
		return nil
	}
	for _, v := range command {
		switch v {
		case 'R':
			if s.horizontal < 1 {
				s.horizontal++
			}
		case 'L':
			if s.horizontal > -1 {
				s.horizontal--
			}
		case 'U':
			if s.vertical < 1 {
				s.vertical++
			}
		case 'D':
			if s.vertical > -1 {
				s.vertical--
			}
		}
	}
	if s.vertical == -1 && s.horizontal == -1 {
		s.num = 7
	}
	if s.vertical == -1 && s.horizontal == 0 {
		s.num = 8
	}
	if s.vertical == -1 && s.horizontal == 1 {
		s.num = 9
	}
	if s.vertical == 0 && s.horizontal == -1 {
		s.num = 4
	}
	if s.vertical == 0 && s.horizontal == 0 {
		s.num = 5
	}
	if s.vertical == 0 && s.horizontal == 1 {
		s.num = 6
	}
	if s.vertical == 1 && s.horizontal == -1 {
		s.num = 1
	}
	if s.vertical == 1 && s.horizontal == 0 {
		s.num = 2
	}
	if s.vertical == 1 && s.horizontal == 1 {
		s.num = 3
	}
	fmt.Printf("%d", s.num)
	return part1
}

func part2(s *state, c <-chan string) stateFn {
	command := <-c
	if command == "O" {
		return nil
	}
	for _, v := range command {
		switch v {
		case 'R':
			if s.vertical == 0 {
				if s.horizontal < 2 {
					s.horizontal++
				}
			}
			if abs(s.vertical) == 1 {
				if s.horizontal < 1 {
					s.horizontal++
				}
			}
		case 'L':
			if s.vertical == 0 {
				if s.horizontal > -2 {
					s.horizontal--
				}
			}
			if abs(s.vertical) == 1 {
				if s.horizontal > -1 {
					s.horizontal--
				}
			}
		case 'U':
			if s.horizontal == 0 {
				if s.vertical < 2 {
					s.vertical++
				}
			}
			if abs(s.horizontal) == 1 {
				if s.vertical < 1 {
					s.vertical++
				}
			}
		case 'D':
			if s.horizontal == 0 {
				if s.vertical > -2 {
					s.vertical--
				}
			}
			if abs(s.horizontal) == 1 {
				if s.vertical > -1 {
					s.vertical--
				}
			}
		}
	}
	switch s.horizontal {
	case 0:
		switch s.vertical {
		case 0:
			s.key = "7"
		case -1:
			s.key = "B"
		case -2:
			s.key = "D"
		case 1:
			s.key = "3"
		case 2:
			s.key = "1"
		}
	case -1:
		switch s.vertical {
		case 0:
			s.key = "6"
		case 1:
			s.key = "2"
		case -1:
			s.key = "A"
		}
	case -2:
		s.key = "5"
	case 1:
		switch s.vertical {
		case 1:
			s.key = "4"
		case 0:
			s.key = "8"
		case -1:
			s.key = "C"
		}
	case 2:
		s.key = "9"
	}
	fmt.Printf("%s", s.key)

	return part2
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
