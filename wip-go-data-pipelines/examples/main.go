package main

import (
	"context"
	"fmt"
	"log"
	"taciogt.com/pipelines/pipelines"
	"time"
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	intCh, errCh := pipelines.Pipeline(ctx)

	for n := range intCh {
		fmt.Printf("result: %d\n", n)
	}

	if err := <-errCh; err != nil {
		log.Fatal(err)
	}

	fmt.Println("successfully executed program")
}
