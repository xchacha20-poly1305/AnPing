# AnPing

AnPing is a network probing tool, supporting ICMP, TCP and UDP ping.

# Usage

example:

```shell
anping -c 1 1.1.1.1 # default to use ICMP
anping icmp 1.1.1.1
```

For more document, please see `anping -h`.

# Build

```shell
CGO_ENABLED=0 make
```

# Credits

* [prometheus-community/pro-bing](https://github.com/prometheus-community/pro-bing)
* [i3h/tcping](https://github.com/i3h/tcping)
* [wangyu-/UDPping](https://github.com/wangyu-/UDPping)