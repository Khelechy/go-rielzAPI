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

// CreateHouse godoc
// @Summary Create House for landlord
// @Accept  json
// @Produce  json
// @Router /api/houses [post]
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

// GetHouses godoc
// @Summary Get All Houses
// @Accept  json
// @Produce  json
// @Router /api/houses [get]
func (a *App) GetHouses(w http.ResponseWriter, r *http.Request) {
    houses, err := models.GetHouses(a.DB)
    if err != nil {
        responses.ERROR(w, http.StatusInternalServerError, err)
        return
    }
    responses.JSON(w, http.StatusOK, houses)
    return
}

// GetHouses By Landlord godoc
// @Summary Get All Houses By Landlord
// @Accept  json
// @Produce  json
// @Router /api/houses/landlord [get]
func (a *App) GetHousesByLandlord(w http.ResponseWriter, r *http.Request) {
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

// GetHouses By Landlord godoc
// @Summary Get All Houses By Landlord
// @Accept  json
// @Produce  json
// @Router /api/houses/landlord/id [get]
func (a *App) GetHousesByLandlordId(w http.ResponseWriter, r *http.Request) {

    vars := mux.Vars(r)

    id, _ := strconv.Atoi(vars["id"])

    userID := uint(id)
    houses, err := models.GetHousesByLandLord(userID, a.DB)
    if err != nil {
        responses.ERROR(w, http.StatusInternalServerError, err)
        return
    }
    responses.JSON(w, http.StatusOK, houses)
    return
}


// GetHouses By state godoc
// @Summary Get All Houses By state
// @Accept  json
// @Produce  json
// @Router /api/houses/state [get]
func (a *App) GetHousesByState(w http.ResponseWriter, r *http.Request) {

    state := mux.Vars(r)["state"]
    houses, err := models.GetHousesByState(state, a.DB)
    if err != nil {
        responses.ERROR(w, http.StatusInternalServerError, err)
        return
    }
    responses.JSON(w, http.StatusOK, houses)
    return
}

// GetHouses
// @Summary Get All Houses By Id
// @Accept  json
// @Produce  json
// @Router /api/houses/id [get]
func (a *App) GetHouseById(w http.ResponseWriter, r *http.Request){

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

// AddTenant By Landlord godoc
// @Summary Add Tenant By Landlord
// @Accept  json
// @Produce  json
// @Router /api/houses/tenant [post]
func (a *App) AddTenant(w http.ResponseWriter, r *http.Request){
    var resp = map[string]interface{}{"status": "success", "message": "Tenant added successfully"}

    tenant := &models.Tenant{}
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        responses.ERROR(w, http.StatusBadRequest, err)
        return
    }

    err = json.Unmarshal(body, &tenant)
    if err != nil {
        responses.ERROR(w, http.StatusBadRequest, err)
        return
    }

    tenant.Prepare()
    err = tenant.Validate()

    if err != nil {
        responses.ERROR(w, http.StatusBadRequest, err)
        return
    }

    house, err := models.GetHouseById(tenant.HouseId, a.DB)
    if err != nil{
        responses.ERROR(w, http.StatusInternalServerError, err)
        return
    }

    if house.AvailableRooms <= 0 {
        resp["status"] = "failed"
        resp["message"] = "There are no available rooms"
        responses.JSON(w, http.StatusInternalServerError, resp)
        return
    }

    tenantCreated, err := tenant.SaveTenant(a.DB)
    if err != nil {
        responses.ERROR(w, http.StatusBadRequest, err)
        return
    }

    house.AvailableRooms = house.AvailableRooms - 1

    _, err = house.UpdateHouse(tenant.HouseId, a.DB)
    if err != nil {
        responses.ERROR(w, http.StatusInternalServerError, err)
        return
    }

    resp["tenant"] = tenantCreated
    responses.JSON(w, http.StatusCreated, resp)
    return
}

// UpdateHouse By Landlord godoc
// @Summary Update Houses By Landlord
// @Accept  json
// @Produce  json
// @Router /api/houses/id [put]
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

// DeleteHouse By Landlord godoc
// @Summary Delete house By Landlord
// @Accept  json
// @Produce  json
// @Router /api/houses/id [delete]
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