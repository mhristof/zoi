package docker

import (
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/stretchr/testify/assert"
)

func TestAlpine(t *testing.T) {
	var cases = []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "single package with dependencies",
			in: heredoc.Doc(`
				Sending build context to Docker daemon  996.9kB
				Step 1/2 : FROM alpine
				 ---> d6e46aa2470d
				Step 2/2 : RUN apk add htop
				 ---> Running in 3f7c87692ab9
				fetch http://dl-cdn.alpinelinux.org/alpine/v3.12/main/x86_64/APKINDEX.tar.gz
				fetch http://dl-cdn.alpinelinux.org/alpine/v3.12/community/x86_64/APKINDEX.tar.gz
				(1/3) Installing ncurses-terminfo-base (6.2_p20200523-r0)
				(2/3) Installing ncurses-libs (6.2_p20200523-r0)
				(3/3) Installing htop (2.2.0-r0)
				Executing busybox-1.31.1-r19.trigger
				OK: 6 MiB in 17 packages
				Removing intermediate container 3f7c87692ab9
				 ---> 7bdf00db2c0c
				Successfully built 7bdf00db2c0c
				Successfully tagged foo:latest
			`),
			out: "htop=2.2.0-r0",
		},
		{
			name: "single package with fixed version",
			in: heredoc.Doc(`
				Sending build context to Docker daemon  996.9kB
				Step 1/2 : FROM alpine
				 ---> d6e46aa2470d
				Step 2/2 : RUN apk add htop=2.2.0-r0
				 ---> Running in 3f7c87692ab9
				fetch http://dl-cdn.alpinelinux.org/alpine/v3.12/main/x86_64/APKINDEX.tar.gz
				fetch http://dl-cdn.alpinelinux.org/alpine/v3.12/community/x86_64/APKINDEX.tar.gz
				(1/3) Installing ncurses-terminfo-base (6.2_p20200523-r0)
				(2/3) Installing ncurses-libs (6.2_p20200523-r0)
				(3/3) Installing htop (2.2.0-r0)
				Executing busybox-1.31.1-r19.trigger
				OK: 6 MiB in 17 packages
				Removing intermediate container 3f7c87692ab9
				 ---> 7bdf00db2c0c
				Successfully built 7bdf00db2c0c
				Successfully tagged foo:latest
			`),
			out: "",
		},
		{
			name: "multiple apk packages",
			in: heredoc.Doc(`
				Sending build context to Docker daemon  1.031MB
				Step 1/6 : FROM alpine:3.5
				3.5: Pulling from library/alpine
				8cae0e1ac61c: Pull complete
				Digest: sha256:66952b313e51c3bd1987d7c4ddf5dba9bc0fb6e524eed2448fa660246b3e76ec
				Status: Downloaded newer image for alpine:3.5
				 ---> f80194ae2e0c
				Step 2/6 : RUN apk add --update-cache py2-pip ca-certificates py2-certifi py2-lxml                            python-dev cython cython-dev libusb-dev build-base                            eudev-dev linux-headers libffi-dev openssl-dev                            jpeg-dev zlib-dev freetype-dev lcms2-dev openjpeg-dev                            tiff-dev tk-dev tcl-dev
				 ---> Running in 565cc2239e79
				fetch http://dl-cdn.alpinelinux.org/alpine/v3.5/main/x86_64/APKINDEX.tar.gz
				fetch http://dl-cdn.alpinelinux.org/alpine/v3.5/community/x86_64/APKINDEX.tar.gz
				(1/97) Installing binutils-libs (2.27-r1)
				(2/97) Installing binutils (2.27-r1)
				(3/97) Installing gmp (6.1.1-r0)
				(4/97) Installing isl (0.17.1-r0)
				(5/97) Installing libgomp (6.2.1-r1)
				(6/97) Installing libatomic (6.2.1-r1)
				(7/97) Installing pkgconf (1.0.2-r0)
				(8/97) Installing libgcc (6.2.1-r1)
				(9/97) Installing mpfr3 (3.1.5-r0)
				(10/97) Installing mpc1 (1.0.3-r0)
				(11/97) Installing libstdc++ (6.2.1-r1)
				(12/97) Installing gcc (6.2.1-r1)
				(13/97) Installing make (4.2.1-r0)
				(14/97) Installing musl-dev (1.1.15-r8)
				(15/97) Installing libc-dev (0.7-r1)
				(16/97) Installing fortify-headers (0.8-r0)
				(17/97) Installing g++ (6.2.1-r1)
				(18/97) Installing build-base (0.4-r1)
				(19/97) Installing ca-certificates (20161130-r1)
				(20/97) Installing libbz2 (1.0.6-r5)
				(21/97) Installing expat (2.2.0-r1)
				(22/97) Installing libffi (3.2.1-r2)
				(23/97) Installing gdbm (1.12-r0)
				(24/97) Installing ncurses-terminfo-base (6.0_p20171125-r1)
				(25/97) Installing ncurses-terminfo (6.0_p20171125-r1)
				(26/97) Installing ncurses-libs (6.0_p20171125-r1)
				(27/97) Installing readline (6.3.008-r4)
				(28/97) Installing sqlite-libs (3.15.2-r2)
				(29/97) Installing python2 (2.7.15-r0)
				(30/97) Installing cython (0.25.1-r0)
				(31/97) Installing python2-dev (2.7.15-r0)
				(32/97) Installing py-pgen (2.7.10-r0)
				(33/97) Installing cython-dev (0.25.1-r0)
				(34/97) Installing udev-init-scripts (30-r6)
				Executing udev-init-scripts-30-r6.post-install
				(35/97) Installing eudev-libs (3.2.1-r1)
				(36/97) Installing libuuid (2.28.2-r1)
				(37/97) Installing libblkid (2.28.2-r1)
				(38/97) Installing xz-libs (5.2.2-r1)
				(39/97) Installing kmod (23-r1)
				(40/97) Installing eudev (3.2.1-r1)
				(41/97) Installing eudev-dev (3.2.1-r1)
				(42/97) Installing libpng (1.6.25-r0)
				(43/97) Installing freetype (2.7-r2)
				(44/97) Installing zlib-dev (1.2.11-r0)
				(45/97) Installing libpng-dev (1.6.25-r0)
				(46/97) Installing freetype-dev (2.7-r2)
				(47/97) Installing libjpeg-turbo (1.5.3-r2)
				(48/97) Installing libjpeg-turbo-dev (1.5.3-r2)
				(49/97) Installing jpeg-dev (8-r6)
				(50/97) Installing tiff (4.0.9-r6)
				(51/97) Installing tiff-dev (4.0.9-r6)
				(52/97) Installing lcms2 (2.8-r1)
				(53/97) Installing lcms2-dev (2.8-r1)
				(54/97) Installing libffi-dev (3.2.1-r2)
				(55/97) Installing libusb (1.0.20-r0)
				(56/97) Installing libusb-dev (1.0.20-r0)
				(57/97) Installing linux-headers (4.4.6-r1)
				(58/97) Installing openjpeg (2.3.0-r0)
				(59/97) Installing openjpeg-dev (2.3.0-r0)
				(60/97) Installing libcrypto1.0 (1.0.2q-r0)
				(61/97) Installing libssl1.0 (1.0.2q-r0)
				(62/97) Installing openssl-dev (1.0.2q-r0)
				(63/97) Installing py2-certifi (2016.9.26-r0)
				(64/97) Installing libgpg-error (1.24-r0)
				(65/97) Installing libgcrypt (1.7.10-r0)
				(66/97) Installing libxml2 (2.9.8-r1)
				(67/97) Installing libxslt (1.1.29-r1)
				(68/97) Installing py2-lxml (3.6.4-r0)
				(69/97) Installing py-setuptools (29.0.1-r0)
				(70/97) Installing py2-pip (9.0.0-r1)
				(71/97) Installing tcl (8.6.6-r0)
				(72/97) Installing tcl-dev (8.6.6-r0)
				(73/97) Installing libxau (1.0.8-r1)
				(74/97) Installing xproto (7.0.31-r0)
				(75/97) Installing libxau-dev (1.0.8-r1)
				(76/97) Installing xcb-proto (1.12-r0)
				(77/97) Installing libxdmcp (1.1.2-r2)
				(78/97) Installing libxcb (1.12-r0)
				(79/97) Installing libpthread-stubs (0.3-r3)
				(80/97) Installing libxdmcp-dev (1.1.2-r2)
				(81/97) Installing libxcb-dev (1.12-r0)
				(82/97) Installing xextproto (7.3.0-r1)
				(83/97) Installing xf86bigfontproto-dev (1.2.0-r3)
				(84/97) Installing xtrans (1.3.5-r0)
				(85/97) Installing inputproto (2.3.2-r0)
				(86/97) Installing libx11 (1.6.6-r0)
				(87/97) Installing kbproto (1.0.7-r1)
				(88/97) Installing libx11-dev (1.6.6-r0)
				(89/97) Installing libxrender (0.9.10-r1)
				(90/97) Installing fontconfig (2.12.1-r0)
				(91/97) Installing libxft (2.3.2-r1)
				(92/97) Installing expat-dev (2.2.0-r1)
				(93/97) Installing fontconfig-dev (2.12.1-r0)
				(94/97) Installing renderproto (0.11.1-r2)
				(95/97) Installing libxrender-dev (0.9.10-r1)
				(96/97) Installing libxft-dev (2.3.2-r1)
				(97/97) Installing tk-dev (8.6.6-r1)
				Executing busybox-1.25.1-r2.trigger
				Executing ca-certificates-20161130-r1.trigger
				OK: 305 MiB in 108 packages
				Removing intermediate container 565cc2239e79
				 ---> 29469422cd39
			`),
			out: heredoc.Doc(`
				build-base=0.4-r1
				ca-certificates=20161130-r1
				cython-dev=0.25.1-r0
				cython=0.25.1-r0
				eudev-dev=3.2.1-r1
				freetype-dev=2.7-r2
				jpeg-dev=8-r6
				lcms2-dev=2.8-r1
				libffi-dev=3.2.1-r2
				libusb-dev=1.0.20-r0
				linux-headers=4.4.6-r1
				openjpeg-dev=2.3.0-r0
				openssl-dev=1.0.2q-r0
				py2-certifi=2016.9.26-r0
				py2-lxml=3.6.4-r0
				py2-pip=9.0.0-r1
				tcl-dev=8.6.6-r0
				tiff-dev=4.0.9-r6
				tk-dev=8.6.6-r1
				zlib-dev=1.2.11-r0`),
		},
	}

	for _, test := range cases {
		assert.Equal(t, test.out, alpine(strings.Split(test.in, "\n")), test.name)
	}
}
