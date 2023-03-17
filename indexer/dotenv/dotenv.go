package dotenv

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// a very minimal env parser
func Load(path string) {
	if path == "" {
		path = ".env"
	}
	log.Printf("Loading env from path: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error loading env: %s\n", err)
	}
	defer f.Close()

	// scan the file line by line
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") {
			// split the line by the first =
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				// set the env variable
				os.Setenv(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error loading env: %s\n", err)
	}

	log.Println("Env loaded")
}

// retrieve env with a default value
func Get(key, def string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	return val
}
