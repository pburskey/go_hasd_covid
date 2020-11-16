package utility

import (
	"time"
)

//
//func main() {
//	/*
//	covid_data_20200929130339.csv
//	 */
//	fmt.Print(parseYYYYMMDDHH24MiSS("20200929130339"))
//}

func ParseYYYYMMDDHH24MiSS(aDateAsString string) (aTime time.Time, err error) {
	//	Mon Jan 2 15:04:05 -0700 MST 2006
	aTime, err = time.Parse("20060102150405", aDateAsString)
	return aTime, err
}

func AsYYYYMMDDHH24MiSS(aTime time.Time) (aDateAsString string) {
	//	Mon Jan 2 15:04:05 -0700 MST 2006
	aDateAsString = aTime.Format("20060102150405")
	return aDateAsString
}

func AsYYYY_MM_DD_HH24(aTime time.Time) (aDateAsString string) {
	//	Mon Jan 2 15:04:05 -0700 MST 2006
	aDateAsString = aTime.Format("2006-01-02 15")
	return aDateAsString
}
