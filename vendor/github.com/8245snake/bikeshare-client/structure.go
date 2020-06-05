package bikeshareapi

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/8245snake/bikeshare_api/src/lib/static"
)

//JsonTimeLayout 時刻フォーマット
const JsonTimeLayout = "2006/01/02 15:04"

/////////////////////////////////////////////////////////////////////////////////////////////////////////
//  JSONマーシャリング構造体
/////////////////////////////////////////////////////////////////////////////////////////////////////////

//GetSpotInfoByJCount GetSpotInfo構造体を返す
func GetSpotInfoByJCount(body static.JCountsBody) SpotInfo {
	var spotinfo SpotInfo
	spotinfo.Area = body.Area
	spotinfo.Spot = body.Spot
	spotinfo.Name = body.Name
	spotinfo.Description = body.Description
	if lat, err := strconv.ParseFloat(body.Lat, 64); err == nil {
		spotinfo.Lat = lat
	}
	if lon, err := strconv.ParseFloat(body.Lon, 64); err == nil {
		spotinfo.Lon = lon
	}
	for _, item := range body.Counts {
		if s, err := NewBikeCount(item.Datetime, item.Count); err == nil {
			spotinfo.Counts = append(spotinfo.Counts, s)
		}
	}
	return spotinfo
}

//GetSpotInfoListByPlaces GetSpotInfo構造体を返す
func GetSpotInfoListByPlaces(body static.JPlacesBody) []SpotInfo {
	var spotinfoList []SpotInfo
	for _, item := range body.Items {
		var spotinfo SpotInfo
		spotinfo.Area = item.Area
		spotinfo.Spot = item.Spot
		spotinfo.Name = item.Name
		spotinfo.Description = item.Description
		if lat, err := strconv.ParseFloat(item.Lat, 64); err == nil {
			spotinfo.Lat = lat
		}
		if lon, err := strconv.ParseFloat(item.Lon, 64); err == nil {
			spotinfo.Lon = lon
		}
		if count, err := NewBikeCount(item.Recent.Datetime, item.Recent.Count); err == nil {
			spotinfo.Counts = append(spotinfo.Counts, count)
		}
		spotinfoList = append(spotinfoList, spotinfo)
	}
	return spotinfoList
}

//GetSpotInfoListByDistance GetSpotInfo構造体を返す
func GetSpotInfoListByDistance(body static.JDistancesBody) []SpotInfo {
	var spotinfoList []SpotInfo
	for _, item := range body.Items {
		var spotinfo SpotInfo
		spotinfo.Area = item.Area
		spotinfo.Spot = item.Spot
		spotinfo.Name = item.Name
		spotinfo.Description = item.Description
		if lat, err := strconv.ParseFloat(item.Lat, 64); err == nil {
			spotinfo.Lat = lat
		}
		if lon, err := strconv.ParseFloat(item.Lon, 64); err == nil {
			spotinfo.Lon = lon
		}
		if count, err := NewBikeCount(item.Recent.Datetime, item.Recent.Count); err == nil {
			spotinfo.Counts = append(spotinfo.Counts, count)
		}
		spotinfoList = append(spotinfoList, spotinfo)
	}
	return spotinfoList
}

//GetDistanceList GetSpotInfo構造体を返す
func GetDistanceList(body static.JDistancesBody) []string {
	var distances []string
	for _, item := range body.Items {
		distances = append(distances, item.Distance)
	}
	return distances
}

