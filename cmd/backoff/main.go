package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/cenkalti/backoff"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 0*time.Second)
	defer cancel()
	b := backoff.WithMaxRetries(backoff.NewConstantBackOff(2*time.Second), 5)
	err := backoff.Retry(func() error {
		fmt.Println("hello")
		return errors.New("error")
	}, backoff.WithContext(b, ctx))
	if err != nil {
		fmt.Println(err)
	}
}
