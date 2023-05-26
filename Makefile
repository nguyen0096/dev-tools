CGO_LDFLAGS ?= $(filter -g -L% -l% -O%,${LDFLAGS})
export CGO_LDFLAGS

EXE = ndv
ifeq ($(GOOS),windows)
EXE = ndv.exe
endif

LDFLAGS_ARG =-ldflags ${CGO_CPPFLAGS}
ifeq ($(CGO_CPPFLAGS),)
LDFLAGS_ARG =
endif

.PHONY: build
build:
	@rm -f ./bin/${EXE}
	@printf "Compiling ${EXE}...\n"
	@go build -trimpath ${LDFLAGS_ARG} -o bin/${EXE} ./main.go

.PHONY: aws_mfa
aws_mfa: build
	@printf "Running ${EXE} aws mfa...\n"
	@./bin/${EXE} aws mfa -t ${TOKEN}