## Go-Core

Current V3

```sh 
# project root
cd "C:\Eliezer\workspace\github.com\go-inventory-v2"

# start project
go mod init .

# install
go mod tidy

# Initialize an empty workspace
go work init .

# Add a single module
go work use ../go-core/v3

# sync files
go work sync

# run
go run .
```