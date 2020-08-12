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

const INPUT_LOCATION = "./properties.txt"

type Property struct {
	id         int
	address    string
	town       string
	value_date string
	value      int
}

// Property constructor, mostly used for ASCII to int conversion
func newProperty(list []string) Property {
	p := Property{
		address:    list[1],
		town:       list[2],
		value_date: list[3],
	}

	id, err_id := strconv.Atoi(list[0])
	value, err_value := strconv.Atoi(list[4])

	if err_id == nil && err_value == nil {
		p.id = id
		p.value = value
	} else {
		p.id = -1
	}
	return p
}

// Pretty output for property struct
func (p Property) PropertyInfo() string {
	fmt.Printf("--ID: %d----------------------------", p.id)
	fmt.Printf("\nAddress: %s, %s", p.address, p.town)
	fmt.Printf("\nValue: $%d - Date: %s", p.value, p.value_date)
	return ""
}

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
		p.value_date == p2.value_date
}

// Checks property list against this property, returns
//   index of first duplicate property found, or 0 if none
func (p Property) dupeCheck(properties []Property) int {
	for index, prop := range properties {
		if prop.Equals(p) {
			return index
		}
	}
	return 0
}

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
	file, err := os.Open(INPUT_LOCATION)

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
		property_data := strings.Split(line, "\t")

		// Looking for lines with five tab-sep'd fields
		if len(property_data) == 5 {
			p := newProperty(property_data)

			// only valid properties are added to list
			if p.IsValid() {
				duplicate := p.dupeCheck(properties)
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
	// Test 4 - Extra Credit - running four goroutines concurrently
  
  var PropBuffer = make(chan []Property)
  
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

  for _, p := range Props {
    
    go func() {
      p = MinCostFilter(p, 40000)
      p = PretentiousFilter(p)
      p = NthRecordFilter(p, 10)
      props := append(<-PropBuffer, p...)
      PropBuffer <- props
    }()
  }
  
  go func() {
    fmt.Printf("Total props: %d", len(<-PropBuffer))
  }()

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
