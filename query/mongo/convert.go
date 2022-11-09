package mongo

import (
	"gitlab.com/dipper-iot/shared/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"time"
)

func StringTimestampToTime(str string) string {
	timeStamp, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		logger.Error("convert timeStamp: ", err)
		return str
	}
	timeData := time.Unix(timeStamp/1000, 0)
	return timeData.String()
}

func GetObjectIdToString(mapData map[string]interface{}) string {
	item, exits := mapData["_id"]
	if !exits {
		return ""
	}
	return item.(primitive.ObjectID).Hex()
}

func GetIdToString(mapData map[string]interface{}) string {
	item, exits := mapData["_id"]
	if !exits {
		return ""
	}
	return item.(string)
}
