package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type MongoInstance struct {
	Client *mongo.Client      // These are the reference to the client
	Db     *mongo.Database    // These are the refernce to the database . 
}                             // Go lang doesnt understand the json .
                              
var mg MongoInstance

const dbName = "fiber-hrms"
const mongoURI = " mongodb://127.0.0.1:27017/" + dbName

type Employee struct {
	ID   string    `json:"id,omitempty" bson:"_id,omitempty"`
	Name string    `json:"name"`
	Salary float64  `json:"salary"`
	Age float64      `json:"age"`
}
//  Connect function
func Connect() error {
 client,err :=mongo.NewClient(options.Client().ApplyURI( mongoURI))
 ctx,cancel:=context.WithTimeout(context.Background(),30*time.Second)
 defer cancel ()
 err =client.Connect(ctx)
 db :=client.Database(dbName)


 if err!=nil {
	return err
 } 
// We have defined the mongoInstance here .
 mg = MongoInstance{
	Client: client,
	Db: db,

 }

 return nil 




}



func main () {
	if err := Connect();err !=nil {
		log.Fatal(err)
	}
	app :=fiber.New()

	// These are our route handlers 
	app.Get("/employee",func(c*fiber.Ctx)error {
    
		 query := bson.D {{}}
		cursor,err :=mg.Db.Collection("employees").Find (c.Context(),query)
        if err!= nil {
			return c.Status(500).SendString(err.Error())

		} 
		 var employess []Employee =make ([]Employee, 0)

        if err := cursor.All(c.Context(),&employess);
		err! = nil {
         return c.Status(500).SendString(err.Error())
		
		}
		return c.JSON(employess)




	})
    app.Post("/employee", func (c *fiber.Ctx)errors{
		collection := mg.Db.Collection("employees")


		employee := new(Employee)
		if err := c.BodyParser(employee);  // IT parses the body and pass it to the variable .
	    err !=nil {
			return c.Status(400).SendString(err.Error())
		}

		employee.ID = ""

		insertionResult,err := collection.InsertOne(c.Context(),employee)
	    if err := c.BodyParser(employee);err !=nill {
			retur c. Status (400).SendString(err.Error ())
		}

		employee.ID = ""
		insertionResult, err := collection.InsterOne (c.Context (),employee)
        if err !=nil {
			return c.Status(500).SendString (err.Error())
		}	
		filter :=bson.D{{Key: "_id",Value: insertionResult.InsertedID}}
		createdRecord := collection.FindOne(c.Context(),filter)


		createdEmployee := &Employee {}
		createdRecord.Decode(createdEmployee)

		return c.Status(201),json (createdEmployee)
	})



	app.Put("/employee/:id")
	app.Delete("/employee/:id")



}