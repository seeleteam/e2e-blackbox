# Makefile to build the command lines and tests in Seele project.
# This Makefile doesn't consider Windows Environment. If you use it in Windows, please be careful.

all:	run

run:
	go build -o ./build/run ./run 
	@echo "Done run building"



.PHONY: run
