package services

import (
	"literary-lions/models"
	"literary-lions/utils"
	"log"
)

func GetAllCategories() (output []models.Category) {
	var c models.Category
	row, err := utils.Getdata().Query(`SELECT id, title FROM categories`)
	if err != nil {
		log.Printf("Error querying categories: %s", err)
	}
	for row.Next() {
		err := row.Scan(&c.ID, &c.Name)
		if err != nil {
			log.Printf("Error querying categories: %s", err)
		} else {
			output = append(output, c)
		}
	}
	return output
}
