package servers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"restaurants/src/models"
	"restaurants/src/proto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/emptypb"
)

type IRestaurantServer interface {
	GetRestaurants(context.Context, *emptypb.Empty) (*proto.GetRestaurantsResponse, error)
	AddRestaurant(context.Context, *proto.AddRestaurantRequest) (*proto.AddRestaurantResponse, error)
	DeleteRestaurantRequest(context.Context, *proto.DeleteRestaurantRequest) (*proto.DeleteRestaurantResponse, error)
	UpdateRestaurant(context.Context, *proto.UpdateRestaurantRequest) (*proto.UpdateRestaurantResponse, error)
}

type RestaurantServer struct {
	proto.UnimplementedRestaurantServiceServer
	coll *mongo.Collection
}

func NewRestaurantServer(coll *mongo.Collection) *RestaurantServer {
	return &RestaurantServer{coll: coll}
}

func (srv *RestaurantServer) GetRestaurants(
	context.Context,
	*emptypb.Empty,
) (*proto.GetRestaurantsResponse, error) {
	var restaurants []*models.Restaurant

	// Filter query
	// TODO: Filter based on user's location
	condition := bson.M{
		"deleted_at": bson.M{
			"$exists": false,
		},
	}

	cur, err := srv.coll.Find(context.Background(), condition)
	if err != nil {
		log.Fatal(err)
	}

	// Loop cursor and retrieve data
	if err := cur.All(context.Background(), &restaurants); err != nil {
		log.Fatal(err)
	}

	// Map the results to expected response type
	var restaurants_grpc []*proto.Restaurant

	for _, v := range restaurants {
		restaurant_mapped := &proto.Restaurant{
			Id:    v.Id.Hex(),
			Name:  v.Name,
			Image: v.Image,
		}

		restaurants_grpc = append(restaurants_grpc, restaurant_mapped)
	}

	return &proto.GetRestaurantsResponse{Data: restaurants_grpc}, nil
}

func (srv *RestaurantServer) AddRestaurant(
	ctx context.Context,
	args *proto.AddRestaurantRequest,
) (*proto.AddRestaurantResponse, error) {

	// Map arg to mongo collection
	now := time.Now()

	openHour, err := time.Parse(time.RFC3339, args.OpenHour)
	closeHour, err := time.Parse(time.RFC3339, args.CloseHour)

	// Uplaoad img to AWS S3 if provided
	var image = "https://example.com/pizza.jpg"

	if args.Image != nil {
		image = "..."
	}

	restaurant := models.Restaurant{
		Id:        primitive.NewObjectID(),
		Name:      args.Name,
		Image:     image,
		OpenHour:  openHour,
		CloseHour: closeHour,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Insert the document into the restaurant collection
	res, err := srv.coll.InsertOne(ctx, restaurant)
	if err != nil {
		log.Fatal(err)
		return nil, fmt.Errorf("Error creating the new restaurant")
	}

	// Map result to grpc response type
	id := res.InsertedID.(primitive.ObjectID).Hex()
	return &proto.AddRestaurantResponse{RestaurantId: id}, nil
}

func (srv *RestaurantServer) DeleteRestaurant(
	ctx context.Context,
	args *proto.DeleteRestaurantRequest,
) (*proto.DeleteRestaurantResponse, error) {

	// Get active orders from a restaurant
	// TODO: throw error if the restaurant has active orders

	// Find the selected restaurant
	objectId, err := primitive.ObjectIDFromHex(args.RestaurantId)

	if err != nil {
		return nil, fmt.Errorf("Id provided not valid")
	}

	filter := bson.M{
		"_id": objectId,
		"deleted_at": bson.M{
			"$exists": false,
		},
	}

	// Update field deletedAt to considerate it disable.
	now := time.Now()

	update := bson.M{
		"$set": bson.M{
			"updated_at": now,
			"deleted_at": now,
		},
	}

	// Finally update in the mongodb collection
	res := srv.coll.FindOneAndUpdate(context.TODO(), filter, update)

	// Validate common error cases
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("Restaurant not found")
		}

		return nil, fmt.Errorf("Error deleting the restaurant")
	}

	return &proto.DeleteRestaurantResponse{RestaurantId: objectId.Hex()}, nil
}

func (srv *RestaurantServer) UpdateRestaurant(
	ctx context.Context,
	args *proto.UpdateRestaurantRequest,
) (*proto.UpdateRestaurantResponse, error) {

	// Find the selected restaurant
	objectId, err := primitive.ObjectIDFromHex(args.RestaurantId)

	if err != nil {
		return nil, fmt.Errorf("The Id provided is not valid")
	}

	filter := bson.M{
		"_id": objectId,
		"deleted_at": bson.M{
			"$exists": false,
		},
	}

	// TODO: UPLOAD IMG to AWS IF PROVIDED

	// Update fields based on arguments provided
	openHour, err := time.Parse(time.RFC3339, args.OpenHour)
	closeHour, err := time.Parse(time.RFC3339, args.CloseHour)

	update := bson.M{
		"$set": bson.M{
			"name":       args.Name,
			"close_hour": closeHour,
			"open_hour":  openHour,
			"updated_at": time.Now(),
		},
	}

	res := srv.coll.FindOneAndUpdate(context.Background(), filter, update)

	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("Restaurant not found")
		}

		log.Printf(res.Err().Error())
		return nil, fmt.Errorf("Error updating the restaurant")
	}

	return &proto.UpdateRestaurantResponse{RestaurantId: args.RestaurantId}, nil
}
