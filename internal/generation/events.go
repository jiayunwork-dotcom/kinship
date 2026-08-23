package generation

import (
	"fmt"
	"sort"

	"kinship/internal/genealogy"
)

// Event represents a life event in the family timeline.
type Event struct {
	Year   int
	Person string
	Type   string // "birth", "parent"
}

// Timeline builds a chronological list of all known events in the family.
func Timeline(f *genealogy.Family) []Event {
	var events []Event
	for _, name := range f.Names() {
		p, _ := f.Person(name)
		events = append(events, Event{Year: p.Birth, Person: name, Type: "birth"})
	}
	// parent events: year child was born, from parent's perspective
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

// YearRange returns min and max years from events.
func YearRange(events []Event) (int, int) {
	if len(events) == 0 {
		return 0, 0
	}
	return events[0].Year, events[len(events)-1].Year
}

// EventsInRange returns events within a year range [from, to].
func EventsInRange(events []Event, from, to int) []Event {
	var result []Event
	for _, e := range events {
		if e.Year >= from && e.Year <= to {
			result = append(result, e)
		}
	}
	return result
}

// EventsForPerson returns all events involving the named person.
func EventsForPerson(events []Event, name string) []Event {
	var result []Event
	for _, e := range events {
		if e.Person == name {
			result = append(result, e)
		}
	}
	return result
}

// PopulationAtYear returns how many persons were already born by the given year.
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
