
ARCH      := "`uname -s`"
LINUX     := "Linux"
MAC       := "Darwin"

all:
	@if [ $(ARCH) = $(LINUX) ]; \
	then \
		echo "make in $(LINUX) platform"; \
		GOOS=linux go build -o ./docker/dbtest  ./main.go; \
	elif [ $(ARCH) = $(MAC) ]; \
	then \
		echo "make in $(MAC) platform"; \
		GOOS=linux GOARCH=arm64 go build -o ./docker/dbtest  ./main.go; \
	else \
		echo "ARCH unknown"; \
	fi
