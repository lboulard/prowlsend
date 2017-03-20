# Send Prowl notification from command line

Simple tool that permit sending short notification to yourself.

## Install

```shell
go get gitub.com/lboulard/prowlsend
```

You will find a executable in `$GOPATH/bin` named `prowlsend`.

## Configuration

This command line utility requires that you have an API key at
<https://www.prowlapp.com>.

Write this API key into configuration file in your home folder at
`$HOME/.config/prowl/prowl.toml` file:

```toml
apikey = "0123456789abcdedfgh"
```

Replace text after `apikey=` with your API key from Prowl.

Make sure this file is not world readable with
 `chmod 600 $HOME/.config/prowl/prowl.toml` command.

I suppose you have now installed Prowl App on your phone.

## Usage

Sending a simple message:

```shell
prowlsend "Hello World!"
```

You will receive a notification with _Prowlsend on ..._ with _..._ replaced by
hostname of computer. In the notification body, you have the message text set
from command argument.

Change default application name _Prowlsend_ with `-a` option. If you do not
want the hostname appended to application name, add `-o=false` to options.

Message can have a title with an event using `-e` option. An optional URL can
be set as information field with `-u` option. You can change message priority
default of 1 with `-p` option. 

`prowlsend` can read a configuration file from another location with `-c`
option.

Finally, see a sum up of options with `prowlsend -h`.
