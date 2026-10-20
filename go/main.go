package main

import (
	"context"
	"fmt"
	"time"

	"github.com/klevente-coralogix/hiring-notificationsender/go/thirdparty"
)

func main() {
	api := thirdparty.NewMockAPI()
	api.SetDelay(100 * time.Millisecond)

	err := api.Post(context.Background(), "Hello, world!")
	fmt.Printf("post error: %v, posted: %v\n", err, api.Posted())
}
