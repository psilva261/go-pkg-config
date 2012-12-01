# go-pkg-config

**go-pkg-config** is a Go clone of pkg-config with no external dependencies.

### Example usage on OS X:

Inside your ~/.profile add:

	APPS=$HOME/Apps
	export PKG_CONFIG_PATH=/usr/lib/pkgconfig:/X11/lib/pkgconfig:$APPS/lib/pkgconfig
	export PATH=$APPS/bin:$PATH
	export DYLD_LIBRARY_PATH=$APPS/lib
	export C_INCLUDE_PATH=$APPS/include
	go build -o pkg-config main.go
	cp pkg-config $APPS/bin

Now when installing some software that uses pkg-config, just use $HOME/Apps as prefix:

	./configure --prefix=$HOME/Apps
	make
	make install
