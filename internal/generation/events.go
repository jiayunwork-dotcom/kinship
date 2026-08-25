package generation

import (
	"fmt"
	"sort"

	"kinship/internal/genealogy"
)

type Event struct {
	Year   int
	Person string
	Type   string
}

func Timeline(f *genealogy.Family) []Event {
	var events []Event
	for _, name := range f.Names() {
		p, _ := f.Person(name)
		events = append(events, Event{Year: p.Birth, Person: name, Type: "birth"})
	}
	for _, child := range f.Names() {
		cp, _ := f.Person(child)
		parents, _ := f.Parents(child)
		for _, pName := range parents {
			events = append(events, Event{
				Year:   cp.Birth,
				Person: pName,
				Type:   fmt.Sprintf("parent of %s", child),
			})
		}
	}
	sort.Slice(events, func(i, j int) bool {
		if events[i].Year != events[j].Year {
			return events[i].Year < events[j].Year
		}
		return events[i].Person < events[j].Person
	})
	return events
}

func YearRange(events []Event) (int, int) {
	if len(events) == 0 {
		return 0, 0
	}
	return events[0].Year, events[len(events)-1].Year
}

func EventsInRange(events []Event, from, to int) []Event {
	var result []Event
	for _, e := range events {
		if e.Year >= from && e.Year <= to {
			result = append(result, e)
		}
	}
	return result
}

func EventsForPerson(events []Event, name string) []Event {
	var result []Event
	for _, e := range events {
		if e.Person == name {
			result = append(result, e)
		}
	}
	return result
}

func PopulationAtYear(f *genealogy.Family, year int) int {
	count := 0
	for _, name := range f.Names() {
		p, _ := f.Person(name)
		if p.Birth <= year {
			count++
		}
	}
	return count
}
