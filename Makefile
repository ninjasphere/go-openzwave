#
# Makefile that builds the required library dependency, then installs the go module
#
GENERATED=NT/NT.go CC/CC.go LOG_LEVEL/LOG_LEVEL.go CODE/CODE.go VT/VT.go

all: build

build: here

here: deps
	mkdir -p NT
	go install

libs:
	cd openzwave && make

clean:
	cd openzwave && make clean 
	go clean -i
	rm -rf $(GENERATED) 

deps:	libs
	scripts/GenerateNT.sh
	scripts/GenerateCODE.sh
	scripts/GenerateCC.sh
	scripts/GenerateLOG_LEVEL.sh
	scripts/GenerateVT.sh

