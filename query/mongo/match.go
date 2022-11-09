package mongo

import (
	"fmt"
	"gitlab.com/dipper-iot/shared/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func OnTimeBuildQuery(field string, start *time.Time, end *time.Time) *bson.M {

	var section bson.M

	if end != nil && start != nil {
		section = bson.M{
			"$match": bson.M{
				field: bson.M{
					"$gte": *start,
					"$lt":  *end,
				},
			},
		}
		return &section
	}

	if end == nil {
		section = bson.M{
			"$match": bson.M{
				field: bson.M{
					"$gte": start,
				},
			},
		}
	}

	if start == nil {
		section = bson.M{
			"$match": bson.M{
				field: bson.M{
					"$lt": end,
				},
			},
		}
	}

	return &section
}

func OnTimeStringBuildQuery(field string, start string, end string, layout string) *bson.M {

	if layout == "" {
		layout = "2006-01-02"
	}

	if start != "" && end != "" {
		startTime, _ := time.Parse(layout, start)
		endTime, _ := time.Parse(layout, end)
		return OnTimeBuildQuery(field, &startTime, &endTime)
	}

	if end == "" {
		startTime, _ := time.Parse(layout, start)
		return OnTimeBuildQuery(field, &startTime, nil)
	}

	endTime, _ := time.Parse(layout, end)
	return OnTimeBuildQuery(field, nil, &endTime)
}

func MatchTextBuildQuery(field string, match string) bson.M {
	return bson.M{
		"$match": bson.M{
			field: bson.M{"$regex": ".*" + match + ".*"},
		},
	}
}

func QueryText(field string, subMatch string) bson.M {
	if subMatch == "" {
		return nil
	}

	return bson.M{
		field: bson.M{"$regex": ".*" + subMatch + ".*", "$options": "i"},
	}
}

func QueryMatch(field string, match string) bson.M {
	if match == "" {
		return nil
	}

	return bson.M{
		field: match,
	}
}

func QueryTextBegin(field string, subMatch string) bson.M {
	if subMatch == "" {
		return nil
	}

	return bson.M{
		field: bson.M{"$regex": primitive.Regex{Pattern: "^" + subMatch + ".*", Options: "s"}},
	}
}

func QueryIsText(field string, textMatch string) bson.M {
	if textMatch == "" {
		return nil
	}

	return bson.M{
		field: textMatch,
	}
}

func QueryIsNumber(field string, nMatch int32, isNot int32) bson.M {
	if nMatch == isNot {
		return nil
	}

	return bson.M{
		field: nMatch,
	}
}

func QueryIsNotNumber(field string, nMatch int32, isNot int32) bson.M {
	if nMatch == isNot {
		return nil
	}

	return bson.M{
		field: bson.M{
			"$ne": nMatch,
		},
	}
}

func PipeMatch(pipe bson.A, query bson.M) bson.A {
	if query != nil {
		pipe = append(pipe, bson.M{
			"$match": query,
		})
	}

	return pipe
}

func QueryTextInArray(field string, listMatch []string) bson.M {

	if listMatch == nil || len(listMatch) == 0 {
		return nil
	}

	var list []bson.M

	for _, match := range listMatch {
		list = append(list, QueryText(field, match))
	}

	return MatchOr(list...)
}

func QueryRegexInArray(field string, listMatch []string) bson.M {

	if listMatch == nil || len(listMatch) == 0 {
		return nil
	}

	if len(listMatch) == 1 {
		return bson.M{
			field: bson.M{
				"$regex": primitive.Regex{Pattern: fmt.Sprintf("%s.*", listMatch[0]), Options: "i"},
			},
		}
	}

	var list []primitive.Regex

	for _, match := range listMatch {
		list = append(list, primitive.Regex{Pattern: fmt.Sprintf("%s.*", match), Options: "i"})
	}

	return bson.M{
		field: bson.M{
			"$in": list,
		},
	}
}

func QueryRegexFullInArray(field string, listMatch []string) bson.M {

	if listMatch == nil || len(listMatch) == 0 {
		return nil
	}

	if len(listMatch) == 1 {
		return bson.M{
			field: bson.M{
				"$regex": primitive.Regex{Pattern: fmt.Sprintf("%s", listMatch[0]), Options: "i"},
			},
		}
	}

	var list []primitive.Regex

	for _, match := range listMatch {
		list = append(list, primitive.Regex{Pattern: fmt.Sprintf("%s", match), Options: "i"})
	}

	return bson.M{
		field: bson.M{
			"$in": list,
		},
	}
}

