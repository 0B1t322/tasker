package tasker_test

import (
	pb "github.com/0B1t322/tasker/tasker"
	"github.com/0B1t322/tasker/tasker_server"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"


	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)


func ConnectToDB() (*sql.DB, error) {
	return taskerserver.ConnectToDB()
}

func clearUsers() {
	db, _ := ConnectToDB()

	_, err := db.Exec("delete from users")
	if err != nil {
		log.Error(err)
	}
	db.Close()
}

func clearTasks() {
	db, _ := ConnectToDB()

	_, err := db.Exec("delete from tasks")
	if err != nil {
		log.Error()
	}
	db.Close()
}

func TestFunc_SuccsecCreateUser(t *testing.T) {
	clearUsers()
	clearTasks()

	if db, err := ConnectToDB(); err != nil {
		t.Log(err)
		t.Fail()
	} else {
		_, err := db.Exec("insert into users (id, token) values ($1, $2)", 0, "admin")
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		db.Close()

	}
	req := &pb.TaskRequest{
		Name:        "скачать",
		Description: "с мыла игру",
		CreatesTime: time.Now().Format(time.Stamp),
		Token:       "admin",
	}

	expectedResponse := &pb.TaskResponse{Error: ""}

	s := &taskerserver.TaskerServer{}

	res, err := s.CreateTask(context.Background(), req)

	if res.Error != expectedResponse.Error {
		t.Log(err)
		t.Log("failed assert")
		t.Fail()
	}

	expectedResponse = &pb.TaskResponse{Error: taskerserver.ErrTaskExsist.Error()}
	res, err = s.CreateTask(context.Background(), req)

	if res.Error != expectedResponse.Error {
		t.Log(res.Error)
		t.Log("faield assert")
		t.Fail()
	}

	expectedResponse = &pb.TaskResponse{Error: taskerserver.ErrUknowUser.Error()}
	req.Token = "-"

	res, _ = s.CreateTask(context.Background(), req)

	if res.Error != expectedResponse.Error {
		t.Log("faield assert")
		t.Fail()
	}
}

func TestFunc_MarkTask(t *testing.T) {
	clearUsers()
	clearTasks()
	if db, err := ConnectToDB(); err != nil {
		t.Fatal(err)
		t.Fail()
	} else {
		_, err := db.Exec("insert into users (id, token) values ($1, $2)", 0, "admin")
		if err != nil {
			t.Fatal(err)
			t.Fail()
		}
		db.Close()

	}

	s := &taskerserver.TaskerServer{}

	s.CreateTask(context.Background(), &pb.TaskRequest{
		Name:        "скачать",
		Description: "с мыла игру",
		CreatesTime: time.Now().Format(time.Stamp),
		Token:       "admin",
	})

	req := &pb.MarkRequest{
		ID:    "1",
		Token: "admin",
		Done:  true,
	}

	res, err := s.MarkTask(context.Background(), req)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if res.Error != "" {
		t.Log("Faield assert")
		t.Fail()
	}

}

func TestFunc_GetAllTasks(t *testing.T) {
	clearTasks()
	clearUsers()
	if db, err := ConnectToDB(); err != nil {
		t.Fatal(err)
		t.Fail()
	} else {
		_, err := db.Exec("insert into users (id, token) values ($1, $2)", 0, "admin")
		if err != nil {
			t.Fatal(err)
			t.Fail()
		}
		db.Close()

	}

	s := &taskerserver.TaskerServer{}

	for i := 0; i < 10; i++ {
		s.CreateTask(context.Background(), &pb.TaskRequest{
			Name:        "скачать " + fmt.Sprint(i),
			Description: "с мыла игру ",
			CreatesTime: time.Now().Format(time.Stamp),
			Token:       "admin",
		})
	}

	req := &pb.GetAllTaskRequest{Token: "admin"}

	res, err := s.GetAllTasks(context.Background(), req)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if res.Error != "" {
		t.Log(res.Error)
		t.Log("error assert")
		t.Fail()
	}

	if len(res.Task) != 10 {
		t.Fail()
	}

	data, _ := json.Marshal(res.Task)

	t.Log(string(data))

}

func toBool(num int) bool {
	if num == 0 {
		return false
	} else {
		return true
	}
}

func TestFunc_Filer(t *testing.T) {
	tasks := []*pb.Task{}
	for i := 0; i < 10; i++ {
		tasks = append(tasks, &pb.Task{
			ID:          fmt.Sprint(i),
			UID:         "0",
			Name:        "task " + fmt.Sprint(i),
			Description: "some desc",
			CreatesTime: time.Now().AddDate(0, 0, i).Format(time.Stamp),
			Done:        toBool(i % 2),
		})
	}
	filter := taskerserver.NewFilter("all:")
	if filter.Type != taskerserver.ALL {
		t.Log(filter.Type)
		t.Fail()
	}
	if len(filter.FilterTasks(tasks)) != 10 {
		t.Log("faield assert")
		t.Fail()
	}

	filter = taskerserver.NewFilter("done: false")
	if filter.Type != taskerserver.DONE {
		t.Log("faield assert")
		t.Fail()
	}

	if len(filter.FilterTasks(tasks)) != 5 {
		t.Log("failed assert")
		t.Fail()
	}

	filter = taskerserver.NewFilter("done: true")
	if filter.Type != taskerserver.DONE {
		t.Log("faield assert")
		t.Fail()
	}

	if len(filter.FilterTasks(tasks)) != 5 {
		t.Log("failed assert")
		t.Fail()
	}

	filter = taskerserver.NewFilter("period: " + tasks[0].CreatesTime + " - " + tasks[7].CreatesTime)

	if filter.Type != taskerserver.PERIOD {
		t.Log("faield assert")
		t.Fail()
	}

	if len(filter.FilterTasks(tasks)) != 8 {
		t.Log("faield assert")
		t.Fail()
	}
}

func TestFunc_sss(t *testing.T) {
	f := func(strs ...string) string {
		var mass []interface{}
		for _, str := range strs {
			mass = append(mass, str)
		}

		return fmt.Sprintf("%s,", mass...)
	}

	t.Log(f("a", "b"))
}

func TestFunc_ArchiveTask(t *testing.T) {
	clearTasks()
	clearUsers()

	if db, err := ConnectToDB(); err != nil {
		t.Fatal(err)
		t.Fail()
	} else {
		_, err := db.Exec("insert into users (id, token) values ($1, $2)", 0, "admin")
		if err != nil {
			t.Fatal(err)
			t.Fail()
		}
		db.Close()

	}

	s := &taskerserver.TaskerServer{}

	for i := 0; i < 10; i++ {
		s.CreateTask(context.Background(), &pb.TaskRequest{
			Name:        "скачать " + fmt.Sprint(i),
			Description: "с мыла игру ",
			CreatesTime: time.Now().Format(time.Stamp),
			Token:       "admin",
		})
	}

	req := &pb.ArchiveRequest{ID: "1", Token: "admin"}

	res, err := s.ArchiveTask(context.Background(), req)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if res.Error != "" {
		t.Log(res.Error)
		t.Fail()
	}

}

func TestFunc_Format(t *testing.T) {
	formatFunc := func(n int) string {
		str := "("

		for i := 0; ; i++ {
			if i == n-1 {
				str += "%s)"
				break
			}
			str += "%s, "
		}
		// now str look like: (field, field, ...., field)

		str += " values ("

		for i := 0; ; i++ {
			if i == n-1 {
				str += "$" + fmt.Sprint(i+1) + ")"
				break
			}
			str += "$" + fmt.Sprint(i+1) + ", "
		}

		return str
	}

	t.Log(formatFunc(5))
}
