package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var part = flag.Int("p", 1, "-p [12]")

func main() {
	flag.Parse()
	buf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	strnums := strings.Fields(string(buf))
	var fn stateFn
	switch *part {
	case 1:
		fn = part1
	case 2:
		fn = part2
	}
	program := &state{}
	data := gen(strnums)
	for fn != nil {
		fn = fn(program, data)
	}
	fmt.Println(program.answer)
}
func gen(nums []string) <-chan int {
	c := make(chan int)
	go func() {
		for _, nstr := range nums {
			num, err := strconv.Atoi(nstr)
			if err != nil {
				panic(err)
			}
			c <- num
		}
		c <- -1 // terminal
		close(c)
	}()
	return c
}

type state struct {
	state  int
	answer int
	buffer []int
}

type stateFn func(*state, <-chan int) stateFn

func part1(s *state, c <-chan int) stateFn {
	n := <-c
	if n < 0 {
		return nil
	}
	switch s.state {
	case 0, 1:
		s.buffer = append(s.buffer, n)
		s.state = (s.state + 1) % 3
	case 2:
		a := s.buffer[0]
		b := s.buffer[1]
		s.answer += triangle(a, b, n)
		s.buffer = s.buffer[:0]
		s.state = 0
	}
	return part1
}
func part2(s *state, c <-chan int) stateFn {
	n := <-c
	if n < 0 {
		return nil
	}
	switch s.state {
	case 0, 1, 2, 3, 4, 5, 6, 7:
		s.buffer = append(s.buffer, n)
		s.state = (s.state + 1) % 9
	case 8:
		a0 := s.buffer[0]
		a1 := s.buffer[1]
		a2 := s.buffer[2]
		b0 := s.buffer[3]
		b1 := s.buffer[4]
		b2 := s.buffer[5]
		c0 := s.buffer[6]
		c1 := s.buffer[7]
		s.answer += triangle(a0, b0, c0)
		s.answer += triangle(a1, b1, c1)
		s.answer += triangle(a2, b2, n)
		s.buffer = s.buffer[:0]
		s.state = 0
	}
	return part2
}

func triangle(a, b, c int) int {
	if a+b > c && a+c > b && b+c > a {
		return 1
	}
	return 0
}
