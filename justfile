alias r := run

# compile and run with debug flag
run:
    go run main.go --debug

runp:
    go run main.go p --debug

# run all tests, disable caching
test:
    go test ./... -count=1

build:
    go build .

# run race detector
race: 
    go run -race main.go --debug 2> race.log

# start debug server
debug-server: 
    dlv debug --headless --api-version=2 --listen=127.0.0.1:43000 . -- --debug

# connect to debug server
debug-connect:
     dlv connect 127.0.0.1:43000
