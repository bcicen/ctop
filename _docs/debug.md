# Debug Mode

`ctop` comes with a built-in logging facility and local socket server to simplify debugging at run time. Debug mode can be enabled via the `CTOP_DEBUG` environment variable:

```bash
CTOP_DEBUG=1 ./ctop
```

While `ctop` is running, you can connect to the logging socket via socat or similar tools:
```bash
socat unix-connect:/tmp/ctop.sock stdio
```

example output:
```
15:06:43.881 ▶ NOTI 002 logger initialized
15:06:43.881 ▶ INFO 003 loaded config param: "filterStr": ""
15:06:43.881 ▶ INFO 004 loaded config param: "sortField": "state"
15:06:43.881 ▶ INFO 005 loaded config switch: "sortReversed": false
15:06:43.881 ▶ INFO 006 loaded config switch: "allContainers": true
15:06:43.881 ▶ INFO 007 loaded config switch: "enableHeader": true
15:06:43.883 ▶ INFO 008 collector started for container: 7120f83ca...
...
```
