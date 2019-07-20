package hl

import (
	"bytes"
	"github.com/buaazp/fasthttprouter"
	"github.com/ravlio/highloadcup2018/account"
	"github.com/ravlio/highloadcup2018/errors"
	"github.com/ravlio/highloadcup2018/gojay"
	"github.com/ravlio/highloadcup2018/metrics"
	"github.com/ravlio/highloadcup2018/requests"
	"github.com/ravlio/highloadcup2018/requests/likes"
	"github.com/ravlio/highloadcup2018/utils"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
	"time"
)

type HTTPServer struct {
	DB   *DB
	Addr string
}

var (
	strApplicationJSON = []byte("application/json")
	filterRespStart    = []byte(`{"accounts":`)
	filterRespEnd      = []byte(`}`)
	emptyAccountsResp  = []byte(`{"accounts":[]}`)
	groupsRespStart    = []byte(`{"groups":`)
	groupsRespEnd      = []byte(`}`)
	emptyGroupsResp    = []byte(`{"groups":[]}`)
	notFoundResp       = []byte(`{}`)
)

var filterDur = metrics.NewDurarion("filter")
var groupDur = metrics.NewDurarion("group")

func (s *HTTPServer) Filter(ctx *fasthttp.RequestCtx) {
	t := time.Now()
	defer filterDur.Write(t)

	ctx.Response.Header.SetContentTypeBytes(strApplicationJSON)
	f := &requests.Filter{}
	err, status := f.FillRequest(ctx)
	if err != nil {
		s.Error(ctx, err, status)
		return
	}
	// постфактум заполняем словари. Сделано для того, чтобы сократить расходы, если окажется, что в пришедшем запросе
	// последний аргумент невалидный, а мы до этого уже сделали кучу выборок из словарей
	err = f.FillWithDicts()
	if err != nil {
		if errors.Cause(err) == requests.ErrNotFoundInDict { // если в каком-то из словарей не найдено значение, значит это 404
			// ctx.Response.SetStatusCode(http.StatusNotFound)
			ctx.Response.SetStatusCode(http.StatusOK)
			ctx.Write(emptyAccountsResp)
		} else {
			s.Error(ctx, err, fasthttp.StatusBadRequest)
		}
		return
	}

	accs, err := s.DB.FilterAccounts(f)
	if err != nil {
		s.Error(ctx, errors.Wrap(err, "error filter accounts"), fasthttp.StatusInternalServerError)
		return
	}

	// Оборачиваем json-массив объектом

	if len(accs) > 0 {
		ctx.Write(filterRespStart)

		cont := &account.AccountsContainer{
			Accounts: accs,
			Fields:   f.FilledFields(),
		}

		b, err := gojay.MarshalJSONArray(cont)
		if err != nil {
			s.Error(ctx, errors.Wrap(err, "error marshal accounts response"), fasthttp.StatusInternalServerError)
		}

		ctx.Write(b)
		ctx.Write(filterRespEnd)
	} else {
		ctx.Write(emptyAccountsResp)
	}

	ctx.Response.SetStatusCode(fasthttp.StatusOK)
}

func (s *HTTPServer) Group(ctx *fasthttp.RequestCtx) {
	t := time.Now()
	defer groupDur.Write(t)
	ctx.Response.Header.SetContentTypeBytes(strApplicationJSON)

	g := &requests.Group{}
	err, status := g.FillRequest(ctx)
	if err != nil {
		s.Error(ctx, err, status)
		return
	}

	err = g.FillWithDicts()
	if err != nil {
		if errors.Cause(err) == requests.ErrNotFoundInDict { // если в каком-то из словарей не найдено значение, значит это 404
			// ctx.Response.SetStatusCode(http.StatusNotFound)
			ctx.Response.SetStatusCode(http.StatusOK)
			ctx.Write(emptyGroupsResp)
		} else {
			s.Error(ctx, err, fasthttp.StatusBadRequest)
		}
		return
	}

	groups, err := s.DB.GroupAccounts(g)
	if err != nil {
		s.Error(ctx, errors.Wrap(err, "error group accounts"), fasthttp.StatusInternalServerError)
		return
	}

	// Оборачиваем json-массив объектом

	if len(groups) == 0 {
		ctx.Write(emptyGroupsResp)
		return
	}

	ctx.Write(groupsRespStart)

	cont := &account.GroupsContainer{
		Groups: groups,
	}

	b, err := gojay.MarshalJSONArray(cont)
	if err != nil {
		s.Error(ctx, errors.Wrap(err, "error marshal groups response"), fasthttp.StatusInternalServerError)
	}

	ctx.Write(b)
	ctx.Write(groupsRespEnd)

	ctx.Response.SetStatusCode(fasthttp.StatusOK)
}

func (s *HTTPServer) Create(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetContentTypeBytes(strApplicationJSON)

	for _, arg := range utils.GetURIArgs(ctx) {
		if bytes.Compare(arg.Key, []byte("query_id")) > 0 {
			s.Error(ctx, errors.New("unknown query arg in create"), fasthttp.StatusBadRequest)
			return
		}
	}

	req := &requests.AccountRequest{}

	err := gojay.UnmarshalJSONObject(ctx.Request.Body(), req)
	if err != nil {
		s.Error(ctx, errors.Wrap(err, "error unmarshal request"), fasthttp.StatusBadRequest)
		return
	}

	if !req.ID.IsSet {
		s.Error(ctx, errors.New("empty id"), fasthttp.StatusBadRequest)
		return
	}

	err = s.DB.CreateAccount(req)

	if err != nil {
		if err == account.ErrEmailAlreadyExists || err == account.ErrPhoneAlreadyExists {
			s.Error(ctx, errors.Wrap(err, "error during adding account"), fasthttp.StatusBadRequest)
		} else {
			s.Error(ctx, errors.Wrap(err, "error during adding account"), fasthttp.StatusInternalServerError)
		}
		return
	}

	ctx.Write(notFoundResp)

	ctx.Response.SetStatusCode(fasthttp.StatusCreated)
}

