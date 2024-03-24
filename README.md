# AnPing

AnPing is a tool to ping.

| OS    | Protocol  |
|-------|-----------|
| Linux | ICMP, TCP |
| Stub  | TCP       |

# Usage

example:

```shell
anping -c 1 1.1.1.1 # default to use ICMP

anping icmp 1.1.1.1
```

For more document, please see `anping -h`.

# Build

```shell
./scripts/build.sh
```

Get anping in `./build/anping`

# Credits

* [prometheus-community/pro-bing](https://github.com/prometheus-community/pro-bing)

* [i3h/tcping](https://github.com/i3h/tcping)