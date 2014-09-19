#
# Makefile that builds the required library dependency, then installs the go module
#

all: build

build: here

here: libs 
	mkdir -p NT
	scripts/GenerateNT.sh
	go install

libs:
	cd openzwave && make

clean:
	cd openzwave && make clean 
	go clean -i

deps:	libs
	
