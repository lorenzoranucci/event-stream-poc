# Event streams POC in Golang

## Quick start 

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

Try to change the value of `MESSAGING_PROTOCOL` between `json` and `protobuf` MESSAGING_PROTOCOLin the `.env.dist` file.
