package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type stringSet map[string]bool

type food struct {
	ingredients stringSet
	allergens   stringSet
}

func (f food) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Food:\n"))
	sb.WriteString("\tIngredients:")
	for i := range f.ingredients {
		sb.WriteString(fmt.Sprintf(" %s", i))
	}
	sb.WriteString("\n")
	sb.WriteString("\tAllergens:")
	for a := range f.allergens {
		sb.WriteString(fmt.Sprintf(" %s", a))
	}
	return sb.String()
}

var foods = make([]food, 0)
var allergenToFood = make(map[string][]int, 0)
var allergenToIngredients = make(map[string]stringSet, 0)
var ingredientTotalCounts = make(map[string]int, 0)
var ingredientsWithAllergens = make(stringSet)
var ingredientsNoAllergens = make(stringSet)
var allergensUniq = make([]string, 0)

func parseAllergens(alg string) stringSet {
	var allergens = make(stringSet)
	for _, a := range strings.Split(alg[0:len(alg)-1], ",") {
		allergens[strings.TrimSpace(a)] = true
	}
	return allergens
}

func parseIngredients(ing string) stringSet {
	var ingredients = make(stringSet)
	for _, i := range strings.Fields(ing) {
		ing := strings.TrimSpace(i)
		if _, found := ingredientTotalCounts[ing]; !found {
			ingredientTotalCounts[ing] = 0
		}
		ingredients[ing] = true
		ingredientTotalCounts[ing] = ingredientTotalCounts[ing] + 1
	}
	return ingredients
}

func parseFood(line string) {
	foodSplits := strings.Split(line, "(contains")
	ingredients := parseIngredients(foodSplits[0])
	var allergens stringSet
	if len(foodSplits) > 1 {
		allergens = parseAllergens(foodSplits[1])
	}
	for a := range allergens {
		if _, found := allergenToFood[a]; !found {
			allergenToFood[a] = make([]int, 0)
			allergensUniq = append(allergensUniq, a)
		}
		allergenToFood[a] = append(allergenToFood[a], len(foods))
	}
	foods = append(foods, food{ingredients, allergens})
}

func intersect(left stringSet, right stringSet) stringSet {
	if len(left) > len(right) {
		return intersect(right, left)
	}
	res := make(stringSet)
	for l := range left {
		if _, found := right[l]; found {
			res[l] = true
		}
	}
	return res
}
func getFoodIntersection(foodIndices []int) stringSet {
	intersection := make(stringSet)
	for i := range foods[foodIndices[0]].ingredients {
		intersection[i] = true
	}
	for _, fid := range foodIndices[1:] {
		intersection = intersect(intersection, foods[fid].ingredients)
	}
	return intersection
}

func solveAllergenToIngredients() {
	for a, f := range allergenToFood {
		allergenToIngredients[a] = getFoodIntersection(f)
		for ing := range allergenToIngredients[a] {
			ingredientsWithAllergens[ing] = true
		}
	}
}

func determineIngredientsNoAllergens() int {
	total := 0
	for i, c := range ingredientTotalCounts {
		if _, found := ingredientsWithAllergens[i]; !found {
			total += c
			ingredientsNoAllergens[i] = true
		}
	}
	return total
}
func getRemainingIng(allergen string) string {
	for k := range allergenToIngredients[allergen] {
		return k
	}
	return ""
}

func solveAllergenMapping() {
	var queue []string
	for a, i := range allergenToIngredients {
		if len(i) == 1 {
			queue = append(queue, a)
		}
	}

	for len(queue) > 0 {
		next := queue[0]
		allergen := getRemainingIng(next)
		queue = queue[1:]
		for a, i := range allergenToIngredients {
			if len(i) > 1 {
				delete(i, allergen)
				if len(i) == 1 {
					queue = append(queue, a)
				}
			}
		}
	}
}

func getDangerList() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s", getRemainingIng(allergensUniq[0])))
	for _, a := range allergensUniq[1:] {
		sb.WriteString(fmt.Sprintf(",%s", getRemainingIng(a)))
	}
	return sb.String()
}

func main() {
	file, err := os.Open("day21_input.txt")
	if err != nil {
		log.Fatalf("failed to open")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		parseFood(line)
	}
	solveAllergenToIngredients()
	totalNoAllergens := determineIngredientsNoAllergens()
	solveAllergenMapping()

	fmt.Println("Foods:")
	for _, f := range foods {
		fmt.Println(f)
	}
	fmt.Println("Allergen To Foods:")
	for a, f := range allergenToFood {
		fmt.Println(a, "=>", f)
	}
	fmt.Println("Allergen To Ingredients:")
	for a, i := range allergenToIngredients {
		fmt.Println(a, "=>", i)
	}
	fmt.Println("Ingredients Possibly With Allergens:")
	for i := range ingredientsWithAllergens {
		fmt.Println(i)
	}
	fmt.Println("Total Ingredients Counts:")
	for i, c := range ingredientTotalCounts {
		fmt.Println(i, "=>", c)
	}
	fmt.Println("Ingredients Without Allergens:")
	for i := range ingredientsNoAllergens {
		fmt.Println(i)
	}
	fmt.Println("Total Without Allergens:", totalNoAllergens)

	fmt.Println("Allergens To Ingredients:")
	sort.Strings(allergensUniq)
	for _, a := range allergensUniq {
		fmt.Println(a, "=>", allergenToIngredients[a], len(allergenToIngredients[a]))
	}
	fmt.Println("Danger List:", getDangerList())
}
