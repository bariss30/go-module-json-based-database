package dbmodule

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Column struct {
	Name      string json:"name"
	DataType  string json:"data_type"
	IsPrimary bool   json:"is_primary"
}

type Table struct {
	TableName string          json:"table_name"
	Columns   []Column        json:"columns"
	Rows      [][]interface{} json:"rows"
}

func getUserInput(prompt string) string {
	var input string
	fmt.Print(prompt)
	fmt.Scanln(&input)
	return input
}

// JSON dosyasını oluştur
func CreateDatabaseJsonFile(filePath string, table Table) error {
	if err := os.MkdirAll("Database", os.ModePerm); err != nil {
		log.Fatal("Dizin oluşturulurken hata oluştu:", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonData, err := json.Marshal(table)
	if err != nil {
		return err
	}

	_, err = file.Write(jsonData)
	return err
}

// Kullanıcıdan veri al
func AddRowFromUser(filePath string) error {
	tableFile, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Dosya okunamadı:", err)
		return err
	}

	var table Table
	if err := json.Unmarshal(tableFile, &table); err != nil {
		fmt.Println("JSON parse hatası:", err)
		return err
	}

	primaryKeys := []interface{}{}
	otherValues := []interface{}{}

	fmt.Println("Yeni veri ekleniyor:", table.TableName)

	for _, column := range table.Columns {
		for {
			input := getUserInput(fmt.Sprintf("%s (%s): ", column.Name, column.DataType))

			var value interface{}
			var convErr error

			switch column.DataType {
			case "int":
				value, convErr = strconv.Atoi(input)
			case "float64":
				value, convErr = strconv.ParseFloat(input, 64)
			case "bool":
				value, convErr = strconv.ParseBool(input)
			default:
				value = input
			}

			if convErr == nil {
				if column.IsPrimary {
					primaryKeys = append(primaryKeys, value) // Öncelikli olarak primary keyleri ekle
				} else {
					otherValues = append(otherValues, value)
				}
				break
			} else {
				fmt.Println("Geçersiz giriş, lütfen tekrar girin.")
			}
		}
	}

	newRow := append(primaryKeys, otherValues...) // Primary keyleri önce ekleyip ardından diğer verileri ekleme

	table.Rows = append(table.Rows, newRow)

	updatedJson, err := json.Marshal(table)
	if err != nil {
		fmt.Println("JSON oluşturulamadı:", err)
		return err
	}

	return os.WriteFile(filePath, updatedJson, 0644)
}

// JSON dosyasından tabloyu okuma ve ekrana yazdırma fonksiyonu
func ReadTable(filePath string) error {
	tableFile, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var table Table
	if err := json.Unmarshal(tableFile, &table); err != nil {
		return err
	}

	fmt.Println("\nTablo Adı:", table.TableName)
	fmt.Println(strings.Repeat("-", 40))

	for _, column := range table.Columns {
		fmt.Printf("%-15s", column.Name)
	}
	fmt.Println("\n" + strings.Repeat("-", 40))

	for i, row := range table.Rows {
		fmt.Printf("%-5d", i+1)
		for _, value := range row {
			fmt.Printf("%-15v", value)
		}
		fmt.Println()
	}

	fmt.Println(strings.Repeat("-", 40))
	return nil
}

// Veri silme fonksiyonu
func DeleteRowFromUser(filePath string) error {
	tableFile, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var table Table
	if err := json.Unmarshal(tableFile, &table); err != nil {
		return err
	}

	if len(table.Rows) == 0 {
		fmt.Println("Tabloda silinecek veri bulunmamaktadır.")
		return nil
	}

	fmt.Println("Silinecek satırı seçin:")
	for i := 0; i < len(table.Rows); i++ {
		fmt.Printf("%d: ", i+1)
		for _, value := range table.Rows[i] {
			fmt.Printf("%v ", value)
		}
		fmt.Println()
	}

	rowToDelete, err := strconv.Atoi(getUserInput("Silmek istediğiniz satır numarasını girin: "))
	if err != nil || rowToDelete < 1 || rowToDelete > len(table.Rows) {
		fmt.Println("Geçersiz satır numarası.")
		return nil
	}

	// Satırı sil
	table.Rows = append(table.Rows[:rowToDelete-1], table.Rows[rowToDelete:]...)

	updatedJson, err := json.Marshal(table)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, updatedJson, 0644)
}

