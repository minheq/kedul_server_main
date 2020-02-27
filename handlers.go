package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/minheq/kedul_server_main/app"
	"github.com/minheq/kedul_server_main/auth"
	"github.com/minheq/kedul_server_main/errors"
)

type phoneNumberVerifyRequest struct {
	PhoneNumber string `json:"phone_number"`
	CountryCode string `json:"country_code"`
}

func (p *phoneNumberVerifyRequest) Bind(r *http.Request) error {
	return nil
}

type phoneNumberVerifyResponse struct {
	VerificationID string `json:"verification_id"`
}

func (rd *phoneNumberVerifyResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *server) handleLoginVerify(authService auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &phoneNumberVerifyRequest{}

		if err := render.Bind(r, data); err != nil {
			s.respondError(w, r, err)
			return
		}

		verificationID, err := authService.LoginVerify(r.Context(), data.PhoneNumber, data.CountryCode)

		if err != nil {
			s.respondError(w, r, err)
			return
		}

		render.Render(w, r, &phoneNumberVerifyResponse{VerificationID: verificationID})
	}
}

type phoneNumberCheckRequest struct {
	VerificationID string `json:"verification_id"`
	Code           string `json:"code"`
}

func (p *phoneNumberCheckRequest) Bind(r *http.Request) error {
	return nil
}

type loginCheckResponse struct {
	AccessToken string `json:"access_token"`
}

func (rd *loginCheckResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *server) handleLoginCheck(authService auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &phoneNumberCheckRequest{}

		if err := render.Bind(r, data); err != nil {
			s.respondError(w, r, err)
			return
		}

		accessToken, err := authService.LoginCheck(r.Context(), data.VerificationID, data.Code)

		if err != nil {
			s.respondError(w, r, err)
			return
		}

		render.Render(w, r, &loginCheckResponse{AccessToken: accessToken})
	}
}

type userResponse struct {
	ID                    string    `json:"id"`
	FullName              string    `json:"full_name"`
	PhoneNumber           string    `json:"phone_number"`
	CountryCode           string    `json:"country_code"`
	IsPhoneNumberVerified bool      `json:"is_phone_number_verified"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

func newUserResponse(user *auth.User) *userResponse {
	return &userResponse{
		ID:                    user.ID,
		FullName:              user.FullName,
		PhoneNumber:           user.PhoneNumber,
		CountryCode:           user.CountryCode,
		IsPhoneNumberVerified: user.IsPhoneNumberVerified,
		CreatedAt:             user.CreatedAt,
		UpdatedAt:             user.UpdatedAt,
	}
}

func (rd *userResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *server) handleGetCurrentUser(authService auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser, _ := r.Context().Value(userCtxKey).(*auth.User)

		render.Render(w, r, newUserResponse(currentUser))
	}
}

func (s *server) handleUpdateUserProfile(authService auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser, _ := r.Context().Value(userCtxKey).(*auth.User)

		input := &auth.UpdateUserProfileInput{}

		if err := s.decode(w, r, input); err != nil {
			s.respondError(w, r, err)
			return
		}

		user, err := authService.UpdateUserProfile(r.Context(), input, currentUser)

		if err != nil {
			s.respondError(w, r, err)
			return
		}

		render.Render(w, r, newUserResponse(user))
	}
}

func (s *server) handleUpdatePhoneNumberVerify(authService auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser, _ := r.Context().Value(userCtxKey).(*auth.User)
		data := &phoneNumberVerifyRequest{}

		if err := render.Bind(r, data); err != nil {
			s.respondError(w, r, err)
			return
		}

		verificationID, err := authService.UpdatePhoneNumberVerify(r.Context(), data.PhoneNumber, data.CountryCode, currentUser)

		if err != nil {
			s.respondError(w, r, err)
			return
		}

		render.Render(w, r, &phoneNumberVerifyResponse{VerificationID: verificationID})
	}
}

func (s *server) handleUpdatePhoneNumberCheck(authService auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser, _ := r.Context().Value(userCtxKey).(*auth.User)
		data := &phoneNumberCheckRequest{}

		if err := render.Bind(r, data); err != nil {
			s.respondError(w, r, err)
			return
		}

		user, err := authService.UpdatePhoneNumberCheck(r.Context(), data.VerificationID, data.Code, currentUser)

		if err != nil {
			s.respondError(w, r, err)
			return
		}

		render.Render(w, r, newUserResponse(user))
	}
}

type pageInfo struct {
	HasNextPage     bool `json:"has_next_page"`
	HasPreviousPage bool `json:"has_previous_page"`
}

type businessListResponse struct {
	TotalCount int                 `json:"total_count,omitempty"`
	PageInfo   *pageInfo           `json:"page_info,omitempty"`
	Data       []*businessResponse `json:"data,omitempty"`
}

func newBusinessListResponse(businesses []*app.Business) *businessListResponse {
	data := []*businessResponse{}

	for _, business := range businesses {
		data = append(data, newBusinessResponse(business))
	}

	return &businessListResponse{
		Data: data,
	}
}

func (rd *businessListResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *server) handleGetBusinessesByUserID(businessService app.BusinessService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.handleGetBusinessesByUserID"
		currentUser, _ := r.Context().Value(userCtxKey).(*auth.User)

		userID := chi.URLParam(r, "userID")

		if userID == "" {
			s.respondError(w, r, errors.Invalid(op, "missing param"))
		}

		businesses, err := businessService.GetBusinessesByUserID(r.Context(), userID, currentUser)

		if err != nil {
			s.respondError(w, r, err)
			return
		}

		render.Render(w, r, newBusinessListResponse(businesses))
	}
}

type businessResponse struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	Name           string    `json:"name"`
	ProfileImageID string    `json:"profile_image_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func newBusinessResponse(business *app.Business) *businessResponse {
	return &businessResponse{
		ID:             business.ID,
		UserID:         business.UserID,
		Name:           business.Name,
		ProfileImageID: business.ProfileImageID,
		CreatedAt:      business.CreatedAt,
		UpdatedAt:      business.UpdatedAt,
	}
}

