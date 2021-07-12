DIST_DIR	= dist
CMD 		= sshw
VER 		= 1.0


install:
	go install

test:
	cd sshlib && go test -v

package: clean
	sh build.sh $(DIST_DIR) $(CMD) $(VER)


clean:
	rm -rf $(DIST_DIR)