// Veri güncelleme fonksiyonu
func UpdateRowFromUser(filePath string) error {
	tableFile, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var table Table
	if err := json.Unmarshal(tableFile, &table); err != nil {
		return err
	}

	if len(table.Rows) == 0 {
		fmt.Println("Tabloda güncellenecek veri bulunmamaktadır.")
		return nil
	}

	fmt.Println("Güncellenecek satırı seçin:")
	for i := 0; i < len(table.Rows); i++ {
		fmt.Printf("%d: ", i+1)
		for _, value := range table.Rows[i] {
			fmt.Printf("%v ", value)
		}
		fmt.Println()
	}

	rowToUpdate, err := strconv.Atoi(getUserInput("Güncellemek istediğiniz satır numarasını girin: "))
	if err != nil || rowToUpdate < 1 || rowToUpdate > len(table.Rows) {
		fmt.Println("Geçersiz satır numarası.")
		return nil
	}

	// Güncelleme işlemi
	fmt.Println("Hangi sütunu güncellemek istersiniz?")
	for i, column := range table.Columns {
		fmt.Printf("%d: %s\n", i+1, column.Name)
	}
	columnToUpdate, err := strconv.Atoi(getUserInput("Güncellemek istediğiniz sütun numarasını girin: "))
	if err != nil || columnToUpdate < 1 || columnToUpdate > len(table.Columns) {
		fmt.Println("Geçersiz sütun numarası.")
		return nil
	}

	newValue := getUserInput(fmt.Sprintf("Yeni değeri girin (%s): ", table.Columns[columnToUpdate-1].DataType))

	// Sütundaki değeri güncelle
	var updatedValue interface{}
	switch table.Columns[columnToUpdate-1].DataType {
	case "int":
		updatedValue, _ = strconv.Atoi(newValue)
	case "float64":
		updatedValue, _ = strconv.ParseFloat(newValue, 64)
	case "bool":
		updatedValue, _ = strconv.ParseBool(newValue)
	default:
		updatedValue = newValue
	}

	table.Rows[rowToUpdate-1][columnToUpdate-1] = updatedValue

	updatedJson, err := json.Marshal(table)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, updatedJson, 0644)
}

func CreateTable() {
	tableName := getUserInput("Tablo adını girin: ")
	filePath := fmt.Sprintf("Database/%s.json", tableName) //string format

	if _, err := os.Stat(filePath); err == nil { // https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go
		action := getUserInput("1: Verileri görüntüle\n2: Yeni satır ekle\n3: Satır sil\n4: Satır güncelle\nSeçiminiz: ")

		if action == "1" {
			if err := ReadTable(filePath); err != nil {
				log.Fatal(err)
			}
		} else if action == "2" {
			if err := AddRowFromUser(filePath); err != nil {
				log.Fatal(err)
			}
		} else if action == "3" {
			if err := DeleteRowFromUser(filePath); err != nil {
				log.Fatal(err)
			}
		} else if action == "4" {
			if err := UpdateRowFromUser(filePath); err != nil {
				log.Fatal(err)
			}
		}
		return
	}

	// Yeni tablo oluşturma
	numColumns, _ := strconv.Atoi(getUserInput("Kaç sütun olacak?: ")) //https://stackoverflow.com/questions/4278430/convert-string-to-integer-type-in-go
	columns := []Column{}

	for i := 0; i < numColumns; i++ {
		colName := getUserInput(fmt.Sprintf("%d. sütun adı: ", i+1))
		colType := getUserInput("Veri tipi (string, int, float64, bool): ")
		isPrimary := getUserInput("Birincil anahtar mı? (evet/hayır): ") == "evet"
		columns = append(columns, Column{Name: colName, DataType: colType, IsPrimary: isPrimary})
	}

	table := Table{TableName: tableName, Columns: columns, Rows: [][]interface{}{}}

	if err := CreateDatabaseJsonFile(filePath, table); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Tablo oluşturuldu!")

	for {
		if getUserInput("Yeni satır eklemek ister misiniz? (evet/hayır): ") != "evet" {
			break
		}
		if err := AddRowFromUser(filePath); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("İşlem tamamlandı! Tablonuz kaydedildi.")
} 