alias r := run

# compile and run with debug flag
run:
    go run main.go --debug

# run all tests, disable caching
test:
    go test ./... -count=1

build:
    go build .

# run race detector
race: 
    go run -race main.go --debug 2> race.log
