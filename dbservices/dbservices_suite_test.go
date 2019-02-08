package dbservices_test

import (
	"database/sql"
	"testing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	_ "github.com/lib/pq"
)

var realDB *sql.DB

func TestDbservices(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dbservices Suite")
}

var _ = BeforeSuite(func() {
	realDB = initializeDb()
})

var _ = AfterSuite(func() {
	realDB.Close()
})

var _ = AfterEach(func() {
	cleanDB()
})

func cleanDB() {
	_, err := realDB.Exec("TRUNCATE TABLE coupons")
	Expect(err).NotTo(HaveOccurred())
}

func initializeDb() *sql.DB {
	connectionString := "user=testing password=testingtesting123 dbname=coupons_test sslmode=disable"

	db, err := sql.Open("postgres", connectionString)
	Expect(err).NotTo(HaveOccurred())

	return db
}