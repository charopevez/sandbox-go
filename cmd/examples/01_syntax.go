// =============================================================
// Go Syntax Refresher for PHP Developers
// Run: go run cmd/examples/01_syntax.go
// =============================================================
package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// -----------------------------------------------------------
//  1. VARIABLES & TYPES
//     PHP: $name = "Go";       (dynamic typing)
//     Go:  name := "Go"        (type inferred, but STATIC)
//
// -----------------------------------------------------------
func variables() {
	fmt.Println("=== VARIABLES ===")

	// Short declaration (most common inside functions)
	name := "Go"
	age := 3
	price := 19.99
	active := true

	// Explicit type (useful for zero values or clarity)
	var count int    // defaults to 0
	var label string // defaults to ""

	fmt.Println(name, age, price, active, count, label)

	// Constants (like PHP's const / define)
	const maxRetries = 3
	const baseURL = "https://api.example.com"
	fmt.Println(maxRetries, baseURL)

	// Multiple assignment (no list() needed like PHP)
	x, y := 10, 20
	x, y = y, x // swap — no temp variable!
	fmt.Println("swapped:", x, y)
}

// -----------------------------------------------------------
//  2. STRINGS
//     PHP: substr, strpos, explode, implode, sprintf
//     Go:  strings package + fmt.Sprintf
//
// -----------------------------------------------------------
func stringOps() {
	fmt.Println("\n=== STRINGS ===")

	s := "Hello, World"

	// PHP: strlen($s)        → Go: len(s) — gives bytes, not chars!
	fmt.Println("length:", len(s))

	// PHP: strpos($s, "World") → Go:
	fmt.Println("contains World:", strings.Contains(s, "World"))
	fmt.Println("index of World:", strings.Index(s, "World"))

	// PHP: strtoupper / strtolower
	fmt.Println("upper:", strings.ToUpper(s))

	// PHP: explode(",", $s) → Go:
	parts := strings.Split(s, ", ")
	fmt.Println("split:", parts)

	// PHP: implode(" | ", $arr) → Go:
	fmt.Println("join:", strings.Join(parts, " | "))

	// PHP: sprintf("Hi %s, age %d", $name, $age) → Go: same!
	msg := fmt.Sprintf("Hi %s, you have %d tasks", "Alice", 5)
	fmt.Println(msg)

	// Multiline strings (like PHP heredoc)
	query := `
		SELECT id, name
		FROM users
		WHERE active = true
	`
	fmt.Println("raw string:", query)
}

// -----------------------------------------------------------
//  3. ARRAYS, SLICES, MAPS
//     PHP arrays do EVERYTHING → Go splits into slice and map
//
// -----------------------------------------------------------
func collections() {
	fmt.Println("\n=== SLICES & MAPS ===")

	// SLICE — like PHP indexed array (but typed + fixed element type)
	// PHP: $fruits = ["apple", "banana", "cherry"];
	fruits := []string{"apple", "banana", "cherry"}
	fmt.Println("fruits:", fruits)

	// PHP: $fruits[] = "date";  → Go:
	fruits = append(fruits, "date")
	fmt.Println("after append:", fruits)

	// PHP: array_slice($fruits, 1, 2) → Go: slicing syntax
	fmt.Println("slice [1:3]:", fruits[1:3]) // "banana", "cherry"

	// PHP: count($fruits) → Go:
	fmt.Println("length:", len(fruits))

	// Iterating — PHP: foreach ($fruits as $i => $f)
	for i, f := range fruits {
		fmt.Printf("  [%d] %s\n", i, f)
	}

	// MAP — like PHP associative array
	// PHP: $scores = ["alice" => 95, "bob" => 87];
	scores := map[string]int{
		"alice": 95,
		"bob":   87,
	}

	// PHP: $scores["charlie"] = 92;
	scores["charlie"] = 92

	// PHP: isset($scores["alice"])
	val, exists := scores["alice"]
	fmt.Printf("alice: %d, exists: %v\n", val, exists)

	// PHP: unset($scores["bob"]);
	delete(scores, "bob")

	// PHP: foreach ($scores as $k => $v)
	for k, v := range scores {
		fmt.Printf("  %s → %d\n", k, v)
	}
}

