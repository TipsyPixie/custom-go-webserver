TEST_OPTS=-cover -timeout 10s -count 1
FIND=find -maxdepth 1 -type d -regextype posix-extended -regex '^\.\/problem[0-9]{1,3}$$'
VET_OPTS=
FMT_OPTS=-l -s -w
FIX_OPTS=
FOR_EACH_MODULE=echo -n 'config database httpError migration router' | xargs -d ' ' -i

FORCE:

test: FORCE
	! $(FOR_EACH_MODULE) go test $(TEST_OPTS) ./{} | grep -q 'FAIL'

vet: FORCE
	$(FOR_EACH_MODULE) go vet $(VET_OPTS) ./{}

lint: FORCE
	$(FOR_EACH_MODULE) gofmt $(FMT_OPTS) ./{}

fix: FORCE
	$(FOR_EACH_MODULE) go fix $(FIX_OPTS) ./{}

report: FORCE
	curl -s -d 'repo=github.com%2FTipsyPixie%2Fcustom-go-webserver' https://goreportcard.com/checks >/dev/null

precommit: test analyze format fix
