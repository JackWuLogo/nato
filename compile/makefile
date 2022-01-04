# Makefile Include Sample

# Set Version
VERSION=v2.0

# Get Revision
# Get Svn Revision
REVISION=$(shell LANGUAGE=en_US svn info | grep Revision | cut -d " " -f 2)

# Get Git Revision
REVISION=$(shell git show -s --pretty=format:%h)

# Get Go Compile File
COMPILE=$(shell GOFLAGS=${GO_FLAGS} go list micro-libs/compile)

# Project Build Time
BUILD_TIME=$(shell date "+%G%m%d_%H%M%S")

# Go Make ldflags
GO_LDFLAGS=-X $(COMPILE).version=$(VERSION) -X $(COMPILE).revision=$(REVISION) -X $(COMPILE).name=$(BINARY) -X $(COMPILE).buildTime=$(BUILD_TIME) -s -w

# Go Flags
GO_FLAGS=''