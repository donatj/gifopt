# gifopt

[![Go Report Card](https://goreportcard.com/badge/github.com/donatj/gifopt)](https://goreportcard.com/report/github.com/donatj/gifopt)
[![CI](https://github.com/donatj/gifopt/actions/workflows/ci.yml/badge.svg)](https://github.com/donatj/gifopt/actions/workflows/ci.yml)
[![GoDoc](https://godoc.org/github.com/donatj/gifopt?status.svg)](https://godoc.org/github.com/donatj/gifopt)

Simple Interframe Gif Optimizer

## Install

### From Source

```bash
go install github.com/donatj/gifopt/cmd/gifopt@latest
```

### From Binary

see: [Releases](https://github.com/donatj/gifopt/releases).

## Usage

```
Usage of gifopt [options] <gif>:
  -o string
    	Where to save the optimized gif (default "<orig>.opt.gif")
  -t float
    	Max interframe color diff percent threshold (default 0.03492566239471114)
```

## Example `t` value results 

### Base GIF

3.2mb

![input](https://user-images.githubusercontent.com/133747/192029927-35a8cb02-c274-4f75-aa16-0d40c403e6a3.gif)

| `gifopt -o output.05.gif -t .05 input.gif`                                                                      | `gifopt -o output.1.gif -t .1 input.gif`                                                                        | 
|-----------------------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------| 
| <img src="https://user-images.githubusercontent.com/133747/192030222-decd4010-fbe4-40ce-a0f5-519102b0dfa7.gif"> | <img src="https://user-images.githubusercontent.com/133747/192030383-98e69a7d-e3f5-40ff-892b-d5dbb3fb6d07.gif"> | 
| 2.6mb                                                                                                           | 1.7mb                                                                                                           | 
---

| `gifopt -o output.5.gif -t .5 input.gif`                                                                        | `gifopt -o output2.gif -t 2 input.gif`                                                                          | 
|-----------------------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------| 
| <img src="https://user-images.githubusercontent.com/133747/192030560-b2729412-88b3-4132-90ba-d1735bcc650e.gif"> | <img src="https://user-images.githubusercontent.com/133747/192030705-3f7aa332-c9ef-406e-b982-5466f3535309.gif"> | 
| 1.1mb                                                                                                           | .5mb                                                                                                            | 


