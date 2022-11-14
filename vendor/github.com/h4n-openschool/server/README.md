# openschool/http

An HTTP server implementation, built to the standards defined in
[RFC2616](https://www.rfc-editor.org/rfc/rfc2616).

## Setting up

To set up, use the `deps` [Makefile](./Makefile) task.

```shell
make deps
```

## Testing

To run tests, use the `test` [Makefile](./Makefile) task.

```shell
make test
```

## Example Server

There is an example server in [./example](./example/). To run it, run `make`.

```shell
$ make
2022/10/30 18:10:46 opening tcp socket...
2022/10/30 18:10:46 beginning accept loop...
...
```
