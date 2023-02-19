# Go Data Pipelines

> How to use Go pipelines to abstract iteration complexity while properly testing them

Pipelines are a Go concurrency pattern for creating a continuous stream of data and enables the responsibilities of producing and consuming information to be very loosely coupled.
But at the same time that concurrency is a beautiful word for software engineering, it is dreaded by the challenges and new concerns it brings to the craft. 

This article tries to address these challenges by providing some examples that hopefully will shed a light when walking this treacherous path. Those names and definitions aren't some standard or consensus across the community, but rather the author own conventions while testing and learning those ideas.        

## What is a pipeline

It is not the intention of this discussion to dive deep in the pipeline pattern details, there are [better sources](https://go.dev/blog/pipelines) for learning that and it is recommended (but not required) to understand at least the basics about [channels](https://go.dev/doc/effective_go#channels).
Even though this section isn't supposed to say something new, it is important to set some ground rules and define the what is the objects being mentioned.

### Data Source

A _data source_ is a method responsible for returning a list of items of the same type using some pagination parameters to limit the results. 
It is usually related to I/O interaction, like retrieving data from a local database or an external API. 
For simplification purposes, this method will be striped down to the bare minimum capable of showcase the ideas and can be something as simple as that:

```go
// ListNumbers emulates a typical data source creating a fictitious result considering the Pagination parameter
func ListNumbers(_ context.Context, p Pagination) ([]int, error) {
	result := make([]int, 0, p.PageSize)

	start := (p.PageNumber - 1) * p.PageSize
	finish := p.PageNumber * p.PageSize

	for i := start; i < finish; i++ {
		result = append(result, i)
	}
	return result, nil
}
```

Its anatomy could be seen as something resembling:
```go
func ListItemsOfTypeInt(_ context.Context, p Pagination) ([]int, error) {
	...
}
```

#### Page quantity

In addition to this method, it is also useful (maybe even necessary) to have a method responsible for knowing the total amount of data that can be retrieved.
Something like:

```go
func NumberCount(_ context.Context) (int, error) {
	...
}
``` 

These methods can be bundled in a `struct` and/or `interface` to always pair them together, but it not the intention of this article to dive into these details. 

### Data Pipe

The data pipe itself is a stream of data that _returns_ an unbounded list of items. 
The quantity of items is not known and shouldn't matter, so an important characteristic of this structure is that it shouldn't rely on holding all possible items in memory at the same moment.

Different languages provide different structures for that, and the idiomatic way for doing something like that in Go is a simple channel like `chan int`. 
And to follow the patterns for error handling in the language, would be natural to pair this channel with a `chan error`.

There isn't so much caveats when looking into the data pipes on isolation. It is just a simple language construct for sharing memory. 
The language shines when connecting the data source with the pipes to allow decoupled consuming later on.

### Pipeline

Pipeline is a function responsible for creating the pipeline and connecting the information from data source to the pipeline. 
It organizes all data production by connecting all moving parts on that side. 

[//]: # (```go)

[//]: # (func Pipeline&#40;ctx context.Context&#41; &#40;chan int, chan error&#41; {)

[//]: # (	resultInt := make&#40;chan int&#41;)

[//]: # (	resultErr := make&#40;chan error&#41;)

[//]: # ()
[//]: # (	go func&#40;&#41; {)

[//]: # (		// while there's data available in the data source send it to the returned channels)

[//]: # (	}&#40;&#41;)

[//]: # ()
[//]: # (	return resultInt, resultErr)

[//]: # (})

[//]: # (```)

## Implementation

Knowing what is a data source and pipe, it is possible to dive into the details on how to connect them with a pipeline.
Even though the first tests and developments were kind of erratic and experimental, the structure of the pipeline can be presented in a more structured and intentional step by step development for better understand. 

### First scratch

Its anatomy will be synthesized starting with the function signature and returning values:

```go
func Pipeline(ctx context.Context) (chan int, chan error) {
	resultInt := make(chan int)
	resultErr := make(chan error)

	// ...

	return resultInt, resultErr
}
```

Despite being quite simple at first glance, this block of code has some decisions worth noticing:
* The function has a single `context.Context` parameter. It doesn't require more than that since all data will be retrieved, and this `ctx` variable will be required to handle communications between go routines;
* There are two distinct channels being returned, the last one reserved for errors. 
  * It is possible to create a struct holding successful values and errors, enabling the return of a single value. 
  The chosen approach looks closer to most go signatures that return an error
* The result values are unbuffered `chan` values. 
  * One could argue that a buffered channel could improve performance by optimizing for I/O throughput, but it would lead to increased complexity. 
  Buffered channels would requiring checking the buffer length before closing to avoid discarding relevant data.  
  * The choice of unbuffered channels ensure that every message send to a channel is received by someone.
  Increasing in throughput can still be achieved by spawning more channel receiver (using the _fan out pattern_)

### Channel management

After the channels creation, it is necessary to ensure they are closed once the functions finishes to avoid locking any resources. 

```go
func Pipeline(ctx context.Context) (chan int, chan error) {
	resultInt := make(chan int)
	resultErr := make(chan error)

	go func() {
	    defer close(resultInt)	
	    defer close(resultErr)

		g, ctx := errgroup.WithContext(ctx)

        // ...
		
        if err := g.Wait(); err != nil {
            resultErr <- err
        }
    }()

	return resultInt, resultErr
}
```

This next step for the pipeline starts a go routine responsible for in fact sending data to the channels created.
The first thing to keep in mind, is reminding to close these channels mentioned and it should be done o function exit using the `defer` statement.

After that, it is created an `errorgroup.Group` using the experimental package [errgroup](https://pkg.go.dev/golang.org/x/sync/errgroup@v0.1.0#pkg-overview).
This group `g` will be responsible for executing the go routines responsible for each page of data, and the `g.Wait()` statement ensures that an error returned by any of these routines is delivered to the correct channel. 

### Iteration through pages

With the channels creation and closing solved, now the data source can be used to retrieve the data using the available pagination.  

```go
func Pipeline(ctx context.Context) (chan int, chan error) {
    resultInt := make(chan int)
    resultErr := make(chan error)
    
    go func () {
        defer close(resultInt)
        defer close(resultErr)
    
        g, ctx := errgroup.WithContext(ctx)

        totalNumbers, err := NumbersCount(ctx)
        if err != nil {
            resultErr <- err
            return
        }
        
        pageSize := 10  // arbitrary value
        pagesCount := totalNumbers / pageSize
        
        for pageNumber := 1; pageNumber <= pagesCount; pageNumber++ {
            pagination := Pagination{
                PageNumber: pageNumber,
                PageSize:   pageSize,
            }
            g.Go(func() error {
                //...
            })
        }
		
        if err := g.Wait(); err != nil {
			resultErr <- err
        }
    }()
    
    return resultInt, resultErr
}
```

The first part of this new code, from the `totalNumber` variable creation to `pagesCount` calculation hasn't so many things worth noticing, but it is interesting to see how the rest of the function made this bit a little easier.
If there is any error returned by the `NumberCount(...)` call, this value is send to a channel and the function returns. 
Due to the deferred closing the `if` block doesn't require more than that.
And the `pageSize` value is completely arbitrary and can be tweaked for optimization purposes without affecting anything outside the pipeline.

The second part starts a `for` loop that has an interesting detail for those not familiar in dealing with go routines.
The `pagination` variable is created to avoid using directly the variable `pageNumber` inside the go routine started inside the loop. 
The latter value changes over the same reference, while the former creates a new reference every iteration. Ignoring this could (which means that eventually will) lead to race conditions and unexpected behavior.
Another strategy to avoid that is passing the used values as parameters for the go routine, but it wasn't possible using the `errgroup` package.

### Send data to the channels

With everything set up, the final piece of the puzzle is calling the data source and send the retrieved data to the returned channels

```go
func Pipeline(ctx context.Context) (chan int, chan error) {
    resultInt := make(chan int)
    resultErr := make(chan error)
    
    go func () {
        defer close(resultInt)
        defer close(resultErr)
    
        g, ctx := errgroup.WithContext(ctx)

        totalNumbers, err := NumbersCount(ctx)
        if err != nil {
            resultErr <- err
            return
        }
        
        pageSize := 10  // arbitrary value
        pagesCount := totalNumbers / pageSize
        
        for pageNumber := 1; pageNumber <= pagesCount; pageNumber++ {
            pagination := Pagination{
                PageNumber: pageNumber,
                PageSize:   pageSize,
            }
            g.Go(func() error {
                ns, err := ListNumbers(ctx, pagination)
                if err != nil {
					return err
                }
                for _, n := range ns {
                    select {
					case resultInt <- n:
					case <-ctx.Done():
                        return ctx.Err()
                    }
				}
                return nil
            })
        }
		
        if err := g.Wait(); err != nil {
			resultErr <- err
        }
    }()
    
    return resultInt, resultErr
}
```

This last go routine started has the core of the Pipeline function, connection all the pieces, but since the explanation was broken down there is little left to describe.

The first detail worth noticing is that by using the `errgrop` package it is not necessary to use the `resultErr` channel inside this go routine.
Just returning the error in this function is enough because leaves for the `g.Wait()` call this responsibility.

The only thing left to do is going through each value returned by the data source and send to the `resultInt` channel. 
This is done by a `select/case` block that is also capable of exiting the whole routine when `ctx` is done. 
Another worthy detail is the `<-ctx.Done()` statement, which enables both external cancellation 
and timeout capabilities.

## Consuming from a pipeline

## Considerations