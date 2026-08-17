# go-core

  A reusable go module with commonly used tools.

  The current version V3

## Enable it locally workspace

```sh
# project root
cd "C:\Eliezer\workspace\github.com\go-inventory-v2"
go work init .
go work use ../go-core

#usage
go work sync
go run .
```

## Installation

```
go get -u github.com/eliezerraj/go-core
```
## Test

```
go test -v -run "^TestGoCore_Kafka_Producer$"
```

## v2/event/kafka 
```sh
#install C compiler
sudo apt-get update
sudo apt-get install build-essential

#enable CGO
go env -w CGO_ENABLED=1

#check and the result must be 1
go env CGO_ENABLED

1
```