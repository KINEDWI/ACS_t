package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"image"
	"gocv.io/x/gocv"
	"image/color"

	"github.com/kinedwi/ACS_t/internal/db"
	"github.com/kinedwi/ACS_t/internal/face"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: add_user <db_path> <cascade_xml>")
		return
	}
	dbpath := os.Args[1]
	cascadePath := os.Args[2]

	cascade, err := face.LoadCascade(cascadePath)
	if err != nil {
		fmt.Println("Cascade load error:", err)
		return
	}
	defer cascade.Close()

	database, err := db.New(dbpath)
	if err != nil {
		panic(err)
	}

	webcam, err := gocv.OpenVideoCapture(0)
	if err != nil {
		fmt.Println("cant open camera:", err)
		return
	}
	defer webcam.Close()

	window := gocv.NewWindow("Add user - press Space to capture, ESC to quit")
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	fmt.Print("Введите имя нового пользователя: ")
	reader := bufio.NewReader(os.Stdin)
	nameRaw, _ := reader.ReadString('\n')
	name := strings.TrimSpace(nameRaw)

	fmt.Println("Поднесите лицо к камере. Нажмите Space для захвата.")

	var lastRect image.Rectangle
	var found bool

	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Println("Ошибка чтения с камеры")
			continue
		}
		if img.Empty() {
			continue
		}
		// покажем прямоугольник лица (если найден)
		rect, ok := face.DetectFirstFace(img, cascade)
		if ok {
			lastRect = rect
			found = true
			gocv.Rectangle(&img, rect, color.RGBA{0, 255, 0, 0}, 2)
		} else {
			found = false
		}
		window.IMShow(img)
		k := window.WaitKey(1)
		if k == 27 { // ESC
			fmt.Println("Отмена")
			return
		}
		if k == 32 { // SPACE
			if !found {
				fmt.Println("Лицо не найдено, попробуй ещё раз")
				continue
			}
			faceMat := img.Region(lastRect)
			desc := face.ComputeDescriptor(faceMat)
			// добавляем в БД
			if err := database.AddUser(name, desc); err != nil {
				fmt.Println("Ошибка добавления пользователя:", err)
			} else {
				fmt.Println("Пользователь добавлен.")
			}
			faceMat.Close()
			return
		}
	}
}
