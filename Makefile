install:
	go install

test:
	cd sshlib && go test -v


DIST_DIR = dist

package:
	rm -rf $(DIST_DIR) && mkdir $(DIST_DIR)
	env GOOS=darwin GOARCH=amd64 sh -c 'go build -o $(DIST_DIR)/sshw-$$GOOS-$$GOARCH'
	env GOOS=linux GOARCH=amd64 sh -c 'go build -o $(DIST_DIR)/sshw-$$GOOS-$$GOARCH'
	env GOOS=windows GOARCH=amd64 sh -c 'go build -o $(DIST_DIR)/sshw-$$GOOS-$$GOARCH'
