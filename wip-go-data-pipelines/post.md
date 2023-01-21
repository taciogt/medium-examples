# Go Data Pipelines

> How to use Go pipelines to abstract iteration complexity while properly testing them

Pipelines are a Go concurrency pattern for creating a continuous stream of data and enables the responsibilities of producing and consuming information to be very loosely coupled. But at the same time that concurrency is a beautiful word for software engineering, it is dreaded by the challenges and new concerns it brings to the craft. 

This article tries to address these challenges by providing some examples that hopefully will shed a light when walking this treacherous path.      

## What is a pipeline 

It is not the intention of this discussion to dive deep in the pipeline pattern details, there are [better sources](https://go.dev/blog/pipelines) for learning that and it is recommended (but not required) to understand at least the basics about [channels](https://go.dev/doc/effective_go#channels).
Even though this section isn't supposed to say something new, it is important to set some ground rules and define the what is the objects of study.

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
func ListItemsOfTypeT(_ context.Context, p Pagination) ([]T, error) {
	...
}
```

### Pipeline

The pipeline itself is a stream of data that _returns_ a list of items. 
The quantity of items is not known and shouldn't matter, so an important characteristic of this structure is that it shouldn't rely on holding all possible items in memory at the same moment.

Different languages provide different structures for that, and the idiomatic way for doing something like that in Go is a simple channel like `chan int`. 
And to follow the patterns for error handling in the language, would be natural to pair this channel with a `chan error`.

There isn't so much caveats when looking into the pipeline on isolation, but the language shines when connecting the data source with the pipeline and when consuming those data later on.

### Pipeline Producer

Pipeline producer is a function responsible for creating the pipeline and connecting the information from data source to the pipeline. Its

```go
func Pipeline(ctx context.Context) (chan int, chan error) {
	resultInt := make(chan int)
	resultErr := make(chan error)

	go func() {
		// while there's data available in the data source send it to the returned channels
	}()

	return resultInt, resultErr
}
```