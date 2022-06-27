package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

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

func main() {
	fmt.Println("Starting application...")

	r := mux.NewRouter()

	r.HandleFunc("/examples/machines", machinesHandler)
	r.HandleFunc("/available/models", modelsHandler)
	r.HandleFunc("/available/units", unitsHandler)
	r.HandleFunc("/available/cachepredictors", cachepredictorsHandler)
	r.HandleFunc("/examples/kernels", kernelsHandler)
	r.HandleFunc("/examples/kernels/{name}", kernelHandler)

	http.Handle("/", r)
	fmt.Println("Running server on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func machinesHandler(w http.ResponseWriter, r *http.Request) {
	var machineFiles []string
	files, err := ioutil.ReadDir("./kerncraft/examples/machine-files")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		arr, _ := json.Marshal(machineFiles)
		w.Write(arr)
		return
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".yml") {
			machineFiles = append(machineFiles, strings.TrimSuffix(file.Name(), ".yml"))
		}
	}
	w.WriteHeader(http.StatusOK)
	arr, _ := json.Marshal(machineFiles)
	w.Write(arr)
}

func modelsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	arr, _ := json.Marshal([]string{"ECM", "ECMData", "ECMCPU", "RooflineASM", "LC", "PerformanceModel"})
	w.Write(arr)
}

func unitsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	arr, _ := json.Marshal([]string{"cy/CL", "cy/It", "It/s", "FLOP/s"})
	w.Write(arr)
}

func cachepredictorsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	arr, _ := json.Marshal([]string{"LC", "SIM"})
	w.Write(arr)
}

func kernelsHandler(w http.ResponseWriter, r *http.Request) {
	var kernelFiles []string
	files, err := ioutil.ReadDir("./kerncraft/examples/kernels")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		arr, _ := json.Marshal(kernelFiles)
		w.Write(arr)
		return
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".c") {
			kernelFiles = append(kernelFiles, strings.TrimSuffix(file.Name(), ".c"))
		}
	}
	w.WriteHeader(http.StatusOK)
	arr, _ := json.Marshal(kernelFiles)
	w.Write(arr)
}

func kernelHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bytes, err := ioutil.ReadFile("./kerncraft/examples/kernels/" + vars["name"] + ".c")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	content := string(bytes)
	w.WriteHeader(http.StatusOK)
	arr, _ := json.Marshal(content)
	w.Write(arr)
}

/*
// Run
func runAnalyzer(runConfiguration RunConfiguration) (int, string) {
	// TODO: validate input

	err := os.WriteFile("./tmp/kern.c", []byte(runConfiguration.Kernel), 0755)
	if err != nil {
		panic(err)
	}

	var defines []string = []string{}
	for i := 0; i < len(runConfiguration.VarNames); i++ {
		defines = append(defines, "-D", runConfiguration.VarNames[i], strconv.Itoa(runConfiguration.VarValues[i]))
	}

	var args []string = []string{}
	args = append(append(append(args, "./tmp/kern.c",
		"--machine", "./kerncraft/examples/machine-files/"+runConfiguration.Machine+".yml",
		"--pmodel", runConfiguration.PerformanceModel),
		defines...),
		"-vvv",
		"--pointer-increment", "auto",
		"-u", runConfiguration.Unit,
		"-c", strconv.Itoa(runConfiguration.Cores),
		"--clean-intermediates",
		"-P", runConfiguration.CachePredictor,
		"-C", "gcc")

	cmd := exec.Command("kerncraft", args...)
	fmt.Println(cmd.String())
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	return http.StatusOK, outb.String() + "\n\n\n" + errb.String()
}
*/
