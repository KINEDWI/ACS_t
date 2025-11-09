package face

import (
	"errors"
	"gocv.io/x/gocv"
	"image"
	"math"
)

// Descriptor length (кол-во бинов)
const HistBins = 64

type Descriptor []float32

// LoadCascade загружает xml классификатор Haar
func LoadCascade(cascadePath string) (gocv.CascadeClassifier, error) {
	c := gocv.NewCascadeClassifier()
	if !c.Load(cascadePath) {
		return c, errors.New("не удалось загрузить cascade xml: " + cascadePath)
	}
	return c, nil
}

// DetectFirstFace возвращает прямоугольник первой найденной области лица
func DetectFirstFace(img gocv.Mat, cascade gocv.CascadeClassifier) (image.Rectangle, bool) {
	rects := cascade.DetectMultiScale(img)
	if len(rects) == 0 {
		return image.Rectangle{}, false
	}
	return rects[0], true
}

// ComputeDescriptor: crop -> gray -> resize(64x64) -> histogram (HistBins) -> L2-normalize
func ComputeDescriptor(faceMat gocv.Mat) Descriptor {
	// приводим к серому
	gray := gocv.NewMat()
	gocv.CvtColor(faceMat, &gray, gocv.ColorBGRToGray)
	defer gray.Close()

	// resize до фиксированного размера
	size := image.Pt(64, 64)
	resized := gocv.NewMat()
	gocv.Resize(gray, &resized, size, 0, 0, gocv.InterpolationLinear)
	defer resized.Close()

	// вычислим гистограмму: значения 0..255, делим на HistBins
	binCount := HistBins
	binSize := 256 / binCount
	hist := make([]float32, binCount)

	for y := 0; y < resized.Rows(); y++ {
		for x := 0; x < resized.Cols(); x++ {
			pixel := resized.GetUCharAt(y, x)
			bin := int(pixel) / binSize
			if bin >= binCount {
				bin = binCount - 1
			}
			hist[bin] += 1.0
		}
	}

	// L2-нормализация
	var sumSq float32
	for i := 0; i < binCount; i++ {
		sumSq += hist[i] * hist[i]
	}
	norm := float32(math.Sqrt(float64(sumSq)))
	if norm == 0 {
		norm = 1
	}
	for i := 0; i < binCount; i++ {
		hist[i] /= norm
	}
	return Descriptor(hist)
}

// Distance (Euclidean) между двумя дескрипторами
func Distance(a, b Descriptor) float32 {
	if len(a) != len(b) {
		return float32(math.Inf(1))
	}
	var s float32
	for i := 0; i < len(a); i++ {
		diff := a[i] - b[i]
		s += diff * diff
	}
	return float32(math.Sqrt(float64(s)))
}
