package scann

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

// PrintFile : test
func PrintFile(name string, dir bool) {
	fmt.Printf("%v ::%v\n", name, dir)
}

// MoveFile : test
func MoveFile(inc chan int, mu *sync.Mutex, dirOrigin string, name string, dir bool, dirTarget string) {

	var t = time.Now()

	//Lee las propiedades del archivo
	f, t, err := decodeFile(dirOrigin, name)

	if err == nil {
		t2, _ := f.DateTime()
		if year, _ := strconv.Atoi(t2.Format("2006")); year > 1900 {
			t = t2
		}
	}

	mu.Lock()
	defer mu.Unlock()

	//fmt.Printf("From: %v (%v) To: %v \n", path.Join(dirOrigin, name), t, setPathTarget(inc, t, dirTarget, name))
	moveFileTarget(path.Join(dirOrigin, name), setPathTarget(inc, t, dirTarget, name), dirTarget)

}

//func decodeFile(dirOrigin string, name string) (*exif.Exif, error) {
func decodeFile(dirOrigin string, name string) (*exif.Exif, time.Time, error) {
	f, err := os.Open(path.Join(dirOrigin, name))
	if err != nil || f == nil {
		log.Fatal(err)
	}

	var st syscall.Stat_t
	if err = syscall.Stat(path.Join(dirOrigin, name), &st); err != nil {
		log.Fatal(err)
	}

	t := time.Unix(int64(st.Ctimespec.Sec), int64(st.Ctimespec.Nsec))

	x, err := exif.Decode(f)

	return x, t, err
}

func setPathTarget(inc chan int, t time.Time, dirTarget string, name string) string {
	//Valida el nombre que se deberá de colocar en el archivo

	//val := t.Format("2006-01-02-15-04-05-") + strconv.Itoa(count)
	//Valida si el archivo ya exíste, si ya exíste le agrega el contador
	filename := t.Format("2006-01-02-15-04-05")
	dirname := path.Join(dirTarget, filename[:7])

	//Valida que el directorio exísta
	_, err := os.Stat(dirname)
	if err != nil {
		//Crea el directorio
		os.MkdirAll(dirname, os.ModePerm)
	}

	filename = path.Join(dirname, filename)
	_, err = os.Stat(filename + path.Ext(name))

	if err == nil {
		//El archivo existe, hay que renombrarlo
		count := <-inc
		count++
		inc <- count

		filename = fmt.Sprintf("%v_%v", filename, count)
	}

	// le coloca la extensión
	filename = fmt.Sprintf("%v%v", filename, path.Ext(name))

	return filename
}

func moveFileTarget(fileOrigin string, fileTarget string, dirError string) {
	err := os.Rename(fileOrigin, fileTarget)

	if err != nil {
		err = os.Rename(fileOrigin, path.Join(dirError, "ERROR", fileOrigin))
		fmt.Printf("Error From: %v \n", fileOrigin)
	}
}
