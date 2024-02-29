package delivery

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/JuFnd/go-proxy/configs"
	proxy "github.com/JuFnd/go-proxy/internal/app/proxy/server"
	"github.com/JuFnd/go-proxy/internal/app/server/usecase"
	scanner "github.com/JuFnd/go-proxy/pkg"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type API struct {
	requestUseCase usecase.IUseCase
	proxyServer    *proxy.ProxyServer
	cfg            *configs.HTTPSrvConfig
	lg             *logrus.Logger
	mx             *mux.Router
}

func GetApi(requestUseCase usecase.IUseCase, proxyServer *proxy.ProxyServer, cfg *configs.HTTPSrvConfig, lg *logrus.Logger) *API {
	api := &API{
		requestUseCase: requestUseCase,
		proxyServer:    proxyServer,
		cfg:            cfg,
		lg:             lg,
		mx:             mux.NewRouter(),
	}

	api.mx.HandleFunc("/requests", api.GetRequests)
	api.mx.HandleFunc("/requests/{id:[0-9]+}", api.GetRequest)
	api.mx.HandleFunc("/scan/{id:[0-9]+}", api.ScanRequest)
	api.mx.HandleFunc("/repeat/{id:[0-9]+}", api.RepeatRequest)

	return api
}

func (a *API) ListenAndServe() error {
	a.lg.Infof("start application-server listening at: " + a.cfg.WebPort)

	err := http.ListenAndServe(":"+a.cfg.WebPort, a.mx)
	if err != nil {
		a.lg.Error("Listen and serve error: ", err.Error())
		return err
	}

	return nil
}

func (a *API) ProxyRequest(w http.ResponseWriter, r *http.Request) {
	a.proxyHttpOrHttps(w, r, true)
}

func (a *API) proxyHttpOrHttps(w http.ResponseWriter, r *http.Request, save bool) {
	var err error

	if r.Method == http.MethodConnect {
		a.lg.Infof("HTTPS....")
		a.proxyServer.ProxyHTTPS(w, r)
		return
	}

	a.proxyServer.ProxyHTTP(w, r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *API) GetRequest(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 32)
	selectedRequest, err := a.requestUseCase.GetRequestDataById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	answer, err := json.Marshal(selectedRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(answer)
}

func (a *API) GetRequests(w http.ResponseWriter, r *http.Request) {
	selectedRequests, err := a.requestUseCase.GetAllRequestsData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	answer, err := json.Marshal(selectedRequests)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(answer)
}

func (a *API) RepeatRequest(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 32)
	selectedRequest, err := a.requestUseCase.GetRequestById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	a.proxyHttpOrHttps(w, &http.Request{
		Method: selectedRequest.Method,
		URL: &url.URL{
			Scheme: selectedRequest.Scheme,
			Host:   selectedRequest.Host,
			Path:   selectedRequest.Path,
		},
		Header: selectedRequest.Headers,
		Body:   ioutil.NopCloser(strings.NewReader(selectedRequest.Body)),
		Host:   r.Host,
	}, false)
}

func (a *API) ScanRequest(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 32)
    selectedRequest, err := a.requestUseCase.GetRequestById(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    request := &http.Request{
        Method: selectedRequest.Method,
        URL: &url.URL{
            Scheme: selectedRequest.Scheme,
            Host:   selectedRequest.Host,
            Path:   selectedRequest.Path,
        },
        Header: selectedRequest.Headers,
        Body:   ioutil.NopCloser(strings.NewReader(selectedRequest.Body)),
        Host:   r.Host,
    }

	rootDir, _ := os.Getwd()
    dictFilePath := rootDir + "/pkg/dicc.txt"
    scanResults := scanner.Scan(request, dictFilePath)

    var foundPaths []string

    for path, isFound := range scanResults {
        if isFound {
            foundPaths = append(foundPaths, path)
        }
    }

    jsonData, err := json.Marshal(foundPaths)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")

    w.Write(jsonData)
}
