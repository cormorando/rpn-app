package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
	"unicode"
	"unicode/utf8"
)

type RPCFunc uint8

func (*RPCFunc) Parse(arg *string, result *int) error {
	*result = calculate([]byte(*arg))

	return nil
}

func skipSpaces(s []byte) []byte {
	c, w := utf8.DecodeRune(s)
	for w > 0 && unicode.IsSpace(c) {
		s = s[w:]
		c, w = utf8.DecodeRune(s)
	}

	return s
}

func readDigits(s []byte) (numStr, remain []byte) {
	numStr = s
	totalW := 0
	c, w := utf8.DecodeRune(s)
	for w > 0 && unicode.IsDigit(c) {
		s = s[w:]
		totalW += w
		c, w = utf8.DecodeRune(s)
	}

	return numStr[:totalW], s
}

func pop(stack []int) (int, []int) {

	return stack[len(stack)-1], stack[:len(stack)-1]
}

func calculate(s []byte) int {
	stack := make([]int, 0)
	var a, b int
	var token []byte

	s = skipSpaces(s)
	for len(s) > 0 {
		c, w := utf8.DecodeRune(s)
		switch {
		case unicode.IsDigit(c):
			token, s = readDigits(s)
			num, err := strconv.Atoi(string(token))
			if err != nil {
				fmt.Println(err)
			} else {
				stack = append(stack, num)
			}
		case c == '+':
			b, stack = pop(stack)
			a, stack = pop(stack)
			stack = append(stack, a+b)
			s = s[w:]
		case c == '-':
			b, stack = pop(stack)
			a, stack = pop(stack)
			stack = append(stack, a-b)
			s = s[w:]
		case c == '*':
			b, stack = pop(stack)
			a, stack = pop(stack)
			stack = append(stack, a*b)
			s = s[w:]
		case c == '/':
			b, stack = pop(stack)
			a, stack = pop(stack)
			stack = append(stack, a/b)
			s = s[w:]
		case c == '%':
			b, stack = pop(stack)
			a, stack = pop(stack)
			stack = append(stack, a%b)
			s = s[w:]
		default:
			fmt.Println("unknown character:", c)
			s = s[w:]
		}
		s = skipSpaces(s)
	}

	return stack[0]
}

func main() {
	log.Print("starting server")
	l, err := net.Listen("tcp", "localhost:1234")
	defer l.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("listening on: ", l.Addr())
	rpc.Register(new(RPCFunc))
	for {
		log.Print("waiting for connections ...")
		conn, err := l.Accept()
		if err != nil {
			log.Printf("accept error: %s", conn)
			continue
		}
		log.Printf("connection started: %v", conn.RemoteAddr())
		go jsonrpc.ServeConn(conn)
		log.Printf("connection ended: %v", conn.RemoteAddr())
	}
}
