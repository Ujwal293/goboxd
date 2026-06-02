package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type RunRequest struct {
	Language string `json:"language"`
	Source   string `json:"source"`
	Stdin    string `json:"stdin"`
}

type RunResponse struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

func runHandler(w http.ResponseWriter, r *http.Request) {

	var req RunRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tempDir, err := os.MkdirTemp("", "goboxd-*")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer os.RemoveAll(tempDir)

	var cmd *exec.Cmd

	if req.Language == "py3" {

		filePath := tempDir + "/solution.py"

		err = os.WriteFile(filePath, []byte(req.Source), 0644)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cmd = exec.Command("python", filePath)

	} else if req.Language == "cpp" {

		cppPath := tempDir + "/main.cpp"

		exePath := tempDir + "/main.exe"

		err = os.WriteFile(cppPath, []byte(req.Source), 0644)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		compileCmd := exec.Command("g++", cppPath, "-o", exePath)

		compileOutput, err := compileCmd.CombinedOutput()

		if err != nil {

			response := RunResponse{
				Stdout: "",
				Stderr: string(compileOutput),
			}

			json.NewEncoder(w).Encode(response)
			return
		}

		cmd = exec.Command(exePath)

	} else {

		http.Error(w, "unsupported language", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel()

	cmd = exec.CommandContext(ctx, cmd.Path, cmd.Args[1:]...)
	cmd.Stdin = strings.NewReader(req.Stdin)
	output, err := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {

		response := RunResponse{
			Stdout: "",
			Stderr: "time limit exceeded",
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	response := RunResponse{
		Stdout: string(output),
	}

	if err != nil {
		response.Stderr = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}
func healthz(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}
func readyz(w http.ResponseWriter, r *http.Request) {

	pythonErr := exec.Command("python3", "--version").Run()

	gppErr := exec.Command("g++", "--version").Run()

	response := map[string]string{
		"python": "ok",
		"g++":    "ok",
		"status": "ready",
	}

	if pythonErr != nil {
		response["python"] = "missing"
		response["status"] = "not ready"
	}

	if gppErr != nil {
		response["g++"] = "missing"
		response["status"] = "not ready"
	}

	json.NewEncoder(w).Encode(response)
}
func main() {
	http.HandleFunc("/healthz", healthz)
	http.HandleFunc("/readyz", readyz)
	http.HandleFunc("/run", runHandler)

	fmt.Println("Server running on port 8080")

	http.ListenAndServe(":8080", nil)
}
