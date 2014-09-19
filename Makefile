#
# Makefile that builds the required library dependency, then installs the go module
#

all: build

build: here

here: libs 
	mkdir -p NT
	scripts/GenerateNT.sh
	scripts/GenerateCODE.sh
	scripts/GenerateCC.sh
	scripts/GenerateLOG_LEVEL.sh
	scripts/GenerateVT.sh
	go install

libs:
	cd openzwave && make

clean:
	cd openzwave && make clean 
	go clean -i

deps:	libs
	
