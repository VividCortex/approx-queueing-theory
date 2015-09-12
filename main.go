package main

// This program computes queue length for M/M/m queues using Erlang C and two
// approximations of it, for a variety of utilizations and number of servers.

import "fmt"
import "math"

func main() {

	utils := make([]float64, 0)
	for u := 0.05; u < 1; u += 0.05 {
		utils = append(utils, u)
	}
	utils = append(utils, 0.99)
	utils = append(utils, 0.999)

	fmt.Println("servers\tutil\terlang\tgunther\tsakasegawa")

	// m is number of servers
	for m := 1.0; m <= 64; m *= 2 {
		for _, u := range utils {
			e := erlang(u, m)
			g := gunther(u, m)
			s := sakasegawa(u, m)
			fmt.Printf("%.f\t%.3f\t%f\t%f\t%f\n", m, u, e, g, s)
		}
	}

}

// The Erlang C formula. u=utilization and m=servers.
// See erlang.pl in Analyzing Computer System Performance with Perl::PDQ.
func erlang(u, m float64) float64 {
	erlangs := u * m
	erlangB := erlangs / (1.0 + erlangs)
	for i := 2.0; i <= m; i++ {
		eb := erlangB
		erlangB = eb * erlangs / (i + (eb * erlangs))
	}

	erlangC := erlangB / (1.0 - u + (u * erlangB))
	return erlangC / (m * (1.0 - u))
}

// Sakasegawa's approximation to the queue length.
// See Sakasegawa (1976) or (1982).
func sakasegawa(u, m float64) float64 {
	return (math.Pow(u, math.Sqrt(2*(m+1.0)))) / (1 - u)
}

// Neil Gunther's generalization of the "stretch factor" formula to m servers.
// See eq 2.63 in Analyzing Computer System Performance with Perl::PDQ.
func gunther(u, m float64) float64 {
	return 1/(1-math.Pow(u, m)) - 1
}
