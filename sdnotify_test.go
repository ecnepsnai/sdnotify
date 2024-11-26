package sdnotify_test

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/ecnepsnai/sdnotify"
)

func startSocket(buf *bytes.Buffer, socketPath string) {
	l, err := net.ListenPacket("unixgram", socketPath)
	if err != nil {
		panic(err)
	}
	for {
		b := make([]byte, 1024)
		n, _, _ := l.ReadFrom(b)
		buf.Write(b[0:n])
	}
}

func expectMesssage(t *testing.T, send func() error, message string) {
	socketPath := path.Join(t.TempDir(), "notify.sock")
	sdnotify.NotifySocketPath = socketPath

	buf := &bytes.Buffer{}
	go startSocket(buf, socketPath)

	// Wait for the server to start
	time.Sleep(5 * time.Millisecond)

	if err := send(); err != nil {
		t.Errorf("Unexpected error sending message: %s", err.Error())
	}

	// Wait for the message to be recieved
	time.Sleep(5 * time.Millisecond)

	if buf.String() != message {
		t.Errorf("Unexpected data sent to socket. Expected: %#v Got: %#v", message, buf.String())
	}
}

func TestMain(m *testing.M) {
	if runtime.GOOS != "linux" {
		fmt.Println("Not running tests on non-linux OS")
		return
	}

	os.Exit(m.Run())
}

func TestReady(t *testing.T) {
	expectMesssage(t, func() error {
		return sdnotify.Ready()
	}, "READY=1\n")
}

func TestReloading(t *testing.T) {
	expectMesssage(t, func() error {
		return sdnotify.Reloading()
	}, "RELOADING=1\n")
}

func TestStopping(t *testing.T) {
	expectMesssage(t, func() error {
		return sdnotify.Stopping()
	}, "STOPPING=1\n")
}

func TestStatus(t *testing.T) {
	expectMesssage(t, func() error {
		return sdnotify.Status("hello")
	}, "STATUS=hello\n")
}

func TestStatusMultiline(t *testing.T) {
	err := sdnotify.Status("STATUS=Hello\nWorld!")
	if err == nil {
		t.Error("No error seen when one expected")
	}
	if !strings.Contains(err.Error(), "status must be a single line") {
		t.Errorf("Unexpected error when sending invalid message: %s", err.Error())
	}
}

func TestWatchdog(t *testing.T) {
	expectMesssage(t, func() error {
		return sdnotify.Watchdog()
	}, "WATCHDOG=1\n")
}

func TestWatchdogTrigger(t *testing.T) {
	expectMesssage(t, func() error {
		return sdnotify.WatchdogTrigger()
	}, "WATCHDOG=trigger\n")
}

func TestCustom(t *testing.T) {
	expectMesssage(t, func() error {
		return sdnotify.Custom("HELLO=WORLD\n")
	}, "HELLO=WORLD\n")
}
