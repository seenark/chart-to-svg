package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CoinKLine struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Symbol      string             `json:"symbol" bson:"symbol"`
	ClosePrices []float64          `json:"closePrices" bson:"closePrices"`
	Svg         string             `json:"svg" bson:"svg"`
}

type ICoinKLineRepository interface {
	GetMultiple([]string) ([]CoinKLine, error)
	GetBySymbol(string) (*CoinKLine, error)
	Create(CoinKLine) error
	Update(string, CoinKLine) error
	Delete(string) error
}

type KLineDb struct {
	cl  *mongo.Collection
	ctx context.Context
}

func NewKLineRepository(collection *mongo.Collection, ctx context.Context) ICoinKLineRepository {
	return &KLineDb{
		cl:  collection,
		ctx: ctx,
	}
}

func (k *KLineDb) GetMultiple(symbols []string) ([]CoinKLine, error) {
	fmt.Printf("symbols: %v\n", symbols)
	all := []CoinKLine{}
	filters := bson.M{}
	newSymbols := []string{}
	for _, v := range symbols {
		if v != "" {
			newSymbols = append(newSymbols, v)
		}
	}
	if len(newSymbols) > 0 {
		symbolArr := []string{}
		symbolArr = append(symbolArr, symbols...)
		filters["symbol"] = bson.M{"$in": symbolArr}
	}
	cur, err := k.cl.Find(k.ctx, filters)
	if err != nil {
		return nil, err
	}
	for cur.Next(k.ctx) {
		kline := CoinKLine{}
		err = cur.Decode(&kline)
		if err != nil {
			continue
		}
		all = append(all, kline)
	}
	// fmt.Println(all)
	return all, nil
}

func (k *KLineDb) GetBySymbol(symbol string) (*CoinKLine, error) {
	filter := bson.D{primitive.E{Key: "symbol", Value: symbol}}
	kline := CoinKLine{}
	err := k.cl.FindOne(k.ctx, filter).Decode(&kline)
	if err != nil {
		return nil, err
	}
	return &kline, nil

}
func (k *KLineDb) Create(kline CoinKLine) error {
	kline.Id = primitive.NewObjectID()

	_, err := k.cl.InsertOne(k.ctx, kline)
	if err != nil {
		return err
	}
	return nil
}
func (k *KLineDb) Update(symbol string, kline CoinKLine) error {
	filter := bson.D{primitive.E{Key: "symbol", Value: symbol}}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{
			Key:   "closePrices",
			Value: kline.ClosePrices,
		},
	}}}
	_, err := k.cl.UpdateOne(k.ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
func (k *KLineDb) Delete(symbol string) error {
	filter := bson.D{primitive.E{Key: "symbol", Value: symbol}}
	_, err := k.cl.DeleteOne(k.ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
