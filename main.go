package main

import (
	"fmt"
	"os"
	"sync"

	"developer.globalsiag.com/hector/sortPhoto/scann"
)

var dirOrigin string
var dirTarget string

func init() {
	//dirOrigin = "/Users/hector/Pictures/Fanny/2015-08-22"
	dirOrigin = "/Users/hector/Pictures/TEST"
	dirTarget = "/Users/hector/Downloads/test"
}

func main() {

	//Si hay valores por parámetros toma estos
	args := os.Args[1:]
	if len(args) == 2 {
		dirOrigin = args[0]
		dirTarget = args[1]
	} else if len(args) != 0 {
		fmt.Println("Error en los parámetros:  Debe incluír primero el path  Origen y luego el path Destino")
	}

	fmt.Printf("Origen: %v     Destino: %v\n", dirOrigin, dirTarget)

	//Inicia la Rutina,  levanta una rutina por directorio encontrado
	//
	var wg sync.WaitGroup
	var mu sync.Mutex

	semaphore := make(chan struct{}, 10) //Inicializa hasta un máximo de 10 concurrencias
	inc := make(chan int, 1)
	inc <- 1

	wg.Add(1)

	go scann.ReadPath(semaphore, inc, &wg, &mu, dirOrigin, dirTarget)
	wg.Wait()

}
