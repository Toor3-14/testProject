package app

import (
	"os"
	"sync"
	"time"
	"net/http"
	"os/signal"
	"log"
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/Toor3-14/testProject/internal/repo"
)
const panicMSG = "Application was not configure"

var (
	once sync.Once
	instance *application
 )

type application struct {
	serverAddr string
	handler http.Handler

	errorLog *log.Logger // for fatal errors and warnings
	infoLog  *log.Logger

	userRepo *repo.UserRepo
	redis *redis.Client
}

func Configure(saddr string, handler http.Handler, errlog, infolog *log.Logger, repo *repo.UserRepo, client *redis.Client) {
   once.Do(func() {
      instance = &application{
		 serverAddr: saddr,
		 handler: handler,
         errorLog: errlog,
         infoLog: infolog,
         userRepo: repo,
         redis: client,
      }
   })
}
func ServAddr() string {
	if instance == nil { panic(panicMSG) }
	return instance.serverAddr
}
func Handler() http.Handler {
	if instance == nil { panic(panicMSG) }
	return instance.handler
}
func ErrLog() *log.Logger {
   if instance == nil { panic(panicMSG) }
   return instance.errorLog
}
func InfoLog() *log.Logger {
   if instance == nil { panic(panicMSG) }
   return instance.infoLog
}
func UserRepo() *repo.UserRepo {
   if instance == nil { panic(panicMSG) }
   return instance.userRepo
}
func Redis() *redis.Client {
   if instance == nil { panic(panicMSG) }
   return instance.redis
}
func Start() {
	if instance == nil { panic(panicMSG) }
	srv := &http.Server{
		Addr: instance.serverAddr,
		ErrorLog: instance.errorLog,
		Handler: instance.handler,
	}
	
	go func() {
		instance.errorLog.Fatal(srv.ListenAndServe())
	}()

	// Wait for an interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	
	// Correct shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	instance.errorLog.Fatal(srv.Shutdown(ctx))
}