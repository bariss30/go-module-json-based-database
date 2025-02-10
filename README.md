# Go JSON Veritabanı Modülü

Bu modül, JSON tabanlı veritabanı işlemleri sağlar. Kullanıcılar veritabanı oluşturabilir, tablo ekleyebilir ve veri işlemleri yapabilir.

## Özellikler

### Veritabanı İşlemleri
- Veritabanı oluşturma.
- Her veritabanı, kendi klasöründe saklanır.

### Tablo İşlemleri
- Tablo oluşturma.
- Veri ekleme, silme, güncelleme ve görüntüleme.
- Her tablo bir JSON dosyasına kaydedilir.

### Veri Türleri
- INT, STRING, BOOL veri türleri desteklenir.
- Her tablonun bir **Primary Key (PK)** olmalıdır.

### Kullanılan Yapılar
- Struct ve Interface kullanımı.

## Fonksiyonlar

- **`createDatabaseJsonFile(filePath string) error`**: Veritabanı oluşturur.
- **`addRowToJsonFile(newRow []interface{}, filePath string) error`**: Tabloya veri ekler.
- **`deleteRowFromJsonFile(idToDelete int, filePath string) error`**: Tablo verisini siler.
- **`readTableFromJsonFile(filePath string) (Table, error)`**: Tabloyu okur.
- **`printTableData(table Table)`**: Tabloyu yazdırır.
- **`generatePrimaryKey() interface{}`**: Birincil anahtar oluşturur.

## Kullanım Örneği

```go
filePath := "Database/database.json"
err := createDatabaseJsonFile(filePath)
if err != nil {
	log.Fatal(err)
}


addRowToJsonFile([]interface{}{"user1", "user1@example.com", true}, filePath)
addRowToJsonFile([]interface{}{"user2", "user2@example.com", false}, filePath)
addRowToJsonFile([]interface{}{"user3", "user3@example.com", true}, filePath)


table, err := readTableFromJsonFile(filePath)
if err != nil {
	log.Fatal(err)
}
printTableData(table)


if err := deleteRowFromJsonFile(2, filePath); err != nil {
	log.Fatal(err)
}



table, err = readTableFromJsonFile(filePath)
if err != nil {
	log.Fatal(err)
}
printTableData(table)
