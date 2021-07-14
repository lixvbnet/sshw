DIST_DIR	= dist
CMD 		= sshw
VER 		= 1.0
LD_FLAGS	= -X 'main.CMD=$(CMD)' -X 'main.Version=$(VER)'

build:
	go build -ldflags="$(LD_FLAGS)"

install:
	go install -ldflags="$(LD_FLAGS)"

test:
	cd sshlib && go test -v

package: clean
	sh build.sh $(DIST_DIR) $(CMD) $(VER) "-ldflags=\"$(LD_FLAGS)\""


clean:
	rm -rf $(DIST_DIR)

t:
	echo ""
