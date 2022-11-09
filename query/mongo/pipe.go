package mongo

import (
	"context"
	"gitlab.com/dipper-iot/shared/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"sync"
)

type TotalItem struct {
	Count int64 `json:"count" bson:"count,omitempty"`
}

type ResultList struct {
	Data  primitive.A `json:"data" bson:"data,omitempty"`
	Total []TotalItem `json:"total" bson:"total,omitempty"`
}

func LimitBuildQuery(limit int64) bson.M {

	return bson.M{
		"$limit": limit,
	}
}

func SkipBuildQuery(skip int64) bson.M {
	return bson.M{
		"$skip": skip,
	}
}

const (
	SORT_ASCENDING  int16 = 1
	SORT_DESCENDING int16 = -1
)

func SortBuildQuery(field bson.M) bson.M {
	return bson.M{
		"$sort": field,
	}
}

type CallBackData = func(data primitive.D) error

func ExecQuery(ctx context.Context, collection *mongo.Collection, pipe bson.A, pipeData bson.A, callBackData CallBackData) (int64, error) {

	errs := make([]error, 0)
	var resultCount TotalItem
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		pipeCount := append(pipe,
			bson.M{
				"$count": "count",
			},
		)

		curC, err := collection.Aggregate(ctx, pipeCount)

		if err != nil {
			logger.Error("Aggregate count: ", err)
			errs = append(errs, err)
			return
		}

		curC.Next(ctx)

		err = curC.Decode(&resultCount)
		if err != nil {
			logger.Error("Aggregate data: ", err)
			errs = append(errs, err)
			return
		}
		curC.Close(ctx)

	}()

	go func() {
		defer wg.Done()
		pipeData = append(pipe, pipeData...)

		cur, err := collection.Aggregate(ctx, pipeData)

		if err != nil {
			logger.Error("Aggregate data: ", err)
			errs = append(errs, err)
			return
		}

		var result primitive.D
		for cur.Next(ctx) {
			cur.Decode(&result)
			callBackData(result)
		}
		cur.Close(ctx)

	}()

	wg.Wait()

	if len(errs) > 0 {
		for _, err := range errs {
			if err == io.EOF && len(errs) == 1 {
				return 0, nil
			}
			return 0, err
		}
	}

	return resultCount.Count, nil
}
