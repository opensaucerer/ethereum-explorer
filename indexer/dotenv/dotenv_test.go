package dotenv

import (
	"log"
	"os"
	"testing"
)

func setup() {
	// create .env.test file
	f, err := os.Create(".env.test")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// write some variables
	_, err = f.WriteString("TEST_VAR1=1\n")
	if err != nil {
		log.Fatal(err)
	}

	// write some variables
	_, err = f.WriteString("#TEST_VAR2=2\n")
	if err != nil {
		log.Fatal(err)
	}
}

func teardown() {
	// remove .env.test file
	err := os.Remove(".env.test")
	if err != nil {
		log.Fatal(err)
	}
}

func TestDotenv(t *testing.T) {

	setup()

	t.Run("Should load variables into environment", func(t *testing.T) {
		Load(".env.test")
	})

	t.Run("Should retrieve variables from environment", func(t *testing.T) {
		Load(".env.test")
		if Get("TEST_VAR1", "") != "1" {
			t.Error("Expected TEST_VAR1 to be 1 but got", Get("TEST_VAR1", ""))
		}
	})

	t.Run("Should retrieve default value if variable is not set", func(t *testing.T) {
		Load(".env.test")
		if Get("TEST_VAR2", "null") != "null" {
			t.Error("Expected TEST_VAR2 to be null but got", Get("TEST_VAR2", "null"))
		}
	})

	teardown()
}
