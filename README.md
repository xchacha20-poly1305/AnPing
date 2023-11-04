# AnPing

一款多功能 Ping 工具。

# USAGE

## CMD

```
$ ./AnPing -h
NAME:
   AnPing - Ping whatever you like.

USAGE:
   AnPing [global options] command [command options] [arguments...]

VERSION:
   Unknow

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --count value, -c value  (default: -1)
   --help, -h               show help
   --version, -v            print the version
```

输入 ping 的目标地址即可，默认 ICMP，也可以使用别的协议。

使用 ICMP 协议时无需指定端口，使用其他协议时默认使用 443 作为端口。

例如：`AnPing tcp://1.1.1.1:443`

支持的协议列表：

| Protocol |
| - |
| ICMP |
| TCP |

# TODO

* [x] TCPing.

* [ ] Fix ICMP count.

* [ ] HTTP Ping.

# FAQ

* 为什么需要在 Linux 下需要请求权限？

  我们必须有 raw socket 权限才能发送 ICMP 数据包。如果您拒绝授予权限，我们只会使用 UDP 来进行 Ping。放心，我们只会要求必要的 `cap_net_raw=ep` 权限。

# Credits

* [urfave/cli](https://github.com/urfave/)

* [prometheus-community/pro-bing](https://github.com/prometheus-community/pro-bing)

* [i3h/tcping](https://github.com/i3h/tcping)