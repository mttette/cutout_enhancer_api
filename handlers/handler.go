package handlers

import (
	"cutout_enhancer/models"
	"encoding/json"
	"fmt"
	"net/http"
)

type RequestPayload struct {
	URL string `json:"url"`
}

type ResponseInstance struct {
    DATA string `json:"data"`
}

func GetEnhanceUrl(url string,tries int) (string,error) {
    dataBaseUrl, err := models.UploadByUrl(url); if err != nil {return "", err}

    taskId, err := models.EnhanceRequest(dataBaseUrl); if err != nil {return "", err}

    downloadUrl, err := models.GetHDdownloadUrl(taskId); if err != nil {return "", err}

    if tries < 2 && downloadUrl == "4002" {
        tries += 1
        return GetEnhanceUrl(url,tries)
    }else if downloadUrl == "4002" {
        return "", &models.NewError{Message: "Unknown error Happend!"}
    }

    return downloadUrl, nil
}
func PostEnhanceHandler(w http.ResponseWriter, r *http.Request) {
    var payload RequestPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

    if payload.URL == "" {
		http.Error(w, "Missing 'url' in payload", http.StatusBadRequest)
		return
	}

    url, err  := models.ReplaceAmpersand(payload.URL); if err != nil {
        http.Error(w, err.Error(),http.StatusInternalServerError)
        return
    }

    responseUrl, err := GetEnhanceUrl(url,0)
    if err != nil {
        http.Error(w, err.Error(),http.StatusInternalServerError)
        return
    }

    //response := ResponseInstance{DATA: responseUrl}
    //w.Header().Set("Content-Type", "application/text")
	//w.WriteHeader(http.StatusOK)
	//json.NewEncoder(w).Encode(response)
    fmt.Fprint(w,responseUrl)

}
