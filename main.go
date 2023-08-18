package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"update-roa/env"
	"update-roa/s3"
)

func init() {
	env.Init()
	s3.Init()

}
func main() {
	ReadPerline()
	for _, object := range ObjectList {
		prefixList := BGPQ3OutPut(object)
		//name: object-v4.juniper.update
		//name: object-v6.juniper.update
		//upload to minio
		//string to file
		v4File := StrToTempFile(prefixList.IPv4List, prefixList.IPv4PrefixName+".juniper.update")
		v6File := StrToTempFile(prefixList.IPv6List, prefixList.IPv6PrefixName+".juniper.update")
		//upload to minio
		file := &s3.File{
			Name: prefixList.IPv4PrefixName + ".juniper.update",
			Path: v4File.Name(),
		}
		file.Upload()
		file = &s3.File{
			Name: prefixList.IPv6PrefixName + ".juniper.update",
			Path: v6File.Name(),
		}
		file.Upload()

	}

}

var ObjectList []string

type PrefixList struct {
	IPv4PrefixName string
	IPv6PrefixName string
	IPv4List       string
	IPv6List       string
}

func ReadPerline() {
	//read ./update.list
	file, err := os.Open("./update.lists")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	//read file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ObjectList = append(ObjectList, scanner.Text())
	}

}

func BGPQ3OutPut(object string) PrefixList {
	//exec bgpq3
	var prefixList PrefixList

	if IsASN(object) {
		prefixList.IPv4PrefixName = asntoasset(object) + "-v4"
		prefixList.IPv6PrefixName = asntoasset(object) + "-v6"
	} else {
		prefixList.IPv4PrefixName = object + "-v4"
		prefixList.IPv6PrefixName = object + "-v6"
	}
	v4Cmd, err := exec.Command("bgpq3", "-4", "-J", "-l", prefixList.IPv4PrefixName, object).Output()
	if err != nil {
		fmt.Println(err)
	}
	v6Cmd, err := exec.Command("bgpq3", "-6", "-J", "-l", prefixList.IPv6PrefixName, object).Output()
	if err != nil {
		fmt.Println(err)
	}
	prefixList.IPv4List = string(v4Cmd)
	prefixList.IPv6List = string(v6Cmd)
	return prefixList

}

func StrToTempFile(str string, name string) *os.File {
	//create temp file
	tmpFile, err := os.CreateTemp("", name)
	if err != nil {
		panic(err)
	}
	//write string to temp file
	_, err = tmpFile.WriteString(str)
	if err != nil {
		panic(err)
	}

	return tmpFile
}

func asntoasset(asn string) string {
	//if string was like as1234 convert to as-1234
	asnRegex := "as([0-9]{1,6})"
	re := regexp.MustCompile(asnRegex)
	return re.ReplaceAllString(asn, "as-$1")
}

func IsASN(asn string) bool {
	//if string was like as1234 convert to as-1234
	asnRegex := "as([0-9]{1,6})"
	re := regexp.MustCompile(asnRegex)
	return re.MatchString(asn)
}
