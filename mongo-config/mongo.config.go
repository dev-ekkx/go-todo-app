package mongoController

import (
	"context"
	"fmt"
	"github.com/dev-ekks/go-todo-app/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

const dbName = "todo"
const colName = "todoList"

var collection *mongo.Collection

func ConnectDb() {
	// Connect to MongoDB
	mongoUser := os.Getenv("MONGO_USER")
	mongoPassword := os.Getenv("MONGO_PASSWORD")
	mongoAppName := os.Getenv("MONGO_APP_NAME")

	mongoUri := fmt.Sprintf("mongodb+srv://%s:%s@todo.vl9hi.mongodb.net/?retryWrites=true&w=majority&appName=%s", mongoUser, mongoPassword, mongoAppName)

	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(mongoUri))
	if err != nil {
		log.Fatal(err)
	}
	//defer func() {
	//	if err := client.Disconnect(context.TODO()); err != nil {
	//		panic(err)
	//	}
	//}()
	fmt.Println("Connected to mongo DB")
	collection = client.Database(dbName).Collection(colName)

	fmt.Println("Collection instance ready...", collection.Name())
}

func GetAllTodos() ([]*model.Todo, error) {
	//Retrieve documents
	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	// Unpacks the cursor into a slice
	var todos []*model.Todo
	for cur.Next(context.Background()) {
		var todo model.TodoStruct
		err := cur.Decode(&todo)
		if err != nil {
			log.Fatal(err)
		}
		newTodoItem := model.Todo{
			Done: todo.Done,
			Text: todo.Text,
			ID:   todo.ID.Hex(),
		}
		todos = append(todos, &newTodoItem)
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(cur, context.Background())
	return todos, nil
}

func CreateTodo(input model.NewTodo) (*model.Todo, error) {
	newTodo := &model.TodoStruct{
		Text: input.Text,
		Done: false,
	}

	result, err := collection.InsertOne(context.Background(), newTodo)
	if err != nil {
		log.Fatal(err)
	}

	newTodo.ID = result.InsertedID.(primitive.ObjectID)
	createdTodo := model.Todo{
		Done: newTodo.Done,
		Text: newTodo.Text,
		ID:   newTodo.ID.Hex(),
	}
	return &createdTodo, nil
}
