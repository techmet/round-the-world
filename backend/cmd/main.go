package main

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/techmet/round-the-world/data"
)

func main() {
	app := fiber.New()

	log.Println("Populating the Initial data..")

	data.PopulateNeighboursForEachCity()

	log.Println("Succesfull populated the Initial data")

	app.Use(cors.New(cors.Config{}))

	app.Get("/cities/search", func(c *fiber.Ctx) error {
		cities := data.GetCitiesForSearch()
		citiesStr, _ := json.Marshal(cities)
		return c.Send(citiesStr)
	})

	app.Get("/cities/trip/:cityId", func(c *fiber.Ctx) error {
		cityId := c.Params("cityId")
		cities := data.GetRoundTrip(cityId)
		citiesStr, _ := json.Marshal(cities)
		return c.Send(citiesStr)
	})

	log.Fatal(app.Listen(":5001"))
}
