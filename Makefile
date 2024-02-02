.ONESHELL:

clean:
	rm -rf build

builddir:
	mkdir -p build

build/showvideo:
	go build -o ./build/showvideo ./demo/showvideo/main.go

showvideo: builddir build/showvideo
	./build/showvideo 2

showphone:
	scrcpy

dronevideo: builddir build/showvideo
	./build/showvideo 4

tinyscan:
	tinygo flash -size short -target=pybadge ./demo/tinyscan/

cubelife:
	tinygo flash -size short -target=metro-m4-airlift -opt=2 -monitor ./demo/cubelife/

minidrone:
	tinygo flash -size short -target=pybadge -ldflags="-X main.DeviceAddress=E0:14:DC:85:3D:D1" ./demo/minidrone/

tello:
	tinygo flash -size short -target=pybadge -ldflags="-X main.ssid=TELLO-AA5556" ./demo/tello/

