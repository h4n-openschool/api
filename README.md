# openschool/classes

The Class microservice of OpenSchool.

## Development Prerequisites

You will need:

- [Docker](https://docker.com)
- A Linux development environment
- An installation of [Nix](https://nixos.org/download)

Once you have those, you need to run:

```shell
mkdir -p $HOME/.config/nix
cat "experimental-features = nix-command flakes" >> $HOME/.config/nix/nix.conf
```

## Running

Clone a copy of the repository, then drop into a development shell via the
included [flake](./flake.nix).

```shell
nix develop
```

This process could take a while the first time, as it will be downloading copies
of required development tooling.

Then, you need to start a copy of RabbitMQ in Docker.

```shell
docker run \
    -d \             # run in the background
    -p 15672:15672 \ # the web ui port
    -p 5672:5672 \   # the amqp port
    rabbitmq
```

You should be able to browse to http://127.0.0.1:15672/ and log in as
`guest:guest`.

Now, you can run the server software (still inside the nix shell).

```shell
go run . serve
```

By default, it will be running on port 8001, bound only to the local interface.

