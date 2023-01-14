# Go Data Pipelines

> How to use Go pipelines to abstract iteration complexity while properly testing them

Pipelines are a Go concurrency pattern for creating a continuous stream of data and enables the responsibilities of producing and consuming information to be very loosely coupled. But at the same time that concurrency is a beautiful word for software engineering, it is dreaded by the challenges and new concerns it brings to the craft. 

This article tries to address these challenges by providing some examples that hopefully will shed a light when walking this treacherous path.      

## What is a pipeline 

It is not the intention of the discussion
