/*
Package sdnotify provides a pure-go alternative to the [sd_notify] C function, allowing a go process
to send messages to systemd's service manager.

On non-linux platforms, the methods of this package are simply a no-op, making it safe for
multi-platform applications.

This package works by connecting to the notify socket created by the service manager. When
a process is launched as a systemd unit, the path to the notify socket is passed through the
`NOTIFY_SOCKET` environment variable automatically.

[sd_notify]: https://www.freedesktop.org/software/systemd/man/249/sd_notify.html
*/
package sdnotify

import (
	"errors"
	"net"
	"os"
	"runtime"
	"strings"
)

// The path to the notify socket to send signals to. The default value is an empty string
// which will use the value of the [NOTIFY_SOCKET] environment variable, populated by systemd
// automatically. This can be overridden if needed.
//
// [NOTIFY_SOCKET]: https://www.freedesktop.org/software/systemd/man/249/sd_notify.html#Environment
var NotifySocketPath = ""

// Inform the invoking service manager about service start-up or configuration reload completion.
func Ready() error {
	return Custom("READY=1\n")
}

// Inform the invoking service manager about the beginning of a configuration reload cycle.
func Reloading() error {
	return Custom("RELOADING=1\n")
}

// Inform the invoking service manager about the beginning of the shutdown phase of the service.
func Stopping() error {
	return Custom("STOPPING=1\n")
}

// Send a free-form human readable status string for the daemon to the service manager. The status
// string must be a single line.
func Status(status string) error {
	if strings.ContainsRune(status, '\n') {
		return errors.New("status must be a single line")
	}

	return Custom("STATUS=" + status + "\n")
}

// Tells the service manager to update the watchdog timestamp. This is the keep-alive ping that
// services need to issue in regular intervals if WatchdogSec= is enabled for it.
func Watchdog() error {
	return Custom("WATCHDOG=1\n")
}

// Tells the service manager that the service detected an internal error that should be handled by
// the configured watchdog options. This will trigger the same behaviour as if WatchdogSec= is
// enabled and the service did not send "WATCHDOG=1" in time.
func WatchdogTrigger() error {
	return Custom("WATCHDOG=trigger\n")
}

// Send a custom message to the service manager. Should take the format of the [state parameter] in
// a call to `sd_notify` and must end with a newline.
//
// [state parameter]: https://www.freedesktop.org/software/systemd/man/249/sd_notify.html#Description
func Custom(message string) error {
	if runtime.GOOS != "linux" {
		return nil
	}

	socketPath := NotifySocketPath
	if socketPath == "" {
		socketPath = os.Getenv("NOTIFY_SOCKET")
	}
	if socketPath == "" {
		return errors.New("no notify socket path defined")
	}

	conn, err := net.Dial("unixgram", socketPath)
	if err != nil {
		return errors.Join(errors.New("unable to connect to notify socket"), err)
	}
	defer conn.Close()

	if _, err := conn.Write([]byte(message)); err != nil {
		return errors.Join(errors.New("unable to send message to notify socket"), err)
	}

	return nil
}
