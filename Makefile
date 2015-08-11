GIT_VER := $(shell git describe --tags)
DATE    := $(shell date +%Y-%m-%dT%H:%M:%S%z)
OWNER   := "mix3"
REPO    := "go-irc2mail"

clean:
	rm -rf pkg/*

binary: clean
	gox -osarch="linux/amd64 darwin/amd64" \
		-output "pkg/{{.Dir}}-${GIT_VER}-{{.OS}}-{{.Arch}}" \
		-ldflags "-X main.version ${GIT_VER} -X main.buildDate ${DATE}"

package: binary
	cd ./pkg && find . -name "*${GIT_VER}*" -type f \
		-exec mkdir -p go-irc2mail \; \
		-exec cp {} go-irc2mail/go-irc2mail \; \
		-exec zip -r {}.zip go-irc2mail \; \
		-exec rm -rf go-irc2mail \;
