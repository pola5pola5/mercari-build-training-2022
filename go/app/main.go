package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

const (
	ImgDir = "image"
)

type Response struct {
	Message string `json:"message"`
}

type Item struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}

type ItemSlice struct {
	Items []Item `json:"items"`
}

func root(c echo.Context) error {
	res := Response{Message: "Hello, world!"}
	return c.JSON(http.StatusOK, res)

}

func getItem(c echo.Context) error {
	var itemSlice ItemSlice

	fp, err := os.ReadFile("./sample.json")
	if err != nil {
		//log.Println("ReadError: ", err)
		fmt.Println("ReadError: ", err)
		return err
	}

	err = json.Unmarshal(fp, &itemSlice)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, itemSlice)
}

func addItem(c echo.Context) error {
	var itemSlice ItemSlice
	// todo: ファイルの存在確認
	fp, err := os.ReadFile("./sample.json")
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = json.Unmarshal(fp, &itemSlice)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusNotFound, err)
	}

	// Get form data
	name := c.FormValue("name")
	category := c.FormValue("category")
	item := Item{name, category}
	// add item to slice
	itemSlice.Items = append(itemSlice.Items, item)

	// save as json
	f, err := os.Create("sample.json")
	if err != nil {
		return err
	}
	err = json.NewEncoder(f).Encode(itemSlice)
	if err != nil {
		return err
	}
	c.Logger().Infof("Receive item: %s %s", item.Name, item.Category)
	message := fmt.Sprintf("item received: %s", item.Name)
	res := Response{Message: message}

	return c.JSON(http.StatusOK, res)
}

func getImg(c echo.Context) error {
	// Create image path
	imgPath := path.Join(ImgDir, c.Param("itemImg"))

	if !strings.HasSuffix(imgPath, ".jpg") {
		res := Response{Message: "Image path does not end with .jpg"}
		return c.JSON(http.StatusBadRequest, res)
	}
	if _, err := os.Stat(imgPath); err != nil {
		c.Logger().Debugf("Image not found: %s", imgPath)
		imgPath = path.Join(ImgDir, "default.jpg")
	}
	return c.File(imgPath)
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Logger.SetLevel(log.INFO)

	front_url := os.Getenv("FRONT_URL")
	if front_url == "" {
		front_url = "http://localhost:3000"
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{front_url},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// Routes
	e.GET("/", root)
	e.GET("/items", getItem)
	e.POST("/items", addItem)
	e.GET("/image/:itemImg", getImg)

	// Start server
	e.Logger.Fatal(e.Start(":9000"))
}
