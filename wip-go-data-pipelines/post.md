# Go Data Pipelines

> How to use Go pipelines to abstract iteration complexity while properly testing them

Pipelines are a Go concurrency pattern for creating a continuous stream of data and enables the responsibilities of producing and consuming information to be very loosely coupled. But at the same time that concurrency is a beautiful word for software engineering, it is dreaded by the challenges and new concerns it brings to the craft. 

This article tries to address these challenges by providing some examples that hopefully will shed a light when walking this treacherous path.      

## What is a pipeline 

It is not the intention of this discussion to dive deep in the pipeline pattern details, there are [better sources](https://go.dev/blog/pipelines) for learning that and it is recommeded (but not required) to understand at least the basics about [channels](https://go.dev/doc/effective_go#channels).
Even though this section isn't supposed to say something new, it is important to set some ground rules and define the what is the objects of study.

### Data Source

A _data source_ is a method responsible for returning a list of items of the same type considering some kind of pagination parameters. 
It is usually related to I/O interaction, like retrieving data from a local database or an external API. 
For simplification purposes, this method will be striped down to the bare minimum capable of showcase the ideas and can be something as simple as that:
