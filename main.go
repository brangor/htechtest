package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// InputLocation - location of data input file
const InputLocation = "./properties.txt"

// Property struct stores property valuation data
type Property struct {
	id        int
	address   string
	town      string
	valueDate string
	value     int
}

// Property constructor, mostly used for ASCII to int conversion
func newProperty(list []string) Property {
	p := Property{
		address:   list[1],
		town:      list[2],
		valueDate: list[3],
	}

	id, errID := strconv.Atoi(list[0])
	value, errValue := strconv.Atoi(list[4])

	if errID == nil && errValue == nil {
		p.id = id
		p.value = value
	} else {
		p.id = -1
	}
	return p
}

// PropertyInfo function - Pretty output for property struct
func (p Property) PropertyInfo() string {
	fmt.Printf("--ID: %d----------------------------", p.id)
	fmt.Printf("\nAddress: %s, %s", p.address, p.town)
	fmt.Printf("\nValue: $%d - Date: %s", p.value, p.valueDate)
	return ""
}

// IsValid function
// Definition of validity, could change if more restrictions added
// Currently only checks if it has an ID in the valid range
func (p Property) IsValid() bool {
	return p.id >= 1
}

// Equals function based on wording in tech test:
// "A duplicate is a row that has the same address and same date."
// Assumed that didn't include the 'town' field, but they come out
//   the same with/without, so didn't worry about it.
func (p Property) Equals(p2 Property) bool {
	return strings.ToUpper(p.address) == strings.ToUpper(p2.address) &&
		p.valueDate == p2.valueDate
}

// DupeCheck checks property list against this property,
//  returns index of first duplicate property found or 0 if none
func (p Property) DupeCheck(properties []Property) int {
	for index, prop := range properties {
		if prop.Equals(p) {
			return index
		}
	}
	return 0
}

// MinCostFilter function
// Test 4.1 - Filtering out cheap properties
// Takes lower limit and returns an array of properties at or above
//  that value
func MinCostFilter(properties []Property, min int, wg *sync.WaitGroup) []Property {
	defer wg.Done()

	var FilteredList []Property
	for _, p := range properties {
		if p.value >= min {
			FilteredList = append(FilteredList, p)
		}
	}
	return FilteredList
}

// PretentiousFilter function
// Test 4.2 - Filtering out pretentious properties
// Removes properties from list whose addresses include AVE, CRES, or PL
func PretentiousFilter(properties []Property, wg *sync.WaitGroup) []Property {
	defer wg.Done()

	var FilteredList []Property

	for _, p := range properties {

		matched, err := regexp.MatchString(`(AVE|CRES|PL).*`, strings.ToUpper(p.address))

		if err == nil {

			if matched == false {
				FilteredList = append(FilteredList, p)
			} else {
				continue
			}

		}
	}
	return FilteredList
}

// NthRecordFilter function
// Test 4.3 - Filtering out every 10th Property
// Removes every 10th property from list
func NthRecordFilter(properties []Property, n int, wg *sync.WaitGroup) []Property {
	defer wg.Done()

	var FilteredList []Property
	for i, p := range properties {
		if i%n != 0 {
			FilteredList = append(FilteredList, p)
		}
	}
	return FilteredList
}

func main() {
	file, err := os.Open(InputLocation)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Initialise properties array
	var properties []Property

	for scanner.Scan() {
		line := scanner.Text()

		// Each set of data contained on a line
		//  w/ tab-sep'd fields
		propertyData := strings.Split(line, "\t")

		// Looking for lines with five tab-sep'd fields
		if len(propertyData) == 5 {
			p := newProperty(propertyData)

			// only valid properties are added to list
			if p.IsValid() {
				duplicate := p.DupeCheck(properties)
				// Test 3 - no instances of duplicate entered at all
				// Removing found duplicates
				if duplicate != 0 {
					//replacing property at index `duplicate` with the one at the end
					properties[duplicate] = properties[len(properties)-1]
					// then slicing the array to just before the last element
					properties = properties[:len(properties)-1]
				} else {
					properties = append(properties, p)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Test 4 - Extra Credit - running four goroutines concurrently
	var wg sync.WaitGroup

	//Tried to do some waitgroup things here, but got all mixed up.
	// Does not work, but leaving my work for reference.
	// I'll admit, I'm not too used to multi-threaded applications!

	// Splitting array into fourths
	TotalProps := len(properties)

	var Props [4][]Property

	for i, p := range properties {
		var PropsIndex int
		if i < TotalProps/4 {
			PropsIndex = 0
		} else if i < TotalProps/2 {
			PropsIndex = 1
		} else if i < TotalProps/4*3 {
			PropsIndex = 2
		} else {
			PropsIndex = 3
		}
		Props[PropsIndex] = append(Props[PropsIndex], p)
	}

	fmt.Printf("Props 1: %d %d %d %d", len(Props[0]), len(Props[1]), len(Props[2]), len(Props[3]))

	go func() {
		wg.Add(3)
		Props[0] = MinCostFilter(Props[0], 40000, &wg)
		Props[0] = PretentiousFilter(Props[0], &wg)
		Props[0] = NthRecordFilter(Props[0], 10, &wg)
	}()

	go func() {
		wg.Add(3)
		Props[1] = MinCostFilter(Props[1], 40000, &wg)
		Props[1] = PretentiousFilter(Props[1], &wg)
		Props[1] = NthRecordFilter(Props[1], 10, &wg)
	}()

	go func() {
		wg.Add(3)
		Props[2] = MinCostFilter(Props[2], 40000, &wg)
		Props[2] = PretentiousFilter(Props[2], &wg)
		Props[2] = NthRecordFilter(Props[2], 10, &wg)
	}()

	go func() {
		wg.Add(3)
		Props[3] = MinCostFilter(Props[3], 40000, &wg)
		Props[3] = PretentiousFilter(Props[3], &wg)
		Props[3] = NthRecordFilter(Props[3], 10, &wg)
	}()

	wg.Wait()

	go func() {
		var allProps []Property
		for _, props := range Props {
			allProps = append(allProps, props...)
		}
		for _, p := range allProps {
			fmt.Println(p.PropertyInfo())
		}
	}()

}
