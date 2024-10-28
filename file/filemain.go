package filemain

import (
	"fmt"
	"log"
	"os"
)

func CreatFile() {
	fmt.Printf("++++++++++write file+++++++++++")
	file, err := os.Create("FILE.TXT")

	if err != nil {
		log.Fatalf("ERROR in creating file")
	}

	defer file.Close()

	len, err := file.WriteString("hi my name is saksham, this is my file" + " how does it look")

	if err != nil {
		log.Fatalf("ERROR in writing")
	}

	fmt.Printf("\n File Name: %s", file.Name())
	fmt.Printf("\n Length of file: %d", len)
}

func ReadFile() {
	fmt.Printf("\n++++++++++read file+++++++++++")
	fileName := "FILE.TXT"
	data, err := os.ReadFile("FILE.TXT")

	if err != nil {
		log.Fatal("ERROR in reading")
	}

	fmt.Printf("\n File name: %s", fileName)
	fmt.Printf("\n Length of filr : %d", len(data))
	fmt.Printf("\n DATA: %s", data)

}

func Main() {
	CreatFile()
	ReadFile()
}
