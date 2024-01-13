package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)


var (
    DOMAIN = "https://restapi.cutout.pro"
    HEADERS = map[string]string{
        "Host":                 "restapi.cutout.pro",
		"User-Agent":           "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0",
		"Accept":               "application/json, text/plain, */*",
		"Accept-Language":      "en-US,en;q=0.5",
		"Accept-Encoding":      "",
		"Token":                "",
		"Origin":               "https://www.cutout.pro",
		"Dnt":                  "1",
		"Referer":              "https://www.cutout.pro/",
		"Sec-Fetch-Dest":       "empty",
		"Sec-Fetch-Mode":       "cors",
		"Sec-Fetch-Site":       "same-site",
		"Te":                   "trailers",
	}
    CLIENT = &http.Client{}
)

type NewError struct {
    Message string
}

func (i *NewError) Error() string {
    return i.Message
}

func ReplaceAmpersand(urlString string) (string,error) {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return "",err
	}
	parsedURL.RawQuery = strings.ReplaceAll(parsedURL.RawQuery, "&", "%26")

	return parsedURL.String(),nil
}

func UploadByUrl(url string) (string,error) {
    goodUrl,err := ReplaceAmpersand(url);if err != nil {
        return "",err
    }
    requestUrl := DOMAIN + "/oss/uploadByUrl?url=" + goodUrl

    resp, err := http.Get(requestUrl)
    if err != nil {
        return "",err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return "",&NewError{Message: "CE : error happend in uplaodUrl function :("}
    }

    var data map[string]interface{}

    decoder := json.NewDecoder(resp.Body)
    if err := decoder.Decode(&data); err != nil {
        return "",err
    }

    newUrl, ok := data["data"].(string)
    if !ok {
        return "", &NewError{Message: "CE : couldnt get 'data' from json ? "}
    }

    return newUrl,nil
}

type cutoutResponse struct {
    Code int `json:"code"`
    Msg string `json:"msg"`
    Data interface{} `json:"data"`
}

func EnhanceRequest(url string) (string,error) {
    goodUrl, err := ReplaceAmpersand(url);if err != nil {
        return "", err
    }

    token, err  := ReadFirstToken()
    if err != nil {
        return "",err
    }

    requestUrl := fmt.Sprintf(
        DOMAIN + "/webMatting/photoEnhancer/submitTaskByUrl?token=%s&imageUrl=%s",
        token,
        goodUrl,
    )

    req, err := http.NewRequest("GET",requestUrl,nil)
    if err != nil {
        return "", err
    }

    for key,value := range HEADERS {
        req.Header.Add(key, value)
    }

    req.Header.Set("Token", token)

    resp, err := CLIENT.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var data cutoutResponse

    decoder := json.NewDecoder(resp.Body)
    if err := decoder.Decode(&data);err != nil {return "", err}

    if data.Code != 0 && resp.StatusCode != 200 {
        return "", &NewError{Message: data.Msg}
    }

    return data.Data.(string) , nil
}

func GetHDdownloadUrl(id string) (string, error) {
    token, err := ReadFirstToken();if err != nil {
        return "", err
    }

    url := fmt.Sprintf(DOMAIN + "/webMatting/photoEnhancer/download?token=%s&id=%s",token,id)

    req, err := http.NewRequest("GET",url,nil)
    if err != nil {
        return "", err
    }

    for key, value := range HEADERS {
        req.Header.Add(key,value)
    }

    req.Header.Set("Token",token)

    resp, err := CLIENT.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var data cutoutResponse

    decoder := json.NewDecoder(resp.Body)
    if err := decoder.Decode(&data);err != nil {return "", err}

    if data.Code == 4001 {
        if err := DeleteFirstLine(); err != nil {
            return "", err
        }
        return GetHDdownloadUrl(id)
    }else if data.Code == 5003{
        time.Sleep(time.Second * 3)
        return GetHDdownloadUrl(id)
    }else if data.Code == 4002{
        code := fmt.Sprint(data.Code)
        return code, nil

    }else if data.Code != 0 && resp.StatusCode != 200 {
        return "", &NewError{Message: data.Msg}
    }

    return data.Data.(string), nil
}


func Test(picUrl string) (string, error) {
    fmt.Println("Uploading Picture ...")
    newUrl, err := UploadByUrl(picUrl)
    if err != nil {return "", err }

    fmt.Println("Enhancing ...")
    taskId, err := EnhanceRequest(newUrl)
    if err != nil {return "", err }

    fmt.Println("Getting the Download Url ...")
    downloadUrl, err := GetHDdownloadUrl(taskId)
    if err != nil {return "", err}
    if downloadUrl == "4002" {
        return Test(picUrl)
    }

    fmt.Println("Everything is Done! this is the enhanced Image -> " + downloadUrl )

    return downloadUrl,nil
}
