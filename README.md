# go-core

   A reusable Go module with commonly used tools.

   The current version V2

## Installation

  go get -u github.com/eliezerraj/go-core

## Test

  go test -v -run "^TestGoCore_Kafka_Producer$"

## v2/event/kafka

  install C compiler

    sudo apt-get update
    sudo apt-get install build-essential

  enable CGO

    go env -w CGO_ENABLED=1

  check

    go env CGO_ENABLED


