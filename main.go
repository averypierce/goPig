package main

/*
-llport #####
-lrport #####
-rrport #####
-rlport #####

-rraddr
-lraddr

-nohead
-notail

-loopr
-loopl
*/

type settings struct {
	localPort, remotePort        int
	remoteAddr                   string
	nohead, notail, loopr, loopl bool
}

var ropts, lopts settings

func init() {
	//create setting profiles for left and right with default options
	lopts = settings{localPort: 36751, remotePort: 36751, remoteAddr: "localhost", nohead: false, notail: false, loopr: false, loopl: false}
	ropts = lopts

}

func parse(args string) {

}

func main() {
	MouseDemoMain()

}
