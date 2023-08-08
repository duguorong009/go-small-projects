package main

import (
	"errors"

	"github.com/gomodule/redigo/redis"
)

// Declare a pool variable to hold the pool of Redis connections.
var pool *redis.Pool

var ErrNoAlbum = errors.New("no album found")

// Define a custom struct to hold Album data.
type Album struct {
	Title  string  `redis:"title"`
	Artist string  `redis:"artist"`
	Price  float64 `redis:"price"`
	Likes  int     `redis:"likes"`
}

func FindAlbum(id string) (*Album, error) {
	// Use the connection pool's Get() method to fetch a single Redis
	// connection from the pool.
	conn := pool.Get()

	// Importantly, use defer and the connection's Close() method to
	// ensure that the connection is always returned to the pool before
	// FindAlbum() exits.
	defer conn.Close()

	// Fetch the details of a specific album. If no album is found
	// the given id, the []interface{} slice returned by redis.Values
	// will have a length of zero. So check fo this and return an
	// ErrNoAlbum error as necessary.
	values, err := redis.Values(conn.Do("HGETALL", "album:"+id))
	if err != nil {
		return nil, err
	} else if len(values) == 0 {
		return nil, ErrNoAlbum
	}

	var album Album
	err = redis.ScanStruct(values, &album)
	if err != nil {
		return nil, err
	}

	return &album, nil
}

func IncrementLikes(id string) error {
	conn := pool.Get()
	defer conn.Close()

	// Before we do anything else, check that an album with the given
	// id exists. The EXISTS command returns 1 if a specific key exists
	// in the database, and 0 if it doesn't.
	exists, err := redis.Int(conn.Do("EXISTS", "album:"+id))
	if err != nil {
		return err
	} else if exists == 0 {
		return ErrNoAlbum
	}

	// Use the MULTI command to inform Redis that we are starting a new
	// transaction. The conn.Send() method writes the command to the
	// connection's output buffer -- it doesn't actually send it to the
	// Redis server.
	err = conn.Send("MULTI")
	if err != nil {
		return err
	}

	// Increment the number of likes in the album hash by 1. Because it
	// follows a MULTI command, this HINCRBY command is not executed but
	// it is QUEUED as part of the transaction. We still need to check
	// the reply's Err filed at this point in case there was a problem
	// queueing the command.
	err = conn.Send("HINCRBY", "album:"+id, "likes", 1)
	if err != nil {
		return err
	}

	// And we do the same with the increment on our sorted test.
	err = conn.Send("ZINCRBY", "likes", 1, id)
	if err != nil {
		return err
	}

	// Execute both commands in our transaction together as an atomic
	// group. EXEC returns the replies from both commands but, because
	// we're not interested in either reply in this example, it
	// suffices to simply check for any errors. Note that calling the
	// conn.Do() method flushes the previous commands from the
	// connection output buffer and sends them to the Redis server.
	_, err = conn.Do("EXEC")
	if err != nil {
		return err
	}

	return nil
}
