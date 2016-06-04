package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

var databasePath string = "/Users/Chris/uptime/bolt.db"

func WebsitesToMonitor() error {

	db, err := bolt.Open(databasePath, 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("websites"))
		if b != nil {
			b.ForEach(func(k, v []byte) error {

				var website Website
				if err := json.Unmarshal(v, website); err != nil {
					return err
				}

				Websites = append(Websites, website)
				return err
			})
		} else {
			b, err = tx.CreateBucket([]byte("websites"))
			fmt.Println("'websites' bucket was created. Please add websites to it now.")
		}
		return err
	})

	return err
}

func AddToDatabase(website Website) error {
	db, err := bolt.Open(databasePath, 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	key, err := IntToBytes(website.Id)
	if err != nil {
		return err
	}
	value, err := json.Marshal(&website)
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("websites"))
		if err != nil {
			return err
		}

		err = bucket.Put(key, value)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func RemoveFromDatabase(website Website) error {
	return nil
}

func IntToBytes(i int) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, i)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
