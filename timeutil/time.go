package timeutil

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"time"
)

func s(x float64) string {
	if int(x) == 1 {
		return ""
	}
	return "s"
}

func TimeElapsed(now time.Time, then time.Time, full bool) string {
	var parts []string
	var text string

	year2, month2, day2 := now.Date()
	hour2, minute2, second2 := now.Clock()

	year1, month1, day1 := then.Date()
	hour1, minute1, second1 := then.Clock()

	year := math.Abs(float64(int(year2 - year1)))
	month := math.Abs(float64(int(month2 - month1)))
	day := math.Abs(float64(int(day2 - day1)))
	hour := math.Abs(float64(int(hour2 - hour1)))
	minute := math.Abs(float64(int(minute2 - minute1)))
	second := math.Abs(float64(int(second2 - second1)))

	week := math.Floor(day / 7)

	if year > 0 {
		parts = append(parts, strconv.Itoa(int(year))+" year"+s(year))
	}

	if month > 0 {
		parts = append(parts, strconv.Itoa(int(month))+" month"+s(month))
	}

	if week > 0 {
		parts = append(parts, strconv.Itoa(int(week))+" week"+s(week))
	}

	if day > 0 {
		parts = append(parts, strconv.Itoa(int(day))+" day"+s(day))
	}

	if hour > 0 {
		parts = append(parts, strconv.Itoa(int(hour))+" hour"+s(hour))
	}

	if minute > 0 {
		parts = append(parts, strconv.Itoa(int(minute))+" minute"+s(minute))
	}

	if second > 0 {
		parts = append(parts, strconv.Itoa(int(second))+" second"+s(second))
	}

	if now.After(then) {
		text = " ago"
	} else {
		text = " after"
	}

	if len(parts) == 0 {
		return "just now"
	}

	if full {
		return strings.Join(parts, ", ") + text
	}
	return parts[0] + text
}

// ParseTime convert date string to time.Time
func ParseTime(s interface{}, layouts ...string) (t time.Time, err error) {
	var layout string
	str := ""
	if len(layouts) > 0 { // custom layout
		layout = layouts[0]
	} else {
		switch s := s.(type) {
		case string:
			str = s
			switch len(s) {
			case 8:
				layout = "20060102"
			case 10:
				layout = "2006-01-02"
			case 13:
				layout = "2006-01-02 15"
			case 16:
				layout = "2006-01-02 15:04"
			case 19:
				layout = "2006-01-02 15:04:05"
			case 20: // time.RFC3339
				layout = "2006-01-02T15:04:05Z07:00"
			}
			break
		case int:
			return time.Unix(int64(s), 0), nil
		case int64:
			return time.Unix(s, 0), nil
		}
	}
	if layout == "" {
		err = errors.New("invalid params")
		return
	}

	// has 'T' eg: "2006-01-02T15:04:05"
	if strings.ContainsRune(str, 'T') {
		layout = strings.Replace(layout, " ", "T", -1)
	}

	// eg: "2006/01/02 15:04:05"
	if strings.ContainsRune(str, '/') {
		layout = strings.Replace(layout, "-", "/", -1)
	}

	t, err = time.Parse(layout, str)
	// t, err = time.ParseInLocation(layout, s, time.Local)
	return
}
