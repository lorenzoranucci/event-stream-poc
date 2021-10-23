# Event streams POC in Golang

## Quick start 

Fetch vendors:
```shell
make mod-vendor
```

Start the system:
```shell
make run-system
```

Watch the consumer:
```shell
make log-consumer
```

Send reviews:
```shell
make post-review
```

Try to change the value of `MESSAGING_PROTOCOL` between `json` and `protobuf` in the `.env.dist` file.
