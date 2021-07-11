## ProxyChecker 

Note:
File Format is new line every address, like below:

```
213.232.127.164:8080
50.192.195.69:52018
177.55.207.38:8080
200.155.139.242:3128
```

### Help
```
❯ ./proxychecker --help
NAME:
   proxychecker - fight the loneliness!

USAGE:
   proxychecker [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --url value   check proxy list from url with extension file is text
   --file value  check proxy from file
   --ip value    checker ip proxy
   --out value   output file (default: "active_proxy.txt")
   --help, -h    show help (default: false)
```


### Check IP Proxy
```
./proxychecker --ip "192.189.10.9:8080"
```

### Check IP Proxy from file
```
./proxychecker --file ~/home/user/Downloads/proxy.txt
```

### Check IP Proxy from url
```
./proxychecker --url https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/socks5.txt
```

