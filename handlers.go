package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"net/http"
)

func getMeHandler(w http.ResponseWriter, r *http.Request, u User, users UserRepository) {
	w.Write([]byte(u.FavoriteCake + "\n"))
	w.Write([]byte(u.Email))
}

func getCakeHandler(w http.ResponseWriter, r *http.Request, u User, users UserRepository) {
	w.Write([]byte(u.FavoriteCake))
}

func getEmailHandler(w http.ResponseWriter, r *http.Request, u User, users UserRepository) {
	w.Write([]byte(u.Email))
}

func (uServ UserService) updateCakeHandler(w http.ResponseWriter, r *http.Request, u User, users UserRepository) {
	params := &UserRegisterParams{}
	err := json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		handleError(errors.New("could not read params"), w)
		return
	}
	if err := validateCake(params); err != nil {
		handleError(err, w)
		return
	}
	jwtAuth, _ := JwtService.getJWT(w, r)
	newCake, err := users.Get(jwtAuth.Email)
	if err != nil {
		handleError(err, w)
		return
	}
	newCake.FavoriteCake = params.FavoriteCake
	err = uServ.repository.Update(jwtAuth.Email, newCake)
	if err != nil {
		handleError(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("cake updated"))
}

func (uServ UserService) updateEmailHandler(w http.ResponseWriter, r *http.Request, u User, users UserRepository) {
	params := &UserRegisterParams{}
	err := json.NewDecoder(r.Body).Decode(params)

	if err != nil {
		handleError(errors.New("could not read params"), w)
		return
	}
	if err := validateEmail(params); err != nil {
		handleError(err, w)
		return
	}

	jwtAuth, _ := JwtService.getJWT(w, r)
	newEmail, err := users.Get(jwtAuth.Email)
	if err != nil {
		handleError(err, w)
		return
	}
	newEmail.Email = params.Email
	uServ.repository.Delete(u.Email)
	uServ.repository.Add(params.Email, newEmail)
	if err != nil {
		handleError(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("email updated"))
}

func (uServ UserService) updatePasswordHandler(w http.ResponseWriter, r *http.Request, u User, users UserRepository) {
	params := &UserRegisterParams{}
	err := json.NewDecoder(r.Body).Decode(params)

	if err != nil {
		handleError(errors.New("could not read params"), w)
		return
	}
	if err := validatePassword(params); err != nil {
		handleError(err, w)
		return
	}

	jwtAuth, _ := JwtService.getJWT(w, r)
	newPass, err := users.Get(jwtAuth.Email)
	if err != nil {
		handleError(err, w)
		return
	}
	newPass.PasswordDigest = string(md5.New().Sum([]byte(params.Password)))

	err = uServ.repository.Update(jwtAuth.Email, newPass)
	if err != nil {
		handleError(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("password updated"))
}
