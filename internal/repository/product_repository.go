package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"store-management/internal/datasource"
	"store-management/internal/model"
	"strings"
	"time"
)

type ProductRepository interface {
	CreateProduct(storeId int64, product *model.Product) (int64, error)
	UpdateProduct(storeId int64, product *model.Product) error
	DeleteProduct(storeId, productId int64) error
	FindProduct(storeId, productId int64) (*model.Product, error)
	FindProductsWithPagination(storeId int64, cursor int64, limit int64) ([]*model.Product, error)
	SearchProducts(storeId int64, query string) ([]*model.Product, error)
}
type productRepositoryImpl struct {
	writer      datasource.SQL
	reader      datasource.SQL
	transaction datasource.Transaction
}

func NewProductRepository(writer, reader datasource.SQL, transaction datasource.Transaction) ProductRepository {
	return &productRepositoryImpl{
		writer:      writer,
		reader:      reader,
		transaction: transaction,
	}
}

var initialKoreanRunes = [19]rune{'ㄱ', 'ㄲ', 'ㄴ', 'ㄷ', 'ㄸ', 'ㄹ', 'ㅁ', 'ㅂ', 'ㅃ', 'ㅅ', 'ㅆ', 'ㅇ', 'ㅈ', 'ㅉ', 'ㅊ', 'ㅋ', 'ㅌ', 'ㅍ', 'ㅎ'}

func (p *productRepositoryImpl) extractAbstractKoreanName(text string) string {
	var abstractName strings.Builder
	for _, r := range text {
		if r >= '가' && r <= '힣' {
			initialRune := initialKoreanRunes[(r-'가')/588] // 588 = 21(중성 갯수) * 28(종성 갯수)
			abstractName.WriteRune(initialRune)
		} else {
			abstractName.WriteRune(r)
		}
	}
	return abstractName.String()
}

func (p *productRepositoryImpl) CreateProduct(storeId int64, product *model.Product) (int64, error) {
	tx := p.transaction.MustBegin()
	product.AbstractName = p.extractAbstractKoreanName(product.Name)
	res := tx.MustExec("INSERT INTO product (category, price, cost, name, abstract_name, description, barcode, expiry_date, size) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		product.Category, product.Price, product.Cost, product.Name, product.AbstractName, product.Description, product.Barcode, product.ExpiryDate.Format("2006-01-02 15:04:05"), product.Size)
	lastInsertId, err := res.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	tx.MustExec("INSERT INTO store_product (store_id, product_id) VALUES (?, ?)", storeId, lastInsertId)
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil
}

func (p *productRepositoryImpl) UpdateProduct(storeId int64, product *model.Product) error {
	tx := p.transaction.MustBegin()
	var count int
	if err := tx.Get(&count, "SELECT COUNT(id) FROM store_product WHERE store_id = ? AND product_id = ?", storeId, product.ID); err != nil || count == 0 {
		_ = tx.Rollback()
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return err
			}
			panic(err)
		} else {
			return sql.ErrNoRows
		}
	}

	fields := reflect.ValueOf(product).Elem()
	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i)
		if field.IsZero() {
			continue
		}
		fieldName := fields.Type().Field(i).Tag.Get("db")
		if fieldName == "ID" {
			continue
		}
		fieldValue := field.Interface()
		if fields.Type().Field(i).Type.String() == "time.Time" {
			if fieldValue.(time.Time).Unix() == 0 {
				continue
			}
			fieldValue = fieldValue.(time.Time).Format("2006-01-02 15:04:05")
		}

		_ = tx.MustExec(fmt.Sprintf("UPDATE product SET %s = ? WHERE id = ?", fieldName), fieldValue, product.ID)
		if fieldName == "name" {
			_ = tx.MustExec("UPDATE product SET abstract_name = ? WHERE id = ?", p.extractAbstractKoreanName(fieldValue.(string)), product.ID)
		}
	}

	err := tx.Commit()
	if err != nil {
		panic(err)
	}
	return nil
}

func (p *productRepositoryImpl) DeleteProduct(storeId, productId int64) error {
	tx := p.transaction.MustBegin()
	_, err := tx.Exec("DELETE FROM product WHERE id = ?", productId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	res := tx.MustExec("DELETE FROM store_product WHERE store_id = ? AND product_id = ?", storeId, productId)
	affectedRows, err := res.RowsAffected()
	if affectedRows == 0 || err != nil {
		_ = tx.Rollback()
		return errors.New("product not found")
	}
	return tx.Commit()
}

func (p *productRepositoryImpl) FindProduct(storeId, productId int64) (*model.Product, error) {
	var count int
	err := p.reader.Get(&count, "SELECT count(id) FROM store_product WHERE store_id = ? AND product_id = ?", storeId, productId)
	if err != nil {
		panic(err)
	}
	if count == 0 {
		return nil, datasource.ErrNoRows
	}

	var product model.Product
	err = p.reader.Get(&product, "SELECT * FROM product WHERE id = ?", productId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, datasource.ErrNoRows
		}
		panic(err)
	}
	return &product, nil
}

func (p *productRepositoryImpl) FindProductsWithPagination(storeId int64, cursor int64, limit int64) ([]*model.Product, error) {
	var products []*model.Product

	var cursorCondition string
	if cursor > 0 {
		cursorCondition = "AND p.id < ?"
	}

	query := fmt.Sprintf(`
        SELECT p.* FROM product p
        INNER JOIN store_product sp ON p.id = sp.product_id
        WHERE sp.store_id = ? %s
        ORDER BY p.id DESC
        LIMIT ?
    `, cursorCondition)

	var err error
	if cursor > 0 {
		err = p.reader.Select(&products, query, storeId, cursor, limit)
	} else {
		err = p.reader.Select(&products, query, storeId, limit)
	}

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		panic(err)
	}
	return products, nil
}

func (p *productRepositoryImpl) SearchProducts(storeId int64, keyword string) ([]*model.Product, error) {
	var products []*model.Product
	query := `
        SELECT p.* FROM product p
        INNER JOIN store_product sp ON p.id = sp.product_id
        WHERE sp.store_id = ? AND (name LIKE ? OR abstract_name LIKE ?)
        ORDER BY p.id DESC
    `
	err := p.reader.Select(&products, query, storeId, "%"+keyword+"%", "%"+p.extractAbstractKoreanName(keyword)+"%")
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		panic(err)
	}
	return products, nil
}
