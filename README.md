

```shell
go get github.com/shuxs/fproxy
```

```shell
git clone https://github.com/shuxs/fproxy.git
cd fproxy
go build

./fproxy 
```

```toml
[app]
listen  =":3000"

[[proxy]]
name    ="backend"
prefix  ="/api"
target  ="http://192.168.1.22:8080/"

[[proxy]]
name    ="front"
prefix  ="/"
target  ="www"
```

```shell
./fproxy -c fproxy.toml
```