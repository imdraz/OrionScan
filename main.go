package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// localIP devuelve la IP de la red local de esta máquina (por ejemplo, 192.168.1.23).
// No envía paquetes: solo pregunta al sistema qué interfaz usaría para salir a Internet.
// Si no hay red, devuelve 127.0.0.1.
func localIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()

	addr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok {
		return "127.0.0.1"
	}
	return addr.IP.String()
}

// parsePorts convierte "80,443" o "1-1024" en una lista de puertos.
func parsePorts(spec string) ([]int, error) {
	var ports []int
	for _, part := range strings.Split(spec, ",") {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "-") {
			bounds := strings.SplitN(part, "-", 2)
			start, err1 := strconv.Atoi(bounds[0])
			end, err2 := strconv.Atoi(bounds[1])
			if err1 != nil || err2 != nil || start < 1 || end > 65535 || start > end {
				return nil, fmt.Errorf("rango inválido: %q", part)
			}
			for p := start; p <= end; p++ {
				ports = append(ports, p)
			}
		} else {
			p, err := strconv.Atoi(part)
			if err != nil || p < 1 || p > 65535 {
				return nil, fmt.Errorf("puerto inválido: %q", part)
			}
			ports = append(ports, p)
		}
	}
	return ports, nil
}

func main() {
	host := flag.String("host", "", "equipo a revisar (vacío = IP de tu red local detectada automáticamente)")
	portSpec := flag.String("p", "1-1024", "puertos: lista (80,443) o rango (1-1024)")
	timeout := flag.Duration("t", 500*time.Millisecond, "tiempo máximo de espera por puerto")
	workers := flag.Int("w", 500, "conexiones simultáneas")
	flag.Parse()

	// Si no se indicó un equipo, usamos la IP de la máquina donde corre el programa.
	target := *host
	if target == "" {
		target = localIP()
	}

	ports, err := parsePorts(*portSpec)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	fmt.Printf("Objetivo: %s\n", target)
	fmt.Printf("Puertos a revisar: %d\n", len(ports))
	start := time.Now()

	jobs := make(chan int)
	var mu sync.Mutex
	var open []int
	var wg sync.WaitGroup

	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range jobs {
				addr := net.JoinHostPort(target, strconv.Itoa(port))
				conn, err := net.DialTimeout("tcp", addr, *timeout)
				if err == nil {
					conn.Close()
					mu.Lock()
					open = append(open, port)
					mu.Unlock()
				}
			}
		}()
	}

	for _, p := range ports {
		jobs <- p
	}
	close(jobs)
	wg.Wait()

	sort.Ints(open)

	if len(open) == 0 {
		fmt.Println("\nNo se encontraron puertos abiertos en el rango revisado.")
	} else {
		fmt.Println("\nPuertos abiertos:")
		for _, p := range open {
			fmt.Printf("  %5d/tcp  abierto\n", p)
		}
	}
	fmt.Printf("\nListo en %s.\n", time.Since(start).Round(time.Millisecond))
}