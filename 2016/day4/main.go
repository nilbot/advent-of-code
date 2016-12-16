package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"strconv"
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
func gen(strs []string) <-chan string {
	c := make(chan string)
	go func() {
		for _, str := range strs {
			c <- str
		}
		c <- "QED" // terminal
		close(c)
	}()
	return c
}

type alphabet byte

const (
	la alphabet = alphabet('a')
	lz alphabet = alphabet('z')
	ua alphabet = alphabet('A')
	uz alphabet = alphabet('Z')
)

type numeric byte

const (
	one  numeric = numeric('0')
	nine numeric = numeric('9')
)

type separator byte

const (
	dash separator = separator('-')
	lb   separator = separator('[')
	rb   separator = separator(']')
)

type state struct {
	input     string
	histogram map[byte]int
	pos       int
	sectorId  int
	lastCount int
	lastChar  byte

	answer int
}

type stateFn func(*state, <-chan string) stateFn
type miniStateFn func(*state) miniStateFn

func part1(s *state, c <-chan string) stateFn {
	str := <-c
	if str == "QED" {
		return nil
	}
	s.input = str
	s.histogram = make(map[byte]int)
	s.pos = 0
	s.lastChar = 0
	s.lastCount = len(str) + 1
	s.sectorId = 0
	fun := lexName // new start
	for fun != nil {
		fun = fun(s)
	}
	return part1
}

func lexName(s *state) miniStateFn {
	b := s.input[s.pos]
	if isAlphabet(b) {
		if num, ok := s.histogram[b]; ok {
			s.histogram[b] = num + 1
		} else {
			s.histogram[b] = 1
		}
		s.pos++
		return lexName
	}
	if separator(b) == dash {
		s.pos++
		return lexName
	}
	if isNumeric(b) {
		return lexNumeric
	}
	panic(fmt.Sprintf("state lexName error, input was: %v, pos was: %v, "+
		"byte as string: %v, byte raw value: %v",
		s.input, s.pos, string(b), b))
}

func lexNumeric(s *state) miniStateFn {
	b := s.input[s.pos]
	if separator(b) == lb {
		s.pos++
		return lexChecksum
	}
	if isNumeric(b) {
		i := int(numeric(b) - one)
		s.sectorId = s.sectorId*10 + i
		s.pos++
		return lexNumeric
	}
	panic(fmt.Sprintf("state lexNumeric error, input was: %v, pos was: "+
		"%v, byte as string: %v",
		s.input, s.pos, string(s.input[s.pos])))
}

func lexChecksum(s *state) miniStateFn {
	b := s.input[s.pos]
	if isAlphabet(b) {
		if _, ok := s.histogram[b]; !ok {
			// totally decoy
			return nil
		}
		cnt := s.histogram[b]
		if cnt > s.lastCount {
			// false, greedy exit
			return nil
		}
		if cnt == s.lastCount && b < s.lastChar {
			// false, greedy exit
			return nil
		}
		s.lastCount = cnt
		s.lastChar = b
		s.pos++
		return lexChecksum
	}
	if separator(b) == rb {
		s.answer += s.sectorId
	}
	return nil
}


func isAlphabet(b byte) bool {
	if (alphabet(b) >= la && alphabet(b) <= lz) || (alphabet(b) >= ua &&
		alphabet(b) <= uz) {
		return true
	}
	return false
}

func isNumeric(b byte) bool {
	if numeric(b) >= one && numeric(b) <= nine {
		return true
	}
	return false
}


func part2(s *state, c <-chan string) stateFn {
	str := <-c
	if str == "QED" {
		return nil
	}

	for i:=1;i<=25;i++ {
		decrypted := caesar(str, i)
		if strings.Contains(decrypted, "north") {
			s.answer = getSectorId(decrypted)
			return nil
		}
	}
	return part2
}


func caesar(s string, n int) string {
	res := ""
	// brutally specific
	for i := 0; i < len(s)-7; i++ {
		b := s[i]
		if b == '-' {
			res += " "
			continue
		}
		if isAlphabet(b) {
			res += string(shift(b, n))
			continue
		}
		res += string(b)
	}
	return res
}

func shift(b byte, n int) byte {
	if alphabet(b) <= lz {
		return ((b - byte(la) + byte(n)) % 26) + byte(la)
	}
	return ((b - byte(ua) + byte(n)) % 26) + byte(ua)
}

func getSectorId(s string) int {
	strs := strings.Fields(s)
	i,e:=strconv.Atoi(strs[len(strs)-1])
	if e != nil {
		panic(e)
	}
	return i
}