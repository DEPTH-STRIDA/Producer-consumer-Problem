package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// hello
// Структура "человек" - элемент из слайса "люди" (слайс - "динамический массив")
type person struct {
	name string
	age  int
}

// Создание среза "люди" с содержимым из структур "человек"
var persons []person

// Структура "семафор" - определяет мьютекс для безопасного доступа к своему значению. (cond используется для безапасного доступа к Available, которое находится ниже.)
type semaphorWriter struct {
	cond             *sync.Cond
	AvailableWriters int
}
type semaphorReader struct {
	cond             *sync.Cond
	AvailableReaders int
}

// Создание экземпляров структуры. Изначально можно записать 5 людей и, логично, нечего читать из слайса.
var newSemaphorWriter = semaphorWriter{
	cond:             sync.NewCond(&sync.Mutex{}),
	AvailableWriters: 5,
}
var newSemaphorReader = semaphorReader{
	cond:             sync.NewCond(&sync.Mutex{}),
	AvailableReaders: 0,
}

// Мьютекс для доступа к критической области.
var mutex sync.Mutex

// Массив с именами.
var names = [11]string{`Умер`, `Алеша`, `Господин`, `Ирак`, `Порес `, `Вор`, `Тролебузина`, `Перкосрак`, `Владлен`, `Сосика`, `Князь`}

// Уменьшает доступное количество писателей или читателей, если доступное количество становится нулем, усыпляет писателя или читателя.
// Писатели и читатели запускаются как горутины(горутина - функция, запущенная в легковесном потоке. Можно сказать функция в отдельном канале.)
// Усыпление (приостановление выполнения функции) достигается функцией Wait, аналог sleep из java.
func down(typeAction string) {
	if typeAction == `semaphorWriter` {
		newSemaphorWriter.cond.L.Lock()
		if newSemaphorWriter.AvailableWriters == 0 {
			newSemaphorWriter.cond.Wait()
		}
		newSemaphorWriter.AvailableWriters--
		newSemaphorWriter.cond.L.Unlock()
	}
	if typeAction == `semaphorReader` {
		newSemaphorReader.cond.L.Lock()
		if newSemaphorReader.AvailableReaders == 0 {
			newSemaphorReader.cond.Wait()
		}
		newSemaphorReader.AvailableReaders--
		newSemaphorReader.cond.L.Unlock()
	}
}
func up(typeAction string) {
	if typeAction == `semaphorWriter` {
		newSemaphorWriter.cond.L.Lock()
		newSemaphorWriter.AvailableWriters++

		if newSemaphorWriter.AvailableWriters == 1 {
			newSemaphorWriter.cond.Signal()
			newSemaphorWriter.cond.L.Unlock()
			return
		}
		newSemaphorWriter.cond.L.Unlock()

	}
	if typeAction == `semaphorReader` {
		newSemaphorReader.cond.L.Lock()
		newSemaphorReader.AvailableReaders++

		if newSemaphorReader.AvailableReaders == 1 {
			newSemaphorReader.cond.Signal()
			newSemaphorReader.cond.L.Unlock()
			return
		}
		newSemaphorReader.cond.L.Unlock()
	}
}
func writer() {
	down(`semaphorWriter`)
	mutex.Lock()
	fmt.Println(`---------------------------------------`)
	fmt.Println(`|||||||||||| Mutex is lock ||||||||||||`)
	fmt.Println(`---------------------------------------`)
	ageTo := 10 + rand.Intn(20)
	nameIndex := rand.Intn(10)
	toAdd := &person{
		name: names[nameIndex],
		age:  ageTo,
	}
	fmt.Println("Adding: ", toAdd)
	*&persons = append(*&persons, *toAdd)
	fmt.Println(`---------------------------------------`)
	fmt.Println(`||||||||||| Mutex is unlock |||||||||||`)
	fmt.Println(`---------------------------------------`)
	fmt.Println(``)
	fmt.Println(``)
	fmt.Println(`---------------------------------------`)
	fmt.Println(persons)
	fmt.Println(`---------------------------------------`)
	fmt.Println(``)
	fmt.Println(``)
	mutex.Unlock()
	up(`semaphorReader`)
}
func reader() {
	down(`semaphorReader`)
	mutex.Lock()
	fmt.Println(`---------------------------------------`)
	fmt.Println(`||||||||||| Mutex is lock |||||||||||`)
	fmt.Println(`---------------------------------------`)
	if len(persons) > 1 {
		fmt.Println("Reading: ", persons[0])
		a := persons[1:]
		persons = a
	} else {
		if len(persons) == 1 {
			fmt.Println("Reading: ", persons[0])
			mySlice1 := make([]person, 0)
			persons = mySlice1
		}
	}

	fmt.Println(`---------------------------------------`)
	fmt.Println(`||||||||||| Mutex is unlock |||||||||||`)
	fmt.Println(`---------------------------------------`)
	fmt.Println(``)
	fmt.Println(``)
	fmt.Println(`---------------------------------------`)
	fmt.Println(persons)
	fmt.Println(`---------------------------------------`)
	fmt.Println(``)
	fmt.Println(``)
	mutex.Unlock()
	up(`semaphorWriter`)
}

func main() {
	fmt.Println(`||||||||||||||||||||| Starting ||||||||||||||||||||||||`)
	time.Sleep(400 * time.Millisecond)
	fmt.Println(`-------------------------------------------------------`)
	time.Sleep(400 * time.Millisecond)
	fmt.Println(`|||||||||||||||| 5 places available |||||||||||||||||||`)
	time.Sleep(400 * time.Millisecond)
	fmt.Println(``)
	fmt.Println(``)
	//Приставка go - перед функцией, указывает запускать её как горутину.
	go func() {
		for {
			time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
			choice := rand.Intn(2)
			if choice == 1 {
				go writer()
			} else {
				go reader()
			}
		}
	}()
	//Т.к. все функции запускаются как горутины, выполнение main кода сразу переходит сюда.
	//Программа завершится через 30 секунд.
	time.Sleep(30 * time.Second)
}
