package modules

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

var (
	df []Row
)

type Row struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Filepath   string    `json:"filepath"`
	Evaluation int       `json:"evaluation"`
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedAt  time.Time `json:"created_at"`
}

type F struct {
	Name string
	Path string
}

func CreateCSV(filedir string) error {
	var datas []F

	if err := filepath.Walk(filedir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		datas = append(datas, F{Name: info.Name(), Path: path})
		return nil
	}); err != nil {
		return err
	}

	// List 作成
	f, err := os.Create("./data/list-before.csv")
	if err != nil {
		return err
	}
	defer f.Close()

	w := gocsv.DefaultCSVWriter(f)

	df = make([]Row, len(datas))
	w.Write([]string{"id", "name", "filepath", "evaluation", "updated_at", "created_at"})
	for i := range datas {
		df[i] = Row{
			ID:         i,
			Name:       datas[i].Name,
			Filepath:   datas[i].Path,
			Evaluation: 0,
			UpdatedAt:  time.Now(),
			CreatedAt:  time.Now(),
		}

		// Save RAW Data
		w.Write([]string{
			fmt.Sprintf("%d", i),
			datas[i].Name,
			datas[i].Path,
			"0",
			time.Now().Format("2006-01-02"),
			time.Now().Format("2006-01-02")})
		fmt.Printf("%s, %s\n", datas[i].Name, datas[i].Path)
	}
	w.Flush()

	return nil
}

type D struct {
	ID         int `query:"id"`
	Evaluation int `query:"evaluation"`
}

func Receive(c echo.Context) error {
	e := new(D)

	e.ID, _ = strconv.Atoi(c.QueryParam("id"))
	e.Evaluation, _ = strconv.Atoi(c.QueryParam("evaluation"))

	log.Infof("received: id: %d, evaluation: %d", e.ID, e.Evaluation)

	// Update Evaluation
	df[e.ID].Evaluation = e.Evaluation
	df[e.ID].UpdatedAt = time.Now()
	log.Info("updated df")

	if err := update(); err != nil {
		return err
	}

	var d = new(Row)
	if len(df) > e.ID+1 {
		d = &df[e.ID+1]
	} else {
		return fmt.Errorf("not required length")
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return c.JSON(http.StatusOK, d)
}

// Update List
func update() error {
	f, err := os.Create("./data/list-updated.csv")
	if err != nil {
		return err
	}
	defer f.Close()

	if err := gocsv.Marshal(df, f); err != nil {
		return err
	}
	log.Info("updated csv")

	return nil
}
