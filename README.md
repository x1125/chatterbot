# chatterbot

A Golang chatter bot API that supports Cleverbot.

For the Mono/.NET, JAVA, Python and PHP version, take a look at [pierredavidbelanger cleverbot-api implementation](https://github.com/pierredavidbelanger/chatter-bot-api).

## Usage

```go
package main

import (
	"x1125io/chatterbot"
	"fmt"
)

func main() {
	session := chatterbot.NewCleverbot()
	fmt.Println(session.ThinkThrough("Hey"))
	fmt.Println(session.ThinkThrough("How are you?"))
	fmt.Println(session.ThinkThrough("Wow, that's terrible..."))
	fmt.Println(session.ThinkThrough("So what you're doing?"))
	fmt.Println(session.ThinkThrough("I'm glad talking to you"))
}
```
