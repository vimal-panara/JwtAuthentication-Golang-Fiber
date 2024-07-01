package handlers

import (
	enums "JwtAuthentication/Enums"
	"JwtAuthentication/database"
	"JwtAuthentication/models"
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddOneUser(user models.User) models.Response {

	collection, err := database.GetMongoCollection(enums.UsersCollction)

	if err != nil {
		return models.Response{
			StatusCode: fiber.StatusInternalServerError,
			Msg:        "unable to get mongo colletion",
			Data:       nil,
		}
	}

	filter := bson.D{
		{Key: "firstName", Value: user.FirstName},
		{Key: "lastName", Value: user.LastName},
		{Key: "mobile", Value: user.Mobile},
		{Key: "email", Value: user.Email},
	}

	var existingUser bson.D
	if err := collection.FindOne(context.TODO(), filter).Decode(&existingUser); err != nil {

		user.Id = primitive.NewObjectID()
		result, err := collection.InsertOne(context.TODO(), user)

		if err != nil {
			return models.Response{
				StatusCode: fiber.StatusInternalServerError,
				Msg:        "unable to insert into mongo colletion",
				Data:       nil,
			}
		}

		return models.Response{
			StatusCode: fiber.StatusOK,
			Msg:        "User inserted successfully",
			Data:       result,
		}
	} else {
		return models.Response{
			StatusCode: fiber.StatusOK,
			Msg:        "duplicate user",
			Data:       nil,
		}
	}
}

func FindOneUser(user *models.User) (models.User, error) {

	collection, err := database.GetMongoCollection(enums.UsersCollction)
	if err != nil {
		return models.User{}, err
	}

	// fmt.Println(user)
	filter := bson.D{
		{Key: "$or", Value: []interface{}{
			bson.D{
				{Key: "email", Value: user.Email},
			},
			bson.D{
				{Key: "mobile", Value: user.Mobile},
			},
		}},
	}

	var dbUser models.User

	if err := collection.FindOne(context.TODO(), filter).Decode(&dbUser); err != nil {
		return models.User{}, err
	}

	return dbUser, nil
}

func FindUserById(id string) (models.User, error) {
	// fmt.Println("inside finduserbyid")
	collection, err := database.GetMongoCollection(enums.UsersCollction)
	if err != nil {
		return models.User{}, err
	}

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, err
	}

	// filter := bson.D{
	// 	{Key: "_id", Value: objId},
	// }

	// fmt.Println(objId)
	filter := bson.M{
		"_id": objId,
	}

	var dbUser models.User

	if err := collection.FindOne(context.TODO(), filter).Decode(&dbUser); err != nil {
		return models.User{}, err
	}

	return dbUser, nil
}

func GetAllUsers() ([]models.User, error) {
	collection, err := database.GetMongoCollection(enums.UsersCollction)
	if err != nil {
		return nil, err
	}

	filter := bson.D{
		{Key: "isActive", Value: true},
	}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var users []models.User
	if err := cursor.All(context.TODO(), &users); err != nil {
		return nil, err
	}

	if len(users) <= 0 {
		return nil, errors.New("user data not found")
	}

	return users, nil
}

func UpdateUserTokens(user models.User, id primitive.ObjectID) error {

	collection, err := database.GetMongoCollection(enums.UsersCollction)
	if err != nil {
		return errors.New("can not get the required mongo collection")
	}

	filter := bson.M{
		"$or": []bson.M{
			{"_id": id},
			{
				"email":  user.Email,
				"mobile": user.Mobile,
			},
		},
	}

	update := bson.M{
		"$set": bson.M{
			"token":        user.Token,
			"refreshToken": user.RefreshToken,
			"updatedAt":    user.UpdatedAt,
		},
	}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return err
	}

	if updateResult.MatchedCount <= 0 {
		return errors.New("unable to find specified user")
	}

	return nil
}

func RemoveToken(token, refToken string) error {

	collection, err := database.GetMongoCollection(enums.UsersCollction)

	if err != nil {
		return err
	}

	filter := bson.M{
		"$or": []bson.M{
			{"token": token},
			{"refreshToken": refToken},
		},
	}

	update := bson.M{
		"$set": bson.M{
			"token":        "",
			"refreshToken": "",
			"updatedAt":    primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	if updateResult.MatchedCount <= 0 {
		return errors.New("unable to log out")
	}

	return nil
}

func CheckTokenInDb(token, refToken string) error {
	collection, err := database.GetMongoCollection(enums.UsersCollction)

	if err != nil {
		return err
	}

	filter := bson.M{
		"$or": []bson.M{
			{"token": token},
			{"refreshToken": refToken},
		},
	}

	var user models.User
	if err := collection.FindOne(context.TODO(), filter).Decode(&user); err != nil {
		return err
	}

	return nil
}
