#!/bin/bash

protoc greet/greetpb/greet.proto --go_out=. --go-grpc_out=.
protoc calculator/calculatorpb/calculator.proto --go_out=. --go-grpc_out=.                                   ✔  9428  20:28:36
protoc blog/blogpd/blog.proto --go_out=. --go-grpc_out=.                                   ✔  9428  20:28:36