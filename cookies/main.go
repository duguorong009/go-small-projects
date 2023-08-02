package main

import (
	"bytes"
	"cookies-example-project/internal/cookies"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// type Cookie struct {
// 	Name  string
// 	Value string

// 	Path       string
// 	Domain     string
// 	Expires    time.Time
// 	RawExpires string

// 	// MaxAge=0 means no 'Max-Age' attribute specified.
// 	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
// 	// MaxAge>0 means Max-Age attribute present and given in seconds
// 	MaxAge   int
// 	Secure   bool
// 	HttpOnly bool
// 	SameSite SameSite
// 	Raw      string
// 	Unparsed []string
// }

// Declare a global variable to hold the secret key.
var secretKey []byte

// Declare the User type.
type User struct {
	Name string
	Age  int
}

func main() {
	// Importantly, we need to tell the encoding/gob package about the Go type
	// that we want to encode. We do this my passing *an instance* of the type
	// to gob.Register(). In this case we pass a pointer to an initialized (but
	// empty) instance of the User struct.
	gob.Register(&User{})

	var err error

	// Decode the random 64-character hex string to give us a slice containing
	// 32 random bytes. For simplicity, I've hardcoded this hex string but in a
	// real application you should read it in at runtime from a command-line
	// flag or environment variable.
	secretKey, err = hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/set", setCookieHandler)
	mux.HandleFunc("/get", getCookieHandler)

	log.Print("Listening...")
	err = http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func setCookieHandler(w http.ResponseWriter, r *http.Request) {
	// Initialize a User struct containing the data that we want to store in the
	// cookie.
	user := User{Name: "Alice", Age: 21}

	// Initialize a buffer to hold the gob-encoded data.
	var buf bytes.Buffer

	// Gob-encode the user data, storing the encoded output in the buffer.
	err := gob.NewEncoder(&buf).Encode(&user)
	if err != nil {
		log.Println(err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	// Call buf.String() to get the gob-encoded value as a string and set it as
	// the cookie value.
	cookie := http.Cookie{
		Name:     "exampleCookie",
		Value:    buf.String(),
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	// Write an encrypted cookie containing the gob-encoded data as normal.
	err = cookies.WriteEncrypted(w, cookie, secretKey)

	if err != nil {
		log.Println(err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("cookie set!"))
}

func getCookieHandler(w http.ResponseWriter, r *http.Request) {
	// Read the gob-encoded value from the encrypted cookie, handling any errors
	// as necessary.
	gobEncodedValue, err := cookies.ReadEncrypted(r, "exampleCookie", secretKey)

	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "cookie not found", http.StatusBadRequest)
		case errors.Is(err, cookies.ErrInvalidValue):
			http.Error(w, "invalid cookie", http.StatusBadRequest)
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	// Create a new instance of a User type.
	var user User

	// Create an strings.Reader containing the gob-encoded value.
	reader := strings.NewReader(gobEncodedValue)

	// Decode it into the User type. Notice that we need to pass a *pointer* to
	// the Decode() target here?
	if err := gob.NewDecoder(reader).Decode(&user); err != nil {
		log.Println(err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	// Print the user information in the response.
	fmt.Fprintf(w, "Name: %q\n", user.Name)
	fmt.Fprintf(w, "Age: %d\n", user.Age)
}

//
// Version 2. Use the base64 encode/decode
//
// func setCookieHandler(w http.ResponseWriter, r *http.Request) {
// 	// Initialize the cookie as normal.
// 	cookie := http.Cookie{
// 		Name:     "exampleCookie",
// 		Value:    "Hello Zoë!",
// 		Path:     "/",
// 		MaxAge:   3600,
// 		HttpOnly: true,
// 		Secure:   true,
// 		SameSite: http.SameSiteLaxMode,
// 	}

// 	// Write the cookie. If there is an error (due to an encoding failure or it
// 	// being too long) then log the error and send a 500 Internal Server Error
// 	// response.
// 	err := cookies.Write(w, cookie)
// 	if err != nil {
// 		log.Println(err)
// 		http.Error(w, "server error", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Write([]byte("cookie set!"))
// }

// func getCookieHandler(w http.ResponseWriter, r *http.Request) {
// 	// Use the Read() function to retrieve the cookie value, additionally
// 	// checking for the ErrInvalidValue error and handling it as necessary.
// 	value, err := cookies.Read(r, "exampleCookie")
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, http.ErrNoCookie):
// 			http.Error(w, "cookie not found", http.StatusBadRequest)
// 		case errors.Is(err, cookies.ErrInvalidValue):
// 			http.Error(w, "invalid cookie", http.StatusBadRequest)
// 		default:
// 			log.Println(err)
// 			http.Error(w, "server error", http.StatusInternalServerError)
// 		}
// 		return
// 	}

// 	w.Write([]byte(value))
// }

//
// Version 1. Simple Cookie
//
// func setCookieHandler(w http.ResponseWriter, r *http.Request) {
// 	// Initialize a new cookie containing the string "Hello world!" and some
// 	// non-default attributes.
// 	cookie := http.Cookie{
// 		Name:     "exampleCookie",
// 		Value:    "Hello world!",
// 		Path:     "/",
// 		MaxAge:   3600,
// 		HttpOnly: true,
// 		Secure:   true,
// 		SameSite: http.SameSiteLaxMode,
// 	}

// 	// Use the http.SetCookie() function to send the cookie to the client.
// 	// Behind the scenes this adds a `Set-Cookie` header to the response
// 	// containing the necessary cookie data.
// 	http.SetCookie(w, &cookie)

// 	// Write a HTTP response as normal.
// 	w.Write([]byte("cookie set!"))
// }

// func getCookieHandler(w http.ResponseWriter, r *http.Request) {
// 	// Retrieve the cookie from the request using its name(which in our case is
// 	// "exampleCookie"). If no matching cookie is found, this will return a
// 	// http.ErrNoCookie error. We check for this, and return a 400 Bad Request
// 	// response to the client.
// 	cookie, err := r.Cookie("exampleCookie")
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, http.ErrNoCookie):
// 			http.Error(w, "cookie not found", http.StatusBadRequest)
// 		default:
// 			log.Println(err)
// 			http.Error(w, "server error", http.StatusInternalServerError)
// 		}
// 		return
// 	}

// 	// Echo out the cookie value in the response body
// 	w.Write([]byte(cookie.Value))
// }