// -----------------------------------------------------------
// 4. FUNCTIONS
//    PHP: function add($a, $b) { return $a + $b; }
//    Go:  func add(a, b int) int { return a + b }
// -----------------------------------------------------------

// Multiple return values — Go's killer feature vs PHP
// (PHP needs array or object to return multiple values)
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

// Named return values
func minMax(nums []int) (min, max int) {
	min, max = nums[0], nums[0]
	for _, n := range nums {
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}
	return // "naked return" — returns min, max
}

// Variadic function — like PHP's ...$args
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// Closures — same concept as PHP closures but cleaner syntax
// PHP: $double = function($x) { return $x * 2; };
// Go:
func makeMutliplier(factor int) func(int) int {
	return func(x int) int {
		return x * factor
	}
}

func functions() {
	fmt.Println("\n=== FUNCTIONS ===")

	// Multiple returns + error handling
	result, err := divide(10, 3)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Printf("10 / 3 = %.2f\n", result)
	}

	// Error case
	_, err = divide(10, 0)
	if err != nil {
		fmt.Println("caught:", err)
	}

	min, max := minMax([]int{3, 1, 4, 1, 5, 9, 2, 6})
	fmt.Printf("min=%d, max=%d\n", min, max)

	fmt.Println("sum:", sum(1, 2, 3, 4, 5))

	double := makeMutliplier(2)
	triple := makeMutliplier(3)
	fmt.Println("double(5):", double(5))
	fmt.Println("triple(5):", triple(5))
}

// -----------------------------------------------------------
// 5. STRUCTS & METHODS
//    PHP: class UserBase { ... }
//    Go:  no classes! Structs + methods instead
// -----------------------------------------------------------

// PHP equivalent:
// class UserBase {
//     public string $Name;
//     public string $Email;
//     public function __construct($name, $email) { ... }
//     public function greet(): string { ... }
// }

type UserBase struct {
	Name  string
	Email string
	Age   int
}

// Method on UserBase (like a class method in PHP)
// Value receiver — doesn't modify the original
func (u UserBase) Greet() string {
	return fmt.Sprintf("Hi, I'm %s (%s)", u.Name, u.Email)
}

// Pointer receiver — CAN modify the original
// PHP: methods always modify $this, Go you must choose
func (u *UserBase) Birthday() {
	u.Age++
}

// "Constructor" pattern (Go has no constructors)
func NewUserBase(name, email string, age int) *UserBase {
	return &UserBase{
		Name:  name,
		Email: email,
		Age:   age,
	}
}

func structs() {
	fmt.Println("\n=== STRUCTS ===")

	// Create struct — like `new UserBase(...)` in PHP
	u1 := UserBase{Name: "Alice", Email: "alice@test.com", Age: 30}
	fmt.Println(u1.Greet())

	// Using "constructor"
	u2 := NewUserBase("Bob", "bob@test.com", 25)
	fmt.Println(u2.Greet())
	fmt.Println("Bob's age:", u2.Age)
	u2.Birthday()
	fmt.Println("After birthday:", u2.Age)

	// Struct embedding (Go's version of inheritance)
	type Admin struct {
		UserBase // embedded — Admin "inherits" UserBase's fields and methods
		Level    int
	}

	admin := Admin{
		UserBase: UserBase{Name: "Charlie", Email: "charlie@test.com", Age: 40},
		Level:    1,
	}
	// Can call UserBase methods directly on Admin
	fmt.Println(admin.Greet(), "- Admin level:", admin.Level)
}

// -----------------------------------------------------------
// 6. INTERFACES
//    PHP: interface Shape { public function area(): float; }
//    Go:  IMPLICIT — no "implements" keyword!
// -----------------------------------------------------------

type Shape interface {
	Area() float64
	String() string
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) String() string {
	return fmt.Sprintf("Circle(r=%.1f)", c.Radius)
}

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) String() string {
	return fmt.Sprintf("Rect(%.1f×%.1f)", r.Width, r.Height)
}