func (rd *businessResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *server) handleGetBusiness(businessService app.BusinessService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.handleGetBusiness"

		businessID := chi.URLParam(r, "businessID")

		if businessID == "" {
			s.respondError(w, r, errors.Invalid(op, "missing param"))
		}

		business, err := businessService.GetBusinessByID(r.Context(), businessID)

		if err != nil {
			s.respondError(w, r, err)
			return
		}

		render.Render(w, r, newBusinessResponse(business))
	}
}

func (s *server) handleCreateBusiness(businessService app.BusinessService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser, _ := r.Context().Value(userCtxKey).(*auth.User)
		input := &app.CreateBusinessInput{}

		if err := s.decode(w, r, input); err != nil {
			s.respondError(w, r, err)
			return
		}

		business, err := businessService.CreateBusiness(r.Context(), currentUser.ID, input)

		if err != nil {
			s.respondError(w, r, err)
			return
		}

		render.Render(w, r, newBusinessResponse(business))
	}
}

func (s *server) handleUpdateBusiness(businessService app.BusinessService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.handleUpdateBusiness"
		currentUser, _ := r.Context().Value(userCtxKey).(*auth.User)
		input := &app.UpdateBusinessInput{}

		businessID := chi.URLParam(r, "businessID")

		if businessID == "" {
			s.respondError(w, r, errors.Invalid(op, "missing param"))
		}

		if err := s.decode(w, r, input); err != nil {
			s.respondError(w, r, err)
			return
		}

		business, err := businessService.UpdateBusiness(r.Context(), businessID, input, currentUser)

		if err != nil {
			s.respondError(w, r, err)
			return
		}

		render.Render(w, r, newBusinessResponse(business))
	}
}

func (s *server) handleDeleteBusiness(businessService app.BusinessService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.handleDeleteBusiness"
		currentUser, _ := r.Context().Value(userCtxKey).(*auth.User)
		businessID := chi.URLParam(r, "businessID")

		if businessID == "" {
			s.respondError(w, r, errors.Invalid(op, "missing param"))
		}

		business, err := businessService.DeleteBusiness(r.Context(), businessID, currentUser)

		if err != nil {
			s.respondError(w, r, err)
			return
		}

		render.Render(w, r, newBusinessResponse(business))
	}
}

