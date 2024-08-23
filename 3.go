package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"flag"
)

// Функция parseFlags анализирует командные флаги для получения пути к директории.
func parseFlags() string {
	includeCurrentDirPtr := flag.String("dst", "", "Путь к дирректории") // Определение флага для пути к директории
	flag.Parse()                                                          // Парсинг командных аргументов
	if *includeCurrentDirPtr == "" {                                      // Проверка, был ли указан путь
		flag.PrintDefaults()                                              // Вывод значений по умолчанию, если путь не указан
		os.Exit(1)                                                        // Выход из программы с ошибкой
	}
	return *includeCurrentDirPtr                                         // Возвращение прочитанного пути к директории
}

// Функция listFilesInDirectory перечисляет все файлы в заданной директории.
func listFilesInDirectory(dirPath string) ([]os.DirEntry, error) {
	return os.ReadDir(dirPath)                                          // Чтение содержимого директории
}

// Функция printFileDetails выводит детальную информацию о каждом файле в переданном списке объектов os.DirEntry.
func printFileDetails(files []os.DirEntry) {
	fmt.Printf("%-15s %-15s %-30s\n", "Тип", "Размер", "Название")   // Вывод заголовков столбцов
	for _, file := range files {                                        // Итерация по каждому файлу
		info, err := file.Info()                                        // Получение информации о файле
		if err != nil {                                                // Обработка ошибок при получении информации о файле
			fmt.Printf("Не удалось получить информацию о файле '%s': %v\n", file.Name(), err)
			continue
		}
		isDirectory := info.IsDir()                                     // Определение, является ли запись директорией
		size := info.Size()                                             // Получение размера файла/директории
		name := file.Name()                                             // Получение имени файла/директории
		tip := " "                                                       // Тип
		if isDirectory {                                                // Обновление описания типа в зависимости от того, является ли оно директорией
			tip = "Директория"
		} else {
			tip = "Файл"
		}
		formattedSize := formatSize(size)                               // Форматирование размера для вывода
		line := fmt.Sprintf("%-15s %-15s %-30s\n", tip, formattedSize, name) // Форматирование строки для вывода
		fmt.Println(line)                                                // Вывод строки
	}
}

// Функция calculateTotalSize рассчитывает общий размер всех файлов в заданной директории.
func calculateTotalSize(dirPath string) (int64, error) {
	var totalSize int64
	err := filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error { // Перебор директории
		if err != nil {                                                                   // Обработка ошибок при переборе
			return err
		}
		if !info.IsDir() {                                                                 // Пропуск директорий
			totalSize += info.Size()                                                         // Добавление размера текущего файла
		}
		return nil                                                                        // Продолжение перебора
	})
	return totalSize, err                                                               // Возвращение общего размера и возможной ошибки
}

// Функция formatSize форматирует количество байт в человекочитаемую строку.
func formatSize(bytes int64) string {
	switch {
	case bytes >= 1024*1024*1024: // Если размер в гигабайтах или больше
		return fmt.Sprintf("%.2f Гигабайта", float64(bytes)/float64(1024*1024*1024))
	case bytes >= 1024*1024:       // Если размер в мегабайтах
		return fmt.Sprintf("%.2f Мегабайта", float64(bytes)/float64(1024*1024))
	case bytes >= 1024:            // Если размер в килобайтах
		return fmt.Sprintf("%.2f Килобайта", float64(bytes)/float64(1024))
	default:                        // Если размер меньше килобайта
		return fmt.Sprintf("%d Байта", bytes)
	}
}

func main() {
	dirPath := parseFlags()                                                             // Парсинг командных флагов для получения пути к директории

	startTime := time.Now()                                                            // Запись времени начала работы программы
	fmt.Println("Программа выполняется...")

	files, err := listFilesInDirectory(dirPath)                                       // Перечисление всех файлов в директории
	if err != nil {                                                                    // Обработка ошибок при перечислении файлов
		fmt.Printf("ошибка при чтении директории '%s': %v\n", dirPath, err)
		os.Exit(1)
	}

	printFileDetails(files)                                                           // Вывод деталей каждого файла

	totalSize, err := calculateTotalSize(dirPath)                                    // Расчет общего размера всех файлов в директории
	if err != nil {                                                                    // Обработка ошибок при расчете общего размера
		fmt.Printf("ошибка при подсчете общего размера директории '%s': %v\n", dirPath, err)
		os.Exit(1)
	}

	fmt.Printf("\nОбщий размер директории: %s\n", formatSize(totalSize))             // Вывод общего размера

	endTime := time.Now()                                                            // Запись времени окончания работы программы
	duration := endTime.Sub(startTime)                                               // Расчет времени выполнения
	fmt.Printf("Время выполнения программы: %v\n", duration)                       // Вывод времени выполнения
}
