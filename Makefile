DIST_DIR	= dist
NAME 		= $(shell basename $$PWD)
LD_FLAGS	= -X 'main.Name=$(NAME)' -X 'main.GitHash=`git rev-parse --short=8 HEAD`'

build:
	go build -ldflags="$(LD_FLAGS)"

install:
	go install -ldflags="$(LD_FLAGS)"

test:
	cd sshlib && go test -v

package: clean test
	sh build.sh $(DIST_DIR) $(NAME) "-ldflags=\"$(LD_FLAGS)\""


clean:
	rm -rf $(DIST_DIR)


.PHONY: build install test package clean
