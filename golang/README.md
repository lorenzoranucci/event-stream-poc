# Event streams POC in Golang

## Quick start 

Fetch vendors:
```shell
make mod-vendor
```

Set envs:
```shell
cp .env.dist .env
```
_Optionally set "kafka" or "pulsar" as event stream and "json" or "protobuf" as format_.

Start the system:
```shell
make run-system
```

Watch the logs:
```shell
make log-system
```

Send reviews:
```shell
make post-review
```

Try to change the value of `MESSAGING_PROTOCOL` between `json` and `protobuf` in the `.env.dist` file.
