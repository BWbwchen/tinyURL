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

func GenerateShortName(longURL string) string {
	shortName := getUniqueShortName()

	UpdateCounterBase()

	return shortName
}

func getUniqueShortName() string {
	counter := getCounter()
	hash := md5.Sum([]byte(strconv.Itoa(counter)))
	candidateName := hex.EncodeToString(hash[:])

	i := 0
	// TODO: maybe can do better
	// check collision
	for {
		if !db.ShortNameExist(candidateName[i : i+LENGTH]) {
			break
		}
		i += 1
	}
	return candidateName[i : i+LENGTH]
}

func getCounter() int {
	counterNowLock.Lock()
	counter := counterNow + counterBase*counterRange
	counterNow += 1
	counterNowLock.Unlock()
	return counter
}

func UpdateCounterBase() {
	if counterNow == counterRange {
		counterNow = 0
		counterBase = getNewCounterBase()
	}
}

func getNewCounterBase() int {
	// get the counter number
	counterByteArray, _, err := zookeeper.Get(zookeeperPath)
	if err != nil {
		panic(err)
	}
	counter, _ := strconv.Atoi(string(counterByteArray))

	zookeeper.Set(zookeeperPath, []byte(strconv.Itoa(counter+1)), -1)
	return counter
}
