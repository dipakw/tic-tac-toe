package server

import (
	"fmt"
	"math/rand"
)

var namesPool = []string{
	"Alice", "Bob", "Charlie", "David", "Emma",
	"Frank", "Grace", "Henry", "Ivy", "Jack",
	"Kate", "Liam", "Mia", "Noah", "Olivia",
	"Peter", "Quinn", "Rachel", "Sam", "Tina",
	"Uma", "Victor", "Wendy", "Xavier", "Yara",
	"Zach", "Aaron", "Bella", "Caleb", "Diana",
	"Ethan", "Fiona", "George", "Hannah", "Ian",
	"Julia", "Kevin", "Laura", "Marcus", "Nora",
	"Oscar", "Paula", "Ryan", "Sarah", "Thomas",
	"Ursula", "Violet", "William", "Xena", "Yusuf",
	"Zoe", "Adrian", "Bianca", "Connor", "Daisy",
	"Edward", "Freya", "Gabriel", "Hazel", "Isaac",
	"Jasmine", "Kyle", "Luna", "Matthew", "Nina",
	"Oliver", "Penelope", "Quincy", "Ruby", "Steven",
	"Taylor", "Uma", "Vincent", "Willow", "Xander",
	"Yasmine", "Zane", "Amelia", "Brandon", "Chloe",
	"Daniel", "Elena", "Felix", "Gianna", "Hugo",
	"Isla", "Jason", "Kara", "Leo", "Maya",
	"Nathan", "Ophelia", "Patrick", "Riley", "Sophie",
	"Tristan", "Valerie", "Wesley", "Zara", "Alex",
}

func getRandomName() string {
	return namesPool[rand.Intn(len(namesPool))]
}

func randomAvatarUrl() string {
	seed := rand.Intn(100)
	return fmt.Sprintf("https://api.dicebear.com/9.x/avataaars/svg?seed=%d", seed)
}
