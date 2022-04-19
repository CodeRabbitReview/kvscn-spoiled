package client

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	address     = "http://localhost:8080/api"
	getAll      = "get all"
	getByID     = "get by id"
	addOrUpdate = "add or update new"
	delete      = "delete"
)

func SendRequestToServer() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("actions to do: \n1. %s\n2. %s\n3. %s\n4. %s\n",
		getAll, getByID, addOrUpdate, delete)
	for {
		action, err := reader.ReadString('\n')
		action = action[:len([]rune(action))-1]
		if err != nil {
			return err
		}
		switch action {
		case getAll, "1":
			if err := GetAll(); err != nil {
				return err
			}
		case delete, "4":
			param, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			if err := Delete(param); err != nil {
				return err
			}
		case getByID, "2":
			param, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			param = param[:len([]rune(param))-1]
			if err := GetByID(param); err != nil {
				return err
			}
		case addOrUpdate, "3":
			param, err := reader.ReadString('~')
			param = param[:len(param)-1]
			if err != nil {
				return err
			}
			if err := AddOrUpdate(param); err != nil {
				return err
			}
		}
	}
}

func GetAll() error {
	resp, err := http.Get(address + "/")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.StatusCode, string(body))
	return nil
}

func Delete(param string) error {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, address+"/"+param, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.StatusCode, string(body))
	return nil
}

func GetByID(param string) error {
	client := &http.Client{}
	fmt.Println(address + "/" + `` + param + ``)
	req, err := http.NewRequest(http.MethodGet, address+"/"+``+param+``, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.StatusCode, string(body))
	return nil
}

func AddOrUpdate(param string) error {
	buf := bytes.NewBuffer([]byte(param))
	req, err := http.Post(address+"/", "application/json;charset=utf-8", buf)
	if err != nil {
		return err
	}
	defer req.Body.Close()
	fmt.Println(req.StatusCode)
	return nil
}
