package server

import (
	"log"
	"net/http"

	"github.com/Prettyletto/post-dude/cmd/internal/db"
	"github.com/Prettyletto/post-dude/cmd/internal/handler"
	"github.com/Prettyletto/post-dude/cmd/internal/repository"
	"github.com/Prettyletto/post-dude/cmd/internal/service"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(addr string) *Server {
	dataBase, err := db.New()
	if err != nil {
		log.Fatal(err)
	}
	dataBase.Init()

	collectionRepo := repository.NewCollectionRepository(dataBase)
	collectionService := service.NewCollectionService(collectionRepo)
	collectionHandler := handler.NewCollectionHandler(collectionService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /collections", collectionHandler.CreateCollectionHandler)
	mux.HandleFunc("GET /collections", collectionHandler.GetAllCollectionsHandler)
	mux.HandleFunc("GET /collections/{id}", collectionHandler.GetCollectionHandler)
	mux.HandleFunc("PUT /collections/{id}", collectionHandler.UpdateCollectionHandler)
	mux.HandleFunc("DELETE /collections/{id}", collectionHandler.DeleteCollectionHandler)

	s := &Server{httpServer: &http.Server{
		Addr:    addr,
		Handler: mux,
	}}
	return s
}

func (s *Server) Start() {
	go func() {
		s.httpServer.ListenAndServe()
	}()
}

func (s *Server) Stop() {
	s.httpServer.Close()
}
