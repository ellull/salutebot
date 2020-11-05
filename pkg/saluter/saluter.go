package saluter

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// Salute represents a salute with the language it corresponds to
type Salute struct {
	Language string
	Text     string
}

// Saluter sends a random Salute to a chat
type Saluter struct {
	salutes []Salute
}

// ChatMessage is a struct that represents a Google Chat message
type ChatMessage struct {
	Text string `json:"text"`
}

// Salute sends a random Salute
func (s *Saluter) Salute(url string) error {
	// pick a random salute
	salute := s.salutes[rand.Intn(len(s.salutes))]

	// build the chat message
	message := ChatMessage{Text: fmt.Sprintf("%s (in %s)", salute.Text, salute.Language)}

	// marshall the message as JSON
	bs, err := json.Marshal(message)
	if err != nil {
		return nil
	}

	// send the message to the chat
	r, err := http.Post(url, "application/json", bytes.NewBuffer(bs))
	if err != nil {
		return err
	}
	defer r.Body.Close()

	// parse the response
	if r.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}

	return nil
}

// NewFileSaluter creates a Saluter loading its Salutes from a CSV file in the format
// language, text
func NewFileSaluter(filename string) (*Saluter, error) {
	// Open file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Setup CSV reader
	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.Comment = '#'
	reader.FieldsPerRecord = 2

	// Read CSV file
	salutes := make([]Salute, 0)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		salutes = append(salutes, Salute{record[0], record[1]})
	}

	return &Saluter{salutes}, nil
}
