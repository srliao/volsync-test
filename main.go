package main

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/caarlos0/env/v10"
)

type config struct {
	Host    string `env:"HOST"`
	Port    string `env:"PORT" envDefault:"3000"`
	DataDir string `env:"DATA_DIR" envDefault:"./"`
}

type server struct {
	cfg config
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	slog.Info("running with config", "config", cfg)

	s := &server{cfg: cfg}
	http.HandleFunc("/", s.handler())
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%v:%v", cfg.Host, cfg.Port), nil))
}

func (s *server) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(s.cfg.DataDir, "data.txt")
		// try reading from file; if none found write
		if _, err := os.Stat(path); err == nil {
			data, err := os.ReadFile(path)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "file %v exists but read failed with error: %v", path, err)
				return
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "file %v exists; content: %s", path, data)
			return
		} else if errors.Is(err, os.ErrNotExist) {
			// path/to/whatever does *not* exist
			content := time.Now().Format(time.RFC822Z)
			err := os.WriteFile(path, []byte(content), 0666)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "file %v does not exists and write failed with error: %v", path, err)
				return
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "file %v does not exist; new content written: %v", path, content)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("unexpected file may or may not exist"))
			return
		}
	}
}
