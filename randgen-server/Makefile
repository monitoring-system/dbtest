

ARCH      := "`uname -s`"
LINUX     := "Linux"
MAC       := "Darwin"

all:
	@if [ $(ARCH) = $(LINUX) ]; \
	then \
		echo "make in $(LINUX) platform"; \
		GOOS=linux go build -o ./randgen-server ./main.go; \
	elif [ $(ARCH) = $(MAC) ]; \
	then \
		echo "make in $(MAC) platform"; \
		GOOS=mac  go build -o ./randgen-server  ./main.go; \
	else \
		echo "ARCH unknown"; \
	fi

