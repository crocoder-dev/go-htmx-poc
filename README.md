## AIR SETUP

install air

```bash
curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```

add to ~/.bashrc

```bash
alias air='$(go env GOPATH)/bin/air'
```

run it with

```bash
air
```
