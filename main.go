package main

import (
	"MemeCraft/internal/adapter/http"
	"MemeCraft/internal/adapter/storage"
	"MemeCraft/internal/preset"
	"MemeCraft/internal/service/meme"
	"fmt"
	"github.com/alecthomas/kingpin/v2"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
)

var (
	port = kingpin.Flag("port", "http port").Short('p').Default("3000").String()
)

func main() {
	kingpin.Parse()
	presetRegistry := preset.NewRegistry()
	if err := presetRegistry.LoadFromDir("./presets"); err != nil {
		log.Fatal(err)
	}

	app := newFiberApp()
	_ = storage.NewCatboxMoeStorage()                     // catbox.moe
	zeroXzeroSTStorage := storage.NewZeroXZeroSTStorage() // 0x0.st
	memeGenerator := meme.NewGenerator(presetRegistry, zeroXzeroSTStorage)
	handler := http.NewHandler(memeGenerator, zeroXzeroSTStorage)

	app.Get("/presets", handler.GetAllPreset)
	app.Get("/presets/:preset_id", handler.GetPresetById)
	app.Post("/presets/:preset_id/memes", handler.GenerateMeme)
	app.Post("/upload", handler.UploadFile)

	app.Static("/", "./public")

	log.Fatal(app.Listen(fmt.Sprintf(":%s", *port)))
}

func newFiberApp() *fiber.App {
	app := fiber.New(fiber.Config{
		StrictRouting:     true,
		CaseSensitive:     true,
		AppName:           "MemeCraft",
		JSONDecoder:       sonic.Unmarshal,
		JSONEncoder:       sonic.Marshal,
		EnablePrintRoutes: true,
	})
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Powered-By", "github.com/n0paleon/MemeCraft")
		return c.Next()
	})
	app.Use(cors.New(cors.Config{
		Next: func(c *fiber.Ctx) bool {
			return true
		},
	}))

	return app
}
