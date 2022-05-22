export PATH := $(PWD):$(PATH)

.PHONY: test
test:
	go test
	go build
	cd example && HUT_DRY_RUN=1 hut init
	cd example/au && HUT_DRY_RUN=1 hut init
	cd example/au/dev && HUT_DRY_RUN=1 hut init
	cd example/au/dev && HUT_DRY_RUN=1 hut plan
