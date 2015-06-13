all: hsstat/hsstat hsps/hsps

test:
	cd hsperfdata && go test

hsstat/hsstat: hsstat/hsstat.go
	cd hsstat && go build

hsps/hsps: hsps/hsps.go
	cd hsps && go build

clean:
	rm -f hsps/hsps hsstat/hsstat

.PHONY: clean test all