// This function accepts ANY Shape — like PHP type-hinting an interface
func printShape(s Shape) {
	fmt.Printf("  %s → area = %.2f\n", s, s.Area())
}

func interfaces() {
	fmt.Println("\n=== INTERFACES ===")

	shapes := []Shape{
		Circle{Radius: 5},
		Rectangle{Width: 3, Height: 4},
		Circle{Radius: 2.5},
	}

	for _, s := range shapes {
		printShape(s)
	}
}

// -----------------------------------------------------------
// 7. ERROR HANDLING
//    PHP: try/catch/throw
//    Go:  return errors explicitly (no exceptions!)
// -----------------------------------------------------------

type ValidationError struct {
	Field   string
	Message string
}

// Implement the error interface (just needs Error() string)
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on %s: %s", e.Field, e.Message)
}

func validateAge(age int) error {
	if age < 0 {
		return &ValidationError{Field: "age", Message: "cannot be negative"}
	}
	if age > 150 {
		return &ValidationError{Field: "age", Message: "unrealistic value"}
	}
	return nil
}

func errors() {
	fmt.Println("\n=== ERROR HANDLING ===")

	// The Go pattern: always check errors immediately
	for _, age := range []int{25, -1, 200} {
		if err := validateAge(age); err != nil {
			fmt.Printf("  age %d → ERROR: %s\n", age, err)
		} else {
			fmt.Printf("  age %d → OK\n", age)
		}
	}
}

// -----------------------------------------------------------
// 8. CONTROL FLOW
// -----------------------------------------------------------
func controlFlow() {
	fmt.Println("\n=== CONTROL FLOW ===")

	// if with init statement (Go-specific, very useful)
	if n := time.Now().Hour(); n < 12 {
		fmt.Println("Good morning!")
	} else if n < 18 {
		fmt.Println("Good afternoon!")
	} else {
		fmt.Println("Good evening!")
	}

	// switch — no break needed! (opposite of PHP)
	day := time.Now().Weekday()
	switch day {
	case time.Saturday, time.Sunday:
		fmt.Println("Weekend!")
	default:
		fmt.Println("Weekday:", day)
	}

	// switch with no condition (cleaner than if/else chains)
	score := 85
	switch {
	case score >= 90:
		fmt.Println("Grade: A")
	case score >= 80:
		fmt.Println("Grade: B")
	case score >= 70:
		fmt.Println("Grade: C")
	default:
		fmt.Println("Grade: F")
	}

	// for loop — Go's ONLY loop (no while, no do-while)
	// PHP's while: while ($i < 5) → Go: for i < 5
	// PHP's for: for ($i=0; $i<5; $i++) → Go: for i := 0; i < 5; i++
	for i := 0; i < 3; i++ {
		fmt.Printf("  loop %d\n", i)
	}
}

// -----------------------------------------------------------
// 9. POINTERS (PHP doesn't have these explicitly)
// -----------------------------------------------------------
func pointers() {
	fmt.Println("\n=== POINTERS ===")

	// In PHP, objects are always passed by reference
	// In Go, you must be explicit with pointers

	x := 42
	p := &x // p is a pointer to x (stores memory address)
	fmt.Println("value:", x)
	fmt.Println("pointer:", p) // memory address like 0xc0000b4008
	fmt.Println("deref:", *p)  // *p reads the value at that address → 42

	*p = 100                          // modify x through the pointer
	fmt.Println("x after *p=100:", x) // 100

	// Why this matters: function arguments are COPIES by default
	// Use pointers when you need to modify the original
	func(val int) {
		val = 999 // this changes nothing outside
	}(x)
	fmt.Println("after pass-by-value:", x) // still 100

	func(ptr *int) {
		*ptr = 999 // this DOES change x
	}(&x)
	fmt.Println("after pass-by-pointer:", x) // 999
}

// -----------------------------------------------------------
// MAIN — run all examples
// -----------------------------------------------------------
func main() {
	variables()
	stringOps()
	collections()
	functions()
	structs()
	interfaces()
	errors()
	controlFlow()
	pointers()

	fmt.Println("\n✅ All examples done!")
}
