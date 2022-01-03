package controllers

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"

    "github.com/khelechy/rielzapi/api/models"
    "github.com/khelechy/rielzapi/api/responses"
)

// CreateHouse parses request, validates data and saves the new house
func (a *App) CreateHouse(w http.ResponseWriter, r *http.Request) {
    var resp = map[string]interface{}{"status": "success", "message": "House successfully created"}

    user := r.Context().Value("userID").(float64)
    house := &models.House{}
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        responses.ERROR(w, http.StatusBadRequest, err)
        return
    }

    err = json.Unmarshal(body, &house)
    if err != nil {
        responses.ERROR(w, http.StatusBadRequest, err)
        return
    }

    house.Prepare() // strip away any white spaces

    if err = house.Validate(); err != nil {
        responses.ERROR(w, http.StatusBadRequest, err)
        return
    }

    house.UserID = uint(user)

    houseCreated, err := house.Save(a.DB)
    if err != nil {
        responses.ERROR(w, http.StatusBadRequest, err)
        return
    }

    resp["house"] = houseCreated
    responses.JSON(w, http.StatusCreated, resp)
    return
}

func (a *App) GetHouses(w http.ResponseWriter, r http.Request) {
    houses, err := models.GetHouses(a.DB)
    if err != nil {
        responses.ERROR(w, http.StatusInternalServerError, err)
        return
    }
    responses.JSON(w, http.StatusOK, houses)
    return
}

func (a *App) GetHousesByLandlord(w http.ResponseWriter, r http.Request) {
    user := r.Context().Value("userID").(float64)
    userID := uint(user)
    houses, err := models.GetHousesByLandLord(userID, a.DB)
    if err != nil {
        responses.ERROR(w, http.StatusInternalServerError, err)
        return
    }
    responses.JSON(w, http.StatusOK, houses)
    return
}

func (a *App) GetHouseById(w http.ResponseWriter, r http.Request){
    var resp = map[string]interface{}{"status": "success", "message": "House fetched successfully"}
    vars := mux.Vars(r)
    id, _ := strconv.Atoi(vars["id"])

    house, err := models.GetHouseById(id, a.DB)
    if err != nil{
        responses.ERROR(w, http.StatusInternalServerError, err)
        return
    }
    responses.JSON(w, http.StatusOK, house)
    return
}

func (a *App) UpdateHouse(w http.ResponseWriter, r *http.Request) {
    var resp = map[string]interface{}{"status": "success", "message": "House updated successfully"}

    vars := mux.Vars(r)

    user := r.Context().Value("userID").(float64)
    userID := uint(user)

    id, _ := strconv.Atoi(vars["id"])

    house, err := models.GetHouseById(id, a.DB)

    if house.UserID != userID {
        resp["status"] = "failed"
        resp["message"] = "Unauthorized house update"
        responses.JSON(w, http.StatusUnauthorized, resp)
        return
    }

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        responses.ERROR(w, http.StatusBadRequest, err)
        return
    }

    houseUpdate := models.House{}
    if err = json.Unmarshal(body, &houseUpdate); err != nil {
        responses.ERROR(w, http.StatusBadRequest, err)
        return
    }

    houseUpdate.Prepare()

    _, err = houseUpdate.UpdateHouse(id, a.DB)
    if err != nil {
        responses.ERROR(w, http.StatusInternalServerError, err)
        return
    }

    responses.JSON(w, http.StatusOK, resp)
    return
}

func (a *App) DeleteHouse(w http.ResponseWriter, r *http.Request) {
    var resp = map[string]interface{}{"status": "success", "message": "House deleted successfully"}

    vars := mux.Vars(r)

    user := r.Context().Value("userID").(float64)
    userID := uint(user)

    id, _ := strconv.Atoi(vars["id"])

    house, err := models.GetHouseById(id, a.DB)

    if house.UserID != userID {
        resp["status"] = "failed"
        resp["message"] = "Unauthorized house delete"
        responses.JSON(w, http.StatusUnauthorized, resp)
        return
    }

    err = models.DeleteHouse(id, a.DB)
    if err != nil {
        responses.ERROR(w, http.StatusInternalServerError, err)
        return
    }
    responses.JSON(w, http.StatusOK, resp)
    return
}