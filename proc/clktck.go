package proc

/*
   #include <unistd.h>
   #include <sys/types.h>
   #include <pwd.h>
   #include <stdlib.h>
*/
import "C"

// getClkTck returns clock ticks in a second by reading
// the system's configuration with Cgo.
func getClkTck() float64 {
	var sc_clk_tck C.long
	sc_clk_tck = C.sysconf(C._SC_CLK_TCK)
	return float64(sc_clk_tck)
}
