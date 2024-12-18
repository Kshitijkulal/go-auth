package main

import (
    "context"
    // "encoding/json"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "golang.org/x/crypto/bcrypt"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "log"
)

var userCollection *mongo.Collection

func main() {
    app := fiber.New()
    ConnectDB()
    userCollection = client.Database("GoAuth").Collection("users")
    app.Use(cors.New(cors.Config{
        AllowOrigins: "*", 
        AllowHeaders: "Content-Type, Authorization",
    }))
    
    app.Use("/protected", AuthMiddleware)
    app.Get("/", func(c *fiber.Ctx) error {
        return c.Status(200).SendString("Hello World")
    })
    app.Post("/api/auth/register", RegisterUser)
    app.Post("/api/auth/login", LoginUser)

    log.Fatal(app.Listen(":8080"))
}

func RegisterUser(c *fiber.Ctx) error {
    log.Println("RegisterUser endpoint hit")
    var user User
    if err := c.BodyParser(&user); err != nil {
        log.Println("Error parsing body:", err)
        return c.Status(fiber.StatusBadRequest).SendString("Failed to parse request body")
    }

    log.Println("Parsed user:", user)
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        log.Println("Error hashing password:", err)
        return c.Status(fiber.StatusInternalServerError).SendString("Failed to hash password")
    }
    user.Password = string(hashedPassword)

    _, err = userCollection.InsertOne(context.TODO(), user)
    if err != nil {
        log.Println("Error inserting user into DB:", err)
        return c.Status(fiber.StatusInternalServerError).SendString("Failed to register user")
    }
    log.Println("User registered successfully")
    return c.SendString("User registered successfully")
}


func LoginUser(c *fiber.Ctx) error {
    var user User
    if err := c.BodyParser(&user); err != nil {
        return c.Status(fiber.StatusBadRequest).SendString("Failed to parse request body")
    }

    var dbUser User
    err := userCollection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&dbUser)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).SendString("User not found")
    }

    err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).SendString("Invalid credentials")
    }

    token, err := GenerateJWT(user.Email)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString("Failed to generate token")
    }

    return c.JSON(fiber.Map{
        "token": token,
    })
}

func AuthMiddleware(c *fiber.Ctx) error {
    tokenStr := c.Get("Authorization")
    if tokenStr == "" {
        return c.Status(fiber.StatusUnauthorized).SendString("Missing token")
    }

    token, err := ValidateJWT(tokenStr)
    if err != nil || !token.Valid {
        return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
    }

    return c.Next()
}
