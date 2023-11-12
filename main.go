package main

import (
	"fmt"
	"vivarium/enums"
)

func describeSex(sex enums.Sexe) {
	if sex == enums.Male {
		fmt.Println("The sex is Male.")
	} else if sex == enums.Femelle {
		fmt.Println("The sex is Female.")
	} else {
		fmt.Println("The sex is Hermaphrodite.")
	}
}

type Creature struct {
	genre enums.MyInsect
	sexe  enums.Sexe
}

func main() {
	sex := enums.Male
	fmt.Println(sex.String()) // print "Male"

	cre1 := Creature{
		genre: enums.Lombric,
		sexe:  enums.Femelle,
	}
	describeSex(cre1.sexe)

	cre2 := Creature{
		genre: enums.Escargot,
		sexe:  enums.Hermaphrodite,
	}
	describeSex(cre2.sexe)

}
