package main

import (
    "os"
	"fmt"
    "log"
    "time"
    "bufio"
    "runtime"
    "os/exec"
    "io/ioutil"
    "encoding/json"
)

type User struct {
    Name        string      `json:"name"`
    JnT         []string    `json:"jnt"`
    SiCepat     []string    `json:"sicepat"`
    AnterAja    []string    `json:"anteraja"`
    Wahana      []string    `json:"wahana"`
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
    value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
    if ok { //if we defined a clear func for that platform:
        value()  //we execute it
    } else { //unsupported platform
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

func IsJnt(resi string) bool {
    if string(resi[0:4]) == "TJNT" {
        return true
    }

    if len(resi) == 12{
        if string(resi[0:1]) == "J" {
            return true
        }
    }

    return false
}

func CekResi(resi string) string {
    if IsJnt(resi) == true {
        return "jnt"
    }
    return "none"
}

func GetList(path string) []string {
    data, err := os.Open(path)
    if err != nil {
        log.Fatal(err)
    }
    fS := bufio.NewScanner(data)
    fS.Split(bufio.ScanLines)
    var list []string
    for fS.Scan() {
        list = append(list, fS.Text())
    }
    return list
}

func main() {
    jsonFile, err := os.Open("./db/data.json")
    if err != nil {
        log.Fatal(err)
    };defer jsonFile.Close()

    file, err := ioutil.ReadAll(jsonFile)
    var users []User
    _ = json.Unmarshal([]byte(file), &users)

    keyUser := 0
    isChangeUser := false

    user := users[0].Name

    currentTime := time.Now()

    for {
        // File
        folder := fmt.Sprintf("./%s/", currentTime.Format("02-01-2006"))
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

        list := GetList(pathLog)

        // Logging
        logger := log.New(fl, "", log.LstdFlags)

        var input string

        fmt.Printf("User [%d. %s] (%d)\n", keyUser+1, user, len(list))
        fmt.Printf("Input:\n>> ")

        fmt.Scanf("%s", &input)

        userList := []string{"1","2","3","4","5"}
        _, isChangeUser = Find(userList, input)
        if isChangeUser == true {
            keyUser, _ = Find(userList, input)
            user = users[keyUser].Name
            continue
        }else{
            if input == "q"  {
                break
            }else if input == "c" {
                CallClear()
            }else  {
                resi := string(input)
                kurir := CekResi(resi)

                if kurir == "none" {
                    fmt.Println("Courier not found!")
                }else if kurir == "jnt" {
                    _, found := Find(users[keyUser].JnT, input)
                    if !found {
                        users[keyUser].JnT = append(users[keyUser].JnT, input)
                        logger.Println(input)
                    }else {
                        fmt.Printf("!ERROR! Resi already exist!\n")
                    }
                }
            }
        }
    }

    file, _ = json.MarshalIndent(users, "", " ")
    _ = ioutil.WriteFile("./db/data.json", file, 0644)

    total := 0
    for i, _ := range(users) {
        folder := fmt.Sprintf("./%s/", currentTime.Format("02-01-2006"))
        file := fmt.Sprintf("%s", users[i].Name)
        path := fmt.Sprintf("%s%s.log", folder, file)
        list := GetList(path)
        total = total + len(list)
        fmt.Printf("\n%s\t: %d", users[i].Name, len(list))
        if i == len(users)-1 {
            fmt.Printf("\nTotal: %d\n", total)
        }
    }
}
