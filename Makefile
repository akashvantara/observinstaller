ARCH_PC := amd64
ARCH_ARM := arm64

LD_FLAGS := -ldflags="-s -w"

all: pc arm

pc: observinstaller_linux_x64 observinstaller_windows_x64 observinstaller_mac_x64
arm: observinstaller_linux_arm64 observinstaller_windows_arm64 observinstaller_mac_arm64

def: main.go
	go build ${LD_FLAGS} -o observinstaller

observinstaller_linux_x64: main.go
	GOOS=linux GOARCH=${ARCH_PC} go build ${LD_FLAGS} -o $@

observinstaller_windows_x64: main.go
	GOOS=windows GOARCH=${ARCH_PC} go build ${LD_FLAGS} -o $@

observinstaller_mac_x64: main.go
	GOOS=darwin GOARCH=${ARCH_PC} go build ${LD_FLAGS} -o $@

observinstaller_linux_arm64: main.go
	GOOS=linux GOARCH=${ARCH_ARM} go build ${LD_FLAGS} -o $@

observinstaller_windows_arm64: main.go
	GOOS=windows GOARCH=${ARCH_ARM} go build ${LD_FLAGS} -o $@

observinstaller_mac_arm64: main.go
	GOOS=darwin GOARCH=${ARCH_ARM} go build ${LD_FLAGS} -o $@

clean: 
	rm observinstaller*
