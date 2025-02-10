package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Column ve Table
type Column struct {
	Name      string `json:"name"`
	DataType  string `json:"data_type"`
	IsPrimary bool   `json:"is_primary"`
}

// Go'daki struct içindeki değişken isimleri JSON içindeki anahtarlara dönüşüyor

//✅ Web API'lerden gelen JSON verisini Go struct'ına çevirmek
//✅ Go struct'larını JSON formatında saklamak veya göndermek
//✅ JSON anahtarlarını özelleştirerek okunaklı hale getirmek

// json:"..." etiketleri, JSON'daki anahtar isimlerini özelleştirmeye yarar.
// json.Marshal() ile struct'ı JSON'a çevirebiliriz.
// json.Unmarshal() ile JSON'u struct'a geri çevirebiliriz.

type Table struct {
	TableName string          `json:"table_name"`
	Columns   []Column        `json:"columns"`
	Rows      [][]interface{} `json:"rows"`
}

// PK oluşturma fonksiyonu
var pkCounter int = 1

func generatePrimaryKey() interface{} {
	pkCounter++
	return pkCounter - 1
}

// JSON dosyasına yeni satır ekleme fonksiyonu
func addRowToJsonFile(newRow []interface{}, filePath string) error {

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("dosya açılırken hata oluştu: %v", err)
	}
	defer file.Close()

	var table Table
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&table); err != nil && err.Error() != "EOF" {
		return fmt.Errorf("JSON verisi okunurken hata oluştu: %v", err)
	}

	pkValue := generatePrimaryKey()
	newRow = append([]interface{}{pkValue}, newRow...)
	table.Rows = append(table.Rows, newRow)

	file.Seek(0, 0)
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(&table); err != nil {
		return fmt.Errorf("JSON verisi yazılırken hata oluştu: %v", err)
	}

	return nil
}

func deleteRowFromJsonFile(idToDelete int, filePath string) error {

	file, err := os.OpenFile(filePath, os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("dosya açılırken hata oluştu: %v", err)
	}
	defer file.Close()

	var table Table
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&table); err != nil {
		return fmt.Errorf("JSON verisi okunurken hata oluştu: %v", err)
	}

	var newRows [][]interface{}
	for _, row := range table.Rows {
		if len(row) > 0 {
			if id, ok := row[0].(float64); ok && int(id) == idToDelete {
				continue // Bu satırı atlıyoruz (silme işlemi)
			}
		}
		newRows = append(newRows, row)
	}
	table.Rows = newRows

	file.Truncate(0) //Dosyanın içeriğini tamamen temizler.

	file.Seek(0, 0)
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(&table); err != nil {
		return fmt.Errorf("JSON verisi yazılırken hata oluştu: %v", err)
	}

	return nil
}

func readTableFromJsonFile(filePath string) (Table, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Table{}, fmt.Errorf("dosya açılırken hata oluştu: %v", err)
	}
	defer file.Close()

	var table Table
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&table); err != nil {
		return Table{}, fmt.Errorf("JSON verisi okunurken hata oluştu: %v", err)
	}

	return table, nil
}

func printTableData(table Table) {
	fmt.Printf("Tablo Adı: %s\n", table.TableName)
	fmt.Println("Sütunlar:")
	for _, column := range table.Columns {
		fmt.Printf(" - %s (%s)\n", column.Name, column.DataType)
	}

	fmt.Println("\nVeriler:")
	for _, row := range table.Rows {
		fmt.Println(row)
	}
}

// Veritabanı dosyası oluşturma fonksiyonu
func createDatabaseJsonFile(filePath string) error {
	if err := os.MkdirAll("Database", os.ModePerm); err != nil {
		log.Fatal("Dizin oluşturulurken hata oluştu:", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	table := Table{
		TableName: "users",
		Columns: []Column{
			{"id", "int", true},
			{"username", "string", false},
			{"email", "string", false},
			{"is_active", "bool", false},
		},
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(&table); err != nil {
		log.Fatal(err)
	}

	return nil
}

func main() {
	filePath := "Database/database.json"

	// Veritabanı dosyasını oluştur
	if err := createDatabaseJsonFile(filePath); err != nil {
		log.Fatal(err)
	}

	// Yeni satırlar ekleyelim
	addRowToJsonFile([]interface{}{"user1", "user1@example.com", true}, filePath)
	addRowToJsonFile([]interface{}{"user2", "user2@example.com", false}, filePath)
	addRowToJsonFile([]interface{}{"user3", "user3@example.com", true}, filePath)

	// Veriyi ekrana yazdır
	fmt.Println("\n--- EKLENEN KULLANICILAR ---")

	table, err := readTableFromJsonFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	printTableData(table)

	// ID=2 olan satırı silelim
	fmt.Println("\n--- ID=2 OLAN KULLANICI SİLİNİYOR... ---")
	if err := deleteRowFromJsonFile(2, filePath); err != nil {
		log.Fatal(err)
	}

	// Güncellenmiş veriyi tekrar ekrana yazdır
	fmt.Println("\n--- GÜNCELLENMİŞ KULLANICILAR ---")
	table, err = readTableFromJsonFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	printTableData(table)

	fmt.Println("\nSilme işlemi tamamlandı!")
}
