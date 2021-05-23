package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/c2h5oh/datasize"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

var (
	port       = "443"
	password   = "nateiscool"
	extensions = []string{"png", "gif", "mp4", "jpg"}
)

func main() {
	engine := html.New("./templates", ".html")

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		Prefork:               true,
		Views:                 engine,
		BodyLimit:             1024 * 1024 * 15, // 15 MB upload limit
	})

	// Raw image GET
	app.Get("/:file.:extension", func(c *fiber.Ctx) error {
		c.Response().Header.Set("Cache-Control", "public, max-age=2629800")

		// Check if extension is in the slice
		if !contains(extensions, c.Params("extension")) {
			return c.Status(400).Type("json", "utf-8").JSON(fiber.Map{
				"success": false,
				"message": "Invalid extension",
			})
		}

		// Check if file exists
		if _, err := os.Stat(fmt.Sprintf("./media/%s/%s.%s", c.Params("extension"), c.Params("file"), c.Params("extension"))); os.IsNotExist(err) {
			return c.Status(404).Type("json", "utf-8").JSON(fiber.Map{
				"success": false,
				"message": "File not found",
			})
		}

		// Return the media file
		return c.Status(200).Type(c.Params("extension")).SendFile(fmt.Sprintf("./media/%s/%s.%s", c.Params("extension"), c.Params("file"), c.Params("extension")))
	})

	// Metadata GET for Discord
	app.Get("/:extension/:file", func(c *fiber.Ctx) error {
		c.Response().Header.Set("Cache-Control", "public, max-age=2629800")

		// Check if the extension is valid
		if !contains(extensions, c.Params("extension")) {
			return c.Status(400).Type("json", "utf-8").JSON(fiber.Map{
				"success": false,
				"message": "Invalid extension",
			})
		}

		// Check if file exists
		info, err := os.Stat(fmt.Sprintf("./media/%s/%s.%s", c.Params("extension"), c.Params("file"), c.Params("extension")))
		if os.IsNotExist(err) {
			return c.Status(404).Type("json", "utf-8").JSON(fiber.Map{
				"success": false,
				"message": "File not found",
			})
		}

		// Calculate size and time
		size := datasize.ByteSize.HumanReadable(datasize.ByteSize(info.Size()))
		ny, _ := time.LoadLocation("America/New_York")
		timestamp := info.ModTime().In(ny)
		timestring := fmt.Sprintf("%s, %s %s %s %02s:%02s:%02s", timestamp.Weekday().String(), strconv.Itoa(timestamp.Day()), timestamp.Month().String(), strconv.Itoa(timestamp.Year()), strconv.Itoa(timestamp.Hour()), strconv.Itoa(timestamp.Minute()), strconv.Itoa(timestamp.Second()))

		// Will present metadata to the Discord scraper
		if strings.Contains(string(c.Request().Header.UserAgent()), "Discord") {
			return c.Status(200).Type("html").Render("metadata", fiber.Map{
				"Filename":  c.Params("file"),
				"Filesize":  size,
				"Filedate":  timestring,
				"Extension": c.Params("extension"),
			})
		}

		// If not Discord scraper then redirect to the raw file
		return c.Redirect(fmt.Sprintf("/%s.%s", c.Params("file"), c.Params("extension")))
	})

	// Upload POST
	app.Post("/upload", func(c *fiber.Ctx) error {
		// Check if the password matches
		if c.FormValue("password") != password {
			return c.Status(403).Type("json", "utf-8").JSON(fiber.Map{
				"success": false,
				"message": "Incorrect password",
			})
		}

		// Get the media file from the POST request
		file, err := c.FormFile("media")
		if err != nil {
			return c.Status(400).Type("json", "utf-8").JSON(fiber.Map{
				"success": false,
				"message": "Missing file",
			})
		}

		// Check if the extension is in the slice of valid extensions
		extension := file.Filename[len(file.Filename)-3:]
		if !contains(extensions, extension) {
			return c.Status(400).Type("json", "utf-8").JSON(fiber.Map{
				"success": false,
				"message": "Invalid extension",
			})
		}

		// Generate a unique filename that is not in use yet.
		filename := ""
		for {
			filename = randomString(6)
			_, err := os.Stat(fmt.Sprintf("./media/%s/%s.%s", extension, filename, extension))
			if os.IsNotExist(err) {
				break
			}
		}

		// Saves the file and returns the info
		err = c.SaveFile(file, fmt.Sprintf("./media/%s/%s.%s", extension, filename, extension))

		if err != nil {
			return c.Status(500).Type("json", "utf-8").JSON(fiber.Map{
				"success": false,
				"message": "Unknown Error",
			})
		}

		log.Printf("Domain: %s, %s uploaded %s.%s", c.Request().Header.Host(), c.IP(), filename, extension)
		return c.Status(200).Type("json", "utf-8").JSON(fiber.Map{
			"success":   true,
			"filename":  filename,
			"extension": extension,
		})
	})

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).SendString("Page not found!")
	})

	log.Fatal(app.ListenTLS(":"+port, "certificate.pem", "certificate.key"))
}
