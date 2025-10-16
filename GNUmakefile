MAKEFLAGS += -j
build-dir := $(CURDIR)/build

binary := $(build-dir)/c
code := $(shell find $(CURDIR) -name "*.go")
code-deps := $(build-dir)/code-deps.mk
tidy := $(build-dir)/go-mod-tidy.out
generate := $(build-dir)/go-generate.out
test := $(build-dir)/go-test.out
lint := $(build-dir)/golangci-lint.out
sec := $(build-dir)/gosec.out
vuln := $(build-dir)/govulncheck.out
nilaway := $(build-dir)/nilaway.out
errcheck := $(build-dir)/errcheck.out
scc := $(build-dir)/scc.out

.PHONY: all
all\
: $(binary) \
  $(lint) \
  $(test) \
  $(sec) \
  $(vuln) \
  $(scc) \
; cat $(test) \
; cat $(scc) \
; date

.PHONY: deep
deep \
: all \
  $(nilaway) \
  $(errcheck)

# This will need to be restarted when code is added or removed.
.PHONY: watch
watch \
: \
; find $(code) $(MAKEFILE_LIST) | entr $(MAKE)

$(binary) \
: $(code) $(code-deps) $(tidy) $(generate) $(MAKEFILE_LIST) \
| $(build-dir) \
; gofumpt -l -w .; go build -o $@ .

$(test) \
: $(binary) \
| $(build-dir) \
; go test -cover -failfast -parallel=8 -count=2 -shuffle=on > $@ 2>&1 \
; exit_code=$$? \
; test $$exit_code -ne 0 && cat $@ \
; exit $$exit_code

$(generate) \
: $(code) $(code-deps) $(MAKEFILE_LIST) \
| $(build-dir) \
; go generate > $@ 2>&1 \
; exit_code=$$? \
; test $$exit_code -ne 0 && cat $@ \
; exit $$exit_code

$(tidy) \
: $(code) $(code-deps) $(MAKEFILE_LIST) \
| $(build-dir) \
; go mod tidy > $@ 2>&1 \
; exit_code=$$? \
; test $$exit_code -ne 0 && cat $@ \
; exit $$exit_code

$(lint) \
: $(binary) \
| $(build-dir) \
; golangci-lint run > $@ 2>&1 \
; exit_code=$$? \
; test $$exit_code -ne 0 && cat $@ \
; exit $$exit_code

$(sec) \
: $(binary) \
| $(build-dir) \
; gosec ./... > $@ 2>&1 \
; exit_code=$$? \
; test $$exit_code -ne 0 && cat $@ \
; exit $$exit_code

$(vuln) \
: $(binary) \
| $(build-dir) \
; govulncheck > $@ 2>&1 \
; exit_code=$$? \
; test $$exit_code -ne 0 && cat $@ \
; exit $$exit_code

$(nilaway) \
: $(binary) \
| $(build-dir) \
; nilaway ./... > $@ 2>&1 \
; exit_code=$$? \
; test $$exit_code -ne 0 && cat $@ \
; exit $$exit_code

$(errcheck) \
: $(binary) \
| $(build-dir) \
; errcheck ./... > $@ 2>&1 \
; exit_code=$$? \
; test $$exit_code -ne 0 && cat $@ \
; exit $$exit_code

$(scc) \
: $(binary) \
| $(build-dir) \
; scc \
  --no-cocomo \
  --sort code \
  -M json \
  -M css \
  -M gitignore \
  --exclude-dir deprecated \
  --exclude-ext xml \
  --exclude-ext toml \
  --exclude-ext yaml \
  --exclude-ext md \
  --exclude-ext txt \
> $@ 2>&1 \
; exit_code=$$? \
; test $$exit_code -ne 0 && cat $@ \
; exit $$exit_code

# This causes rebuilds on addition or removal of source files. Removing the
# candidate file ensures that it is rebuilt every time to check for new or
# removed files This is superior to making it a .PHONY target because doing so
# would simply cause a fresh rebuild every time obviating the whole dependency
# tree.
$(build-dir)/code-deps.mk \
: $(build-dir)/code-deps-candidate.mk \
; @diff $< $@ > /dev/null 2>&1 || (echo found new or deleted code; cp $< $@) \
; rm -f $<

$(build-dir)/code-deps-candidate.mk\
: \
| $(build-dir) \
; @echo checking for new or deleted code \
; echo $(code) > $@

.PHONY: clean
clean \
: \
; rm -Rf $(build-dir)

$(build-dir)\
:\
; mkdir -p $@
