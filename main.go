package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Body struct {
	ProductIdList []string `json:"productIdList"`
}

type BodyUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type BodyNewUser struct {
	User    Information `json:"user"`
	Password string `json:"password"`
}

type Information struct {
	Information    InformationContent `json:"information"`
}

type InformationContent struct {
	Name	string 	`json:"name"`
	Email    string `json:"email"`
	Phone string `json:"phone"`
	Cpf string `json:"cpf"`
	Sex string `json:"sex"`
	Birthdate string `json:"birthdate"`
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, project-id")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	godotenv.Load()
	r := gin.Default()
	r.Use(CORSMiddleware())
	r.GET("/api", func(c *gin.Context) {
		host := strings.ReplaceAll(c.DefaultQuery("url", ""), "\"", "")
		url := strings.ReplaceAll(host, " ", "%20")

		resp, err := http.Get(url)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.Data(200, "application/json", body)
	})

	r.POST("/api", func(c *gin.Context) {
		host := strings.ReplaceAll(c.DefaultQuery("url", ""), "\"", "")
		url := strings.ReplaceAll(host, " ", "%20")

		var BodyImpl Body

		err := c.ShouldBindJSON(&BodyImpl)

		if err != nil {
			c.JSON(400, gin.H{
				"error": "cannot bind JSON: " + err.Error(),
			})
			return
		}

		postBody, _ := json.Marshal(map[string]interface{}{
			"productIdList": BodyImpl.ProductIdList,
		})

		payload := bytes.NewBuffer(postBody)

		resp, err := http.Post(url, "application/json", payload)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.Data(200, "application/json", body)
	})

	r.POST("/api/user", func(c *gin.Context) {
		client := http.Client{}
		host := strings.ReplaceAll(c.DefaultQuery("url", ""), "\"", "")
		url := strings.ReplaceAll(host, " ", "%20")

		var BodyImpl BodyUser

		err := c.ShouldBindJSON(&BodyImpl)

		if err != nil {
			c.JSON(400, gin.H{
				"error": "cannot bind JSON: " + err.Error(),
			})
			return
		}

		postBody, _ := json.Marshal(map[string]interface{}{
			"email":    BodyImpl.Email,
			"password": BodyImpl.Password,
		})

		payload := bytes.NewBuffer(postBody)

		resp, err := http.NewRequest("POST", url, payload)

		resp.Header = http.Header{
			"Content-Type": {"application/json"},
			"project-id":   {"pglv"},
		}
		res, err := client.Do(resp)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.Data(200, "application/json", body)
	})

	r.POST("/api/newuser", func(c *gin.Context) {
		client := http.Client{}
		host := strings.ReplaceAll(c.DefaultQuery("url", ""), "\"", "")
		url := strings.ReplaceAll(host, " ", "%20")

		var BodyImpl BodyNewUser

		err := c.ShouldBindJSON(&BodyImpl)

		if err != nil {
			c.JSON(400, gin.H{
				"error": "cannot bind JSON: " + err.Error(),
			})
			return
		}

		postBody, _ := json.Marshal(BodyImpl)

		payload := bytes.NewBuffer(postBody)

		resp, err := http.NewRequest("POST", url, payload)

		resp.Header = http.Header{
			"Content-Type": {"application/json"},
			"project-id":   {"pglv"},
		}
		res, err := client.Do(resp)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.Data(200, "application/json", body)
	})
	r.Run()
}
