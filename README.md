# go-core

A reusable Go module with commonly used tools.
- api
    - Make REST call
- aws
    - aws config
    - s3
    - dynamo
    - paramter store
    - secret manager
- cache (redis cluster)
    - Get
    - SetCount
    - GetCount
    - Set
- cert
    - ParsePemToRSAPriv
    - ParsePemToRSAPub
- coreJson
    - Read JSON
    - Write JSON
- database
    - Acquire
    - Release
    - GetConnection
    - CloseConnection
    - StartTx
    - ReleaseTx
- event
    - Producer
    - InitTransactions
    - BeginTransaction
    - CommitTransaction
    - AbortTransaction
    - Consumer
    - Commit
- middleware
    - MiddleWareHandlerHeader
    - MiddleWareErrorHandler
- observability
    - Event
    - Span
    - SpanCtx
- tools
    - ConvertToDate

## Installation

go get -u github.com/eliezerraj/go-core