func QueryTextBeginInArray(field string, listMatch []string) bson.M {

	if listMatch == nil || len(listMatch) == 0 {
		return nil
	}

	var list []bson.M

	for _, match := range listMatch {
		list = append(list, QueryTextBegin(field, match))
	}

	return MatchOr(list...)
}

func QueryFullInArray(field string, listMatch []string) bson.M {

	if listMatch == nil || len(listMatch) == 0 {
		return nil
	}

	list := []bson.M{}

	for _, match := range listMatch {
		list = append(list, bson.M{
			field: match,
		})
	}

	return MatchOr(list...)
}

func QueryItemInArray(field string, listMatch []string) bson.M {
	if len(listMatch) == 0 {
		return nil
	}

	return bson.M{
		field: bson.M{
			"$in": listMatch,
		},
	}
}

func QueryItemInArrayException(field string, listMatch []string, textException string) bson.M {
	if len(listMatch) == 0 {
		return nil
	}

	if util.StringInSlice(textException, listMatch) {
		return nil
	}

	return bson.M{
		field: bson.M{
			"$in": listMatch,
		},
	}
}

func MatchOr(listMath ...bson.M) bson.M {

	list := make([]bson.M, 0)

	for _, item := range listMath {
		if item != nil {
			list = append(list, item)
		}
	}

	if len(list) == 0 {
		return nil
	}

	if len(list) == 1 {
		return bson.M{
			"$match": list[0],
		}
	}

	return bson.M{
		"$match": bson.M{
			"$or": list,
		},
	}
}

func MatchAnd(listMath ...bson.M) bson.M {

	list := make([]bson.M, 0)

	for _, item := range listMath {
		if item != nil {
			list = append(list, item)
		}
	}

	if len(list) == 0 {
		return nil
	}

	if len(list) == 1 {
		return bson.M{
			"$match": list[0],
		}
	}

	return bson.M{
		"$match": bson.M{
			"$and": list,
		},
	}
}

func MatchItemArrayQuery(field string, listMatch []string) bson.M {
	return bson.M{
		"$match": QueryItemInArray(field, listMatch),
	}
}

func ToPipe(field string, listMatch []string) bson.M {
	return bson.M{
		"$match": QueryItemInArray(field, listMatch),
	}
}

func TimeBuild(field string, start *time.Time, end *time.Time) bson.M {

	if end != nil && start != nil {
		return bson.M{
			field: bson.M{
				"$gte": start,
				"$lt":  end,
			},
		}
	}

	if end == nil {
		return bson.M{
			field: bson.M{
				"$gte": start,
			},
		}
	}

	if start == nil {
		return bson.M{
			field: bson.M{
				"$lt": end,
			},
		}
	}

	return nil
}

func TimeBuildId(start *time.Time, end *time.Time) bson.M {

	if end != nil && start != nil {
		return bson.M{
			"_id": bson.M{
				"$gte": primitive.NewObjectIDFromTimestamp(*start),
				"$lt":  primitive.NewObjectIDFromTimestamp(*end),
			},
		}
	}

	if end == nil {
		return bson.M{
			"_id": bson.M{
				"$gte": primitive.NewObjectIDFromTimestamp(*start),
			},
		}
	}

	if start == nil {
		return bson.M{
			"_id": bson.M{
				"$lt": primitive.NewObjectIDFromTimestamp(*end),
			},
		}
	}

	return nil
}

func TimeStringBuild(field string, start string, end string, layout string) bson.M {
	if start == "" && end == "" {
		return nil
	}

	if layout == "" {
		layout = "2006-01-02"
	}

	if start != "" && end != "" {
		startTime, _ := time.Parse(layout, start)
		endTime, _ := time.Parse(layout, end)
		return TimeBuild(field, &startTime, &endTime)
	}

	if end == "" {
		startTime, _ := time.Parse(layout, start)
		return TimeBuild(field, &startTime, nil)
	}

	endTime, _ := time.Parse(layout, end)
	return TimeBuild(field, nil, &endTime)
}

func TimeStringBuildId(start string, end string, layout string) bson.M {
	if start == "" && end == "" {
		return nil
	}

	if layout == "" {
		layout = "2006-01-02"
	}

	if start != "" && end != "" {
		startTime, _ := time.Parse(layout, start)
		endTime, _ := time.Parse(layout, end)
		return TimeBuildId(&startTime, &endTime)
	}

	if end == "" {
		startTime, _ := time.Parse(layout, start)
		return TimeBuildId(&startTime, nil)
	}

	endTime, _ := time.Parse(layout, end)
	return TimeBuildId(nil, &endTime)
}
