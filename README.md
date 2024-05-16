# logger
A lightweight and thread-safe logger for Go based on [log.c](https://github.com/rxi/log.c). This logger is designed to provide a simple yet powerful logging mechanism for Go applications, supporting multiple log levels and custom callbacks.


## Installation

```sh
go get github.com/ayo-ajayi/logger
```

## Usage

Import the package into your code:

```go
import "github.com/ayo-ajayi/logger"
```
```go
func main(){
    log := logger.NewLogger(logger.INFO, false)
    world:="world"
    log.Info("Hello %s", world)
}
```

## License
This library is free software; you can redistribute it and/or modify it under
the terms of the MIT license. See [LICENSE](LICENSE) for details.

##  Author

-   [Ayomide Ajayi](https://github.com/ayo-ajayi)