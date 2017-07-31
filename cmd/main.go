package main

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"time"
)

var (
	ErrTimestamp      = errors.New("timestamp must be greater than the current time")
	ErrEventListEmpty = errors.New("at least one event is required")
)

type Event struct {
	timestamp time.Time
	action    func()
}

func NewEvent(timestamp time.Time, f func()) (event *Event, err error) {
	if timestamp.Before(time.Now()) {
		return nil, ErrTimestamp
	}

	return &Event{
		timestamp: timestamp,
		action:    f,
	}, nil
}

func (event Event) AfterCurrentTime() bool {
	return event.timestamp.After(time.Now())
}

type Cron struct {
	timer  *time.Timer
	stopC  chan struct{}
	events []Event // TODO 链表
}

func NewCron(events ...Event) (*Cron, error) {
	if len(events) <= 0 {
		return nil, ErrEventListEmpty
	}
	// 按时间排序
	sort.Slice(events, func(i, j int) bool {
		return events[i].timestamp.Before(events[j].timestamp)
	})

	c := &Cron{
		events: events,
		timer:  time.NewTimer(events[0].timestamp.Sub(time.Now())),
		stopC:  make(chan struct{}, 1),
	}
	go c.run()
	return c, nil
}

func (c *Cron) run() {
	for i := range c.events {
		select {
		case <-c.timer.C:
			c.events[i].action()

			if i < len(c.events)-1 {
				c.timer.Reset(c.events[i+1].timestamp.Sub(time.Now()))
			}
		case <-c.stopC:
			log.Println("The stop signal has been received")
			break
		}
	}
	log.Println("Stop run")
}

func (c *Cron) AddEvent(event Event) (err error) { // TODO 添加事件方法
	panic("unsupported")
}

func (c *Cron) Stop() bool {
	c.stopC <- struct{}{}
	close(c.stopC)
	return c.timer.Stop()
}

func main() {
	now := time.Now()
	fmt.Printf("Now: %d\n", now.UnixNano())

	e0, _ := NewEvent(now.Add(time.Second*2), func() {
		log.Printf("The event0 has been fired(%d)\n", time.Now().UnixNano())
	})

	e1, _ := NewEvent(now.Add(time.Millisecond*500), func() {
		log.Printf("The event1 has been fired(%d)\n", time.Now().UnixNano())
	})

	e2, _ := NewEvent(now.Add(time.Millisecond*505), func() {
		log.Printf("The event2 has been fired(%d)\n", time.Now().UnixNano())
	})

	e3, _ := NewEvent(now.Add(time.Millisecond*505), func() {
		log.Printf("The event3 has been fired(%d)\n", time.Now().UnixNano())
	})

	e4, _ := NewEvent(now.Add(time.Second*5), func() {
		log.Printf("The event4 has been fired(%d)\n", time.Now().UnixNano())
	})

	c, err := NewCron(*e0, *e1, *e2, *e3, *e4)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(3 * time.Second)

	c.Stop()

	time.Sleep(time.Second)
}
