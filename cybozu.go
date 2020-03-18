package cybozu // import "github.com/synerex/meeting_cybozu"

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/sclevine/agouti"
)

var (
	url       = "https://onlinedemo.cybozu.info/scripts/office10/ag.cgi?"
	loginName = "高橋 健太"
)

func getPageDOM(page *agouti.Page) *goquery.Document {
	// get whole page
	wholePage, err := page.HTML()
	if err != nil {
		fmt.Println("Failed to get whole page:", err)
	}
	// use goquery
	readerOfPage := strings.NewReader(wholePage)
	pageDom, err := goquery.NewDocumentFromReader(readerOfPage)
	if err != nil {
		fmt.Println("Failed to get page dom:", err)
	}
	return pageDom
}

func searchIndex(dates []string, target string) (int, error) {
	index := -1
	for i, date := range dates {
		if date == target {
			index = i
		}
		// fmt.Println(i, date)
	}
	if index == -1 {
		errMsg := "Failed to set parameter: " + target
		return -1, errors.New(errMsg)
	} else {
		return index, nil
	}
}

func login(page *agouti.Page, user string) error {
	// get user list
	usersDom := getPageDOM(page).Find("select[name='_ID']").Children()
	users := make([]string, usersDom.Length())
	usersDom.Each(func(i int, sel *goquery.Selection) {
		users[i] = sel.Text()
		// fmt.Println(i, users[i])
	})
	// search index
	userIndex, err := searchIndex(users, user)
	if err != nil {
		return err
	}
	// set login user
	name := page.FindByName("_ID")
	if _, err := name.Count(); err != nil {
		return err
	}
	name.Select(users[userIndex])
	// click login button
	submitBtn := page.FindByName("Submit")
	if _, err := submitBtn.Count(); err != nil {
		return err
	}
	// click
	if err := submitBtn.Click(); err != nil {
		return err
	}
	fmt.Println("Login complete:", users[userIndex])
	return nil
}

