package main

// This program computes queue length for M/M/m queues using Erlang C and two
// approximations of it, for a variety of utilizations and number of servers.

import "fmt"
import "math"

func main() {

	utils := make([]float64, 0)
	for u := 0.0; u < 0.999999; u += 0.1 {
		utils = append(utils, u)
	}
	utils = append(utils, 0.99)
	utils = append(utils, 0.999)

	// m is number of servers
	servers := make([]float64, 0)
	for m := 1.0; m <= 64; m *= 2 {
		servers = append(servers, m)
	}

	fmt.Println("servers\tutil\terlang\tgunther\tsakasegawa")

	for _, m := range servers {
		for _, u := range utils {
			e := erlang(u, m)
			g := gunther(u, m)
			s := sakasegawa(u, m)
			fmt.Printf("%.f\t%.3f\t%f\t%f\t%f\n", m, u, e, g, s)
		}
	}

}

// The Erlang C formula for queue length. u=utilization and m=servers.
// See erlang.pl in Analyzing Computer System Performance with Perl::PDQ.
func erlang(u, m float64) float64 {
	erlangs := u * m
	erlangB := erlangs / (1.0 + erlangs)
	for i := 2.0; i <= m; i++ {
		eb := erlangB
		erlangB = eb * erlangs / (i + (eb * erlangs))
	}

	erlangC := erlangB / (1.0 - u + (u * erlangB))
	return u * m * erlangC / (m * (1.0 - u))
}

// Sakasegawa's approximation to the queue length.
// See Sakasegawa (1977) or (1982).
func sakasegawa(u, m float64) float64 {
	return math.Pow(u, math.Sqrt(2*(m+1.0))) / (1 - u)
}

// Neil Gunther's generalization of the "stretch factor" formula to m servers,
// solved for queue length.  See eq 2.63 in Analyzing Computer System
// Performance with Perl::PDQ.
func gunther(u, m float64) float64 {
	return u * m * (1/(1-math.Pow(u, m)) - 1)
}
