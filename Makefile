DIST_DIR	= dist
NAME 		= $(shell basename $$PWD)
MAIN_DIR	= .
LD_FLAGS	= -X 'main.Name=$(NAME)' -X 'main.GitHash=`git rev-parse --short=8 HEAD`'

build:
	go build -ldflags="$(LD_FLAGS)"	$(MAIN_DIR)

install:
	go install -ldflags="$(LD_FLAGS)" $(MAIN_DIR)

test:
	cd sshlib && go test -v

package: clean test
	sh build.sh $(DIST_DIR) $(NAME) "-ldflags=\"$(LD_FLAGS)\"" $(MAIN_DIR)


clean:
	rm -rf $(DIST_DIR)


.PHONY: build install test package clean
