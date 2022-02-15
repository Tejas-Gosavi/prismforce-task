package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

type SingleData struct {
	Amount    int       `json:"amount"`
	StartDate time.Time `json:"startDate"`
}

type Data struct {
	ExpenseData []SingleData `json:"expenseData"`
	RevenueData []SingleData `json:"revenueData"`
}

type SingleBalance struct {
	Amount    int    `json:"amount"`
	StartDate string `json:"startDate"`
}

type BalanceSheet struct {
	Balance []SingleBalance `json:"balance"`
}

type dateAmountMAP map[time.Time]int

func ParseJSON(fileBytes []byte) Data {
	fmt.Println("Parseing JSON data...")
	var parsedData Data

	// parsing JSON encoded data(in []byte format) to struct
	err := json.Unmarshal(fileBytes, &parsedData)
	if err != nil {
		// if any error occurs then print that error and stop executing current program
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Print("Completed parseing JSON data!!!\n\n")
	return parsedData
}

func mergeSameDayData(parsedData []SingleData) ([]SingleData, int) {
	fmt.Println("Merging same day data...")
	var newParsedData []SingleData

	// Creating a map with Keys of type time.Time and value of type int
	mp := make(dateAmountMAP, 12)

	// As told in README.md file, year and month will change, but day and time will remain same
	// so, getting first month-year
	// e.g. if first entry in parsedData is 2020-05-01T00:00:00.000Z, then
	// starting month for that year will be 2020-01-01T00:00:00.000Z
	stMonth := parsedData[0].StartDate.AddDate(0, 1-int(parsedData[0].StartDate.Month()), 0)

	// Looping from i = 0 to 12 and adding i to starting month-year to get i'th month-date
	// and setting Amount of i'th month-date to 0
	// for i = 0,
	// 2020-01-01T00:00:00.000Z + 0 -> 2020-01-01T00:00:00.000Z
	// 2020-01-01T00:00:00.000Z = 0
	// for i = 1,
	// 2020-01-01T00:00:00.000Z + 1 -> 2020-02-01T00:00:00.000Z
	// 2020-02-01T00:00:00.000Z = 0
	// for i = 2,
	// 2020-01-01T00:00:00.000Z + 0 -> 2020-03-01T00:00:00.000Z
	// 2020-03-01T00:00:00.000Z = 0
	// and so on...
	for i := 0; i < 12; i++ {
		// AddDate function will give date after (year, month, day)
		// here we need next month with same year and day
		// so will add 0, i, 0
		// follow this link to know more -> https://www.golangprograms.com/add-n-number-of-year-month-day-hour-minute-second-millisecond-microsecond-and-nanosecond-to-current-date-time.html
		mp[stMonth.AddDate(0, i, 0)] = 0
	}

	// looping through parsedData
	// for every startDate we will map it to its Amount in mp
	// here data with multiple entries will also be merged
	// we will also find (mxmonth) to upto which month we need to data add in balance sheet
	// e.g. if in parsedData entries are like
	//   {
	// 		"amount": 50,
	// 		"startDate": "2021-01-01T00:00:00.000Z"
	//   },
	//   {
	// 		"amount": 20,
	// 		"startDate": "2021-02-01T00:00:00.000Z"
	//   },
	//   {
	// 		"amount": 30,
	// 		"startDate": "2021-03-01T00:00:00.000Z"
	//   }
	// then we only need to add data upto 3rd month i.e. march
	// we don't need data after march, so here mxmonth = 3
	var mxmonth int
	for _, d := range parsedData {
		mp[d.StartDate] += d.Amount
		if mxmonth < int(d.StartDate.Month()) {
			mxmonth = int(d.StartDate.Month())
		}
	}

	// converting merged data of type map to struct
	for t, v := range mp {
		newParsedData = append(newParsedData, SingleData{Amount: v, StartDate: t})
	}
	fmt.Print("Completed merging same day data!!!\n\n")
	fmt.Println("Sorting data by date...")

	// As mentioned in README.md file, Sorting that struct by Date in Ascending order
	sort.Slice(newParsedData, func(i, j int) bool {
		return newParsedData[i].StartDate.Before(newParsedData[j].StartDate)
	})
	fmt.Print("Completed sorting data by date!!!\n\n")

	// returning newly unique, merged, sorted parsedData and mxmonth
	return newParsedData, mxmonth
}

func CreateBalanceSheet(parsedData Data) BalanceSheet {
	var mx, mxmonthrd, mxmonthed int

	// getting sorted data with no multiple entries per month, no timestamps may not overlapd
	// and mxmonth for RevenueData and ExpenseData resp.
	parsedData.RevenueData, mxmonthrd = mergeSameDayData(parsedData.RevenueData)
	parsedData.ExpenseData, mxmonthed = mergeSameDayData(parsedData.ExpenseData)

	fmt.Println("Generating Balance sheet...")
	// creating temparary balance sheet and adding data in it
	// using merge two sorted arrays algorithm, know more -> https://www.geeksforgeeks.org/merge-two-sorted-arrays/
	var tempBalance []SingleData
	var i, j int

	// for every i'th entry in ExpenseData and j'th entry in RevenueData
	for i < len(parsedData.ExpenseData) && j < len(parsedData.RevenueData) {

		// if their starting date are same then
		// Amount for that date will be revenue.amount - expense.amount'
		// and append it to temparary balance sheet and increment i and j'
		if parsedData.ExpenseData[i].StartDate == parsedData.RevenueData[j].StartDate {
			tempBalance = append(tempBalance, SingleData{Amount: parsedData.RevenueData[j].Amount - parsedData.ExpenseData[i].Amount,
				StartDate: parsedData.ExpenseData[i].StartDate})
			i++
			j++

			// otherwise less date from expensesData[i] and RevenueData[j]
			// with its Amount will be appended to temparary balance sheet and will increment i or j
		} else if parsedData.ExpenseData[i].StartDate.After(parsedData.RevenueData[j].StartDate) {
			tempBalance = append(tempBalance, SingleData{Amount: parsedData.RevenueData[j].Amount,
				StartDate: parsedData.RevenueData[j].StartDate})
			j++
		} else {
			tempBalance = append(tempBalance, SingleData{Amount: parsedData.ExpenseData[i].Amount,
				StartDate: parsedData.ExpenseData[i].StartDate})
			i++
		}
	}

	// Data from remaining array will be add directly to the temparary balance sheet
	for i < len(parsedData.ExpenseData) {
		tempBalance = append(tempBalance, SingleData{Amount: parsedData.ExpenseData[i].Amount,
			StartDate: parsedData.ExpenseData[i].StartDate})
		i++
	}
	for j < len(parsedData.RevenueData) {
		tempBalance = append(tempBalance, SingleData{Amount: parsedData.RevenueData[j].Amount,
			StartDate: parsedData.RevenueData[j].StartDate})
		j++
	}

	// now we only need to add data from temporary balance sheet to final balance sheet having more mxmonth
	// i.e. if mxmonthrd = 3 and mxmonthed = 5
	// then from mxmonthrd and mxmonthed
	// mxmonthed is greater then we need to add data upto 5'th month i.e. upto may month
	// mx = mxmonthed
	if mxmonthrd > mxmonthed {
		mx = mxmonthrd
	} else {
		mx = mxmonthed
	}

	// looping from 0 to mx
	// append data from temporary balance sheet to final balance sheet
	var bal BalanceSheet
	for i := 0; i < mx; i++ {
		bal.Balance = append(bal.Balance, SingleBalance{Amount: tempBalance[i].Amount, StartDate: strings.Join(strings.Split(tempBalance[i].StartDate.Format(time.RFC3339), "Z"), ".000Z")})
	}

	fmt.Print("Completed generating Balance sheet!!!\n")

	// returning final balance sheet
	return bal
}
