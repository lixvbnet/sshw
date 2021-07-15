DIST_DIR	= dist
NAME 		= sshw
VER 		= 1.0
LD_FLAGS	= -X 'main.Name=$(NAME)' -X 'main.Version=$(VER)' -X 'main.GitHash=`git rev-parse --short=8 HEAD`'

build:
	go build -ldflags="$(LD_FLAGS)"

install:
	go install -ldflags="$(LD_FLAGS)"

test:
	cd sshlib && go test -v

package: clean
	sh build.sh $(DIST_DIR) $(NAME) $(VER) "-ldflags=\"$(LD_FLAGS)\""


clean:
	rm -rf $(DIST_DIR)


.PHONY: build install test package clean
