package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

func main() {
	ch := make(chan int, 1)
	c := make(chan os.Signal, 1)
	chExit := make(chan int, 1)
	chExitMain := make(chan int, 1)
	signal.Notify(c, os.Interrupt)
	ticker := time.NewTicker(1 * time.Second)

	countPing := 0
	countPong := 0
	isShowPong := false

	go func() {
	LOOP1:
		for {
			select {
			case <-ch:
				countPong++
				if isShowPong {
					fmt.Println("pong")
				}
			case s := <-c:
				fmt.Println(s)
				break LOOP1
			default:
			}
			time.Sleep(1000 * time.Millisecond)
		}

		chExit <- 1
		fmt.Println("1 out")
	}()

	go func() {
	LOOP2:
		for {
			select {
			case <-chExit:
				fmt.Println("2 start exit")
				break LOOP2
			case <-ticker.C:
				ch <- 1
				// fmt.Println("ping")
			default:
			}

			time.Sleep(1000 * time.Millisecond)
		}

		chExitMain <- 1
		fmt.Println("2 out")
	}()

	reader := bufio.NewReader(os.Stdin)
LOOP:
	for {
		select {
		case <-chExitMain:
			break LOOP
		default:
			fmt.Print("-> ")
			text, err := reader.ReadString('\n')
			if err != nil {
				panic(err)
			}
			text = strings.Replace(text, "\n", "", -1)

			switch text {
			case "get ping":
				fmt.Printf("ping: %d\n", countPing)
			case "get pong":
				fmt.Printf("pong: %d\n", countPong)
			case "show pong":
				isShowPong = true
			case "hide pong":
				isShowPong = false
			default:
				if strings.HasPrefix(text, "set ping") {
					text = strings.Replace(text, "set ping", "", -1)
					text = strings.Replace(text, " ", "", -1)
					n, err := strconv.ParseInt(text, 10, 64)
					if err != nil {
						fmt.Println(err)
						continue
					}
					fmt.Println(n)
					ticker.Stop()
					ticker = time.NewTicker(time.Duration(n) * time.Second)
					fmt.Println("bingo")
				}
			}
		}
	}
	fmt.Println("main out")
}
