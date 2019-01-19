BINARY=secspy
PACKAGES=`find ./cmd -mindepth 1 -maxdepth 1 -type d`
LIBRARYS=

all: clean man build

clean:
	for p in $(PACKAGES); do rm -f `echo $${p}|cut -d/ -f3`{,.1,.1.gz}; done

build: test
	for p in $(PACKAGES); do go build -ldflags "-w -s" $${p}; done

linux:
	for p in $(PACKAGES); do GOOS=linux go build -ldflags "-w -s" $${p}; done

install: build man
	@echo "If you get errors, you may need sudo."
	GOBIN=/usr/local/bin go install -ldflags "-w -s" ./...
	mkdir -p /usr/local/etc/$(BINARY) /usr/local/share/man/man1
	mv *.1.gz /usr/local/share/man/man1

uninstall:
	@echo "If you get errors, you may need sudo."
	rm -rf /usr/local/{etc,bin}/$(BINARY) /usr/local/share/man/man1/$(BINARY).1.gz


test: lint
	for p in $(PACKAGES) $(LIBRARYS); do go test -race -covermode=atomic $${p}; done

lint:
	goimports -l $(PACKAGES) $(LIBRARYS)
	gofmt -l $(PACKAGES) $(LIBRARYS)
	errcheck $(PACKAGES) $(LIBRARYS)
	golint $(PACKAGES) $(LIBRARYS)
	go vet $(PACKAGES) $(LIBRARYS)

man:
	script/build_manpages.sh ./

deps:
	dep ensure -update
