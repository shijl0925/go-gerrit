**gerritctl** is a gerrit CLI based on [go-gerrit](https://github.com/shijl0925/go-gerrit) library. 🚀

:one: Generate a http password for the username that will manage the gerrit.

:two: Create the `configuration directory` and the `config.json file`
```
$ mkdir -p ~/.config/gerritctl/
$ pushd ~/.config/gerritctl/
    $ vi config.json 
    {
        "Url": "https://gerrit.mydomain.com",
        "Account": "xxx",
        "Password": "yyy"
    }
$ popd
```

:three: Build the gerritctl

```
$ git clone https://github.com/shijl0925/go-gerrit.git
$ cd cli/gerritctl
$ make
```

```
$ ./gerritctl
Client for gerrit, manage resources by the gerrit

Usage:
  gerritctl [command]

Available Commands:
  help        Help about any command
  project     project related commands
  version     get server version

Flags:
      --config string   Path to config file
  -h, --help            help for gerritctl
  -v, --version         version for gerritctl

Use "gerritctl [command] --help" for more information about a command.
```

:rocket: :rocket: :rocket: :rocket:
