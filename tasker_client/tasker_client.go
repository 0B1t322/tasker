package taskerclient

import (
	pb "github.com/0B1t322/tasker/tasker"
	"context"
	"encoding/json"
	"errors"

	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	ErrNotFoundConn = NewError(errors.New("Not found conn"))
)


// Config need to create conn
type Config struct {
	Port string `json:"port"`
}

func MiddlewareTasker(conf Config) mux.MiddlewareFunc {
	return (func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			conn, err := grpc.Dial("127.0.0.1"+conf.Port, []grpc.DialOption{grpc.WithInsecure()}...)
			if err != nil {
				log.Error(err)
			}
	
			r = r.WithContext(
				context.WithValue(
					r.Context(), 
					"conn", 
					conn,
				),
			)
	
			next.ServeHTTP(w,r)
			
			defer conn.Close()
		})
	})
}

//CreateTask ....
func CreateTask(w http.ResponseWriter, r *http.Request) {
	log.Info("Create Task")
	conn, ok := r.Context().Value("conn").(*grpc.ClientConn)
	if !ok || conn == nil {
		log.Error("Not found conn")

		if data, err := ErrNotFoundConn.Marshall(); err != nil {
			log.Error(err)
		} else {
			w.Write(data)
		}

		return
	}

	client := pb.NewTaskerClient(conn)
	req := &pb.TaskRequest{}

	json.NewDecoder(r.Body).Decode(req)

	resp, err := client.CreateTask(context.Background(), req)
	if err != nil {
		log.Error(err)
	}

	log.Info("Error: ",resp.Error)
	log.Info("err: ", err)
	if data, err := json.Marshal(resp); err != nil {
		log.Error(err)
	} else {
		w.Write(data)
	}
}

// MarkTask ....
func MarkTask(w http.ResponseWriter, r *http.Request) {
	conn, ok := r.Context().Value("conn").(*grpc.ClientConn)
	if !ok || conn == nil {
		log.Error("Not found conn")

		if data, err := ErrNotFoundConn.Marshall(); err != nil {
			log.Error(err)
		} else {
			w.Write(data)
		}

		return
	}

	client := pb.NewTaskerClient(conn)
	req := &pb.MarkRequest{}

	json.NewDecoder(r.Body).Decode(req)

	resp, err := client.MarkTask(context.Background(), req)
	if err != nil {
		log.Error(err)
	}

	if data, err := json.Marshal(resp); err != nil {
		log.Error(err)
	} else {
		w.Write(data)
	}
}

// ArchiveTask ....
func ArchiveTask(w http.ResponseWriter, r *http.Request) {
	conn, ok := r.Context().Value("conn").(*grpc.ClientConn)
	if !ok || conn == nil {
		log.Error("Not found conn")

		if data, err := ErrNotFoundConn.Marshall(); err != nil {
			log.Error(err)
		} else {
			w.Write(data)
		}

		return
	}

	client := pb.NewTaskerClient(conn)
	req := &pb.ArchiveRequest{}

	json.NewDecoder(r.Body).Decode(req)

	resp, err := client.ArchiveTask(context.Background(), req)
	if err != nil {
		log.Error(err)
	}

	if data, err := json.Marshal(resp); err != nil {
		log.Error(err)
	} else {
		w.Write(data)
	}
}

// GetTask ...
func GetTask(w http.ResponseWriter, r *http.Request) {
	conn, ok := r.Context().Value("conn").(*grpc.ClientConn)
	if !ok || conn == nil {
		log.Error("Not found conn")

		if data, err := ErrNotFoundConn.Marshall(); err != nil {
			log.Error(err)
		} else {
			w.Write(data)
		}

		return
	}

	client := pb.NewTaskerClient(conn)
	req := &pb.GetTaskRequest{}

	json.NewDecoder(r.Body).Decode(req)

	resp, err := client.GetTask(context.Background(), req)
	if err != nil {
		log.Error(err)
	}

	if data, err := json.Marshal(resp); err != nil {
		log.Error(err)
	} else {
		w.Write(data)
	}
}

// GetAllTasks ....
func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	conn, ok := r.Context().Value("conn").(*grpc.ClientConn)
	if !ok || conn == nil {
		log.Error("Not found conn")

		if data, err := ErrNotFoundConn.Marshall(); err != nil {
			log.Error(err)
		} else {
			w.Write(data)
		}

		return
	}

	client := pb.NewTaskerClient(conn)
	req := &pb.GetAllTaskRequest{}

	json.NewDecoder(r.Body).Decode(req)

	resp, err := client.GetAllTasks(context.Background(), req)
	if err != nil {
		log.Error(err)
	}

	if data, err := json.Marshal(resp); err != nil {
		log.Error(err)
	} else {
		w.Write(data)
	}
}

// NewTaskerRouter return new tasker router
func NewTaskerRouter(conf Config) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/tasker/CreateTask", CreateTask)
	r.HandleFunc("/api/tasker/MarkTask", MarkTask).Methods("PATCH")
	r.HandleFunc("/api/tasker/ArchiveTask", ArchiveTask).Methods("PUT")
	r.HandleFunc("/api/tasker/GetTask", GetTask).Methods("GET")
	r.HandleFunc("/api/tasker/GetAllTasks", GetAllTasks).Methods("GET")
	r.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})
	
	r.Use(MiddlewareTasker(conf))

	return r
}