package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gomodule/redigo/redis"
)

// Define a custom struct to hold Album data.
type Album struct {
	Title  string
	Artist string
	Price  float64
	Likes  int
}

func main() {
	// Establish a connection to the Redis server listening on port
	// 6379 of the local machine. 6379 is the default port, so unless
	// you've already changed the Redis configuration file this should
	// work.
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatal(err)
	}

	// Importantly, use defer to ensure the connection is always
	// properly closed before exiting the main() function.
	defer conn.Close()

	// redisSetExample(conn)

	// redisGetAndConversionExample(conn)

	// Fetch all album fields with the HGETALL command. Because HGETALL
	// returns an array reply, and because the underlying data structure
	// in Redis is a hash, it makes sense to use the Map() helper
	// function to convert the reply to a map[string]string.
	reply, err := redis.StringMap(conn.Do("HGETALL", "album:1"))
	if err != nil {
		log.Fatal(err)
	}

	// Use the populateAlbum helper function to create a new Album
	// object from the map[string]string.
	album, err := populateAlbum(reply)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", album)

}

// Create, populate and return a pointer to a new Album struct, based
// on data from a map[string]string.
func populateAlbum(reply map[string]string) (*Album, error) {
	var err error
	album := new(Album)
	album.Title = reply["title"]
	album.Artist = reply["artist"]
	// We need to use the strconv package to convert the 'price' value
	// from a string to a float64 before assigning it.
	album.Price, err = strconv.ParseFloat(reply["price"], 64)
	if err != nil {
		return nil, err
	}

	// similaryly, we need to convert the 'likes' value from a string to
	// an integer
	album.Likes, err = strconv.Atoi(reply["likes"])
	if err != nil {
		return nil, err
	}

	return album, nil
}

func redisSetExample(conn redis.Conn) {
	// Send our command across the connection. The first parameter to
	// Do() is always the name of the Redis command(in this example
	// HMSET), optionally followed by any necessary arguments (in this
	// example the key, followed by the various hash fields and values).
	_, err := conn.Do("HMSET", "album:2", "title", "Electric Ladyland", "artist", "Jimi Hendrix", "price", 4.95, "likes", 8)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Electric Ladyland added!")
}

func redisGetAndConversionExample(conn redis.Conn) {
	// Issue a HGET command to retrieve the title for a specific album,
	// and use the Str() helper method to convert the reply to a string.
	title, err := redis.String(conn.Do("HGET", "album:1", "title"))
	if err != nil {
		log.Fatal(err)
	}

	// Similarly, get the artist and convert it to a string.
	artist, err := redis.String(conn.Do("HGET", "album:1", "artist"))
	if err != nil {
		log.Fatal(err)
	}

	// And the price as a float64...
	price, err := redis.Float64(conn.Do("HGET", "album:1", "price"))
	if err != nil {
		log.Fatal(err)
	}

	// And the number of likes as an integer
	likes, err := redis.Int(conn.Do("HGET", "album:1", "likes"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s by %s: £%.2f [%d likes]\n", title, artist, price, likes)
}
