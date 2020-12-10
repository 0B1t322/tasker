package taskerserver

import (
	"github.com/0B1t322/statements"
	pb "github.com/0B1t322/tasker/tasker"
	"net"
	"google.golang.org/grpc"
	
	"context"
	"database/sql"
	"errors"

	"time"

	log "github.com/sirupsen/logrus"
)

// Errors
var (
	ErrTaskExsist 		= 		errors.New("Task exsist")
	ErrTaskNotExsist	=		errors.New("Task is not exsist")
	ErrUknowUser 		= 		errors.New("Uknown user")
	ErrInternalServer 	=		errors.New("Server internal error")
)

// TaskerServer .....
type TaskerServer struct {
	pb.UnimplementedTaskerServer
}



type Config struct {
	Port string `json:"port"`
	Protocol string `"json:"protocol"`
}

// NewTaskerServer listen this microservice
func NewTaskerServer(conf Config, opts []grpc.ServerOption) error {
	listener, err := net.Listen(conf.Protocol, conf.Port)
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer(opts...)

	pb.RegisterTaskerServer(grpcServer, &TaskerServer{})
	if err := grpcServer.Serve(listener); err != nil {
		return err
	}

	return nil
}

/*
	ctx should be withValue "db" *sql.DB
*/
func getTaskByName(ctx context.Context, name string) (*pb.Task, error) {
	task := &pb.Task{}

	stmt, err := statements.NewGetAllStmt(ctx, "tasks" ,"name")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(name)

	err = row.Scan(&task.ID, &task.UID, &task.Name, &task.Description, &task.CreatesTime, &task.Done)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func getTaskByid(ctx context.Context, ID string) (*pb.Task, error) {
	stmt, err := statements.NewGetAllStmt(ctx, "tasks", "id")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	task := &pb.Task{}
	row := stmt.QueryRow(ID)
	err = row.Scan(&task.ID, &task.UID, &task.Name, &task.Description, &task.CreatesTime, &task.Done)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func getTasksByUID(ctx context.Context, UID int) ([]*pb.Task, error) {
	stmt, err := statements.NewGetAllStmt(ctx, "tasks", "user_id")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	tasks := []*pb.Task{}

	row, err := stmt.Query(UID)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	for row.Next() {
		task := &pb.Task{}

		err := row.Scan(
			&task.ID, 
			&task.UID, 
			&task.Name, 
			&task.Description, 
			&task.CreatesTime, 
			&task.Done,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil;
}

/*
getUserIDByToken:
	description: return a id by given token
	return: -1 id of some err and can return sql.ErrNoRows
*/
func getUserIDByToken(ctx context.Context, token string) (int, error) {
	stmt, err := statements.NewGetStmt(ctx,"id","users","token")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	var UID int
	row := stmt.QueryRow(token)

	if err := row.Scan(&UID); err != nil {
		log.Error(err)
		return -1, err
	}

	return UID, nil
}

// CreateTask ...
func (ts *TaskerServer) CreateTask(
	ctx context.Context, 
	req *pb.TaskRequest,
	) (*pb.TaskResponse, error) {
	
	db, err := ConnectToDB()
	if err != nil {
		return &pb.TaskResponse{Error: ErrInternalServer.Error()}, err
	}
	defer db.Close()
	// провреяем токен
	var UID int
	if id, err := getUserIDByToken(context.WithValue(ctx,"db", db), req.Token); err == sql.ErrNoRows {
		return &pb.TaskResponse{Error: ErrUknowUser.Error()}, nil
	} else if err != nil {
		return &pb.TaskResponse{Error: ErrInternalServer.Error()}, err
	} else {
		UID = id
	}

	// проверим есть ли уже задание с таким именем либо уже прошел 1 день
	task, err := getTaskByName( context.WithValue(ctx,"db", db), req.Name)
	if err != sql.ErrNoRows{
		createDate, errTime := time.Parse(time.Stamp, task.CreatesTime)
		if errTime != nil {
			return &pb.TaskResponse{Error: ErrInternalServer.Error()}, err
		}

		if after :=	time.Now().AddDate(0,0,-1).Add(-time.Minute); !createDate.After(after) {
				return &pb.TaskResponse{Error: "Task exsist"}, nil
		}
	} else if err != nil && err != sql.ErrNoRows {
		return &pb.TaskResponse{Error: ErrInternalServer.Error()}, err
	}
	
	if _, err := db.Exec("insert into tasks (user_id, name, description, create_time, done) values ($1, $2, $3, $4, $5)", 
			UID, req.Name, req.Description, req.CreatesTime, false,
	); err != nil {
		return &pb.TaskResponse{Error: ErrInternalServer.Error()}, err
	}

	return &pb.TaskResponse{Error: ""}, nil
}


// MarkTask ......
func (ts *TaskerServer) MarkTask(
		ctx context.Context, 
		req *pb.MarkRequest,
	) (*pb.TaskResponse, error) {
		db, err := ConnectToDB()
		if err != nil {
			return &pb.TaskResponse{Error: ErrInternalServer.Error()}, err
		}
		defer db.Close()

		dbCtx := context.WithValue(ctx, "db", db)
		

		// проверка токена
		if _, err := getUserIDByToken(dbCtx, req.Token); err == ErrUknowUser {
			return &pb.TaskResponse{Error: err.Error()}, nil
		} else if err != nil {
			return &pb.TaskResponse{Error: ErrInternalServer.Error()}, err
		}
		
		res, err := db.Exec("update tasks set done = $1 where id = $2", req.Done, req.ID)
		rowsAffected, _ := res.RowsAffected()
		if err != nil {
			return &pb.TaskResponse{Error: ErrInternalServer.Error()}, err
		}

		if rowsAffected == 0 {
			return &pb.TaskResponse{Error: ErrTaskNotExsist.Error()}, nil
		}

		return &pb.TaskResponse{Error: ""}, nil
}

// GetAllTasks ....
func (ts *TaskerServer) GetAllTasks(
	ctx context.Context, 
	req *pb.GetAllTaskRequest,
	) (*pb.GetTaskResponse, error) {
		
		db, err := ConnectToDB()
		if err != nil {
			return &pb.GetTaskResponse{
				Task: nil, 
				Error: ErrInternalServer.Error(),
				}, err 
		}
		defer db.Close()

		dbCtx := context.WithValue(ctx, "db", db)
		uid, err := getUserIDByToken(dbCtx, req.Token)
		switch err {
			case ErrUknowUser:
				return &pb.GetTaskResponse{
					Task: nil,
					Error: ErrUknowUser.Error(),
				}, nil

			case nil:

			default:
				return &pb.GetTaskResponse{
					Task: nil,
					Error: ErrInternalServer.Error(),
				}, err
		}

		tasks, err := getTasksByUID(dbCtx, uid)
		if err != nil {
			return &pb.GetTaskResponse{
				Task: nil,
				Error: ErrInternalServer.Error(),
			}, err
		}

		return &pb.GetTaskResponse{
			Task: tasks,
			Error: "",
		}, nil
		
}

// GetTask ........
func (ts TaskerServer) GetTask(
	ctx context.Context, 
	req *pb.GetTaskRequest,
	) (*pb.GetTaskResponse, error) {
		
		allTasksRes, err := ts.GetAllTasks(ctx, &pb.GetAllTaskRequest{Token: req.Token})
		if err != nil {
			return &pb.GetTaskResponse{Error: err.Error(), Task: nil}, err
		}

		if allTasksRes.Error != "" {
			return &pb.GetTaskResponse{Error: allTasksRes.Error}, nil
		}

		tasks := NewFilter(req.Filer).FilterTasks(allTasksRes.Task)

		return &pb.GetTaskResponse{Task: tasks, Error: ""}, nil
}

// ArchiveTask ....
func (ts TaskerServer) ArchiveTask(
	ctx context.Context, 
	req *pb.ArchiveRequest,
	) (*pb.TaskResponse, error) {
		db, err := ConnectToDB()
		if err  != nil {
			return &pb.TaskResponse{Error: ErrInternalServer.Error()}, err
		}
		defer db.Close()

		dbCtx := context.WithValue(ctx, "db", db)
		// проверим токен
		_, err = getUserIDByToken(dbCtx, req.Token)
		if err == sql.ErrNoRows {
			return &pb.TaskResponse{Error: ErrUknowUser.Error()}, nil
		} else if err != nil {
			return &pb.TaskResponse{Error: ErrInternalServer.Error()}, err
		}
		
		// мы должны получить таск по ID удалить его

		task, err := getTaskByid(dbCtx, req.ID)
		if err == sql.ErrNoRows {
			return &pb.TaskResponse{Error: ErrUknowUser.Error()}, nil
		} else if err != nil {
			return &pb.TaskResponse{Error: ErrInternalServer.Error()}, err
		}
		
		stmtDel, err := statements.NewDeleteStmt(dbCtx, "tasks", "id")
		if err != nil {
			return &pb.TaskResponse{Error: ErrInternalServer.Error()}, err
		}

		if _, err := stmtDel.Exec(req.ID); err != nil {
			return &pb.TaskResponse{Error: ErrInternalServer.Error()}, err
		}

		// перемещаем его
		stmtIns, err := statements.NewInsertStmt(
			dbCtx, // ctx
			"archive_tasks", // table 
			"user_id", "name", "description", "create_time", "done", // fields
		)
		if err != nil {
			return &pb.TaskResponse{Error: ErrInternalServer.Error()}, err
		}

		if _, err := stmtIns.Exec(
			task.UID,
			task.Name,
			task.Description,
			task.CreatesTime,
			task.Done,
		); err != nil {
			return &pb.TaskResponse{Error: ErrInternalServer.Error()}, err
		}

		return &pb.TaskResponse{Error: ""}, nil
}