type locationResponse struct {
	ID             string    `json:"id"`
	BusinessID     string    `json:"business_id"`
	Name           string    `json:"name"`
	ProfileImageID string    `json:"profile_image_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func newLocationResponse(location *app.Location) *locationResponse {
	return &locationResponse{
		ID:             location.ID,
		BusinessID:     location.BusinessID,
		Name:           location.Name,
		ProfileImageID: location.ProfileImageID,
		CreatedAt:      location.CreatedAt,
		UpdatedAt:      location.UpdatedAt,
	}
}

func (rd *locationResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type locationListResponse struct {
	TotalCount int                 `json:"total_count,omitempty"`
	PageInfo   *pageInfo           `json:"page_info,omitempty"`
	Data       []*locationResponse `json:"data,omitempty"`
}

func newLocationListResponse(locations []*app.Location) *locationListResponse {
	data := []*locationResponse{}

	for _, location := range locations {
		data = append(data, newLocationResponse(location))
	}

	return &locationListResponse{
		Data: data,
	}
}

func (rd *locationListResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *server) handleGetLocationsByUserIDAndBusinessID(locationsService app.LocationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.handleGetLocationsByUserIDAndBusinessID"
		currentUser, _ := r.Context().Value(userCtxKey).(*auth.User)

		userID := chi.URLParam(r, "userID")
		businessID := chi.URLParam(r, "businessID")

		if userID == "" {
			s.respondError(w, r, errors.Invalid(op, "missing param"))
		}

		if businessID == "" {
			s.respondError(w, r, errors.Invalid(op, "missing param"))
		}

		locations, err := locationsService.GetLocationsByUserIDAndBusinessID(r.Context(), userID, businessID, currentUser)

		if err != nil {
			s.respondError(w, r, err)
			return
		}

		render.Render(w, r, newLocationListResponse(locations))
	}
}

func (s *server) handleGetLocation(locationService app.LocationService, permissionsService app.PermissionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.handleGetLocation"
		currentUser, _ := r.Context().Value(userCtxKey).(*auth.User)

		locationID := chi.URLParam(r, "locationID")

		if locationID == "" {
			s.respondError(w, r, errors.Invalid(op, "missing param"))
		}

		actor, err := permissionsService.GetEmployeeActor(r.Context(), currentUser.ID, locationID)

		if err != nil {
			s.respondError(w, r, err)
			return
		}

		location, err := locationService.GetLocationByID(r.Context(), locationID, actor)

		if err != nil {
			s.respondError(w, r, err)
			return
		}

		render.Render(w, r, newLocationResponse(location))
	}
}

func (s *server) handleCreateLocation(locationService app.LocationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser, _ := r.Context().Value(userCtxKey).(*auth.User)
		input := &app.CreateLocationInput{}

		if err := s.decode(w, r, input); err != nil {
			s.respondError(w, r, err)
			return
		}

		location, err := locationService.CreateLocation(r.Context(), input, currentUser)

		if err != nil {
			s.respondError(w, r, err)
			return
		}

		render.Render(w, r, newLocationResponse(location))
	}
}

func (s *server) handleUpdateLocation(locationService app.LocationService, permissionsService app.PermissionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.handleUpdateLocation"
		currentUser, _ := r.Context().Value(userCtxKey).(*auth.User)
		input := &app.UpdateLocationInput{}

		locationID := chi.URLParam(r, "locationID")

		if locationID == "" {
			s.respondError(w, r, errors.Invalid(op, "missing param"))
		}

		if err := s.decode(w, r, input); err != nil {
			s.respondError(w, r, err)
			return
		}

		actor, err := permissionsService.GetEmployeeActor(r.Context(), currentUser.ID, locationID)

		if err != nil {
			s.respondError(w, r, err)
			return
		}

		location, err := locationService.UpdateLocation(r.Context(), locationID, input, actor)

		if err != nil {
			s.respondError(w, r, err)
			return
		}

		render.Render(w, r, newLocationResponse(location))
	}
}

func (s *server) handleDeleteLocation(locationService app.LocationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.handleDeleteLocation"
		currentUser, _ := r.Context().Value(userCtxKey).(*auth.User)
		locationID := chi.URLParam(r, "locationID")

		if locationID == "" {
			s.respondError(w, r, errors.Invalid(op, "missing param"))
		}

		location, err := locationService.DeleteLocation(r.Context(), locationID, currentUser)

		if err != nil {
			s.respondError(w, r, err)
			return
		}

		render.Render(w, r, newLocationResponse(location))
	}
}