//GetGraphInfoByJGraphResponse GetSpotInfo構造体を返す
func GetGraphInfoByJGraphResponse(body static.JGraphResponse) GraphInfo {
	var graphInfo GraphInfo
	graphInfo.Height = body.Height
	graphInfo.Width = body.Width
	graphInfo.URL = body.URL
	graphInfo.Title = body.Title

	var spotinfo SpotInfo
	spotinfo.Area = body.Item.Area
	spotinfo.Spot = body.Item.Spot
	spotinfo.Name = body.Item.Name
	spotinfo.Description = body.Item.Description
	if lat, err := strconv.ParseFloat(body.Item.Lat, 64); err == nil {
		spotinfo.Lat = lat
	}
	if lon, err := strconv.ParseFloat(body.Item.Lon, 64); err == nil {
		spotinfo.Lon = lon
	}
	if s, err := NewBikeCount(body.Item.Recent.Datetime, body.Item.Recent.Count); err == nil {
		spotinfo.Counts = append(spotinfo.Counts, s)
	}
	graphInfo.SpotInfo = spotinfo

	return graphInfo
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////
//   共通構造体
/////////////////////////////////////////////////////////////////////////////////////////////////////////

//BikeCount 台数
type BikeCount struct {
	Time  time.Time
	Count int
}

//NewBikeCount JSONから変換するため
func NewBikeCount(datetime string, count string) (BikeCount, error) {
	t, err := time.Parse(JsonTimeLayout, datetime)
	if err != nil {
		return BikeCount{}, err
	}
	c, err := strconv.Atoi(count)
	if err != nil {
		return BikeCount{}, err
	}
	return BikeCount{Time: t, Count: c}, nil
}

//SpotInfo 駐輪場情報
type SpotInfo struct {
	Area, Spot, Name, Description string
	Lat, Lon                      float64
	Counts                        []BikeCount
}

//DistanceInfo 指定地点からの距離情報
type DistanceInfo struct {
	BaseLat, BaseLon float64
	Spots            []struct {
		SpotInfo SpotInfo
		Distance string
	}
}

//SpotName 駐輪場の名前
type SpotName struct {
	Area, Spot, Name string
}

//GraphInfo グラフ画像の情報
type GraphInfo struct {
	Title    string
	Width    string
	Height   string
	URL      string
	SpotInfo SpotInfo
}

//Users ユーザ情報
type Users struct {
	LineID    string
	SlackID   string
	Favorites []string
	Notifies  []string
	Histories []string
}

//ServiceStatus システム稼働状況
type ServiceStatus struct {
	Status     string
	Connection string
	Scraping   string
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////
//  リクエストパラメータ構造体
/////////////////////////////////////////////////////////////////////////////////////////////////////////

//SearchPlacesOption 駐輪場検索
type SearchPlacesOption struct {
	Area, Spot, Query string
	Places            []string
	Limit             int
	Sort              string
}

//GetQuery 検索条件作成
func (option SearchPlacesOption) GetQuery() string {
	var params []string
	if option.Area != "" {
		param := fmt.Sprintf("area=%s", option.Area)
		params = append(params, param)
	}
	if option.Spot != "" {
		param := fmt.Sprintf("spot=%s", option.Spot)
		params = append(params, param)
	}
	if option.Query != "" {
		param := fmt.Sprintf("q=%s", option.Query)
		params = append(params, param)
	}
	if len(option.Places) > 0 {
		param := fmt.Sprintf("places=%s", strings.Join(option.Places, ","))
		params = append(params, param)
	}
	if option.Sort != "" {
		param := fmt.Sprintf("sort=%s", option.Sort)
		params = append(params, param)
	}
	if option.Limit > 0 {
		param := fmt.Sprintf("limit=%d", option.Limit)
		params = append(params, param)
	}
	return strings.Join(params, `&`)
}

//SearchCountsOption 台数検索
type SearchCountsOption struct {
	Area, Spot, Day string
}

//GetQuery 検索条件作成
func (option SearchCountsOption) GetQuery() string {
	var params []string
	if option.Area != "" {
		param := fmt.Sprintf("area=%s", option.Area)
		params = append(params, param)
	}
	if option.Spot != "" {
		param := fmt.Sprintf("spot=%s", option.Spot)
		params = append(params, param)
	}
	if option.Day != "" {
		param := fmt.Sprintf("day=%s", option.Day)
		params = append(params, param)
	}
	return strings.Join(params, `&`)
}

//SearchDistanceOption 近いスポット検索
type SearchDistanceOption struct {
	Lat, Lon float64
}

//GetQuery 検索条件作成
func (option SearchDistanceOption) GetQuery() string {
	var params []string
	if option.Lat != 0 {
		param := fmt.Sprintf("lat=%v", option.Lat)
		params = append(params, param)
	}
	if option.Lon != 0 {
		param := fmt.Sprintf("lon=%v", option.Lon)
		params = append(params, param)
	}
	return strings.Join(params, `&`)
}

//SearchGraphOption グラフ検索
type SearchGraphOption struct {
	Area, Spot  string
	Property    string
	Days        []string
	DrawTitle   bool
	UploadImgur bool
}

//GetQuery 検索条件作成
func (option SearchGraphOption) GetQuery() string {
	var params []string
	if option.Area != "" {
		param := fmt.Sprintf("area=%s", option.Area)
		params = append(params, param)
	}
	if option.Spot != "" {
		param := fmt.Sprintf("spot=%s", option.Spot)
		params = append(params, param)
	}
	if option.Property != "" {
		param := fmt.Sprintf("property=%s", option.Property)
		params = append(params, param)
	}
	if len(option.Days) > 0 {
		param := fmt.Sprintf("days=%s", strings.Join(option.Days, `,`))
		params = append(params, param)
	}
	if option.DrawTitle {
		params = append(params, "title=yes")
	}
	if option.UploadImgur {
		params = append(params, "imgur=yes")
	}
	return strings.Join(params, `&`)
}
