package functions

import (
	"Toy_Cryptocurrency/models"
	"image/color"
	"io/ioutil"
	"math/rand"
	"runtime"
	"strconv"
	"time"

	"github.com/fogleman/gg"
)

func CreateSecurityCode(user models.User) (int, string) {
	// Generar 4 dígitos de código de seguridad
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	firstDigit := 1 + generator.Intn(9)
	secondDigit := 1 + generator.Intn(9)
	thirdDigit := 1 + generator.Intn(9)
	fourthDigit := 1 + generator.Intn(9)
	stringCode := strconv.Itoa(firstDigit) + strconv.Itoa(secondDigit) + strconv.Itoa(thirdDigit) + strconv.Itoa(fourthDigit)

	// Guardar código en un archivo para verificar luego
	securityCode := []byte(stringCode)
	securityCodeTxtRoute := ""
	if runtime.GOOS == "windows" {
		securityCodeTxtRoute = "security_codes/users_codes_texts/" + user.Email + ".txt"
	} else {
		securityCodeTxtRoute = "/home/piero/Toy_Cryptocurrency_Backend/security_codes/users_codes_texts/" + user.Email + ".txt"
	}
	err := ioutil.WriteFile(securityCodeTxtRoute, securityCode, 0777)
	if err != nil {
		return 1, ""
	}

	// Generar imagen con código de seguridad
	securityCodeImageRoute := ""
	if runtime.GOOS == "windows" {
		securityCodeImageRoute = "security_codes/base_image.jpg"
	} else {
		securityCodeImageRoute = "/home/piero/Toy_Cryptocurrency_Backend/security_codes/base_image.jpg"
	}
	bgImage, err := gg.LoadImage(securityCodeImageRoute)
	if err != nil {
		return 2, ""
	}
	imgWidth := bgImage.Bounds().Dx()
	imgHeight := bgImage.Bounds().Dy()
	dc := gg.NewContext(imgWidth, imgHeight)
	dc.DrawImage(bgImage, 0, 0)
	securityCodeFontRoute := ""
	if runtime.GOOS == "windows" {
		securityCodeFontRoute = "security_codes/amasis_MT_bold.ttf"
	} else {
		securityCodeFontRoute = "/home/piero/Toy_Cryptocurrency_Backend/security_codes/amasis_MT_bold.ttf"
	}
	if err := dc.LoadFontFace(securityCodeFontRoute, 200); err != nil {
		return 3, ""
	}
	x := float64(imgWidth / 2)
	y := float64(imgHeight / 2)
	maxWidth := float64(imgWidth) - 60.0
	dc.SetColor(color.White)
	dc.DrawStringWrapped(stringCode, x, y, 0.5, 0.5, maxWidth, 1.5, gg.AlignCenter)
	securityCodeUserImageRoute := ""
	if runtime.GOOS == "windows" {
		securityCodeUserImageRoute = "security_codes/users_codes_images/" + user.Email + ".png"
	} else {
		securityCodeUserImageRoute = "/home/piero/Toy_Cryptocurrency_Backend/security_codes/users_codes_images/" + user.Email + ".png"
	}
	_ = dc.SavePNG(securityCodeUserImageRoute)

	return 0, securityCodeUserImageRoute
}
