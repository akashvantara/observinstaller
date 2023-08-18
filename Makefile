ARCH_PC := amd64
ARCH_ARM := arm64

all: pc arm

pc: linux_x64 windows_x64 mac_x64
arm: linux_arm64 windows_arm64 mac_arm64

linux_x64: observinstaller_linux_x64
windows_x64: observinstaller_windows_x64
mac_x64: observinstaller_mac_x64
linux_arm64: observinstaller_linux_arm64
windows_arm64: observinstaller_windows_arm64
mac_arm64: observinstaller_mac_arm64

observinstaller_linux_x64: main.go
	GOOS=linux GOARCH=${ARCH_PC} go build -o $@

observinstaller_windows_x64: main.go
	GOOS=windows GOARCH=${ARCH_PC} go build -o $@

observinstaller_mac_x64: main.go
	GOOS=darwin GOARCH=${ARCH_PC} go build -o $@

observinstaller_linux_arm64: main.go
	GOOS=linux GOARCH=${ARCH_ARM} go build -o $@

observinstaller_windows_arm64: main.go
	GOOS=windows GOARCH=${ARCH_ARM} go build -o $@

observinstaller_mac_arm64: main.go
	GOOS=darwin GOARCH=${ARCH_ARM} go build -o $@

clean: 
	rm observinstaller_*
