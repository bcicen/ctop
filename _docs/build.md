# Build

To build `ctop` from source, simply clone the repo and run:

```bash
make build
```

To build a minimal Docker image containing only `ctop`:
```bash
make image
```

Now you can run your local image:

```bash
docker run --rm -ti \
  --name ctop \
  -v /var/run/docker.sock:/var/run/docker.sock \
  ctop:latest
```
