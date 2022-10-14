# TL-WPA4220 go package

[![GitHub Super-Linter](https://github.com/angelbarrera92/tlwpa4220/workflows/Lint%20Code%20Base/badge.svg)](https://github.com/marketplace/actions/super-linter)

Go library for the TL-WPA4220

## Usage

### As a library

```go
package main

import (
    "fmt"

    "github.com/angelbarrera92/tlwpa4220/pkg/tlwpa4220"
)

func main() {
    t := tlwpa4220.Client{
        Username: "admin",
        Password: "hopeItsNotAdmin",
        IP:       "192.168.2.1",
    }

    wireless, err := t.WirelessStatistics()
    if err != nil {
        panic(err)
    }

    for _, device := range wireless.Data {
        fmt.Println(device.Mac)
    }

    err = t.Reboot()
    if err != nil {
        fmt.Println(err)
        panic(err)
    }
}
```

### As a CLI

To be developed

### As an API

To be developed

## License

[GPLv3](LICENSE)
