#
# Makefile that builds the required library dependency, then installs the go module
#

all: install

install: libs 
	go install

libs:
	cd openzwave && make

clean:
	cd openzwave && make clean 
	go clean
