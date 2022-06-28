package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type KerncraftConfig struct {
	PerformanceModels []string `json:"performanceModels"`
	CachePredictors   []string `json:"cachePredictors"`
	Units             []string `json:"units"`
}

var kerncraftConfig KerncraftConfig

type RunConfiguration struct {
	Kernel           string   `json:"kernel"`
	Machine          string   `json:"machine"`
	PerformanceModel string   `json:"performanceModel"`
	CachePredictor   string   `json:"cachePredictor"`
	Cores            int      `json:"cores"`
	VarNames         []string `json:"varNames"`
	VarValues        []int    `json:"varValues"`
	Unit             string   `json:"unit"`
}

func populateKerncraftConfigDefault() {

}

func main() {
	fmt.Println("Starting application...")

	fmt.Println("Loading Kerncraft config...")
	{
		data, err := ioutil.ReadFile("./kerncraftConfig.json")
		if err != nil {
			fmt.Println("Error reading kerncraft config file; falling back to default values... Error: " + err.Error())
			kerncraftConfig.PerformanceModels = []string{"ECM"}
			kerncraftConfig.CachePredictors = []string{"LC"}
			kerncraftConfig.Units = []string{"cy/CL"}
		} else {
			err = json.Unmarshal(data, &kerncraftConfig)
			if err != nil {
				fmt.Println("Error parsing JSON from ./kerncraftConfig.json; falling back to default values... Error: " + err.Error())
				populateKerncraftConfigDefault()
			}
		}
	}

	fmt.Println("Registering routes...")
	r := mux.NewRouter()
	r.HandleFunc("/examples/machines", machinesHandler)
	r.HandleFunc("/available/models", modelsHandler)
	r.HandleFunc("/available/units", unitsHandler)
	r.HandleFunc("/available/cachepredictors", cachepredictorsHandler)
	r.HandleFunc("/examples/kernels", kernelsHandler)
	r.HandleFunc("/examples/kernels/{name}", kernelHandler)
	r.HandleFunc("/analyze/{folderID}", runAnalyzer)

	fmt.Println("Registering CORS...")
	http.Handle("/", r)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080", "*"},
		AllowCredentials: true,
	})
	fmt.Println("Running server on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", c.Handler(r)))
}

func machinesHandler(w http.ResponseWriter, _ *http.Request) {
	var machineFiles []string
	files, err := ioutil.ReadDir("./kerncraft/examples/machine-files")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		arr, _ := json.Marshal([]string{})
		_, err = w.Write(arr)
		if err != nil {
			fmt.Println("Error writing machines: " + err.Error())
		}
		return
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".yml") {
			machineFiles = append(machineFiles, strings.TrimSuffix(file.Name(), ".yml"))
		}
	}
	w.WriteHeader(http.StatusOK)
	arr, _ := json.Marshal(machineFiles)
	_, err = w.Write(arr)
	if err != nil {
		fmt.Println("Error writing machines: " + err.Error())
	}
}

func modelsHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	arr, _ := json.Marshal(kerncraftConfig.PerformanceModels)
	_, err := w.Write(arr)
	if err != nil {
		fmt.Println("Error writing models: " + err.Error())
	}
}

func unitsHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	arr, _ := json.Marshal(kerncraftConfig.Units)
	_, err := w.Write(arr)
	if err != nil {
		fmt.Println("Error writing units: " + err.Error())
	}
}

func cachepredictorsHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	arr, _ := json.Marshal(kerncraftConfig.CachePredictors)
	_, err := w.Write(arr)
	if err != nil {
		fmt.Println("Error writing cache predictors: " + err.Error())
	}
}

func kernelsHandler(w http.ResponseWriter, _ *http.Request) {
	var kernelFiles []string
	files, err := ioutil.ReadDir("./kerncraft/examples/kernels")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		arr, _ := json.Marshal(kernelFiles)
		_, err = w.Write(arr)
		if err != nil {
			fmt.Println("Error writing kernel files: " + err.Error())
		}
		return
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".c") {
			kernelFiles = append(kernelFiles, strings.TrimSuffix(file.Name(), ".c"))
		}
	}
	w.WriteHeader(http.StatusOK)
	arr, _ := json.Marshal(kernelFiles)
	_, err = w.Write(arr)
	if err != nil {
		fmt.Println("Error writing kernel files: " + err.Error())
	}
}

func kernelHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	data, err := ioutil.ReadFile("./kerncraft/examples/kernels/" + vars["name"] + ".c")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	content := string(data)
	w.WriteHeader(http.StatusOK)
	arr, _ := json.Marshal(content)
	_, err = w.Write(arr)
	if err != nil {
		fmt.Println("Error writing kernel: " + err.Error())
	}
}

func runAnalyzer(w http.ResponseWriter, r *http.Request) {
	data, _ := io.ReadAll(r.Body)
	var configuration RunConfiguration
	_ = json.Unmarshal(data, &configuration)

	vars := mux.Vars(r)
	folderID := vars["folderID"]

	{
		// validate data
		_, err := strconv.Atoi(folderID)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}

	err := os.MkdirAll("./tmp/" + folderID, os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	file, err := os.Create("./tmp/" + folderID + "/kernel.c")
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}(file)

	err = os.WriteFile("./tmp/" + folderID + "/kernel.c", []byte(configuration.Kernel), 0755)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("File written (supposedly)")

	var defines []string
	for i := 0; i < len(configuration.VarNames); i++ {
		defines = append(defines, "-D", configuration.VarNames[i], strconv.Itoa(configuration.VarValues[i]))
	}

	var args []string
	args = append(append(append(args, "./tmp/" + folderID + "/kernel.c",
		"--machine", "./kerncraft/examples/machine-files/"+configuration.Machine+".yml",
		"--pmodel", configuration.PerformanceModel),
		defines...),
		"-vvv",
		"--pointer-increment", "auto",
		"-u", configuration.Unit,
		"-c", strconv.Itoa(configuration.Cores),
		"--clean-intermediates",
		"-P", configuration.CachePredictor,
		"-C", "gcc",
		"--json", "./tmp/" + folderID + "/output.json")

	cmd := exec.Command("kerncraft", args...)
	fmt.Println(cmd.String())

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	//data, _ = ioutil.ReadFile("./tmp/" + folderID + "/output.json")
	//_, _ = w.Write(data)
	arr, _ := json.Marshal(outb.String() + "\n\n\n" + errb.String())
	_, _ = w.Write(arr)

	// cleanup
	err = os.RemoveAll("./tmp/" + folderID)
	if err != nil {
		fmt.Println("Error cleaning up temporary folder: " + err.Error())
	}
}