func (s *HTTPServer) Update(ctx *fasthttp.RequestCtx, accID uint32) {
	ctx.Response.Header.SetContentTypeBytes(strApplicationJSON)

	for _, arg := range utils.GetURIArgs(ctx) {
		if bytes.Compare(arg.Key, []byte("query_id")) > 0 {
			s.Error(ctx, errors.New("unknown query arg in update"), fasthttp.StatusBadRequest)
			return
		}
	}

	req := &requests.AccountRequest{}

	err := gojay.UnmarshalJSONObject(ctx.Request.Body(), req)
	if err != nil {
		s.Error(ctx, errors.Wrap(err, "error unmarshal request"), fasthttp.StatusBadRequest)
		return
	}

	if req.ID.IsSet {
		s.Error(ctx, errors.Wrap(err, "id can't be in body"), fasthttp.StatusBadRequest)
		return
	}

	err = s.DB.UpdateAccount(req, accID)

	if err != nil {
		if err == account.ErrUnexistingAccount {
			s.Error(ctx, errors.Wrap(err, "error during updating account"), fasthttp.StatusNotFound)
		} else {
			s.Error(ctx, errors.Wrap(err, "error during updating account"), fasthttp.StatusBadRequest)
		} /*if err==account.ErrEmailAlreadyExists || err==account.ErrPhoneAlreadyExists {
			s.Error(ctx,errors.Wrap(err,"error during updating account"),fasthttp.StatusBadRequest)
		} else {
			s.Error(ctx,errors.Wrap(err,"error during updating account"),fasthttp.StatusInternalServerError)
		}*/
		return
	}

	ctx.Write(notFoundResp)
	ctx.Response.SetStatusCode(202)
}

func (s *HTTPServer) Likes(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetContentTypeBytes(strApplicationJSON)

	for _, arg := range utils.GetURIArgs(ctx) {
		if bytes.Compare(arg.Key, []byte("query_id")) > 0 {
			s.Error(ctx, errors.New("unknown query arg in create"), fasthttp.StatusBadRequest)
			return
		}
	}

	req := &likes.LikesRequest{}

	err := gojay.UnmarshalJSONObject(ctx.Request.Body(), req)
	if err != nil {
		s.Error(ctx, errors.Wrap(err, "error unmarshal request"), fasthttp.StatusBadRequest)
		return
	}

	if len(req.Likes) == 0 {
		s.Error(ctx, errors.New("empty likes requrst"), fasthttp.StatusBadRequest)
	}
	if err := s.DB.ValidateLikes(req.Likes); err != nil {
		s.Error(ctx, errors.Wrap(err, "error during adding likes"), fasthttp.StatusBadRequest)
		return
	}

	s.DB.IndexLikes(req.Likes)

	ctx.Write(notFoundResp)

	ctx.Response.SetStatusCode(fasthttp.StatusAccepted)
}

func (s *HTTPServer) RouterAccounts(ctx *fasthttp.RequestCtx) {
	token, ok := ctx.UserValue("token").(string)

	if ok && token == "new" {
		s.Create(ctx)
		return
	}

	if ok && token == "likes" {
		s.Likes(ctx)
		return
	}

	accID, err := strconv.ParseUint(token, 10, 32)
	if err != nil {
		s.Error(ctx, errors.New("error cast token to uint32"), fasthttp.StatusNotFound) // Должен быть 400, но от решения ожидается 404
		return
	}

	// Update
	if accID <= 0 {
		s.Error(ctx, errors.New("accID can't be <=0"), fasthttp.StatusBadRequest)
		return
	}

	s.Update(ctx, uint32(accID))
	return
}

func (s *HTTPServer) Error(ctx *fasthttp.RequestCtx, err error, status int) {
	ctx.SetStatusCode(status)
	if s.DB.opts.Debug {
		ctx.Response.Header.Add("X-Error", err.Error())
		// ctx.Response.Header.AddBytesV("X-Error-Trace", []byte(fmt.Sprintf("%+v", err)))
	}
}

func (s *HTTPServer) NotFound(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetContentTypeBytes(strApplicationJSON)
	ctx.SetStatusCode(fasthttp.StatusNotFound)
}

func (s *HTTPServer) Run() error {
	router := fasthttprouter.New()
	router.GET("/accounts/filter/", s.Filter)
	router.GET("/accounts/group/", s.Group)
	router.POST("/accounts/:token/", s.RouterAccounts) // Ох уж эти знатоки REST. /accounts/<id>/ и /accounts/new/ — конфликтующие пути

	/*router.GET("/accounts/group/", Index)
	router.GET("/accounts/:id/recommend/", Index)
	router.GET("/accounts/:id/suggest/", Index)*/
	router.NotFound = s.NotFound
	log.Info().Msg("HTTP is up and ready to serve requests!")
	return fasthttp.ListenAndServe(s.Addr, router.Handler)
}

func (s *HTTPServer) Shutdown() error {
	// fasthttp не умеет :(
	return nil
}
