include $(GOROOT)/src/Make.inc

TARG=lists

GOFILES=\
	lists.go\
	caching.go\
	header.go\
	linear.go\
	cyclic.go

include $(GOROOT)/src/Make.pkg