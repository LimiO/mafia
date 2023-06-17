package users

import (
	"fmt"
	"github.com/signintech/gopdf"
	"mafia/internal/db"
)

const (
	pdfTemplate = `<html>
<body>
	<h1 style="color:black;">Users
	%s
</body>
</html>`
	statTemplate = `<div>
	<span>id: %s</span><br>
	<span>name: %s</span><br>
	<span>email: %s</span><br>
	<span>sex: %s</span><br>
	%s
	<span>count_games: %d</span><br>
	<span>count_wins: %d</span><br>
	<span>count_loses: %d</span><br>
	<span>total_time: %d</span><br>
	<hr>
</div>
`
)

type Data struct {
	// 0 - wait
	// 1 - good
	// 2 - exception
	status int
}

type UserData struct {
	stats *db.Stats
	user  *db.User
}

type Request struct {
	Name string
	Data map[string]*UserData
}

type PDFController struct {
	queue chan *Request
	data  map[string]*Data

	dbController *db.Manager
}

func (c *PDFController) Run() {
	for request := range c.queue {
		c.data[request.Name] = &Data{
			status: 0,
		}
		request := request
		go func() {
			err := c.genPdf(request.Name, request.Data)
			if err != nil {
				c.data[request.Name].status = 2
				fmt.Printf("failed to gen pdf: %v", err)
			}
		}()
	}
}

func mergeUserData(stats []*db.Stats, users []*db.User) map[string]*UserData {
	result := map[string]*UserData{}
	fmt.Println(users, 12312)
	for _, user := range users {
		result[user.ID] = &UserData{
			user: user,
		}
	}
	for _, stat := range stats {
		result[stat.ID].stats = stat
	}
	for userID := range result {
		if result[userID].stats == nil {
			delete(result, userID)
		}
	}

	return result
}

func (c *PDFController) genPdf(request string, usersData map[string]*UserData) error {
	pdfg := gopdf.GoPdf{}
	pdfg.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	err := pdfg.AddTTFFont("ttf", "content/font.ttf")
	if err != nil {
		return fmt.Errorf("failed to add ttf: %v", err)
	}
	err = pdfg.SetFont("ttf", "", 20)
	if err != nil {
		return fmt.Errorf("failed to set font: %v", err)
	}
	//var statsTemplates []string
	for _, data := range usersData {
		pdfg.AddPage()
		pdfg.SetXY(10, 10)
		_ = pdfg.Cell(nil, fmt.Sprintf("id: %s\n", data.user.ID))

		pdfg.SetXY(10, 30)
		_ = pdfg.Cell(nil, fmt.Sprintf("name: %s\n", data.user.Name))

		pdfg.SetXY(10, 50)
		_ = pdfg.Cell(nil, fmt.Sprintf("email: %s\n", data.user.Email))

		pdfg.SetXY(10, 70)
		_ = pdfg.Cell(nil, fmt.Sprintf("sex: %s\n", data.user.Sex))

		pdfg.SetXY(10, 90)
		_ = pdfg.Cell(nil, fmt.Sprintf("count games: %d\n", data.stats.CountGames))

		pdfg.SetXY(10, 110)
		_ = pdfg.Cell(nil, fmt.Sprintf("count wins: %d\n", data.stats.CountWins))

		pdfg.SetXY(10, 130)
		_ = pdfg.Cell(nil, fmt.Sprintf("count loses: %d\n", data.stats.CountGames-data.stats.CountWins))

		pdfg.SetXY(10, 150)
		_ = pdfg.Cell(nil, fmt.Sprintf("total time: %d\n", data.stats.TotalTime))
		if data.user.Image != "" {
			_ = pdfg.Cell(nil, "image: \n")
			_ = pdfg.Image(fmt.Sprintf("content/img/%s", data.user.Image), 10, 170, nil)
		}
	}

	err = pdfg.WritePdf(fmt.Sprintf("content/pdf/%s.pdf", request))
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	c.data[request].status = 1
	return nil
}
