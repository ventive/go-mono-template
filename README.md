# go-mono-template

## Usage

Start the application with the following command:

```
docker compose -f docker-compose.local.yml up -d
```

Install the nats cli with the following command:

```
go install github.com/nats-io/natscli/nats@latest
```

Subscribe to the nats topics with the following command:

```
nats sub ventive.service.subtractor.inbox
nats sub ventive.service.subtractor.outbox.default

nats sub ventive.service.adder.inbox
nats sub ventive.service.adder.outbox.default
```

Publish to the nats topics with the following command:

```
nats pub ventive.service.adder.inbox '{"data": {"a": 1, "b": 2}}'
nats pub ventive.service.subtractor.inbox '{"data": {"a": 3, "b": 2}}'
```
