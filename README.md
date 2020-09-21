# WORKERPOOL

_Simple worker pool library_

[![Build Status](https://secure.travis-ci.org/the4thamigo-uk/workerpool.png?branch=master)](https://travis-ci.org/the4thamigo-uk/workerpool?branch=master)
[![Coverage Status](https://coveralls.io/repos/the4thamigo-uk/workerpool/badge.svg?branch=master&service=github)](https://coveralls.io/github/the4thamigo-uk/workerpool?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/the4thamigo-uk/workerpool)](https://goreportcard.com/report/github.com/the4thamigo-uk/workerpool)

## Introduction

It is often said that there is no need for a general purpose worker pool library for `go`, and that you can do everything with the basic `chan` primitive. However, it is relatively easy to
make mistakes in safely catering for edge cases. This library aims to make it simple to embed worker pool logic into your project.

The worker pool simply accepts _work_ in the form of functions `func()`. Parameters are expected to be passed into the work function via closures and return values are expected to be
implemented by the calling code using whatever mechanism is required (i.e. `chan` or otherwise).

An example of usage is provided [here](./example/linkpuller) and the godoc is [here](https://godoc.org/github.com/the4thamigo-uk/workerpool)




