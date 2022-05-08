package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var eng *gin.Engine
var runningSessions []bool

func main() {
	fmt.Println("Starting application...")
	eng = gin.Default()

	eng.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// register getters
	eng.GET("/examples/machines", func(c *gin.Context) {
		c.JSON(getExampleMachines())
	})
	eng.GET("/available/models", func(c *gin.Context) {
		c.JSON(getAvailableModels())
	})
	eng.GET("/available/units", func(c *gin.Context) {
		c.JSON(getAvailableUnits())
	})
	eng.GET("/available/cachepredictors", func(c *gin.Context) {
		c.JSON(getAvailablePredictors())
	})
	eng.GET("/examples/kernels", func(c *gin.Context) {
		c.JSON(getExampleKernels())
	})
	eng.GET("/examples/kernels/:name", func(c *gin.Context) {
		c.JSON(getKernel(c.Param("name")))
	})

	// register posts
	eng.POST("/session", func(c *gin.Context) {
		c.JSON(createSession())
	})

	err := eng.Run("localhost:7248")
	if err != nil {
		log.Fatal(err)
	}
}

// Get configurable data [httpStatus, data]
func getExampleMachines() (int, []string) {
	var machineFiles []string
	files, err := ioutil.ReadDir("./kerncraft/examples/machine-files")
	if err != nil {
		return http.StatusNotFound, machineFiles
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".yml") {
			machineFiles = append(machineFiles, strings.TrimSuffix(file.Name(), ".yml"))
		}
	}
	return http.StatusOK, machineFiles
}
func getAvailableModels() (int, []string) {
	// TODO: Performance model -> required or not? ASK
	return http.StatusOK, []string{"ECM", "ECMData", "ECMCPU", "RooflineASM", "LC", "PerformanceModel"}
}
func getAvailableUnits() (int, []string) {
	return http.StatusOK, []string{"cy/CL", "cy/It", "It/s", "FLOP/s"}
}
func getAvailablePredictors() (int, []string) {
	return http.StatusOK, []string{"LC", "SIM"}
}
func getExampleKernels() (int, []string) {
	var kernelFiles []string
	files, err := ioutil.ReadDir("./kerncraft/examples/kernels")
	if err != nil {
		return http.StatusNotFound, kernelFiles
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".c") {
			kernelFiles = append(kernelFiles, strings.TrimSuffix(file.Name(), ".c"))
		}
	}
	return http.StatusOK, kernelFiles
}
func getKernel(kernelName string) (int, string) {
	bytes, err := ioutil.ReadFile("./kerncraft/examples/kernels/" + kernelName + ".c")
	if err != nil {
		return http.StatusNotFound, ""
	}
	content := string(bytes)
	return http.StatusOK, content
}

// Post session [httpStatus, sessionId]
// TODO: overflow protection
func createSession() (int, int) {
	id := 0
	found := false
	for i := 1; i < len(runningSessions); i++ {
		if !runningSessions[i] {
			runningSessions[i] = true
			found = true
			id = i
		}
	}
	if !found {
		id = len(runningSessions)
		runningSessions = append(runningSessions, true)
	}
	return http.StatusOK, id
}
