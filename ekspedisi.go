package main


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

func IsSicepat(resi string) bool {
    if string(resi[0:2]) == "00" {
        if len(resi) == 12 {
            return true
        }
    }

    return false
}

func CekResi(resi string) string {
    if len(resi) < 5 {
        return "none"
    }

    if IsJnt(resi) == true {
        return "jnt"
    }

    if IsSicepat(resi) == true {
        return "sicepat"
    }
    return "none"
}

