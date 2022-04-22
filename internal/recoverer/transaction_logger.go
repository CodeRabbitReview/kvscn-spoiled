package recoverer

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/mishaprokop4ik/storage/internal/client"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Recovered struct {
	method string
	data   string
}

type actions []string

var defaultActions = actions{"put", "delete"}

func (a *actions) in(action string) bool {
	for i := 0; i < len(*a); i++ {
		if (*a)[i] == action {
			return true
		}
	}
	return false
}

type Recover struct {
	file   *os.File
	logger *log.Logger
	c      chan Recovered
}

func NewRecover(fileName string, l *log.Logger) (*Recover, error) {
	if fileName == "" {
		return &Recover{
			file:   nil,
			logger: l,
			c:      make(chan Recovered),
		}, nil
	}

	var err error
	var f *os.File
	_, err = os.Stat(fileName)
	if errors.Is(err, os.ErrNotExist) {
		f, err = os.Create(fileName)
	} else if err == nil {
		f, err = os.OpenFile(fileName, os.O_RDWR, os.ModeAppend)
	} else {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return &Recover{
		file:   f,
		logger: l,
		c:      make(chan Recovered),
	}, nil
}

func (r *Recover) RecoverData(action, data string) error {
	if r.file == nil {
		return nil
	}
	removeAllSpaces, err := regexp.Compile(`\r|\t|\n|\t| `)
	if err != nil {
		return err
	}
	if !defaultActions.in(action) {
		return fmt.Errorf("incorrect action type: %s; want one of this: %v",
			action, defaultActions)
	}
	data = removeAllSpaces.ReplaceAllString(data, "")
	_, err = r.file.WriteString(fmt.Sprintf("%s\t%s\n", action, data))
	if err != nil {
		return err
	}
	return nil
}

func (r *Recover) takeRecovered() {
	_, err := os.Stat(r.file.Name())
	if err != nil {
		r.logger.Fatal(err)
	}
	f, err := os.OpenFile(r.file.Name(),
		os.O_APPEND, 0644)
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
		switch recoverData[0] {
		case "put":
			r.c <- Recovered{
				method: http.MethodPut,
				data:   recoverData[1],
			}
		case "delete":
			r.c <- Recovered{
				method: http.MethodDelete,
				data:   recoverData[1],
			}
		}
	}
	err = sc.Err()
	if err != nil {
		r.logger.Fatalf("scan file error: %v", err)
	}
	err = os.Truncate(r.file.Name(), 0)
	if err != nil {
		r.logger.Fatal(err)
	}
	r.c <- Recovered{}
}

func (r *Recover) SendRecovered(addr string) {
	go r.takeRecovered()
	client := client.NewAPI(fmt.Sprintf("http://localhost%s", addr))
	defer close(r.c)
	for {
		recovered := <-r.c
		if (Recovered{}) == recovered {
			break
		}
		go func(recovered Recovered) {
			switch recovered.method {
			case http.MethodPut:
				_, err := client.AddOrUpdate(recovered.data)
				if err != nil {
					r.logger.Fatal(err)
				}
			case http.MethodDelete:
				_, err := client.Delete(recovered.data)
				if err != nil {
					r.logger.Fatal(err)
				}
			}
		}(recovered)
	}
}
