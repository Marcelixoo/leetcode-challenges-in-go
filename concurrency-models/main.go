package main

import (
	"fmt"
	"math/rand"
	"time"
)

/*
Restoring sequencing:

1. Send a channel on a channel, making a goroutin wait til its turn.
2. Receive all messages and then enable them again by sending on a private channel.
3. Controller channel (i.e. waitForIt) is shared between all messages
*/
type Message struct {
	str  string
	wait chan bool
}

type QuitableMessage struct {
	str  string
	quit chan bool
}

type QuitableMessageV2 struct {
	str  string
	quit chan string
}

func main() {
	c := quitableV2("Joe")

	var msg QuitableMessageV2
	for i := rand.Intn(10); i >= 0; i-- {
		msg = <-c
		fmt.Println(msg.str)
	}
	msg.quit <- "Terminating..."
	fmt.Printf("Joe last words are: %q\n", <-msg.quit)
}

func quitableV2(msg string) <-chan QuitableMessageV2 {
	c := make(chan QuitableMessageV2)
	quit := make(chan string)
	go func() {
		for i := 0; ; i++ {
			select {
			case c <- QuitableMessageV2{fmt.Sprintf("%s: %d", msg, i), quit}:
				// do nothing
			case <-quit:
				// chance to clean things up!
				quit <- "See you!"
				return
			}
		}
	}()

	return c
}

func main_quitable() {
	c := quitable("Joe")

	var msg QuitableMessage

	for i := rand.Intn(10); i >= 0; i-- {
		msg = <-c
		fmt.Println(msg.str)
	}

	msg.quit <- true
}

func quitable(msg string) <-chan QuitableMessage {
	c := make(chan QuitableMessage)

	quit := make(chan bool)

	go func() {
		for i := 0; ; i++ {
			select {
			case c <- QuitableMessage{fmt.Sprintf("%s: %d", msg, i), quit}:
				// do nothing
			case <-quit:
				return
			}
		}
	}()

	return c
}

func main_timesoutEntireConversation() {
	c := generator("Joe")

	// Constrain the whole loop to up to 3 seconds.
	// After that period, terminates the entire conversation.
	timeout := time.After(3 * time.Second)
	for {
		select {
		case s := <-c:
			fmt.Println(s.str)
			s.wait <- true
		case <-timeout:
			fmt.Println("You talk too much.")
			return
		}
	}
}

func main_timeoutIfChannelTakesTooLong() {
	c := generator("Joe")

	for {
		select {
		case s := <-c:
			fmt.Println(s.str)
		case <-time.After(1 * time.Second):
			fmt.Println("You're too slow.")
			return
		}
	}
}

// 1. Simple example of consuming 5 messages before exiting
func main_boring() {
	c := make(chan string) // a channel of strings
	go boring("Booooring!", c)

	for i := 0; i < 5; i++ {
		fmt.Printf("You say: %q\n", <-c) // receive expression is just a value
	}

	fmt.Println("You're boring; I'm leaving.")
}

func boring(msg string, c chan string) {
	random := func() time.Duration {
		return time.Duration(rand.Intn(1e3))
	}

	for i := 0; ; i++ {
		c <- fmt.Sprintf("%s %d", msg, i) // any suitable value can be sent through the channel
		time.Sleep(random() * time.Millisecond)
	}
}

// 2. Generator: function that returns a channel
func generator(msg string) <-chan Message { // returns a receive-only channel of strings; both are possible
	c := make(chan Message)

	waitForIt := make(chan bool) // shared among all messages

	go func() { // we launch a goroutine from inside the function
		for i := 0; ; i++ {
			c <- Message{fmt.Sprintf("%s %d", msg, i), waitForIt}
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
			<-waitForIt // i think we're simply discarding messages here to keep processing
		}
	}() // this is an anonymous function immediatelly invoked into a goroutine (like IFFEs)
	return c // return the channel to the caller
}

// 3.
func main_ordering() {
	joe := generator("Joe")
	ann := generator("Ann")

	c := fanIn(joe, ann)

	for i := 0; i < 5; i++ {
		msg1 := <-c
		fmt.Println(msg1.str)
		msg2 := <-c
		fmt.Println(msg2.str)
		msg1.wait <- true
		msg2.wait <- true
	}
	fmt.Println("You're both boring; I'm leaving.")
}

func generator_waitForIt(msg string) <-chan Message { // returns a receive-only channel of strings; both are possible
	c := make(chan Message)

	waitForIt := make(chan bool) // shared among all messages

	go func() { // we launch a goroutine from inside the function
		for i := 0; ; i++ {
			c <- Message{fmt.Sprintf("%s %d", msg, i), waitForIt}
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
			// this seems to be necessary to block the code here until the channel
			// is released to move on, which makes me think that any value traveling
			// through these channels would do.
			<-waitForIt
		}
	}() // this is an anonymous function immediatelly invoked into a goroutine (like IFFEs)
	return c // return the channel to the caller
}

// 3.
func fanIn(input1, input2 <-chan Message) <-chan Message {
	c := make(chan Message)

	// we're returning the return value from input 1 and 2 into the multiplexed channel c
	go func() {
		for {
			c <- <-input1
		}
	}()
	go func() {
		for {
			c <- <-input2
		}
	}()

	return c
}

func main_fanInV2() {
	joe := generator("Joe")
	ann := generator("Ann")
	karl := generator("Kalr")

	c := fanInV2(joe, ann, karl)

	for i := 0; i < 5; i++ {
		msg1 := <-c
		fmt.Println(msg1.str)
		msg2 := <-c
		fmt.Println(msg2.str)
		msg3 := <-c
		fmt.Println(msg3.str)
		msg1.wait <- true
		msg2.wait <- true
		msg3.wait <- true
	}
	fmt.Println("You're both boring; I'm leaving.")
}

func fanInV2(input1, input2, input3 <-chan Message) <-chan Message {
	c := make(chan Message)
	// will block until one of the channels respond
	go func() {
		for {
			select {
			case s := <-input1:
				c <- s
			case s := <-input2:
				c <- s
			case s := <-input3:
				c <- s
			}
		}
	}()
	return c
}

func balance() {
	c1, c2, c3 := make(chan string), make(chan string), make(chan string)

	select {
	case v1 := <-c1:
		fmt.Printf("received %v from c1\n", v1)
	case v2 := <-c2:
		fmt.Printf("received %v from c2\n", v2)
	case v1 := <-c3:
		fmt.Printf("received %v from c1\n", v1)
	default:
		fmt.Printf("no one was ready to communicate\n")
	}
}
