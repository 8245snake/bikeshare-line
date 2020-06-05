package bikeshareapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/8245snake/bikeshare_api/src/lib/static"
)

//ApiClient クライアント
type ApiClient struct {
	Client   *http.Client
	CertKey  string
	Endpoint string
}

const (
	urlRoot         = "https://hanetwi.ddns.net/bikeshare/api/v1/"
	urlGetPlaces    = "places?"
	urlGetCounts    = "counts?"
	urlGetDistances = "distances?"
	urlGetAllPlaces = "all_places"
	urlGetStatus    = "status"
	urlGetUser      = "private/users"
	urlPostUser     = "private/user"
	urlGetGraph     = "https://hanetwi.ddns.net/bikeshare/graph?"
)

//NewApiClient コンストラクタ
func NewApiClient() ApiClient {
	var api ApiClient
	api.Client = &http.Client{}
	api.Endpoint = urlRoot
	return api
}

//SetCertKey キーを設定
func (api *ApiClient) SetCertKey(certKey string) {
	api.CertKey = certKey
}

//SetEndpoint エンドポイントを設定（デバッグ用）
func (api *ApiClient) SetEndpoint(endpoint string) {
	api.Endpoint = endpoint
}

//SendGetRequest GETリクエストを送信してレスポンスのバイト配列を得る
func (api *ApiClient) SendGetRequest(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("cert", api.CertKey)
	resp, err := api.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	} else if len(byteArray) == 0 {
		return nil, fmt.Errorf("レスポンスのデータ長が不正です")
	}
	return byteArray, nil
}

//SendPostRequest POSTリクエストを送信してレスポンスのバイト配列を得る
func (api *ApiClient) SendPostRequest(URL string, payload interface{}) ([]byte, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error")
	}
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	//ヘッダを設定
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("cert", api.CertKey)
	resp, err := api.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	} else if len(byteArray) == 0 {
		return nil, fmt.Errorf("レスポンスのデータ長が不正です")
	}
	return byteArray, nil
}

//GetPlaces 駐輪場検索
func (api *ApiClient) GetPlaces(option SearchPlacesOption) ([]SpotInfo, error) {
	url := api.Endpoint + urlGetPlaces + option.GetQuery()
	byteArray, err := api.SendGetRequest(url)
	if err != nil {
		return []SpotInfo{}, err
	}
	var data static.JPlacesBody
	if err := json.Unmarshal(([]byte)(byteArray), &data); err != nil {
		return []SpotInfo{}, err
	}
	return GetSpotInfoListByPlaces(data), nil
}

//GetCounts 台数検索
func (api ApiClient) GetCounts(option SearchCountsOption) (SpotInfo, error) {
	url := api.Endpoint + urlGetCounts + option.GetQuery()
	byteArray, err := api.SendGetRequest(url)
	if err != nil {
		return SpotInfo{}, err
	}
	var data static.JCountsBody
	if err := json.Unmarshal(([]byte)(byteArray), &data); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return SpotInfo{}, err
	}
	return GetSpotInfoByJCount(data), nil
}

//GetDistances 近いスポット検索
func (api ApiClient) GetDistances(option SearchDistanceOption) (DistanceInfo, error) {
	url := api.Endpoint + urlGetDistances + option.GetQuery()
	byteArray, err := api.SendGetRequest(url)
	if err != nil {
		return DistanceInfo{}, err
	}
	var data static.JDistancesBody
	if err := json.Unmarshal(([]byte)(byteArray), &data); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return DistanceInfo{}, err
	}
	spotinfoList := GetSpotInfoListByDistance(data)
	distanceInfo := DistanceInfo{BaseLat: option.Lat, BaseLon: option.Lon}
	distances := GetDistanceList(data)
	for i, item := range spotinfoList {
		distanceInfo.Spots = append(distanceInfo.Spots, struct {
			SpotInfo SpotInfo
			Distance string
		}{SpotInfo: item, Distance: distances[i]})
	}
	return distanceInfo, nil
}

//GetAllSpotNames すべてのスポットの名前だけ検索
func (api ApiClient) GetAllSpotNames() ([]SpotName, error) {
	byteArray, err := api.SendGetRequest(api.Endpoint + urlGetAllPlaces)
	if err != nil {
		return []SpotName{}, err
	}
	var data static.JAllPlacesBody
	if err := json.Unmarshal(([]byte)(byteArray), &data); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return []SpotName{}, err
	}

	var names []SpotName
	for _, item := range data.Items {
		names = append(names, SpotName{Area: item.Area, Spot: item.Spot, Name: item.Name})
	}
	return names, nil
}

//GetGraph グラフ検索
func (api ApiClient) GetGraph(option SearchGraphOption) (GraphInfo, error) {
	url := urlGetGraph + option.GetQuery()
	byteArray, err := api.SendGetRequest(url)
	if err != nil {
		return GraphInfo{}, err
	}
	var data static.JGraphResponse
	if err := json.Unmarshal(([]byte)(byteArray), &data); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return GraphInfo{}, err
	}
	return GetGraphInfoByJGraphResponse(data), nil
}

//GetUsers ユーザ情報取得
func (api ApiClient) GetUsers() ([]Users, error) {
	byteArray, err := api.SendGetRequest(api.Endpoint + urlGetUser)
	if err != nil {
		return nil, err
	}

	var data static.JUsers
	if err := json.Unmarshal(([]byte)(byteArray), &data); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return nil, err
	}
	var users []Users
	for _, jUser := range data.Users {
		users = append(users,
			Users{
				LineID:    jUser.LineID,
				SlackID:   jUser.SlackID,
				Favorites: jUser.Favorites,
				Histories: jUser.Histories,
				Notifies:  jUser.Notifies,
			},
		)
	}
	return users, nil
}

//GetStatus サービス稼働状況取得
func (api ApiClient) GetStatus() (static.JServiceStatus, error) {
	var data static.JServiceStatus
	byteArray, err := api.SendGetRequest(api.Endpoint + urlGetStatus)
	if err != nil {
		return data, err
	}

	if err := json.Unmarshal(([]byte)(byteArray), &data); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return data, err
	}

	return data, nil
}

//UpdateUser ユーザー更新
func (api ApiClient) UpdateUser(user Users) ([]Users, error) {

	juser := static.JUser{LineID: user.LineID, SlackID: user.SlackID, Favorites: user.Favorites, Histories: user.Histories, Notifies: user.Notifies}

	url := api.Endpoint + urlPostUser
	byteArray, err := api.SendPostRequest(url, juser)
	if err != nil {
		return nil, err
	}
	var data static.JUsers
	if err := json.Unmarshal(([]byte)(byteArray), &data); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return nil, err
	}
	var users []Users
	for _, jUser := range data.Users {
		users = append(users,
			Users{
				LineID:    jUser.LineID,
				SlackID:   jUser.SlackID,
				Favorites: jUser.Favorites,
				Histories: jUser.Histories,
				Notifies:  jUser.Notifies,
			},
		)
	}
	return users, nil
}
