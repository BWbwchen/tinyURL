package main

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-zookeeper/zk"
)

var (
	// CharacterSet consists of 62 characters [0-9][A-Z][a-z].
	Base         = 62
	CharacterSet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

var LENGTH int = 7

var URL = os.Getenv("ZOOKEEPER_URL")
var zookeeperPath = "/seed"
var zookeeper *zk.Conn
var counterRange int = 100
var counterNow int = 0
var counterBase int

var counterNowLock sync.Mutex

func InitZookeeper() {
	c, _, err := zk.Connect([]string{URL}, time.Second) //*10)
	zookeeper = c
	if err != nil {
		panic(err)
	}
	// zookeeper register
	data := []byte("0")
	zookeeper.Create(zookeeperPath, data, 0, zk.WorldACL(zk.PermAll))
	counterBase = getCounter()
}

func GetShortName(longURL string) string {
	counterNowLock.Lock()
	counter := counterNow + counterBase*counterRange
	counterNow += 1
	counterNowLock.Unlock()
	hash := md5.Sum([]byte(strconv.Itoa(counter)))

	shortName := hex.EncodeToString(hash[:])
	// TODO: maybe can do better
	// check collision
	i := 0
	for {
		_, err := DatabaseGet(shortName[i : i+LENGTH])
		if err != nil {
			break
		}
		i += 1
	}

	// update counter
	if counterNow == counterRange {
		counterNow = 0
		counterBase = getCounter()
	}

	return shortName[i : i+LENGTH]
}

func getCounter() int {
	// get the counter number
	counterByteArray, _, err := zookeeper.Get(zookeeperPath)
	if err != nil {
		panic(err)
	}
	counter, _ := strconv.Atoi(string(counterByteArray))

	zookeeper.Set(zookeeperPath, []byte(strconv.Itoa(counter+1)), -1)
	return counter
}

/*
func base62(num int) string {
	b := make([]byte, 0)

	// loop as long the num is bigger than zero
	for num > 0 {
		// receive the rest
		r := math.Mod(float64(num), float64(Base))

		// devide by Base
		num /= Base

		// append chars
		b = append([]byte{CharacterSet[int(r)]}, b...)
	}

	return string(b)
}
*/
