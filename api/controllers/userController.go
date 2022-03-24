//userController.go
package controllers

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"

    "github.com/khelechy/rielzapi/api/models"
    "github.com/khelechy/rielzapi/api/responses"
    "github.com/khelechy/rielzapi/utils"
)

// UserSignUp controller for creating new users
func (a *App) UserSignUp(w http.ResponseWriter, r *http.Request) {
    var resp = map[string]interface{}{"status": "success", "message": "Registered successfully"}

    user := &models.User{}
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        responses.ERROR(w, http.StatusBadRequest, err)
        return
    }

    err = json.Unmarshal(body, &user)
    if err != nil {
        responses.ERROR(w, http.StatusBadRequest, err)
        return
    }

    usr, _ := user.GetUser(a.DB)
    if usr != nil {
        resp["status"] = "failed"
        resp["message"] = "User already registered, please login"
        responses.JSON(w, http.StatusBadRequest, resp)
        return
    }

    user.Prepare() // here strip the text of white spaces

    err = user.Validate("") // default were all fields(email, lastname, firstname, password, profileimage) are validated
    if err != nil {
        responses.ERROR(w, http.StatusBadRequest, err)
        return
    }
    userCreated, err := user.SaveUser(a.DB)
    if err != nil {
        responses.ERROR(w, http.StatusBadRequest, err)
        return
    }
    resp["user"] = userCreated
    responses.JSON(w, http.StatusCreated, resp)
    return
}

// Login signs in users
func (a *App) Login(w http.ResponseWriter, r *http.Request) {
    var resp = map[string]interface{}{"status": "success", "message": "logged in"}

    user := &models.User{}
    body, err := ioutil.ReadAll(r.Body) // read user input from request
    if err != nil {
        responses.ERROR(w, http.StatusBadRequest, err)
        return
    }

    err = json.Unmarshal(body, &user)
    if err != nil {
        responses.ERROR(w, http.StatusBadRequest, err)
        return
    }

    user.Prepare() // here strip the text of white spaces

    err = user.Validate("login") // fields(email, password) are validated
    if err != nil {
        responses.ERROR(w, http.StatusBadRequest, err)
        return
    }

    usr, err := user.GetUser(a.DB)
    if err != nil {
        responses.ERROR(w, http.StatusInternalServerError, err)
        return
    }

    if usr == nil { // user is not registered
        resp["status"] = "failed"
        resp["message"] = "Login failed, please signup"
        responses.JSON(w, http.StatusBadRequest, resp)
        return
    }

    err = models.CheckPasswordHash(user.Password, usr.Password)
    if err != nil {
        resp["status"] = "failed"
        resp["message"] = "Login failed, please try again"
        responses.JSON(w, http.StatusForbidden, resp)
        return
    }
    token, err := utils.EncodeAuthToken(usr.ID)
    if err != nil {
        responses.ERROR(w, http.StatusBadRequest, err)
        return
    }

    resp["token"] = token
    responses.JSON(w, http.StatusOK, resp)
    return
}

func (a *App) GetUsers(w http.ResponseWriter, r *http.Request) {
    houses, err := models.GetAllUsers(a.DB)
    if err != nil {
        responses.ERROR(w, http.StatusInternalServerError, err)
        return
    }
    responses.JSON(w, http.StatusOK, houses)
    return
}

func (a *App) GetUserById(w http.ResponseWriter, r *http.Request){

    vars := mux.Vars(r)

    id, _ := strconv.Atoi(vars["id"])

    user, err := models.GetUserById(id, a.DB)
    if err != nil{
        responses.ERROR(w, http.StatusInternalServerError, err)
        return
    }
    responses.JSON(w, http.StatusOK, user)
    return
}

func (a *App) UpdateUser(w http.ResponseWriter, r *http.Request) {
    var resp = map[string]interface{}{"status": "success", "message": "User updated successfully"}

    vars := mux.Vars(r)

    myUser := r.Context().Value("userID").(float64)
    userID := uint(myUser)

    id, _ := strconv.Atoi(vars["id"])

    user, err := models.GetUserById(id, a.DB)

    if user.ID != userID {
        resp["status"] = "failed"
        resp["message"] = "Unauthorized user update"
        responses.JSON(w, http.StatusUnauthorized, resp)
        return
    }

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        responses.ERROR(w, http.StatusBadRequest, err)
        return
    }

    userUpdate := models.User{}
    if err = json.Unmarshal(body, &userUpdate); err != nil {
        responses.ERROR(w, http.StatusBadRequest, err)
        return
    }

    userUpdate.Prepare()

    _, err = userUpdate.UpdateUser(id, a.DB)
    if err != nil {
        responses.ERROR(w, http.StatusInternalServerError, err)
        return
    }

    responses.JSON(w, http.StatusOK, resp)
    return
}
