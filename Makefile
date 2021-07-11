install:
	go install

test:
	cd sshlib && go test -v
