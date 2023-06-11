package users

import (
	"fmt"
	"mafia/internal/db"
	"strings"

	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
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

	for _, user := range users {
		result[user.ID] = &UserData{
			user: user,
		}
	}
	for _, stat := range stats {
		result[stat.ID].stats = stat
	}
	for _, user := range users {
		if result[user.ID].stats == nil {
			delete(result, user.ID)
		}
	}

	return result
}

func (c *PDFController) genPdf(request string, usersData map[string]*UserData) error {
	pdfg, err := wkhtml.NewPDFGenerator()
	if err != nil {
		return fmt.Errorf("failed to make new pdf gen: %v", err)
	}

	var statsTemplates []string
	for _, data := range usersData {
		img := ""
		if data.user.Image != "" {
			img = fmt.Sprintf(
				"<span>image: </span><br><img src=\"/infoserver/content/img/%s\" alt=\"bad pic\", height=\"250\"><br>",
				data.user.Image,
			)
		}
		statsTemplates = append(statsTemplates, fmt.Sprintf(
			statTemplate,
			data.user.ID, data.user.Name, data.user.Email, data.user.Sex, img,
			data.stats.CountGames, data.stats.CountWins, data.stats.CountGames-data.stats.CountWins, data.stats.TotalTime,
		))
	}
	pdfData := fmt.Sprintf(pdfTemplate, strings.Join(statsTemplates, "\n\n"))
	pageReader := wkhtml.NewPageReader(strings.NewReader(pdfData))
	pageReader.EnableLocalFileAccess.Set(true)
	pdfg.AddPage(pageReader)

	err = pdfg.Create()
	if err != nil {
		return fmt.Errorf("failed to create pdf: %v", err)
	}

	err = pdfg.WriteFile(fmt.Sprintf("content/pdf/%s.pdf", request))
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	c.data[request].status = 1
	return nil
}
