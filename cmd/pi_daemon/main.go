package main

import (
	"fmt"
	"time"
	"os"
	"gocv.io/x/gocv"

	"github.com/kinedwi/ACS_t/internal/db"
	"github.com/kinedwi/ACS_t/internal/face"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: daemon <db_path> <cascade_xml>")
		return
	}
	dbpath := os.Args[1]
	cascadePath := os.Args[2]

	cascade, err := face.LoadCascade(cascadePath)
	if err != nil {
		panic(err)
	}
	defer cascade.Close()

	database, err := db.New(dbpath)
	if err != nil {
		panic(err)
	}

	webcam, err := gocv.OpenVideoCapture(0)
	if err != nil {
		panic(err)
	}
	defer webcam.Close()
	fmt.Println("Демон запущен. Нажми Ctrl+C для остановки.")

	img := gocv.NewMat()
	defer img.Close()

	// порог: подобрать опытным путём (меньше = более строгий)
	const threshold = 0.6

	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Println("Не удалось прочитать кадр")
			time.Sleep(200 * time.Millisecond)
			continue
		}
		if img.Empty() {
			continue
		}

		rect, found := face.DetectFirstFace(img, cascade)
		if !found {
			time.Sleep(200 * time.Millisecond)
			continue
		}
		faceMat := img.Region(rect)
		desc := face.ComputeDescriptor(faceMat)
		faceMat.Close()

		name, dist, any, err := database.FindBestMatch(desc)
		if err != nil {
			fmt.Println("DB error:", err)
			time.Sleep(time.Second)
			continue
		}
		if any && dist <= threshold {
			fmt.Printf("Доступ: %s (dist=%.4f)\n", name, dist)
			database.LogEvent(name, fmt.Sprintf("access_granted (d=%.4f)", dist))
			// TODO: включить реле (GPIO)
		} else {
			fmt.Printf("Неизвестно (best=%s d=%.4f)\n", name, dist)
			database.LogEvent("Unknown", fmt.Sprintf("access_denied (d=%.4f)", dist))
			database.AddAlert(fmt.Sprintf("Неизвестное лицо (d=%.4f)", dist))
		}

		time.Sleep(700 * time.Millisecond)
	}
}
