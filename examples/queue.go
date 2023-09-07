package main

func main() {
	q := New[string](1024, 32)
	q.Push("Fucking")
	q.Push("The")
	q.Push("World")

	q.Printf()

	var new_q Queue[string]
	new_q.Push("You")
	new_q.Push("Find Your Worth In The Fucking World")
	new_q.UnShift("May")

	new_q.Printf()
}
