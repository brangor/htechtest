package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

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
func MinCostFilter(properties []Property, min int) []Property {
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
func PretentiousFilter(properties []Property) []Property {
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
func NthRecordFilter(properties []Property, n int) []Property {
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

	// Test 4 - Filtering
	var FilteredProperties []Property
	FilteredProperties = MinCostFilter(properties, 40000)
	FilteredProperties = PretentiousFilter(FilteredProperties)
	FilteredProperties = NthRecordFilter(FilteredProperties, 10)

	fmt.Println("Filtered Properties List:")

	for _, p := range FilteredProperties {
		fmt.Println(p.PropertyInfo())
	}

	fmt.Printf("\n# Properties (unfiltered): %d", len(properties))
	fmt.Printf("\n# Properties (filtered): %d", len(FilteredProperties))

}
