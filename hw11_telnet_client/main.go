package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
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

	go func() {
		if err := client.Send(); err != nil {
			log.Println(err)
		}
	}()

	endReceive := make(chan struct{}, 1)
	go func() {
		if err := client.Receive(); err != nil {
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				fmt.Fprintln(os.Stderr, "Время ожидания соединения истекло")
			} else {
				fmt.Fprintln(os.Stderr, "Соединение было закрыто")
			}
		}
		endReceive <- struct{}{}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-endReceive:
		fmt.Fprintln(os.Stderr, "Соединение было закрыто пользователем")
	case <-quit:
		fmt.Fprintln(os.Stderr, "Соединение было закрыто операционной системой")
	}

	if err := client.Close(); err != nil {
		log.Println(err)
	}
}
