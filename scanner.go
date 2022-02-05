package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/fatih/color"
)

type User struct {
	Name     string   `json:"name"`
	JnT      []string `json:"jnt"`
	SiCepat  []string `json:"sicepat"`
	AnterAja []string `json:"anteraja"`
	Wahana   []string `json:"wahana"`
}

//create a map for storing clear funcs
var clear map[string]func()

func init() {
	// Initialize
	clear = make(map[string]func())
	// Linux
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	// Windows
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	} else {
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func GetList(path string, kurir string) []string {
	data, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	fS := bufio.NewScanner(data)
	fS.Split(bufio.ScanLines)
	var list []string
	for fS.Scan() {
		if kurir == "all" {
			list = append(list, fS.Text())
		} else {
			if fS.Text()[len(fS.Text())-2:] == kurir {
				list = append(list, fS.Text())
			}
		}
	}
	return list
}

func resiExistHandler() bool {
	var answer string

	color.Red("[!!] Resi already exist!")
	fmt.Fprintf(color.Output, "Proceed with scanning?[%s/%s]\n", color.GreenString("y"), color.RedString("n"))
	for {
		fmt.Printf(">> ")
		fmt.Scan(&answer)
		if answer == "n" {
			return false
		} else if answer == "y" {
			return true
		} else {
			color.Red("Answer not valid!")
		}
	}
}

func CreateFile(currentTime time.Time, users []User) {
	for i := 0; i < len(users); i++ {
		folder := fmt.Sprintf("./log/%s/", currentTime.Format("02-01-2006"))
		file := fmt.Sprintf("%s.log", users[i].Name)
		pathLog := fmt.Sprintf("%s%s", folder, file)
		if _, err := os.Stat(folder); os.IsNotExist(err) {
			os.MkdirAll(folder, 0777)
		}

		fl, err := os.OpenFile(pathLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			log.Fatal(err)
		}
		defer fl.Close()
	}
}

func Pause() {
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func PrintTotal(users []User) {
	currentTime := time.Now()

	fmt.Println("\n====================")
	total := 0
	fmt.Println("JnT")
	fmt.Println("====================")
	for i := 0; i < len(users); i++ {
		folder := fmt.Sprintf("./log/%s/", currentTime.Format("02-01-2006"))
		file := fmt.Sprintf("%s", users[i].Name)
		path := fmt.Sprintf("%s%s.log", folder, file)
		list := GetList(path, "01")
		total = total + len(list)
		fmt.Printf("%s\t: %d\n", users[i].Name, len(list))
		if i == len(users)-1 {
			fmt.Printf("\nTotal: %d\n", total)
		}
	}

	total = 0
	fmt.Println("\n====================")
	fmt.Println("SiCepat")
	fmt.Println("====================")
	for i := 0; i < len(users); i++ {
		folder := fmt.Sprintf("./log/%s/", currentTime.Format("02-01-2006"))
		file := fmt.Sprintf("%s", users[i].Name)
		path := fmt.Sprintf("%s%s.log", folder, file)
		list := GetList(path, "02")
		total = total + len(list)
		fmt.Printf("%s\t: %d\n", users[i].Name, len(list))
		if i == len(users)-1 {
			fmt.Printf("\nTotal: %d\n\n", total)
		}
	}
}

func main() {
	var users []User

	jsonFile, err := os.Open("./db/data.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	fl, err := ioutil.ReadAll(jsonFile)
	_ = json.Unmarshal([]byte(fl), &users)

	keyUser := 0
	isChangeUser := false

	user := users[0].Name

	currentTime := time.Now()

	CreateFile(currentTime, users)

	for {
		// File
		folder := fmt.Sprintf("./log/%s/", currentTime.Format("02-01-2006"))
		file := fmt.Sprintf("%s.log", users[keyUser].Name)
		pathLog := fmt.Sprintf("%s%s", folder, file)
		if _, err := os.Stat(folder); os.IsNotExist(err) {
			os.MkdirAll(folder, 0777)
		}

		fl, err := os.OpenFile(pathLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			log.Fatal(err)
		}
		defer fl.Close()

		list := GetList(pathLog, "all")

		// Logging
		logger := log.New(fl, "", log.LstdFlags)

		var input string

		fmt.Printf("User [%d. %s] (%d)\n", keyUser+1, user, len(list))
		fmt.Printf(">> ")

		fmt.Scanf("%s", &input)

		userList := []string{"1", "2", "3", "4", "5"}
		_, isChangeUser = Find(userList, input)
		if isChangeUser == true {
			keyUser, _ = Find(userList, input)
			user = users[keyUser].Name
			continue
		} else {
			if input == "q" || input == "quit" {
				break
			} else if input == "c" || input == "clear" {
				CallClear()
			} else if input == "p" || input == "print" {
                PrintTotal(users)
			} else if input == "w" || input == "write" {
            } else {
				resi := string(input)
				kurir := CekResi(resi)

				if kurir == "none" {
					fmt.Println("Courier not found!")
				} else if kurir == "jnt" {
					_, found := Find(users[keyUser].JnT, input)
					if !found {
						users[keyUser].JnT = append(users[keyUser].JnT, input)
						logger.Println(input, "01")
					} else {
						answer := resiExistHandler()
						if answer == false {
							break
						}
					}
				} else if kurir == "sicepat" {
					_, found := Find(users[keyUser].SiCepat, input)
					if !found {
						users[keyUser].SiCepat = append(users[keyUser].SiCepat, input)
						logger.Println(input, "02")
					} else {
						answer := resiExistHandler()
						if answer == false {
							break
						}
					}
				}
			}
		}
	}

	fl, _ = json.MarshalIndent(users, "", " ")
	_ = ioutil.WriteFile("./db/data.json", fl, 0644)

    PrintTotal(users)

	Pause()
}
