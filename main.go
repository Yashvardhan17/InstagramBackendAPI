package main

import (
	"awesomeProject1/db"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"
)
//Struct defined for user collection
type UserDetails struct {
	Id  int `json: "Id"`
	Name  string    `json: "Name"`
	Email string   `json: "Email"`
	Password string `json: "Password"`
}
var UserArr []UserDetails
//var myClient = &http.Client{Timeout: 10*time.Second}
//handle functions
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Instagram API")
	fmt.Println("Endpoint Hit: homePage")
}

func AddUserFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
//Method POST and get the reponse from http response
		reqBody, _ := ioutil.ReadAll(r.Body)
		fmt.Fprintf(w, "%+v", string(reqBody))
		var users UserDetails
		json.Unmarshal(reqBody, &users)
		// update our global userdetail array to include
		// our new user
		UserArr = append(UserArr, users)

		json.NewEncoder(w).Encode(users)

		fmt.Println("Endpoint Hit: users")
		/*data := UserDetails{}
		json.Unmarshal([]byte(reqBody), data)*/


		dbName := "Instagram_DB" // move to env
		client, err := db.CreateDatabaseConnection(dbName)
		if err != nil {
			fmt.Println("Failed to connect to DB")
			panic(err)
		}
		defer client.Disconnect(context.TODO())
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		col := client.Database(dbName).Collection("users_collection")


		fmt.Println(HashPassword(users.Password))
		result, insertErr := col.InsertOne(ctx,users);

		if insertErr != nil {
			fmt.Println("InsertONE Error:", insertErr)
			os.Exit(1)
		} else {
			fmt.Println("InsertOne() result type: ", reflect.TypeOf(result))
			fmt.Println("InsertOne() api result type: ", result)

			newID := result.InsertedID
			fmt.Println("InsertedOne(), newID", newID)
			fmt.Println("InsertedOne(), newID type:", reflect.TypeOf(newID))

		}

	}
}
func GetUserFunc(w http.ResponseWriter, r *http.Request) {

		keys, ok := r.URL.Query()["id"]

		if !ok || len(keys[0]) < 1 {
			log.Println("Url Param 'key' is missing")
			return
		}

		// Query()["key"] will return an array of items,
		// we only want the single item.
		key := keys[0]

		log.Println("Url Param 'key' is: " + string(key))
		//var userDB UserDetails
	dbName := "Instagram_DB"
	client, err := db.CreateDatabaseConnection(dbName)
	if err != nil {
		fmt.Println("Failed to connect to DB")
		panic(err)
	}
	defer client.Disconnect(context.TODO())
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	col := client.Database(dbName).Collection("users_collection")
	id_key, err := strconv.Atoi(key)
	filterCursor, err := col.Find(ctx, bson.M{"id":id_key })
	fmt.Println(filterCursor)
	if err != nil {
		log.Fatal(err)
	}
	var userFiltered []bson.M
	if err = filterCursor.All(ctx, &userFiltered); err != nil {
		log.Fatal(err)
	}
	fmt.Println(userFiltered)
	json.NewEncoder(w).Encode(userFiltered)
}

func handlerequests(){
	http.HandleFunc("/", homePage)
	http.HandleFunc("/users", AddUserFunc)
	http.HandleFunc("/users/",GetUserFunc)
	http.HandleFunc("/posts",creategitPosts)
	http.HandleFunc("/posts/",PostById)
	http.HandleFunc("/posts/users/",returnAllPosts)

	log.Fatal(http.ListenAndServe(":30000", nil))



}
//password encryption methods
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func main() {

handlerequests()

}