func booking(page *agouti.Page, date string, start string, end string, title string, room string) error {
	reserveButton := page.FindByXPath("//*[@id=\"content-wrapper\"]/div[4]/div/div[1]/table/tbody/tr/td[1]/table/tbody/tr/td[1]/span/span/a")
	_, err := reserveButton.Count()
	if err != nil {
		return err
	}
	reserveButton.Click()

	// set the date
	yearDom := getPageDOM(page).Find("select[name='SetDate.Year']").Children()
	monthDom := getPageDOM(page).Find("select[name='SetDate.Month']").Children()
	dayDom := getPageDOM(page).Find("select[name='SetDate.Day']").Children()
	startHourDom := getPageDOM(page).Find("select[name='SetTime.Hour']").Children()
	startMinuteDom := getPageDOM(page).Find("select[name='SetTime.Minute']").Children()
	endHourDom := getPageDOM(page).Find("select[name='EndTime.Hour']").Children()
	endMinuteDom := getPageDOM(page).Find("select[name='EndTime.Minute']").Children()

	years := make([]string, yearDom.Length())
	months := make([]string, monthDom.Length())
	days := make([]string, dayDom.Length())
	startHours := make([]string, startHourDom.Length())
	startMinutes := make([]string, startMinuteDom.Length())
	endHours := make([]string, endHourDom.Length())
	endMinutes := make([]string, endMinuteDom.Length())

	yearDom.Each(func(i int, g *goquery.Selection) {
		tx := g.Text()
		years[i] = tx
	})
	monthDom.Each(func(i int, g *goquery.Selection) {
		tx := g.Text()
		months[i] = tx
	})
	dayDom.Each(func(i int, g *goquery.Selection) {
		tx := g.Text()
		days[i] = tx
	})
	startHourDom.Each(func(i int, g *goquery.Selection) {
		tx := g.Text()
		startHours[i] = tx
	})
	startMinuteDom.Each(func(i int, g *goquery.Selection) {
		tx := g.Text()
		startMinutes[i] = tx
	})
	endHourDom.Each(func(i int, g *goquery.Selection) {
		tx := g.Text()
		endHours[i] = tx
	})
	endMinuteDom.Each(func(i int, g *goquery.Selection) {
		tx := g.Text()
		endMinutes[i] = tx
	})

	dateSplit := strings.Split(date, "/")
	yearIndex, err := searchIndex(years, dateSplit[0])
	if err != nil {
		return err
	}
	monthIndex, err := searchIndex(months, dateSplit[1])
	if err != nil {
		return err
	}
	dayIndex, err := searchIndex(days, dateSplit[2])
	if err != nil {
		return err
	}

	startSplit := strings.Split(start, ":")
	endSplit := strings.Split(end, ":")
	startHourIndex, err := searchIndex(startHours, startSplit[0]+"時")
	if err != nil {
		return err
	}
	startMinuteIndex, err := searchIndex(startMinutes, startSplit[1]+"分")
	if err != nil {
		return err
	}
	endHourIndex, err := searchIndex(endHours, endSplit[0]+"時")
	if err != nil {
		return err
	}
	endMinuteIndex, err := searchIndex(endMinutes, endSplit[1]+"分")
	if err != nil {
		return err
	}

	yearX := page.FindByName("SetDate.Year")
	_, err = yearX.Count()
	if err != nil {
		return err
	}
	monthX := page.FindByName("SetDate.Month")
	_, err = monthX.Count()
	if err != nil {
		return err
	}
	dayX := page.FindByName("SetDate.Day")
	_, err = dayX.Count()
	if err != nil {
		return err
	}
	startHourX := page.FindByName("SetTime.Hour")
	_, err = startHourX.Count()
	if err != nil {
		return err
	}
	startMinuteX := page.FindByName("SetTime.Minute")
	_, err = startMinuteX.Count()
	if err != nil {
		return err
	}
	endHourX := page.FindByName("EndTime.Hour")
	_, err = endHourX.Count()
	if err != nil {
		return err
	}
	endMinuteX := page.FindByName("EndTime.Minute")
	_, err = endMinuteX.Count()
	if err != nil {
		return err
	}

	err = yearX.Select(years[yearIndex])
	if err != nil {
		return err
	}
	err = monthX.Select(months[monthIndex])
	if err != nil {
		return err
	}
	err = dayX.Select(days[dayIndex])
	if err != nil {
		return err
	}
	err = startHourX.Select(startHours[startHourIndex])
	if err != nil {
		return err
	}
	err = startMinuteX.Select(startMinutes[startMinuteIndex])
	if err != nil {
		return err
	}
	err = endHourX.Select(endHours[endHourIndex])
	if err != nil {
		return err
	}
	err = endMinuteX.Select(endMinutes[endMinuteIndex])
	if err != nil {
		return err
	}

	// set the title
	detail := page.FindByName("Detail")
	if _, err := detail.Count(); err != nil {
		return err
	}
	detail.Fill(title)

	// choose room
	xpath := ""
	switch room {
	case "第一会議室":
		xpath = "//*[@id=\"content-wrapper\"]/div[4]/div/form/div[2]/table/tbody/tr/td/table/tbody/tr[2]/td/div/div[1]/div/table/tbody/tr[7]/td/table/tbody/tr[1]/td[3]/select/option[1]"
	case "第二会議室":
		xpath = "//*[@id=\"content-wrapper\"]/div[4]/div/form/div[2]/table/tbody/tr/td/table/tbody/tr[2]/td/div/div[1]/div/table/tbody/tr[7]/td/table/tbody/tr[1]/td[3]/select/option[2]"
	case "打合せルーム":
		xpath = "//*[@id=\"content-wrapper\"]/div[4]/div/form/div[2]/table/tbody/tr/td/table/tbody/tr[2]/td/div/div[1]/div/table/tbody/tr[7]/td/table/tbody/tr[1]/td[3]/select/option[3]"
	}
	theRoomY := page.FindByXPath(xpath)
	theRoomY.Click()

	time.Sleep(2 * time.Second)

	// submit to make a reservation
	entryButton := page.FindByName("Entry")
	_, err = entryButton.Count()
	if err != nil {
		return err
	}
	entryButton.Click()
	fmt.Println("Booking complete:", years[yearIndex], months[monthIndex], days[dayIndex], startHours[startHourIndex], startMinutes[startMinuteIndex], endHours[endHourIndex], endMinutes[endMinuteIndex])

	return nil
}

