RM=rm -f
RMFORCE=rm -rf
GENERATED_IPC=$(SR_CODE_BASE)/generated/src
SRCS=children.go patricia.go  
COMP_NAME=patriciad
all: exe

exe: $(SRCS)
	 go build -o $(COMP_NAME) $(SRCS)

guard:
ifndef SR_CODE_BASE
	 $(error SR_CODE_BASE is not set)
endif

clean:guard
	 $(RM) $(COMP_NAME)
