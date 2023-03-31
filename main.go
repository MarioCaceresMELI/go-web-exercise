package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"os"
	"github.com/gin-gonic/gin"
)

type Product struct{
	Id int `json:"id"`
	Name string `json:"name"`
	Quantity int `json:"quantity"`
	CodeValue string `json:"code_value"`
	IsPublished bool `json:"is_published"`
	Expiration string `json:"expiration"`
	Price float64 `json:"price"`
}

var(
	ErrorCantOpenFile = errors.New("cant open the requested file")
	ErrorJsonParse = errors.New("cant parse the JSON to struct")
)

func main(){
	products, err := readJsonFile("products.json")
	if err!= nil{
		panic(err)
	}
	startServer(products)

}

func readJsonFile(filename string) ([]Product, error){
	var productList []Product
	file,err := os.ReadFile(filename)
	if err!= nil{
		return productList, ErrorCantOpenFile
	}

	err = json.Unmarshal(file,&productList)
	if err!= nil{
		return productList, ErrorJsonParse
	}

	return productList,nil

}

func startServer(products []Product){
	router:=gin.Default()

	//A. GET para ping-pong
	router.GET("/ping", func(c *gin.Context){
		c.String(http.StatusOK,"pong")
	})

	//B. GET para pretty-products
	router.GET("/pretty-products", func(c *gin.Context){
		c.IndentedJSON(http.StatusOK, products)
	})
	//B2. GET para products
	router.GET("/products", func(c *gin.Context){
		c.JSON(http.StatusOK, products)
	})

	//C. GET product by id
	router.GET("/products/:id", func(c *gin.Context){
		productId,_ := strconv.Atoi(c.Param("id"))
		var requestedProduct Product
		for _, product:=range products{
			if product.Id == productId{
				requestedProduct = product
			}
		}
		c.JSON(http.StatusOK, requestedProduct)
	})

	//D. GET search by parameter priceGt
	router.GET("/products/search", func(c *gin.Context){
		priceGt,_ := strconv.ParseFloat(c.Query("priceGt"),64)
		var validProducts []Product
		for _, product:=range products{
			if product.Price > priceGt{
				validProducts = append(validProducts, product)
			}
		}
		c.IndentedJSON(http.StatusOK, validProducts)
	})

	
	router.Run()
}