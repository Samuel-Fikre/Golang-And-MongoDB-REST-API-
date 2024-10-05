package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"mongodb-golang/models"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	Client *mongo.Client
}

func NewUserController(s *mongo.Client) *UserController {
	return &UserController{Client: s}
}

// You are getting the user ID from the URL parameters using httprouter.Params.
// You are then checking if the user ID is valid using bson.IsObjectIdHex.

func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	// Here, you're checking if the ID is a valid BSON ObjectId. If it's not, you return a 404 status indicating that the resource was not found.

	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// bson.ObjectIdHex(id): Converts the hexadecimal string id into an ObjectId
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		// Handle the error (e.g., return an HTTP error response)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//  Initializes a new User struct. If no fields are provided, they will have their default values (zero for integers, empty strings, etc.)

	u := models.User{}

	// uc.Session: Refers to the MongoDB session associated with the uc object. This uc is likely a struct that holds the MongoDB session.

	// DB("mongo-golang"): Accesses the MongoDB database named mongo-golang.

	// C("users"): Refers to the users collection within the mongo-golang database

	// This line is often used to specify which collection of a database you want to interact with before performing operations like querying, inserting, updating, or deleting data

	//if err := uc.Session.DB("mongo-golang").C("users").FindId(oid).One(&user): This line checks if there’s an error when trying to find a user with the specified ObjectId

	// oid: The ObjectId of the document you are trying to find.
	//One(&user): Retrieves the result into the user variable.

	var user models.User
	if err := uc.Client.Database("mongo-golang").Collection("users").FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&user); err != nil {
		w.WriteHeader(404)
		return
	}

	//  is used to convert a Go struct (u) into a JSON-encoded byte slice

	// u: This is a Go struct or any data structure that you want to serialize (convert) to JSON

	// json.Marshal(u): Converts the struct u into JSON format. This returns a byte slice ([]byte) representing the JSON-encoded data.

	//uj: The variable where the JSON-encoded data is stored. It's a byte slice ([]byte).

	uj, err := json.Marshal(u)

	//  If marshalling is successful, the JSON-encoded data is stored in uj, and if an error occurs, it's stored in err.
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Use fmt.Fprintf to write formatted output to the ResponseWriter (w)
	fmt.Fprintf(w, "%s\n", uj) // uj is the JSON-encoded data

	// This uses fmt.Fprintf to write the JSON-encoded byte slice (uj) to the http.ResponseWriter (w). The uj variable represents the marshalled JSON data (e.g., {"name":"John Doe","email":"john.doe@example.com"}).

}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := models.User{}

	// In this line, the incoming JSON data (from r.Body) is being decoded into the struct u (e.g., User)

	// Decoding refers to the process of converting data from a specific format (such as JSON, XML, or binary data) back into a usable form in your program (like a struct, class, or variable).

	// Encoding: Taking a Go struct and converting it into a format (like JSON) for output (e.g., sending over HTTP).

	// In Go, the & symbol is used to get the memory address of a variable, creating a pointer to that variable. A pointer stores the memory address where the actual data is located rather than a copy of the data itself.

	// In the context of decoding JSON or working with struct data, you use & to pass the pointer to a variable so that changes can be made directly to the original variable (not a copy).

	json.NewDecoder(r.Body).Decode(&u)

	// bson.NewObjectId():

	// This function call generates a new BSON ObjectId. In MongoDB, each document has a unique _id field, and this ObjectId serves as that unique identifier. The NewObjectId() function creates a new ObjectId that is unique across all documents in the database

	u.ID = primitive.NewObjectID()

	// C("users"):

	//This method specifies the collection within the database where you want to perform operations. Here, users is the collection where user documents are stored

	collection := uc.Client.Database("mongo-golang").Collection("users")
	_, err := collection.InsertOne(context.TODO(), u)
	if err != nil {
		// Handle error
	}

	// ... existing code ...
	uj, err := json.Marshal(u)
	if err != nil {
		// handle error
		return
	}

	fmt.Println(string(uj))
	// Use uj here, for example:
	// fmt.Println(string(uj))
	// ... existing code ...

	// This is for sending it to front end-

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	userJSON, err := json.Marshal(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s\n", userJSON)

}

func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	// The snippet you've provided checks if a given string id is a valid MongoDB ObjectID using the primitive.IsValidObjectID function from the go.mongodb.org/mongo-driver/bson/primitive package. This is a good practice to ensure that the ID format is correct before attempting to perform operations on it.

	// Here are some examples of valid MongoDB ObjectIDs:

	// 	507f1f77bcf86cd799439011

	// Invalid 507f1f77bcf86cd79943901z (contains an invalid character z)

	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// The line of code oid := bson.ObjectIDHex(id) is used to convert a string representation of a MongoDB ObjectID into a BSON ObjectID type in Go

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		// Handle the error (e.g., return an error response)
		return
	}

	result, err := uc.Client.Database("mongo-golang").Collection("users").DeleteOne(r.Context(), bson.M{"_id": oid})
	if err != nil || result.DeletedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Deleted user", oid, "\n")

}
