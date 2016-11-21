package db

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

var (
	safari  = []byte("safari")
	members = map[string]string{
		"Hil":   "Warthog",
		"Pete":  "Springbok",
		"Tuck":  "Buffalo",
		"Eli":   "Impala",
		"Liz":   "Zebra",
		"Berto": "Hippo",
	}
)

// Validate checks for a correct name and animal match
func Validate(name, animal string) bool {
	if strings.ToLower(members[name]) != strings.ToLower(animal) {
		return false
	}
	return true
}

// GetPair returns matching animal
func GetPair(name string) string {
	return members[name]
}

// RandomInsert assigns secret pairs
func RandomInsert() error {
	db, err := bolt.Open("./bolt.db", 0644, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(safari)
		if err != nil {
			return err
		}
		nameList := shuffledCopy(members)
		var match string
		for n := range members {
			match, nameList = makeMatch(n, nameList)
			err = bucket.Put([]byte(n), []byte(match))
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

// GetMatch returns the current match for a person
func GetMatch(name string) (match string, err error) {
	db, err := bolt.Open("./bolt.db", 0644, nil)
	if err != nil {
		return "", err
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(safari)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", safari)
		}
		val := bucket.Get([]byte(name))
		match = string(val)
		return nil
	})

	if err != nil {
		return "", err
	}

	return match, nil
}

// GetAll returns the current list inside a bucket
func GetAll() (map[string]string, error) {
	list := make(map[string]string)
	db, err := bolt.Open("./bolt.db", 0644, nil)
	if err != nil {
		return list, err
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(safari)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", safari)
		}

		for n := range members {
			val := bucket.Get([]byte(n))
			list[n] = string(val)
		}
		return nil
	})

	return list, err
}

func shuffledCopy(members map[string]string) []string {
	i := 0
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	length := len(members)
	names := make([]string, length)
	perm := random.Perm(length)
	for n := range members {
		names[perm[i]] = n
		i++
	}
	return names
}

func makeMatch(person string, list []string) (string, []string) {
	for i := 0; i < len(list); i++ {
		if person != list[i] {
			name := list[i]
			list = append(list[:i], list[i+1:]...)
			return name, list
		} else if person == list[i] && len(list) == 1 {
			log.Fatal("Last person got themselves")
		}
	}
	return "", list
}
