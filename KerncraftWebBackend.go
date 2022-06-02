package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html"
	"net/http"
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
	r := mux.NewRouter()
	r.HandleFunc("/products", ProductsHandler)
	r.HandleFunc("/product/{id}", ProductHandler)
	http.Handle("/", r)
}

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func ProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Id: %v\n", vars["id"])
}

/*
var eng *gin.Engine

func main() {
	fmt.Println("Starting application...")
	eng = gin.Default()

	eng.Use(cors.Default())

	// register getters
	registerGetters(eng)

	eng.POST("/test", func(c *gin.Context) {
		var rc RunConfiguration
		decoder := json.NewDecoder(c.Request.Body)
		err := decoder.Decode(&rc)
		if err != nil {
			panic(err)
		}
		c.JSON(runAnalyzer(rc))
	})

	err := eng.Run("localhost:8081")
	if err != nil {
		log.Fatal(err)
	}
}

func registerGetters(engine *gin.Engine) {
	eng.GET("/examples/machines", func(c *gin.Context) { c.JSON(getExampleMachines()) })
	eng.GET("/available/models", func(c *gin.Context) { c.JSON(getAvailableModels()) })
	eng.GET("/available/units", func(c *gin.Context) { c.JSON(getAvailableUnits()) })
	eng.GET("/available/cachepredictors", func(c *gin.Context) { c.JSON(getAvailablePredictors()) })
	eng.GET("/examples/kernels", func(c *gin.Context) { c.JSON(getExampleKernels()) })
	eng.GET("/examples/kernels/:name", func(c *gin.Context) { c.JSON(getKernel(c.Param("name"))) })
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
	// TODO: "Performance model" -> required or not? ASK
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
