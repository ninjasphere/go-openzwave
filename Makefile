#
# Makefile that builds the required library dependency, then installs the go module
#
GENERATED=NT/NT.go CC/CC.go LOG_LEVEL/LOG_LEVEL.go CODE/CODE.go VT/VT.go MF/MF.go

all: build

build: here

here: deps fmt
	mkdir -p NT
	go install

libs:
	cd openzwave && make

clean: clean-src
	cd openzwave && make clean 
	go clean -i
	rm -rf $(GENERATED) 

clean-src:
	find . -name '*~' -exec rm {} \;

fmt:
	gofmt -s -w *.go

deps:	libs
	scripts/GenerateNT.sh
	scripts/GenerateCODE.sh
	scripts/GenerateCC.sh
	scripts/GenerateLOG_LEVEL.sh
	scripts/GenerateVT.sh
	scripts/GenerateMF.sh

control-panel:
	@echo "run 'scripts/start-ozwcp.sh', configure the device as /dev/cu.SLAB_USBtoUART, then hit initialize."