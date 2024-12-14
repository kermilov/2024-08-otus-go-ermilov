package main

import (
	"bytes"
	"io"
	"net"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})
}

// При нажатии Ctrl+D программа должна закрывать сокет и завершаться с сообщением.
func TestTelnetClientCtrlD(t *testing.T) {
	// Создаем мок-сервер
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	// Создаем буферы для ввода и вывода
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}

	// Создаем клиент
	client := NewTelnetClient(listener.Addr().String(), time.Second, io.NopCloser(in), out)
	err = client.Connect()
	require.NoError(t, err)

	// Запускаем горутину для приема соединения
	go func() {
		conn, err := listener.Accept()
		require.NoError(t, err)
		defer conn.Close()

		// Читаем данные из соединения
		_, err = io.Copy(io.Discard, conn)
		require.NoError(t, err)
	}()

	// Симулируем ввод данных
	in.WriteString("Hello\n")
	in.WriteString("World\n")

	// Симулируем EOF (Ctrl+D)
	in.Write([]byte{})

	// Запускаем Send в отдельной горутине
	done := make(chan struct{})
	go func() {
		err := client.Send()
		assert.NoError(t, err)
		close(done)
	}()

	// Ждем завершения Send или таймаута
	select {
	case <-done:
		// Send завершился успешно
	case <-time.After(5 * time.Second):
		t.Fatal("Send не завершился после EOF")
	}

	err = client.Close()
	assert.NoError(t, err)
}

// При получении SIGINT программа должна завершать свою работу.
func TestTelnetClientCtrlC(t *testing.T) {
	t.Skip(`
		Тест не проходит ни на винде (Received unexpected error: not supported by windows), ни на линуксе (signal: interrupt)
		Судя по коду func (p *Process) signal(sig Signal) error в пакете os из го вообще такой сигнал не послать
	`)
	// Создаем мок-сервер
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	// Создаем буферы для ввода и вывода
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}

	// Создаем клиент
	client := NewTelnetClient(listener.Addr().String(), time.Second, io.NopCloser(in), out)
	err = client.Connect()
	require.NoError(t, err)

	// Запускаем горутину для приема соединения
	go func() {
		conn, err := listener.Accept()
		require.NoError(t, err)
		defer conn.Close()

		// Читаем данные из соединения
		_, err = io.Copy(io.Discard, conn)
		require.NoError(t, err)
	}()

	// Запускаем Send и Receive в отдельных горутинах
	sendDone := make(chan struct{})
	receiveDone := make(chan struct{})

	go func() {
		err := client.Send()
		assert.NoError(t, err)
		close(sendDone)
	}()

	go func() {
		err := client.Receive()
		assert.NoError(t, err)
		close(receiveDone)
	}()

	// Даем немного времени для запуска горутин
	time.Sleep(100 * time.Millisecond)

	// Симулируем сигнал SIGINT (Ctrl+C)
	p, err := os.FindProcess(os.Getpid())
	require.NoError(t, err)
	err = p.Signal(syscall.SIGINT)
	require.NoError(t, err)

	// Ждем завершения Send и Receive или таймаута
	select {
	case <-sendDone:
		// Send завершился
	case <-time.After(5 * time.Second):
		t.Fatal("Send не завершился после SIGINT")
	}

	select {
	case <-receiveDone:
		// Receive завершился
	case <-time.After(5 * time.Second):
		t.Fatal("Receive не завершился после SIGINT")
	}

	err = client.Close()
	assert.NoError(t, err)
}

// Если сокет закрылся со стороны сервера, то при следующей попытке отправить сообщение программа
// должна завершаться (допускается завершать программу после "неудачной" отправки нескольких сообщений).
func TestTelnetClientSendAfterServerClose(t *testing.T) {
	// Создаем мок-сервер
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	// Создаем буферы для ввода и вывода
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}

	// Создаем клиент
	client := NewTelnetClient(listener.Addr().String(), time.Second, io.NopCloser(in), out)
	err = client.Connect()
	require.NoError(t, err)

	// Запускаем горутину для закрытия соединения
	go func() {
		conn, err := listener.Accept()
		require.NoError(t, err)
		conn.Close()
	}()

	// Посылаем сообщения некоторое время
	for i := 0; i < 10; i++ {
		in.WriteString("hello\n")
		err = client.Send()
		if err != nil {
			break
		}
	}

	require.NotNil(t, err)
}

// При подключении к несуществующему серверу, программа должна завершаться с ошибкой соединения/таймаута.
func TestClientConnectToNonExistingServer(t *testing.T) {
	address := "127.0.0.1:0" // не запускаем в рамках теста мок-сервер на этом адресе
	timeout := time.Second
	// Создаем буферы для ввода и вывода
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}

	// Создаем клиент
	client := NewTelnetClient(address, timeout, io.NopCloser(in), out)
	err := client.Connect()
	require.NotNil(t, err)
}
