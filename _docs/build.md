# Build

To build `ctop` from source, ensure you have a recent version of [glide](https://github.com/Masterminds/glide) installed and run:

```bash
go get github.com/bcicen/ctop && \
cd $GOPATH/src/github.com/bcicen/ctop && \
make build
```

To build a minimal Docker image containing only `ctop`:
```bash
make image
```

Now you can run your local image:

```bash
docker run -ti --name ctop --rm -v /var/run/docker.sock:/var/run/docker.sock ctop
```
