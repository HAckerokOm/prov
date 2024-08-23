package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)
func calcSumDirect(pathDirectory string) (int64, error) {

    //считывание содержания директории    
	files, err := os.ReadDir(pathDirectory)
    if err != nil {        
		return 0, err
    }
    //суммарный размер директории    
	var sum int64 = 0
    //проход по кадому файлу
    for _, file := range files {
        //формирование нового пути к внутренней директории и рекурсивный вызов        
		//с этим путем в качестве аргумента
        if file.IsDir() {            
            dirSum, err := calcSumDirect(fmt.Sprintf("%s/%s", pathDirectory, file.Name())) 
			if err != nil {          // Обработка ошибок при получении информации о файле
				fmt.Printf("Не удалось получить размер директории'%s': %v\n", file.Name(), err)
				continue  }         
			sum += dirSum
            continue        
		}
        info, _ := file.Info()
        sum += info.Size()
    }   
	return sum, nil}

// Функция parseFlags анализирует командные флаги для получения пути к директории.
func parseFlags() string {
	includeCurrentDirPtr := flag.String("dst", "", "Путь к директории") // Определение флага для пути к директории
	flag.Parse()                                                         // Парсинг командных аргументов
	if *includeCurrentDirPtr == "" {                                     // Проверка, был ли указан путь
		flag.PrintDefaults() // Вывод значений по умолчанию, если путь не указан
		os.Exit(1)           // Выход из программы с ошибкой
	}
	return *includeCurrentDirPtr // Возвращение прочитанного пути к директории
}

// Функция listFilesInDirectory перечисляет все файлы в заданной директории.
func listFilesInDirectory(dirPath string) ([]os.DirEntry, error) {
	return os.ReadDir(dirPath) // Чтение содержимого директории
}

// Функция printFileDetails выводит детальную информацию о каждом файле в переданном списке объектов os.DirEntry.
func printFileDetails(files []os.DirEntry, mdir string) {
	fmt.Printf("%-15s %-15s %-30s\n", "Тип", "Размер", "Название") // Вывод заголовков столбцов
	for _, file := range files {                                   // Итерация по каждому файлу
		info, err := file.Info() // Получение информации о файле
		if err != nil {          // Обработка ошибок при получении информации о файле
			fmt.Printf("Не удалось получить информацию о файле '%s': %v\n", file.Name(), err)
			continue
		}
		isDirectory := info.IsDir() // Определение, является ли запись директорией
		size := info.Size()         // Получение размера файла/директории
		name := file.Name()         // Получение имени файла/директории
		tip := " "                  // Тип
		if isDirectory {            // Обновление описания типа в зависимости от того, является ли оно директорией
			tip = "Директория"

			directsum, err := calcSumDirect(fmt.Sprintf("%s/%s", mdir , file.Name()))
			size+=directsum
			if err != nil {          // Обработка ошибок при получении информации о файле
				fmt.Printf("Не удалось получить размер дирректории '%s': %v\n", file.Name(), err)
				continue}
		} else {
			tip = "Файл"
		}
		formattedSize := formatSize(size)                                    // Форматирование размера для вывода
		line := fmt.Sprintf("%-15s %-15s %-30s\n", tip, formattedSize, name) // Форматирование строки для вывода
		fmt.Println(line)                                                    // Вывод строки
	}
}

// Функция formatSize форматирует количество байт в человекочитаемую строку.
func formatSize(bytes int64) string {
	switch {
	case bytes >= 1000*1000*1000: // Если размер в гигабайтах или больше
		return fmt.Sprintf("%.2f Гигабайта", float64(bytes)/float64(1000*1000*1000))
	case bytes >= 1000*1000: // Если размер в мегабайтах
		return fmt.Sprintf("%.2f Мегабайта", float64(bytes)/float64(1000*1000))
	case bytes >= 1000: // Если размер в килобайтах
		return fmt.Sprintf("%.2f Килобайта", float64(bytes)/float64(1000))
	default: // Если размер меньше килобайта
		return fmt.Sprintf("%d Байта", bytes)
	}
}

func main() {
	dirPath := parseFlags() // Парсинг командных флагов для получения пути к директории

	startTime := time.Now() // Запись времени начала работы программы
	fmt.Println("Программа выполняется...")

	files, err := listFilesInDirectory(dirPath) // Перечисление всех файлов в директории
	if err != nil {                             // Обработка ошибок при перечислении файлов
		fmt.Printf("ошибка при чтении директории '%s': %v\n", dirPath, err)
		os.Exit(1)
	}

	printFileDetails(files,dirPath) // Вывод деталей каждого файла

	endTime := time.Now()                                    // Запись времени окончания работы программы
	duration := endTime.Sub(startTime)                       // Расчет времени выполнения
	fmt.Printf("Время выполнения программы: %v\n", duration) // Вывод времени выполнения
}