func Execute(year string, month string, day string, week string, start string, end string, people string, title string, room string) error {
	// set of Chrome
	driver := agouti.ChromeDriver(agouti.Browser("chrome"))
	if err := driver.Start(); err != nil {
		return err
	}
	defer driver.Stop()

	page, err := driver.NewPage()
	if err != nil {
		return err
	}

	// sample Cybozu
	if err := page.Navigate(url); err != nil {
		return err
	}

	// login
	if err := login(page, loginName); err != nil {
		return err
	}

	// get group list
	groupsDom := getPageDOM(page).Find("select[name='GID']").Children()
	groups := make([]string, groupsDom.Length())
	groupsDom.Each(func(i int, sel *goquery.Selection) {
		groups[i] = sel.Text()
		// fmt.Println(i, groups[i])
	})

	// move to meeting room page
	group := page.FindByName("GID")
	if _, err := group.Count(); err != nil {
		fmt.Println("Cannot find path:", err)
	}
	group.Select(groups[9]) // "(全施設)"

	// make a reservation
	// date := "2019年/4月/23(火)"
	// start := "10:00"
	// end := "11:30"
	date := year + "年/" + month + "月/" + day + week
	if err := booking(page, date, start, end, title, room); err != nil {
		return err
	}

	time.Sleep(3 * time.Second)
	return nil
}

func Schedules(year string, month string, day string, start string, end string, people string) (map[string][]string, error) {

	driver := agouti.ChromeDriver(agouti.Browser("chrome"))
	// driver := agouti.ChromeDriver(
	// 	agouti.ChromeOptions("args", []string{
	// 		"--headless",
	// 		"--disable-gpu",
	// 	}),
	// 	agouti.Debug,
	// )

	if err := driver.Start(); err != nil {
		return nil, err
	}
	defer driver.Stop()

	page, err := driver.NewPage()
	if err != nil {
		return nil, err
	}

	// sample Cybozu
	if err := page.Navigate(url); err != nil {
		return nil, err
	}

	// login
	if err := login(page, loginName); err != nil {
		return nil, err
	}

	// get group list
	groupsDom := getPageDOM(page).Find("select[name='GID']").Children()
	groups := make([]string, groupsDom.Length())
	groupsDom.Each(func(i int, sel *goquery.Selection) {
		groups[i] = sel.Text()
		// fmt.Println(i, groups[i])
	})

	// move to meeting room page
	group := page.FindByName("GID")
	if _, err := group.Count(); err != nil {
		return nil, err
	}
	group.Select(groups[10]) // "会議室"

	// get schedules
	schedulesDom := getPageDOM(page).Find("#redraw > table > tbody").Children()
	rooms := make(map[string][]string, schedulesDom.Length())
	schedulesDom.Each(func(i int, sel *goquery.Selection) {
		if i == 0 {
			sel.Children().Each(func(j int, cc *goquery.Selection) {
				if j == 0 {
					rooms["dates"] = []string{}
				} else {
					st := strings.TrimSpace(cc.Text())
					rooms["dates"] = append(rooms["dates"], st)
				}
			})
		} else {
			roomName := "none"
			sel.Children().Each(func(j int, cc *goquery.Selection) {
				if j == 0 {
					roomName = strings.Trim(cc.Children().First().First().Text(), " \n")
					roomName = strings.TrimSpace(roomName)
					rooms[roomName] = []string{}
				} else {
					st := strings.Trim(cc.Text(), "\n")
					st = strings.TrimSpace(st)
					rooms[roomName] = append(rooms[roomName], st)
				}
			})
		}
	})

	return rooms, nil
}
