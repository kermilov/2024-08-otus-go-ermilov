package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "Время ожидания соединения")
	flag.Parse()

	if flag.NArg() != 2 {
		log.Fatal("Количество аргументов не соответствует паттерну: go-telnet [--timeout=<timeout>] <host> <port>")
	}

	host := flag.Arg(0)
	port := flag.Arg(1)

	address := net.JoinHostPort(host, port)

	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(os.Stderr, "Установлено соединение", address)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := client.Send(); err != nil {
			log.Println(err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := client.Receive(); err != nil {
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				fmt.Fprintln(os.Stderr, "Время ожидания соединения истекло")
			} else {
				fmt.Fprintln(os.Stderr, "Соединение было закрыто")
			}
		}
	}()

	wg.Wait()

	if err := client.Close(); err != nil {
		log.Println(err)
	}
}
