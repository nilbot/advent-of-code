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
	strs := strings.Split(string(buf), ", ")
	strs = append(strs, "O")
	source := gen(strs)
	data := &state{table: make(map[string]int)}
	var fn stateFn
	switch *part {
	case 1:
		fn = part1
	case 2:
		fn = part2
	}
	for fn != nil {
		fn = fn(data, source)
	}
	fmt.Println(data)
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
	facing     direction
	total      int
	table      map[string]int
}

type direction int

const (
	north direction = iota
	east
	south
	west
)

type stateFn func(*state, <-chan string) stateFn

func part1(s *state, c <-chan string) stateFn {
	command := <-c
	if command == "O" {
		return nil
	}
	dir := command[0]
	switch dir {
	case 'R':
		s.facing = (s.facing + 1) % 4
	case 'L':
		s.facing = (s.facing + 3) % 4
	}
	rest := command[1:]
	var howfar int
	fmt.Sscanf(rest, "%d", &howfar)
	switch s.facing {
	case north:
		s.vertical += howfar
	case south:
		s.vertical -= howfar
	case east:
		s.horizontal += howfar
	case west:
		s.horizontal -= howfar
	}
	s.total = abs(s.horizontal) + abs(s.vertical)
	return part1
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func (s state) String() string {
	return fmt.Sprintf("%d", s.total)
}

func part2(s *state, c <-chan string) stateFn {
	command := <-c
	if command == "O" {
		return nil
	}
	dir := command[0]
	switch dir {
	case 'R':
		s.facing = (s.facing + 1) % 4
	case 'L':
		s.facing = (s.facing + 3) % 4
	}
	rest := command[1:]
	var howfar int
	fmt.Sscanf(rest, "%d", &howfar)
	for i := 0; i < howfar; i++ {
		switch s.facing {
		case north:
			s.vertical++
		case south:
			s.vertical--
		case east:
			s.horizontal++
		case west:
			s.horizontal--
		}
		coordinate := fmt.Sprintf("%dx%d", s.horizontal, s.vertical)
		if _, ok := s.table[coordinate]; !ok {
			s.table[coordinate] = 1
		} else {
			s.total = abs(s.horizontal) + abs(s.vertical)
			return nil
		}
	}
	return part2
}
