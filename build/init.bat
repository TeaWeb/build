REM initialize project
set GOPATH=
set GO111MODULE=on
set GOPROXY=direct

REM download
go mod tidy

echo "[done]"