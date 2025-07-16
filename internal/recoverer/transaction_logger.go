package recoverer

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/mishaprokop4ik/storage/internal/client"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

// DefaultSaveFile is name of file where recovered data will be sent
const DefaultSaveFile = "recovered"

// Recovered combines http.Method and data.
// data should be in json format string
type Recovered struct {
	method string
	data   string
}

type Actions []string

var DefaultActions = Actions{"p", "d"}

func (a *Actions) in(action string) bool {
	for i := 0; i < len(*a); i++ {
		if (*a)[i] == action {
			return true
		}
	}
	return false
}

// TransactionLogger contains options to recover data
type TransactionLogger struct {
	fileName string
	logger   *log.Logger
	c        chan Recovered
}

// NewTransactionLogger creates new instance of TransactionLogger
// if fileName is empty string will use DefaultSaveFile
// if there is not a file with input fileName
// NewTransactionLogger will create new file by this name
func NewTransactionLogger(fileName string, l *log.Logger) *TransactionLogger {
	if fileName == "" {
		fileName = DefaultSaveFile
	}

	_, err := os.Stat(fileName)
	if errors.Is(err, os.ErrNotExist) {
		_, err := os.Create(fileName)
		if err != nil {
			l.Fatal(err)
		}
	}

	return &TransactionLogger{
		fileName: fileName,
		logger:   l,
		c:        make(chan Recovered),
	}
}

// RecoverData recovers data into file
// if action does not exist in Actions input
// will send error message.
// In correct way RecoverData saves data in file by format:
// action\tdata\n
func (r *TransactionLogger) RecoverData(action, data string, actions Actions) error {
	if !actions.in(action) {
		return fmt.Errorf(`incorrect action type: %s; want one of this: %v`, action, DefaultActions)
	}
	f, err := os.OpenFile(r.fileName,
		os.O_APPEND|os.O_WRONLY, os.ModeAppend)

	if err != nil {
		r.logger.Fatal(err)
	}
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("%s\t%s\n", action, data))
	if err != nil {
		return err
	}
	return nil
}

func (r *TransactionLogger) takeRecovered() {
	_, err := os.Stat(r.fileName)
	if err != nil {
		r.logger.Fatal(err)
	}
	f, err := os.OpenFile(r.fileName,
		os.O_RDONLY, 0644)
	if err != nil {
		r.logger.Fatal(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		err = sc.Err()
		if err != nil {
			r.logger.Fatalf("scan file error: %v", err)
		}
		recoverData := strings.Split(sc.Text(), "\t")
		r.c <- Recovered{
			method: recoverData[0],
			data:   recoverData[1],
		}
	}
	err = sc.Err()
	if err != nil {
		r.logger.Fatalf("scan file error: %v", err)
	}
	err = os.Truncate(r.fileName, 0)
	if err != nil {
		r.logger.Fatal(err)
	}
	r.c <- Recovered{}
}

// SendRecovered takes recovered data and send it to client
// first requests will put after it delete
func (r *TransactionLogger) SendRecovered(port string) {
	go r.takeRecovered()
	client := client.NewAPI(fmt.Sprintf("http://localhost%s", port))
	defer close(r.c)
	wg := &sync.WaitGroup{}
	var toDelete []string
	for {
		recovered := <-r.c
		if (Recovered{}) == recovered {
			break
		}
		switch recovered.method {
		case "p":
			wg.Add(1)
			go func(data string) {
				defer wg.Done()
				resp, err := client.AddOrUpdate(data)
				if err != nil {
					r.logger.Fatal(err)
				}
				if resp.StatusCode != http.StatusCreated {
					r.logger.Println(resp.StatusCode, string(resp.Body))
				}
			}(recovered.data)
		case "d":
			toDelete = append(toDelete, recovered.data)
		}
	}

	wg.Wait()
	for i := 0; i < len(toDelete); i++ {
		go func(data string) {
			resp, err := client.Delete(data)
			if err != nil {
				r.logger.Fatal(err)
			}
			if resp.StatusCode != http.StatusNoContent {
				r.logger.Println(resp.StatusCode, string(resp.Body))
			}
		}(toDelete[i])
	}
}
