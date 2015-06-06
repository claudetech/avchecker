# Website availability checker [![Build Status](https://travis-ci.org/claudetech/avchecker.svg?branch=master)](https://travis-ci.org/claudetech/avchecker)

Simple website availability checker written in Go.

Periodically sends an HTTP request to the given URL,
and periodically sends the success count and success ratio.

Here is a sample usage to post stats to a Redis queue.

```go
package main

import (
  "github.com/claudetech/avchecker"
  redis "github.com/xuyu/goredis"
)

func main() {
  reporter, err := avchecker.NewRedisQueueReporter("coocreed", &redis.DialConfig{})
  if err != nil {
    panic(err)
  }
  checker := avchecker.New("http://localhost:4567", reporter, &avchecker.Options{})
  checker.StartChecking()
}
```

To POST the stats through HTTP instead of posting them to redis, you could
replace the `reporter, err := ` line by

```go
reporter, err := avchecker.NewHttpReporter("http://mystatscollector.com/stats", "application/json")
```

## CLI usage

CLI binaries are avilable on the [release page](https://github.com/claudetech/avchecker/releases).
See [avchecker-cli](./avchecker-cli) for usage.
