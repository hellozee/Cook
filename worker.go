package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	ps "local.proj/Cook/parser"
)

func structToMap(parsedStruct []entity) {
	for _, item := range parsedStruct {
		oldfileTimings[item.File] = item.Hash
	}
}

func generateList() {
	for _, value := range fileList {
		file, err := os.Stat(value)
		checkErr(err)
		t := file.ModTime()
		hash := hashTime(t.String())

		newfileTimings[value] = hash
		hashJSONnew.Body.Entity = append(hashJSONnew.Body.Entity, entity{File: value, Hash: hash})
	}
}

func compileFirst() {
	//Iteratively generate .o files

	for key, value := range fileList {
		fmt.Println("Compiling " + value)
		cmd := exec.Command(ps.CompilerDetails.Binary, "-c", file, "-o", "Cooking/"+tag+".o")
		checkCommand(cmd)
	}
}

func compareAndCompile() {
	//Compare the file hash with current hash if do not match generate .o file
	//also replace the current hash with the new hash

	for key, value := range fileList {

		file, err := os.Stat(value)

		checkErr(err)
		t := file.ModTime()
		timeStamp := strings.Replace(t.String(), " ", "", -1)

		if !checkTimeStamp(timeStamp, oldfileTimings[value]) {
			fmt.Println("Compiling " + value)
			cmd := exec.Command(ps.CompilerDetails.Binary, "-c", value, "-o", "Cooking/"+key+".o")
			checkCommand(cmd)

			oldfileTimings[value] = hashTime(t.String())
		}

		hashJSONnew.Body.Entity = append(hashJSONnew.Body.Entity, entity{File: value, Hash: oldfileTimings[value]})
	}
}

func linkAll() {

	//Compile all the generated .o files under the Cooking directory
	fmt.Println("Linking files..")
	args := []string{ps.CompilerDetails.Binary, "-o", ps.CompilerDetails.Name, ps.CompilerDetails.Includes, ps.CompilerDetails.OtherFlags,
		"Cooking/*.o", ps.CompilerDetails.LdFlags}
	cmd := exec.Command(os.Getenv("SHELL"), "-c", strings.Join(args, " "))
	checkCommand(cmd)
}
