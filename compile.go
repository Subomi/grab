package grab

import (
	"io"
	"log"
	"os"
)

func Compile(files []string) bool {
	var done bool

	file, err := os.OpenFile(files[0], os.O_APPEND|os.O_WRONLY, 0755)

	if err != nil {
		log.Fatal(err)
	}

	for s := 1; s < len(files); s++ {
		part, err := os.OpenFile(files[s], os.O_RDONLY, 0755)

		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(file, part)

		if err != nil {
			log.Fatalf(err.Error())
		}

		if err = part.Close(); err != nil {
			log.Fatal(err.Error())
		}

		if err = os.Remove(files[s]); err != nil {
			log.Fatal(err.Error())
		}
	}

	done = true
	return done
}
