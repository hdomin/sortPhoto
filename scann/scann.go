package scann

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
)

// Sobrecarga de interface para el ordenamiento por nombre de los archivos
type byNumericalFilename []os.FileInfo

func (nf byNumericalFilename) Len() int      { return len(nf) }
func (nf byNumericalFilename) Swap(i, j int) { nf[i], nf[j] = nf[j], nf[i] }
func (nf byNumericalFilename) Less(i, j int) bool {
	//Obtiene el path completo
	pathA := strings.ToLower(nf[i].Name())
	pathB := strings.ToLower(nf[j].Name())

	return pathA < pathB
}

// ReadPath : Lee recursivamente los directorios
func ReadPath(semaphore chan struct{}, inc chan int, wg *sync.WaitGroup, mu *sync.Mutex, dirOrigin string, dirTarget string) {
	defer wg.Done()

	//semaphore <- struct{}{} // Espera hasta que se libere un proceso
	semaphore <- struct{}{} // Lock
	defer func() {
		<-semaphore // Libera
	}()

	//lee el contenido del directorio Origen
	files, err := ioutil.ReadDir(dirOrigin)

	if err != nil {
		log.Fatal(err)
	}

	sort.Sort(byNumericalFilename(files))

	for _, f := range files {

		if f.IsDir() {
			wg.Add(1)
			go ReadPath(semaphore, inc, wg, mu, path.Join(dirOrigin, f.Name()), dirTarget)
		} else {
			MoveFile(inc, mu, dirOrigin, f.Name(), f.IsDir(), dirTarget)
		}
	}
}
