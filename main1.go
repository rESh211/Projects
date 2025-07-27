package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	K1 = 0.035
	K2 = 0.029
)

var (
	Format     = "20060102 15:04:05" // формат даты и времени
	StepLength = 0.65                // длина шага в метрах
	Weight     = 75.0                // вес кг
	Height     = 1.75                // рост м
	Speed      = 1.39                // скорость м/с
)

// parsePackage разбирает входящий пакет в параметре data.
func parsePackage(data string) (time.Time, int, bool) {
	// 1. Разделите строку на две части по запятой в слайс ds+
	ds := strings.Split(data, ",")
	// 2. Проверьте, чтобы ds состоял из двух элементов+
	if len(ds) != 2 {
		return time.Time{}, 0, false
	}
	// Формируем полную дату в формате YYYYMMDD HH:MM:SS
	// получаем время time.Time
	t, err := time.Parse(Format, ds[0])
	if err != nil {
		return time.Time{}, 0, false
	}
	// получаем количество шагов
	steps, err := strconv.Atoi(ds[1])
	if err != nil || steps < 0 {
		return time.Time{}, 0, false
	}
	return t, steps, true
}

// stepsDay перебирает все записи слайса, подсчитывает и возвращает
// общее количество шагов
func stepsDay(storage []string) int {
	sum := 0
	for _, v := range storage {
		if _, steps, ok := parsePackage(v); ok {
			sum += steps
		}
	}
	return sum
}

// calories возвращает количество килокалорий, которые потрачены на
func calories(distance float64) float64 {
	// Преобразуем скорость из км/ч в м/с
	// Рассчитываем время в движении (в минутах)
	energyMinute := K1*Weight + (Speed*Speed/Height)*K2*Weight
	period := distance / Speed / 60 // время ходьбы в минутах
	return energyMinute * period
}

// achievement возвращает мотивирующее сообщение в зависимости от
// пройденного расстояния в километрах
func achievement(distance float64) string {
	if distance >= 6.5 {
		return "Отличный результат! Цель достигнута."
	}
	if distance >= 3.9 {
		return "Неплохо! День был продуктивный."
	}
	if distance >= 2 {
		return "Завтра наверстаем!"
	}
	return "Лежать тоже полезно. Главное — участие, а не победа!"
}

// showMessage выводит строку и добавляет два переноса строк
func showMessage(s string) {
	fmt.Printf("%s\n\n", s)
}

// AcceptPackage обрабатывает входящий пакет, который передаётся в
// виде строки в параметре data. Параметр storage содержит пакеты за текущий день.
// Если время пакета относится к новым суткам, storage предварительно
// очищается.
// Если пакет валидный, он добавляется в слайс storage, который возвращает
// функция. Если пакет невалидный, storage возвращается без изменений.
func AcceptPackage(data string, storage []string) []string {
	// 1. Используйте parsePackage для разбора пакета
	//    t, steps, ok := parsePackage(data)
	//    выведите сообщение в случае ошибки
	//    также проверьте количество шагов на равенство нулю
	t, steps, ok := parsePackage(data)
	if !ok {
		showMessage(`ошибочный формат пакета`)
		return storage
	}
	if steps == 0 {
		return storage
	}

	// 2. Получите текущее UTC-время и сравните дни
	//    выведите сообщение, если день в пакете t.Day() не совпадает
	//    с текущим днём
	now := time.Now().UTC()
	if t.Day() != now.Day() {
		showMessage(`неверный день`)
		return storage
	}
	if t.After(now) {
		showMessage(`некорректное значение времени`)
		return storage
	}
	if len(storage) > 0 {
		// 3. Достаточно сравнить первые len(Format) символов пакета с
		//    len(Format) символами последней записи storage
		//    если меньше или равно, то ошибка — некорректное значение времени
		if data[:len(Format)] <= storage[len(storage)-1][:len(Format)] {
			showMessage(`некорректное значение времени`)
			return storage
		}
		// если наступили новые сутки
		if data[:8] != storage[len(storage)-1][:8] {
			// то обнуляем слайс с накопленными данными
			storage = storage[:0]
		}
	}
	// остаётся совсем немного
	// 5. Добавить пакет в storage
	storage = append(storage, data)
	// 6. Получить общее количество шагов
	allSteps := stepsDay(storage)
	// distance — пройденное расстояние в метрах
	distance := float64(allSteps) * StepLength
	// 8. Получить потраченные килокалории
	energy := calories(distance)
	// 9. Получить мотивирующий текст
	distance /= 1000 // переводим в километры
	achiev := achievement(distance)
	// 10. Сформировать и вывести полный текст сообщения
	msg := fmt.Sprintf(`Время: %s.
Количество шагов за сегодня: %d.
Дистанция составила %.2f км.
Вы сожгли %.2f ккал.
%s`, t.Format("15:04:05"), allSteps, distance, energy, achiev)

	showMessage(msg)
	return storage
}

func main() {
	now := time.Now().UTC()
	today := now.Format("20060102")

	input := []string{
		"01:41:03,-100",
		",3456",
		"12:40:00, 3456 ",
		"something is wrong",
		"02:11:34,678",
		"02:11:34,792",
		"17:01:30,1078",
		"03:25:59,7830",
		"04:00:46,5325",
		"04:45:21,3123",
	}

	var storage []string
	storage = AcceptPackage("20230720 00:11:33,100", storage)
	for _, v := range input {
		storage = AcceptPackage(today+" "+v, storage)
	}
}
