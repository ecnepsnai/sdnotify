# sdnotify

Package sdnotify provides a pure-go alternative to
[sd_notify](https://www.freedesktop.org/software/systemd/man/249/sd_notify.html), allowing a go
process to send signals to systemd.

On non-linux platforms, the methods of this package are simply a no-op, making it safe for
multi-platform applications.

## Usage

The `sdnotify` package works by connecting to the notify socket created by systemd. When a process
is launched as a systemd unit, the path to the notify socket is passed through the
[NOTIFY_SOCKET](https://www.freedesktop.org/software/systemd/man/249/sd_notify.html#Environment)
environment variable automatically.

### Sending Messages

The `sdnotify` package supports sending the following messages:

```go
sdnotify.Ready() // READY=1
sdnotify.Reloading() // RELOADING=1
sdnotify.Stopping() // STOPPING=1
sdnotify.Status("some status") // STATUS=some status 
sdnotify.Watchdog() // WATCHDOG=1
sdnotify.WatchdogTrigger() // WATCHDOG=trigger
```

Additionally you may call `sdnotify.Custom("")` to send a message not supported by this package.
The format of the message must be the same as the `state` parameter to a `sd_notify` call.
[_Learn more_](https://www.freedesktop.org/software/systemd/man/249/sd_notify.html)

#### Examples

```go
package main

import (
	"github.com/ecnepsnai/sdnotify"
)

func main() {
    sdnotify.Ready()
}
```

### Configuration

By default the `sdnotify` package is designed to work for go processes started by systemd, so no
extra configuration should be required.

The path to the notify socket to send signals to can be contolled with `sdnotify.NotifySocketPath`.
The default value is an empty string which will use the value of the NOTIFY_SOCKET environment
variable, populated by systemd automatically.
