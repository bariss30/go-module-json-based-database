# Go Module: JSON Based Database

Bu proje, JSON tabanlı bir veritabanı modülü oluşturmaya yönelik bir Go modülüdür.

## Kurulum

Öncelikle, projenizi oluşturmak için aşağıdaki adımları takip edin:

```sh
# Yeni bir proje klasörü oluşturun
mkdir test
cd test

# Go modülünü başlatın
go mod init go-module-json-based-database

# Gerekli dosyaları oluşturun
touch main.go
```

Daha sonra, modülü projenize dahil etmek için aşağıdaki komutu çalıştırın:

```sh
go get github.com/bariss30/go-module-json-based-database/dbmodule
```

## Kullanım

Aşağıdaki örnek kod, `dbmodule` modülünü nasıl kullanacağınızı göstermektedir:

```go
package main

import (
	"fmt"
	dbmodule "github.com/bariss30/go-module-json-based-database/dbmodule"
)

func main() {
	fmt.Println("Veritabanı oluşturuluyor...")
	dbmodule.CreateTable() // CreateTable fonksiyonunu çağırıyoruz
}
```

